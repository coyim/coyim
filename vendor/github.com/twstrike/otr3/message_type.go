package otr3

import "bytes"

type messageTypeGuess int

const (
	msgGuessNotOTR messageTypeGuess = iota
	msgGuessTaggedPlaintext
	msgGuessQuery
	msgGuessDHCommit
	msgGuessDHKey
	msgGuessRevealSig
	msgGuessSignature
	msgGuessV1KeyExch
	msgGuessData
	msgGuessError
	msgGuessFragment
	msgGuessUnknown
)

func guessMessageType(msg []byte) messageTypeGuess {
	if bytes.HasPrefix(msg, []byte("?OTR")) {
		switch {
		case bytes.HasPrefix(msg, []byte("?OTR:AAMC")):
			return msgGuessDHCommit
		case bytes.HasPrefix(msg, []byte("?OTR:AAIC")):
			return msgGuessDHCommit

		case bytes.HasPrefix(msg, []byte("?OTR:AAMK")):
			return msgGuessDHKey
		case bytes.HasPrefix(msg, []byte("?OTR:AAIK")):
			return msgGuessDHKey

		case bytes.HasPrefix(msg, []byte("?OTR:AAMR")):
			return msgGuessRevealSig
		case bytes.HasPrefix(msg, []byte("?OTR:AAIR")):
			return msgGuessRevealSig

		case bytes.HasPrefix(msg, []byte("?OTR:AAMS")):
			return msgGuessSignature
		case bytes.HasPrefix(msg, []byte("?OTR:AAIS")):
			return msgGuessSignature

		case bytes.HasPrefix(msg, []byte("?OTR:AAED")):
			return msgGuessData
		case bytes.HasPrefix(msg, []byte("?OTR:AAID")):
			return msgGuessData
		case bytes.HasPrefix(msg, []byte("?OTR:AAMD")):
			return msgGuessData

		case bytes.HasPrefix(msg, []byte("?OTR?")):
			return msgGuessQuery
		case bytes.HasPrefix(msg, []byte("?OTRv")):
			return msgGuessQuery

		case bytes.HasPrefix(msg, []byte("?OTR:AAEK")):
			return msgGuessV1KeyExch

		case bytes.HasPrefix(msg, []byte("?OTR Error:")):
			return msgGuessError

		case bytes.HasPrefix(msg, []byte("?OTR|")):
			return msgGuessFragment
		case bytes.HasPrefix(msg, []byte("?OTR,")):
			return msgGuessFragment
		}
		return msgGuessUnknown
	}
	if bytes.Index(msg, whitespaceTagHeader) != -1 {
		return msgGuessTaggedPlaintext
	}
	return msgGuessNotOTR
}
