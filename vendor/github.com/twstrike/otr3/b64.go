package otr3

import "encoding/base64"

func b64encode(msg []byte) []byte {
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(b64, msg)
	return b64
}

func b64decode(inp []byte) ([]byte, error) {
	msg := make([]byte, base64.StdEncoding.DecodedLen(len(inp)))
	msgLen, err := base64.StdEncoding.Decode(msg, inp)
	return msg[:msgLen], err
}
