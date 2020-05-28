package scram

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/subtle"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/transform"

	"github.com/coyim/coyim/sasl"
)

type state interface {
	next(sasl.Token, sasl.Properties, sasl.AttributeValuePairs) (state, sasl.Token, error)
	finished() bool
}

// RFC 5802, section 5

type start struct {
	hash     func() hash.Hash
	hashSize int
}

func (c start) finished() bool {
	return false
}

// next for start will send the client-first-message to the server if successful
func (c start) next(_ sasl.Token, props sasl.Properties, pairs sasl.AttributeValuePairs) (state, sasl.Token, error) {
	user, ok := props[sasl.AuthID]
	if !ok {
		return c, nil, sasl.PropertyMissingError{Property: sasl.AuthID}
	}

	clientNonce, ok := props[sasl.ClientNonce]
	if !ok {
		return c, nil, sasl.PropertyMissingError{Property: sasl.ClientNonce}
	}

	bare := fmt.Sprintf("n=%s,r=%s", user, clientNonce)
	t := sasl.Token("n,," + bare)

	return expectingServerFirstMessage{[]byte(bare), c.hash, c.hashSize}, t, nil
}

type expectingServerFirstMessage struct {
	firstMessageBare []byte
	hash             func() hash.Hash
	hashSize         int
}

func (c expectingServerFirstMessage) finished() bool {
	return false
}

func (c expectingServerFirstMessage) next(serverMessage sasl.Token, props sasl.Properties, pairs sasl.AttributeValuePairs) (state, sasl.Token, error) {
	serverNonce, ok := pairs["r"]
	if !ok {
		return c, nil, sasl.ErrMissingParameter
	}

	salt, ok := pairs["s"]
	if !ok {
		return c, nil, sasl.ErrMissingParameter
	}

	saltToken, err := sasl.DecodeToken([]byte(salt))
	if err != nil {
		return c, nil, err
	}

	count, ok := pairs["i"]
	if !ok {
		return c, nil, sasl.ErrMissingParameter
	}

	numIterations, err := strconv.Atoi(count)
	if err != nil {
		return c, nil, err
	}

	password, ok := props[sasl.Password]
	if !ok {
		return c, nil, sasl.PropertyMissingError{Property: sasl.Password}
	}

	clientNonce, ok := props[sasl.ClientNonce]
	if !ok {
		return c, nil, sasl.PropertyMissingError{Property: sasl.ClientNonce}
	}

	if !strings.HasPrefix(serverNonce, clientNonce) {
		return c, nil, errors.New("nonce mismatch")
	}

	finalMessageBare := []byte("c=biws,r=" + serverNonce)
	normPass, _ := c.normalizedPassword(password)
	saltedPassword := pbkdf2.Key([]byte(normPass), saltToken, numIterations, c.hashSize, c.hash)

	return c.compose(saltedPassword, finalMessageBare, serverMessage)
}

func (c expectingServerFirstMessage) compose(saltedPassword, finalMessageBare, serverFirstMessage []byte) (state, sasl.Token, error) {
	clientMAC := hmac.New(c.hash, saltedPassword)
	clientMAC.Write([]byte("Client Key"))
	clientKey := clientMAC.Sum(nil)
	storedKeyHash := c.hash()
	storedKeyHash.Write(clientKey)
	storedKey := storedKeyHash.Sum(nil)

	serverMAC := hmac.New(c.hash, saltedPassword)
	serverMAC.Write([]byte("Server Key"))
	serverKey := serverMAC.Sum(nil)

	authMessage := bytes.Join([][]byte{
		c.firstMessageBare,
		serverFirstMessage,
		finalMessageBare,
	}, []byte(","))

	clientSignatureMAC := hmac.New(c.hash, storedKey)
	clientSignatureMAC.Write(authMessage)
	clientSignature := clientSignatureMAC.Sum(nil)

	clientProof := make([]byte, c.hashSize)
	for i := range clientProof {
		clientProof[i] = clientKey[i] ^ clientSignature[i]
	}

	serverSignatureMAC := hmac.New(c.hash, serverKey[:])
	serverSignatureMAC.Write(authMessage)
	serverSignature := serverSignatureMAC.Sum(nil)

	p := []byte(",p=")
	encodedProf := sasl.Token(clientProof).Encode()
	finalMessage := make([]byte, len(finalMessageBare)+len(p)+len(encodedProf))

	n := copy(finalMessage[:len(finalMessageBare)], finalMessageBare)
	n = n + copy(finalMessage[n:n+len(p)], p)
	copy(finalMessage[n:n+len(encodedProf)], encodedProf)

	encodedServerSignature := sasl.Token(serverSignature).Encode()
	serverAuthentication := append([]byte("v="), encodedServerSignature...)

	return expectingServerFinalMessage{serverAuthentication}, sasl.Token(finalMessage), nil
}

func (c expectingServerFirstMessage) normalizedPassword(password string) (string, error) {
	t := transform.NewReader(
		strings.NewReader(password), sasl.Stringprep)
	r := bufio.NewReader(t)

	normalized, _, err := r.ReadLine()
	return string(normalized), err
}

type expectingServerFinalMessage struct {
	serverAuthentication []byte
}

func (c expectingServerFinalMessage) finished() bool {
	return false
}

func (c expectingServerFinalMessage) next(t sasl.Token, _ sasl.Properties, _ sasl.AttributeValuePairs) (state, sasl.Token, error) {
	if subtle.ConstantTimeCompare(t, c.serverAuthentication) != 1 {
		return c, nil, errors.New("server signature mismatch")
	}

	return finished{}, nil, nil
}

type finished struct{}

func (c finished) finished() bool {
	return true
}

func (finished) next(sasl.Token, sasl.Properties, sasl.AttributeValuePairs) (state, sasl.Token, error) {
	return finished{}, nil, nil
}
