package otr3

import (
	"math/big"
	"testing"
)

func Test_generateSMP3_generatesLongerValuesForR4WithProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	smp, err := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.r4, fixtureLong1)
	assertDeepEquals(t, err, nil)
}

func Test_generateSMP3_generatesLongerValuesForR5WithProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.r5, fixtureLong2)
}

func Test_generateSMP3_generatesLongerValuesForR6WithProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.r6, fixtureLong3)
}

func Test_generateSMP3_generatesLongerValuesForR7WithProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.r7, fixtureLong4)
}

func Test_generateSMP3_generatesShorterValuesForR4WithProtocolV2(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.r4, fixtureShort1)
}

func Test_generateSMP3_computesPaCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.pa, fixtureMessage3().pa)
}

func Test_generateSMP3_computesQaCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.qa, fixtureMessage3().qa)
}

func Test_generateSMP3_computesPaPbCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.papb, fixtureSmp3().papb)
}

func Test_generateSMP3_computesQaQbCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.qaqb, fixtureSmp3().qaqb)
}

func Test_generateSMP3_storesG3b(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.g3b, fixtureMessage2().g3b)
}

func Test_generateSMP3_computesCPCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.cp, fixtureMessage3().cp)
}

func Test_generateSMP3_computesD5Correctly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.d5, fixtureMessage3().d5)
}

func Test_generateSMP3_computesD6Correctly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.d6, fixtureMessage3().d6)
}

func Test_generateSMP3_computesRaCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.ra, fixtureMessage3().ra)
}

func Test_generateSMP3_computesCrCorrectly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.cr, fixtureMessage3().cr)
}

func Test_generateSMP3_computesD7Correctly(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, smp.msg.d7, fixtureMessage3().d7)
}

func Test_generateSMP3Parameters_returnsAnErrorIfThereIsntRandomnessToGenerate_r4(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP3Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP3Parameters_returnsAnErrorIfThereIsntRandomnessToGenerate_r5(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP3Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP3Parameters_returnsAnErrorIfThereIsntRandomnessToGenerate_r6(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP3Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP3Parameters_returnsAnErrorIfThereIsntRandomnessToGenerate_r7(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP3Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP3Parameters_returnsOKIfThereIsEnoughRandomnessToGenerateBlindingFactors(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
	})).generateSMP3Parameters()
	assertDeepEquals(t, err, nil)
}

func Test_generateSMP3_returnsAnErrorIfThereIsNotEnoughRandomnessForBlinding(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP3(fixtureSecret(), *fixtureSmp1(), fixtureMessage2())
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_verifySMP3_failsIfPaIsNotInTheGroupForProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	err := otr.verifySMP3(fixtureSmp2(), smp3Message{pa: big.NewInt(1)})
	assertDeepEquals(t, err, newOtrError("Pa is an invalid group element"))
}

func Test_verifySMP3_failsIfQaIsNotInTheGroupForProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	err := otr.verifySMP3(fixtureSmp2(), smp3Message{
		pa: big.NewInt(2),
		qa: big.NewInt(1),
	})
	assertDeepEquals(t, err, newOtrError("Qa is an invalid group element"))
}

func Test_verifySMP3_failsIfRaIsNotInTheGroupForProtocolV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	err := otr.verifySMP3(fixtureSmp2(), smp3Message{
		pa: big.NewInt(2),
		qa: big.NewInt(2),
		ra: big.NewInt(1),
	})
	assertDeepEquals(t, err, newOtrError("Ra is an invalid group element"))
}

func Test_verifySMP3_succeedsForValidZKPS(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	err := otr.verifySMP3(fixtureSmp2(), fixtureMessage3())
	assertDeepEquals(t, err, nil)
}

func Test_verifySMP3_failsIfCpIsNotAValidZKP(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	m := fixtureMessage3()
	m.cp = sub(m.cp, big.NewInt(1))
	err := otr.verifySMP3(fixtureSmp2(), m)
	assertDeepEquals(t, err, newOtrError("cP is not a valid zero knowledge proof"))
}

func Test_verifySMP3_failsIfCrIsNotAValidZKP(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	m := fixtureMessage3()
	m.cr = sub(m.cr, big.NewInt(1))
	err := otr.verifySMP3(fixtureSmp2(), m)
	assertDeepEquals(t, err, newOtrError("cR is not a valid zero knowledge proof"))
}
