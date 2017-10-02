package otr3

import "testing"

func Test_newOtrVersion_returnsTheCorrectOTRVersionForAValidVersionNumber(t *testing.T) {
	v, _ := newOtrVersion(3, policies(allowV3))
	_, ok := v.(otrV3)
	assertEquals(t, ok, true)
}

func Test_newOtrVersion_returnsUnsupportedVersionErrorIfGivenAWrongVersion(t *testing.T) {
	_, err := newOtrVersion(4, policies(allowV3))
	assertEquals(t, err, errUnsupportedOTRVersion)
}

func Test_newOtrVersion_returnsAnErrorIfGivenAVersionThatIsntAllowedByPolicy(t *testing.T) {
	_, err := newOtrVersion(3, policies(allowV2))
	assertEquals(t, err, errInvalidVersion)
}

func Test_checkVersion_returnsErrorIfTheMessageIsCorrupt(t *testing.T) {
	c := &Conversation{}
	e := c.checkVersion([]byte{0x00})
	assertEquals(t, e, errInvalidOTRMessage)
}

func Test_checkVersion_setsTheConversationVersionIfWeHaveNoExistingVersion(t *testing.T) {
	c := &Conversation{Policies: policies(allowV3)}
	c.ourKeys = []PrivateKey{alicePrivateKey}
	e := c.checkVersion([]byte{0x00, 0x03})
	assertEquals(t, e, nil)
	assertDeepEquals(t, c.version, otrV3{})
}

func Test_checkVersion_setsTheConversationVersionIfWeHaveTheCorrectPolicy(t *testing.T) {
	c := &Conversation{Policies: policies(allowV2)}
	c.ourKeys = []PrivateKey{alicePrivateKey}
	e := c.checkVersion([]byte{0x00, 0x02})
	assertEquals(t, e, nil)
	assertDeepEquals(t, c.version, otrV2{})
}

func Test_checkVersion_returnsTheErrorFromNewOtrVersion(t *testing.T) {
	c := &Conversation{Policies: policies(allowV2)}
	c.ourKeys = []PrivateKey{alicePrivateKey}
	e := c.checkVersion([]byte{0x00, 0x03})
	assertEquals(t, e, errUnsupportedOTRVersion)
}

func Test_checkVersion_doesNotSetConversationVersionIfOneIsAlreadySet(t *testing.T) {
	c := &Conversation{Policies: policies(allowV2 | allowV3), version: otrV3{}}
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.checkVersion([]byte{0x00, 0x02})
	assertEquals(t, otrV3{}, c.version)
}

func Test_checkVersion_returnsErrorIfCurrentVersionIsDifferentFromMessageVersion(t *testing.T) {
	c := &Conversation{Policies: policies(allowV2 | allowV3), version: otrV3{}}
	c.ourKeys = []PrivateKey{alicePrivateKey}
	e := c.checkVersion([]byte{0x00, 0x02})
	assertEquals(t, e, errWrongProtocolVersion)
}
