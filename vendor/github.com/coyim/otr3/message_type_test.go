package otr3

import "testing"

func Test_guessMessageType_correctlyIdentifiesAMessageWithWhitespaceTags(t *testing.T) {
	msg := append([]byte("Hello world"), whitespaceTagHeader...)

	assertEquals(t, guessMessageType(msg), msgGuessTaggedPlaintext)
}

func Test_guessMessageType_correctlyIdentifiesAMessageWithNoOTRContent(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("Hello world")), msgGuessNotOTR)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeADHCommit(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR:AAMC")), msgGuessDHCommit)
	assertEquals(t, guessMessageType([]byte("?OTR:AAIC")), msgGuessDHCommit)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeADHKey(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR:AAMK")), msgGuessDHKey)
	assertEquals(t, guessMessageType([]byte("?OTR:AAIK")), msgGuessDHKey)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeARevealSig(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR:AAMR")), msgGuessRevealSig)
	assertEquals(t, guessMessageType([]byte("?OTR:AAIR")), msgGuessRevealSig)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeASignature(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR:AAMS")), msgGuessSignature)
	assertEquals(t, guessMessageType([]byte("?OTR:AAIS")), msgGuessSignature)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeAData(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR:AAMD")), msgGuessData)
	assertEquals(t, guessMessageType([]byte("?OTR:AAID")), msgGuessData)
	assertEquals(t, guessMessageType([]byte("?OTR:AAED")), msgGuessData)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeAQuery(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR?")), msgGuessQuery)
	assertEquals(t, guessMessageType([]byte("?OTRv")), msgGuessQuery)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeAV1KeyExch(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR:AAEK")), msgGuessV1KeyExch)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeAnError(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR Error:")), msgGuessError)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeSomethingElse(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR Weird:")), msgGuessUnknown)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeV3Fragment(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR|")), msgGuessFragment)
}

func Test_guessMessageType_correctlyIdentifiesAMessageThatLooksLikeV2Fragment(t *testing.T) {
	assertEquals(t, guessMessageType([]byte("?OTR,")), msgGuessFragment)
}
