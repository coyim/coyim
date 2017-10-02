package otr3

import "testing"

func Test_receiveQueryMessage_sendDHCommitv3AndTransitToStateAwaitingDHKey(t *testing.T) {
	queryMsg := []byte("?OTRv?23?")

	c := &Conversation{Policies: policies(allowV3)}
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	msg, err := c.receiveQueryMessage(queryMsg)

	assertNil(t, err)
	assertEquals(t, c.ake.state, authStateAwaitingDHKey{})
	assertDeepEquals(t, dhMsgType(msg[0]), msgTypeDHCommit)
	assertDeepEquals(t, dhMsgVersion(msg[0]), uint16(3))
}

func Test_receiveQueryMessageV2_sendDHCommitv2(t *testing.T) {
	queryMsg := []byte("?OTRv?23?")

	c := &Conversation{Policies: policies(allowV2)}
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	msg, err := c.receiveQueryMessage(queryMsg)

	assertNil(t, err)
	assertEquals(t, c.ake.state, authStateAwaitingDHKey{})
	assertDeepEquals(t, dhMsgType(msg[0]), msgTypeDHCommit)
	assertDeepEquals(t, dhMsgVersion(msg[0]), uint16(2))
}

func Test_receiveQueryMessageV2V3_sendDHCommitv3WhenV2AndV3AreAllowed(t *testing.T) {
	queryMsg := []byte("?OTRv?23?")

	c := &Conversation{Policies: policies(allowV2 | allowV3)}
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	msg, err := c.receiveQueryMessage(queryMsg)

	assertNil(t, err)
	assertEquals(t, c.ake.state, authStateAwaitingDHKey{})
	assertDeepEquals(t, dhMsgType(msg[0]), msgTypeDHCommit)
	assertDeepEquals(t, dhMsgVersion(msg[0]), uint16(3))
}

func Test_receiveQueryMessage_StoresRAndXAndGx(t *testing.T) {
	fixture := fixtureConversation()
	fixture.dhCommitMessage()

	cxt := &Conversation{
		version: otrV3{},
		Rand:    fixtureRand(),
	}
	cxt.SetOurKeys([]PrivateKey{bobPrivateKey})

	_, err := cxt.sendDHCommit()

	assertNil(t, err)
	assertDeepEquals(t, cxt.ake.r, fixture.ake.r)
	assertDeepEquals(t, cxt.ake.secretExponent, fixture.ake.secretExponent)
	assertDeepEquals(t, cxt.ake.ourPublicValue, fixture.ake.ourPublicValue)
}

func Test_receiveQueryMessage_signalsMessageEventOnFailure(t *testing.T) {
	queryMsg := []byte("?OTRv3?")

	c := newConversation(nil, fixedRand([]string{"ABCD"}))
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	c.Policies.add(allowV3)
	c.expectMessageEvent(t, func() {
		c.receiveQueryMessage(queryMsg)
	}, MessageEventSetupError, nil, errShortRandomRead)
}

func Test_receiveQueryMessage_returnsErrorIfNoCompatibleVersionCouldBeFound(t *testing.T) {
	c := &Conversation{Policies: policies(allowV3)}
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	_, err := c.receiveQueryMessage([]byte("?OTRv?2?"))
	assertEquals(t, err, errUnsupportedOTRVersion)
}

func Test_receiveQueryMessage_returnsErrorIfDhCommitMessageGeneratesError(t *testing.T) {
	c := &Conversation{
		Policies: policies(allowV2),
		Rand:     fixedRand([]string{"ABCDABCD"}),
	}
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	_, err := c.receiveQueryMessage([]byte("?OTRv2?"))
	assertEquals(t, err, errShortRandomRead)
}

func Test_parseOTRQueryMessage(t *testing.T) {
	var exp = map[string][]int{
		"?OTR?":     []int{1},
		"?OTRv2?":   []int{2},
		"?OTRv23?":  []int{2, 3},
		"?OTR?v2":   []int{1, 2},
		"?OTRv248?": []int{2, 4, 8},
		"?OTR?v?":   []int{1},
		"?OTRv?":    []int{},
	}

	for queryMsg, versions := range exp {
		m := []byte(queryMsg)
		assertDeepEquals(t, parseOTRQueryMessage(m), versions)
	}
}

func Test_extractVersionsFromQueryMessage_returnsNilForUnsupportedVersions(t *testing.T) {
	p := policies(0)
	msg := []byte("?OTR?")
	versions := extractVersionsFromQueryMessage(p, msg)

	assertEquals(t, versions, 0)
}

func Test_extractVersionsFromQueryMessage_acceptsBothV2AndV3IfThePolicyAllows(t *testing.T) {
	msg := []byte("?OTRv32?")
	p := policies(allowV2 | allowV3)
	versions := extractVersionsFromQueryMessage(p, msg)

	assertEquals(t, versions, 1<<2|1<<3)
}

func Test_extractVersionsFromQueryMessage_acceptsOTRV2IfHasOnlyAllowV2Policy(t *testing.T) {
	msg := []byte("?OTRv32?")
	p := policies(allowV2)
	versions := extractVersionsFromQueryMessage(p, msg)

	assertEquals(t, versions, 1<<2)
}

func Test_QueryMessage_returnsARegularQueryMessage(t *testing.T) {
	c := &Conversation{Policies: policies(allowV3)}
	assertEquals(t, string(c.QueryMessage()), "?OTRv3?")
}

func Test_QueryMessage_returnsAQueryMessageWithExtraMessage(t *testing.T) {
	c := &Conversation{Policies: policies(allowV3)}
	c.SetFriendlyQueryMessage("hello foobarium")
	assertEquals(t, string(c.QueryMessage()), "?OTRv3? hello foobarium")
}
