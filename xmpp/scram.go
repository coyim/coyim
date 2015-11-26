package xmpp

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/transform"
)

type scramClient struct {
	user     string
	password string

	clientNonce   string
	serverNonce   string
	salt          []byte
	numIterations int

	firstMessageBare   string
	serverFirstMessage string
}

var (
	errPasswordContainsInvalid = errors.New("password contains invalid characters")
)

func (s *scramClient) normalizedPassword() (string, error) {
	t := transform.NewReader(
		strings.NewReader(s.password), SASLprep)
	r := bufio.NewReader(t)

	normalized, _, err := r.ReadLine()
	return string(normalized), err
}

func (s *scramClient) firstMessage() string {
	if s.clientNonce == "" {
		panic("empty client nonce")
	}

	s.firstMessageBare = fmt.Sprintf("n=%s,r=%s", s.user, s.clientNonce)

	return base64.StdEncoding.EncodeToString(
		[]byte("n,," + s.firstMessageBare),
	)
}

func (s *scramClient) receive(msg string) (err error) {
	encoding := base64.StdEncoding

	var dec []byte
	dec, err = encoding.DecodeString(msg)
	if err != nil {
		return
	}

	s.serverFirstMessage = string(dec)

	for _, p := range strings.Split(s.serverFirstMessage, ",") {
		switch {
		case strings.HasPrefix(p, "r="):
			s.serverNonce = strings.TrimPrefix(p, "r=")
		case strings.HasPrefix(p, "s="):
			saltParam := strings.TrimPrefix(p, "s=")
			if s.salt, err = encoding.DecodeString(saltParam); err != nil {
				return
			}
		case strings.HasPrefix(p, "i="):
			countParam := strings.TrimPrefix(p, "i=")
			if s.numIterations, err = strconv.Atoi(countParam); err != nil {
				return
			}
		}
	}

	if !strings.HasPrefix(s.serverNonce, s.clientNonce) {
		err = ErrAuthenticationFailed
	}

	return
}

func (s *scramClient) reply() (string, string, error) {
	finalMessageBare := "c=biws,r=" + s.serverNonce
	normPass, _ := s.normalizedPassword()
	saltedPassword := pbkdf2.Key([]byte(normPass), s.salt, s.numIterations, sha1.Size, sha1.New)

	clientMAC := hmac.New(sha1.New, saltedPassword)
	clientMAC.Write([]byte("Client Key"))
	clientKey := clientMAC.Sum(nil)
	storedKey := sha1.Sum(clientKey)

	serverMAC := hmac.New(sha1.New, saltedPassword)
	serverMAC.Write([]byte("Server Key"))
	serverKey := serverMAC.Sum(nil)

	authMessage := s.firstMessageBare + "," + s.serverFirstMessage + "," + finalMessageBare

	clientSignatureMAC := hmac.New(sha1.New, storedKey[:])
	clientSignatureMAC.Write([]byte(authMessage))
	clientSignature := clientSignatureMAC.Sum(nil)

	clientProof := make([]byte, sha1.Size)
	for i := range clientProof {
		clientProof[i] = clientKey[i] ^ clientSignature[i]
	}

	serverSignatureMAC := hmac.New(sha1.New, serverKey[:])
	serverSignatureMAC.Write([]byte(authMessage))
	serverSignature := serverSignatureMAC.Sum(nil)

	encoding := base64.StdEncoding
	finalMessage := finalMessageBare + ",p=" + encoding.EncodeToString(clientProof)

	serverAuthentication := "v=" + encoding.EncodeToString(serverSignature)

	return encoding.EncodeToString([]byte(finalMessage)),
		encoding.EncodeToString([]byte(serverAuthentication)),
		nil
}
