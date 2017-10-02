package scram

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/transform"

	"github.com/coyim/coyim/sasl"
)

type state interface {
	challenge(sasl.Token, sasl.Properties, sasl.AttributeValuePairs) (state, sasl.Token, error)
}

// RFC 5802, section 5

type firstMessage struct{}

func (c firstMessage) challenge(_ sasl.Token, props sasl.Properties, pairs sasl.AttributeValuePairs) (state, sasl.Token, error) {
	user, ok := props[sasl.AuthID]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.AuthID}
	}

	clientNonce, ok := props[sasl.ClientNonce]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.ClientNonce}
	}

	bare := fmt.Sprintf("n=%s,r=%s", user, clientNonce)
	t := sasl.Token("n,," + bare)

	return clientFinalMessage{[]byte(bare)}, t, nil
}

type clientFinalMessage struct {
	firstMessageBare []byte
}

func (c clientFinalMessage) challenge(serverMessage sasl.Token, props sasl.Properties, pairs sasl.AttributeValuePairs) (state, sasl.Token, error) {
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
		return c, nil, sasl.PropertyMissingError{sasl.Password}
	}

	clientNonce, ok := props[sasl.ClientNonce]
	if !ok {
		return c, nil, sasl.PropertyMissingError{sasl.ClientNonce}
	}

	if !strings.HasPrefix(serverNonce, clientNonce) {
		return c, nil, errors.New("nonce mismatch")
	}

	finalMessageBare := []byte("c=biws,r=" + serverNonce)
	normPass, _ := c.normalizedPassword(password)
	saltedPassword := pbkdf2.Key([]byte(normPass), saltToken, numIterations, sha1.Size, sha1.New)

	return c.compose(saltedPassword, finalMessageBare, serverMessage)
}

func (c clientFinalMessage) compose(saltedPassword, finalMessageBare, serverFirstMessage []byte) (state, sasl.Token, error) {
	clientMAC := hmac.New(sha1.New, saltedPassword)
	clientMAC.Write([]byte("Client Key"))
	clientKey := clientMAC.Sum(nil)
	storedKey := sha1.Sum(clientKey)

	serverMAC := hmac.New(sha1.New, saltedPassword)
	serverMAC.Write([]byte("Server Key"))
	serverKey := serverMAC.Sum(nil)

	authMessage := bytes.Join([][]byte{
		c.firstMessageBare,
		serverFirstMessage,
		finalMessageBare,
	}, []byte(","))

	clientSignatureMAC := hmac.New(sha1.New, storedKey[:])
	clientSignatureMAC.Write(authMessage)
	clientSignature := clientSignatureMAC.Sum(nil)

	clientProof := make([]byte, sha1.Size)
	for i := range clientProof {
		clientProof[i] = clientKey[i] ^ clientSignature[i]
	}

	serverSignatureMAC := hmac.New(sha1.New, serverKey[:])
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

	return authenticateServer{serverAuthentication}, sasl.Token(finalMessage), nil
}

func (c clientFinalMessage) normalizedPassword(password string) (string, error) {
	t := transform.NewReader(
		strings.NewReader(password), sasl.Stringprep)
	r := bufio.NewReader(t)

	normalized, _, err := r.ReadLine()
	return string(normalized), err
}

type authenticateServer struct {
	serverAuthentication []byte
}

func (c authenticateServer) challenge(t sasl.Token, _ sasl.Properties, _ sasl.AttributeValuePairs) (state, sasl.Token, error) {
	if subtle.ConstantTimeCompare(t, c.serverAuthentication) != 1 {
		return c, nil, errors.New("server signature mismatch")
	}

	return finished{}, nil, nil
}

type finished struct{}

func (finished) challenge(sasl.Token, sasl.Properties, sasl.AttributeValuePairs) (state, sasl.Token, error) {
	return finished{}, nil, nil
}
