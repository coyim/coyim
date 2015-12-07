package sasl

import "encoding/base64"

var encoding = base64.StdEncoding

// Token represents a SASL token exchanged by client and server.
type Token []byte

// Encode encoded the token using base64
func (t Token) Encode() []byte {
	l := encoding.EncodedLen(len(t))
	ret := make([]byte, l)
	encoding.Encode(ret[:], t)

	return ret
}

// DecodeToken decodes a base64 encoded token
func DecodeToken(in []byte) (Token, error) {
	l := encoding.DecodedLen(len(in))
	ret := make([]byte, l)
	n, err := encoding.Decode(ret, in)

	return Token(ret[0:n]), err
}

func (t Token) String() string {
	return string(t)
}
