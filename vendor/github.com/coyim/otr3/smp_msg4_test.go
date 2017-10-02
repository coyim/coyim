package otr3

import (
	"math/big"
	"testing"
)

func Test_generateSMP4_generatesLongerValuesForR7WithProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	smp, err := otr.generateSMP4(fixtureSecret(), *fixtureSmp2(), fixtureMessage3())
	assertDeepEquals(t, smp.r7, fixtureLong1)
	assertDeepEquals(t, err, nil)
}

func Test_generateSMP4Parameters_returnsAnErrorIfThereIsntEnoughRandomnessToGenerateBlindingFactor(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP4Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP4_returnsAnErrorIfGenerationOfFourthParametersFails(t *testing.T) {
	otr := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	}))
	_, err := otr.generateSMP4(fixtureSecret(), *fixtureSmp2(), fixtureMessage3())
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP4_generatesShorterValuesForR7WithProtocolV3(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP4(fixtureSecret(), *fixtureSmp2(), fixtureMessage3())
	assertDeepEquals(t, smp.r7, fixtureShort1)
}

func Test_generateSMP4_computesRbCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP4(fixtureSecret(), *fixtureSmp2(), fixtureMessage3())
	assertDeepEquals(t, smp.msg.rb, fixtureMessage4().rb)
}

func Test_generateSMP4_computesCrCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP4(fixtureSecret(), *fixtureSmp2(), fixtureMessage3())
	assertDeepEquals(t, smp.msg.cr, fixtureMessage4().cr)
}

func Test_generateSMP4_computesD7Correctly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP4(fixtureSecret(), *fixtureSmp2(), fixtureMessage3())
	assertDeepEquals(t, smp.msg.d7, fixtureMessage4().d7)
}

func Test_verifySMP4_succeedsForValidZKPS(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	err := otr.verifySMP4(fixtureSmp3(), fixtureMessage4())
	assertDeepEquals(t, err, nil)
}

func Test_verifySMP4_failsIfRbIsNotInTheGroupForProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	err := otr.verifySMP4(fixtureSmp3(), smp4Message{rb: big.NewInt(1)})
	assertDeepEquals(t, err, newOtrError("Rb is an invalid group element"))
}

func Test_verifySMP4_failsIfCrIsNotACorrectZKP(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	m := fixtureMessage4()
	m.cr = sub(m.cr, big.NewInt(1))
	err := otr.verifySMP4(fixtureSmp3(), m)
	assertDeepEquals(t, err, newOtrError("cR is not a valid zero knowledge proof"))
}
