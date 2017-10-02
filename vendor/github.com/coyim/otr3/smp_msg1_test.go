package otr3

import (
	"math/big"
	"testing"
)

func Test_generatesLongerAandRValuesForOtrV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	smp, err := otr.generateSMP1()
	assertDeepEquals(t, smp.a2, fixtureLong1)
	assertDeepEquals(t, smp.a3, fixtureLong2)
	assertDeepEquals(t, smp.r2, fixtureLong3)
	assertDeepEquals(t, smp.r3, fixtureLong4)
	assertDeepEquals(t, err, nil)
}

func Test_generateSMP1Parameters_ReturnsErrorIfThereIsntEnoughRandomnessForA2(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b"})).generateSMP1Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP1_ReturnsErrorIfGenerateInitialParametersDoesntWork(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b"})).generateSMP1()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP1Parameters_ReturnsErrorIfThereIsntEnoughRandomnessForA3(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP1Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP1Parameters_ReturnsErrorIfThereIsntEnoughRandomnessForR2(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP1Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generateSMP1Parameters_ReturnsErrorIfThereIsntEnoughRandomnessForR3(t *testing.T) {
	_, err := newConversation(otrV2{}, fixedRand([]string{
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b8b",
		"1a2a3a4a5a6a7a8a1b2b3b4b5b6b7b",
	})).generateSMP1Parameters()
	assertDeepEquals(t, err, errShortRandomRead)
}

func Test_generatesShorterAandRValuesForOtrV2(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP1()
	assertDeepEquals(t, smp.a2, fixtureShort1)
	assertDeepEquals(t, smp.a3, fixtureShort2)
	assertDeepEquals(t, smp.r2, fixtureShort3)
	assertDeepEquals(t, smp.r3, fixtureShort4)
}

func Test_computesG2aAndG3aCorrectlyForOtrV3(t *testing.T) {
	otr := newConversation(otrV3{}, fixtureRand())
	smp, _ := otr.generateSMP1()
	assertDeepEquals(t, smp.msg.g2a, fixtureMessage1v3().g2a)
	assertDeepEquals(t, smp.msg.g3a, fixtureMessage1v3().g3a)
}

func Test_computesG2aAndG3aCorrectlyForOtrV2(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP1()
	assertDeepEquals(t, smp.msg.g2a, fixtureMessage1().g2a)
	assertDeepEquals(t, smp.msg.g3a, fixtureMessage1().g3a)
}

func Test_computesC2AndD2CorrectlyForOtrV2(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP1()
	assertDeepEquals(t, smp.msg.c2, fixtureMessage1().c2)
	assertDeepEquals(t, smp.msg.d2, fixtureMessage1().d2)
}

func Test_computesC3AndD3CorrectlyForOtrV2(t *testing.T) {
	otr := newConversation(otrV2{}, fixtureRand())
	smp, _ := otr.generateSMP1()
	assertDeepEquals(t, smp.msg.c3, fixtureMessage1().c3)
	assertDeepEquals(t, smp.msg.d3, fixtureMessage1().d3)
}

func Test_thatVerifySMPStartParametersCheckG2AForOtrV3(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	err := c.verifySMP1(smp1Message{g2a: new(big.Int).SetInt64(1)})
	assertDeepEquals(t, err, newOtrError("g2a is an invalid group element"))
}

func Test_thatVerifySMPStartParametersCheckG3AForOtrV3(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	err := c.verifySMP1(smp1Message{g2a: new(big.Int).SetInt64(3), g3a: p})
	assertDeepEquals(t, err, newOtrError("g3a is an invalid group element"))
}

func Test_thatVerifySMPStartParametersDoesntCheckG2AForOtrV2(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	err := c.verifySMP1(smp1Message{
		g2a: new(big.Int).SetInt64(1),
		g3a: new(big.Int).SetInt64(1),
		c2:  new(big.Int).SetInt64(1),
		c3:  new(big.Int).SetInt64(1),
		d2:  new(big.Int).SetInt64(1),
		d3:  new(big.Int).SetInt64(1),
	})
	assertDeepEquals(t, err, newOtrError("c2 is not a valid zero knowledge proof"))
}

func Test_thatVerifySMPStartParametersDoesntCheckG3AForOtrV2(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	err := c.verifySMP1(smp1Message{
		g2a: new(big.Int).SetInt64(3),
		g3a: new(big.Int).SetInt64(1),
		c2:  new(big.Int).SetInt64(1),
		c3:  new(big.Int).SetInt64(1),
		d2:  new(big.Int).SetInt64(1),
		d3:  new(big.Int).SetInt64(1),
	})
	assertDeepEquals(t, err, newOtrError("c2 is not a valid zero knowledge proof"))
}

func Test_thatVerifySMPStartParametersChecksThatc2IsAValidZeroKnowledgeProof(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	err := c.verifySMP1(smp1Message{
		g2a: new(big.Int).SetInt64(3),
		g3a: new(big.Int).SetInt64(3),
		c2:  new(big.Int).SetInt64(3),
		c3:  new(big.Int).SetInt64(3),
		d2:  new(big.Int).SetInt64(3),
		d3:  new(big.Int).SetInt64(3),
	})
	assertDeepEquals(t, err, newOtrError("c2 is not a valid zero knowledge proof"))
}

func Test_thatVerifySMPStartParametersChecksThatc3IsAValidZeroKnowledgeProof(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	err := c.verifySMP1(smp1Message{
		g2a: fixtureMessage1().g2a,
		g3a: new(big.Int).SetInt64(3),
		c2:  fixtureMessage1().c2,
		c3:  new(big.Int).SetInt64(3),
		d2:  fixtureMessage1().d2,
		d3:  new(big.Int).SetInt64(3),
	})
	assertDeepEquals(t, err, newOtrError("c3 is not a valid zero knowledge proof"))
}

func Test_thatVerifySMPStartParametersIsOKWithAValidParameterMessage(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())

	g2a, _ := new(big.Int).SetString("8a88c345c63aa25dab9815f8c51f6b7b621a12d31c8220a0579381c1e2e85a2275e2407c79c8e6e1f72ae765804e6b4562ac1b2d634313c70d59752ac119c6da5cb95dde3eedd9c48595b37256f5b64c56fb938eb1131447c9af9054b42841c57d1f41fe5aa510e2bd2965434f46dd0473c60d6114da088c7047760b00bc10287a03afc4c4f30e1c7dd7c9dbd51bdbd049eb2b8921cbdc72b4f69309f61e559c2d6dec9c9ce6f38ccb4dfd07f4cf2cf6e76279b88b297848c473e13f091a0f77", 16)
	g3a, _ := new(big.Int).SetString("d275468351fd48246e406ee74a8dc3db6ee335067bfa63300ce6a23867a1b2beddbdae9a8a36555fd4837f3ef8bad4f7fd5d7b4f346d7c7b7cb64bd7707eeb515902c66aa0c9323931364471ab93dd315f65c6624c956d74680863a9388cd5d89f1b5033b1cf232b8b6dcffaaea195de4e17cc1ba4c99497be18c011b2ad7742b43fa9ee3f95f7b6da02c8e894d054eb178a7822273655dc286ad15874687fe6671908d83662e7a529744ce4ea8dad49290d19dbe6caba202a825a20a27ee98a", 16)
	c2, _ := new(big.Int).SetString("d3b6ef5528fa97e983395bec165fa4ced7657bdabf3742d60880965c369c880c", 16)
	d2, _ := new(big.Int).SetString("7fffffffffffffffe487ed5110b4611a62633145c06e0e68948127044533e63a0105df531d89cd9128a5043cc71a026ef7ca8cd9e69d218d98158536f92f8a1ba7f09ab6b6a8e122f242dabb312f3f637a262174d31bf6b585ffae5b7a035bf6f71c35fdad44cfd2d74f9208be258ff324943328f6722d9ee1003e5c50b1df82cc6d241b0e2ae9cd348b1fd47e9267af339d65211b4fcfa466656c89b4217f90102e4aa3ac176a41f6240f32689712b0391c1c659757f4bfb83e6ba66bf8b630", 16)
	c3, _ := new(big.Int).SetString("57d8cfda442854ecb01b28e631aa9165d51d1192f7f464bf17ea7f6665c05030", 16)
	d3, _ := new(big.Int).SetString("7fffffffffffffffe487ed5110b4611a62633145c06e0e68948127044533e63a0105df531d89cd9128a5043cc71a026ef7ca8cd9e69d218d98158536f92f8a1ba7f09ab6b6a8e122f242dabb312f3f637a262174d31bf6b585ffae5b7a035bf6f71c35fdad44cfd2d74f9208be258ff324943328f6722d9ee1003e5c50b1df82cc6d241b0e2ae9cd348b1fd47e9267af8140bb2aa65628bcff455920bba95a1392f2fcb5c115f43a7a828b5bf0393c5c775a17a88506a7893ff509d674cd655c", 16)

	err := c.verifySMP1(smp1Message{
		g2a: g2a,
		g3a: g3a,
		c2:  c2,
		c3:  c3,
		d2:  d2,
		d3:  d3,
	})
	assertDeepEquals(t, err, nil)
}

func Test_thatVerifySMPStartParametersIsOKWithAValidParameterMessageWithProtocolV2(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())

	g2a, _ := new(big.Int).SetString("8a88c345c63aa25dab9815f8c51f6b7b621a12d31c8220a0579381c1e2e85a2275e2407c79c8e6e1f72ae765804e6b4562ac1b2d634313c70d59752ac119c6da5cb95dde3eedd9c48595b37256f5b64c56fb938eb1131447c9af9054b42841c57d1f41fe5aa510e2bd2965434f46dd0473c60d6114da088c7047760b00bc10287a03afc4c4f30e1c7dd7c9dbd51bdbd049eb2b8921cbdc72b4f69309f61e559c2d6dec9c9ce6f38ccb4dfd07f4cf2cf6e76279b88b297848c473e13f091a0f77", 16)
	g3a, _ := new(big.Int).SetString("d275468351fd48246e406ee74a8dc3db6ee335067bfa63300ce6a23867a1b2beddbdae9a8a36555fd4837f3ef8bad4f7fd5d7b4f346d7c7b7cb64bd7707eeb515902c66aa0c9323931364471ab93dd315f65c6624c956d74680863a9388cd5d89f1b5033b1cf232b8b6dcffaaea195de4e17cc1ba4c99497be18c011b2ad7742b43fa9ee3f95f7b6da02c8e894d054eb178a7822273655dc286ad15874687fe6671908d83662e7a529744ce4ea8dad49290d19dbe6caba202a825a20a27ee98a", 16)
	c2, _ := new(big.Int).SetString("d3b6ef5528fa97e983395bec165fa4ced7657bdabf3742d60880965c369c880c", 16)
	d2, _ := new(big.Int).SetString("7fffffffffffffffe487ed5110b4611a62633145c06e0e68948127044533e63a0105df531d89cd9128a5043cc71a026ef7ca8cd9e69d218d98158536f92f8a1ba7f09ab6b6a8e122f242dabb312f3f637a262174d31bf6b585ffae5b7a035bf6f71c35fdad44cfd2d74f9208be258ff324943328f6722d9ee1003e5c50b1df82cc6d241b0e2ae9cd348b1fd47e9267af339d65211b4fcfa466656c89b4217f90102e4aa3ac176a41f6240f32689712b0391c1c659757f4bfb83e6ba66bf8b630", 16)
	c3, _ := new(big.Int).SetString("57d8cfda442854ecb01b28e631aa9165d51d1192f7f464bf17ea7f6665c05030", 16)
	d3, _ := new(big.Int).SetString("7fffffffffffffffe487ed5110b4611a62633145c06e0e68948127044533e63a0105df531d89cd9128a5043cc71a026ef7ca8cd9e69d218d98158536f92f8a1ba7f09ab6b6a8e122f242dabb312f3f637a262174d31bf6b585ffae5b7a035bf6f71c35fdad44cfd2d74f9208be258ff324943328f6722d9ee1003e5c50b1df82cc6d241b0e2ae9cd348b1fd47e9267af8140bb2aa65628bcff455920bba95a1392f2fcb5c115f43a7a828b5bf0393c5c775a17a88506a7893ff509d674cd655c", 16)

	err := c.verifySMP1(smp1Message{
		g2a: g2a,
		g3a: g3a,
		c2:  c2,
		c3:  c3,
		d2:  d2,
		d3:  d3,
	})
	assertDeepEquals(t, err, nil)
}
