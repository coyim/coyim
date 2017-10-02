package otr3

import "testing"

func Test_verifyInstanceTags_ignoresOurInstaceTagIfItIsZero(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}

	err := v.verifyInstanceTags(c, 0x100, 0)

	assertNil(t, err)
}

func Test_verifyInstanceTags_returnsErrorWhenOurInstanceTagIsLesserThan0x100(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}

	err := v.verifyInstanceTags(c, 0x100, 0x99)

	assertEquals(t, err, errInvalidOTRMessage)
}

func Test_verifyInstanceTags_signalsMalformedMessageWhenOurInstanceTagIsLesserThan0x100(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}

	c.expectMessageEvent(t, func() {
		v.verifyInstanceTags(c, 0x100, 0x99)
	}, MessageEventReceivedMessageMalformed, nil, nil)
}

func Test_verifyInstanceTags_returnsErrorWhenOurInstanceDoesNotMatch(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}
	c.theirInstanceTag = 0x100
	c.ourInstanceTag = 0x122

	err := v.verifyInstanceTags(c, c.theirInstanceTag, 0x121)

	assertEquals(t, err, errReceivedMessageForOtherInstance)
}

func Test_verifyInstanceTags_signalsAMessageEventWhenOurInstanceTagDoesNotMatch(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}
	c.theirInstanceTag = 0x100
	c.ourInstanceTag = 0x122

	c.expectMessageEvent(t, func() {
		v.verifyInstanceTags(c, c.theirInstanceTag, 0x121)
	}, MessageEventReceivedMessageForOtherInstance, nil, nil)
}

func Test_verifyInstanceTags_savesTheirInstanceTag(t *testing.T) {
	v := otrV3{}
	c := &Conversation{}

	err := v.verifyInstanceTags(c, 0x101, 0)
	assertNil(t, err)
	assertEquals(t, c.theirInstanceTag, uint32(0x101))
}

func Test_verifyInstanceTags_returnsErrorWhenTheirInstanceTagIsLesserThan0x100(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}

	err := v.verifyInstanceTags(c, 0x99, 0x100)

	assertEquals(t, err, errInvalidOTRMessage)
}

func Test_verifyInstanceTags_returnsErrorWhenTheirInstanceTagIsZero(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}

	err := v.verifyInstanceTags(c, 0, 0x100)

	assertEquals(t, err, errInvalidOTRMessage)
}

func Test_verifyInstanceTags_signalsMalformedMessageWhenTheirInstanceTagIsTooLow(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}

	c.expectMessageEvent(t, func() {
		v.verifyInstanceTags(c, 0, 0x100)
	}, MessageEventReceivedMessageMalformed, nil, nil)
}

func Test_verifyInstanceTags_returnsErrorWhenTheirInstanceTagDoesNotMatch(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}
	c.theirInstanceTag = 0x122

	err := v.verifyInstanceTags(c, 0x121, c.ourInstanceTag)
	assertEquals(t, err, errReceivedMessageForOtherInstance)
}

func Test_verifyInstanceTags_signalsAMessageEventWhenTheirInstanceTagDoesNotMatch(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}
	c.theirInstanceTag = 0x122

	c.expectMessageEvent(t, func() {
		v.verifyInstanceTags(c, 0x121, c.ourInstanceTag)
	}, MessageEventReceivedMessageForOtherInstance, nil, nil)
}

func Test_otrv3_parseMessageHeader_signalsMalformedMessageWhenWeCantParseInstanceTags(t *testing.T) {
	v := otrV3{}
	c := &Conversation{version: v}

	c.expectMessageEvent(t, func() {
		v.parseMessageHeader(c, []byte{0x00, 0x03, 0x02, 0x00, 0x00, 0x01, 0x22, 0x00, 0x00, 0x01})
	}, MessageEventReceivedMessageMalformed, nil, nil)
}

func Test_generateInstanceTag_generatesOurInstanceTag(t *testing.T) {
	rand := fixedRand([]string{"00000099", "00001234"})
	c := &Conversation{Rand: rand}

	err := c.generateInstanceTag()

	assertEquals(t, err, nil)
	assertEquals(t, c.ourInstanceTag, uint32(0x1234))
}

func Test_generateInstanceTag_returnsAnErrorIfFailsToReadFromRand(t *testing.T) {
	rand := fixedRand([]string{"00000099", "00000080"})
	c := &Conversation{Rand: rand}

	err := c.generateInstanceTag()

	assertEquals(t, err, errShortRandomRead)
	assertEquals(t, c.ourInstanceTag, uint32(0))
}

func Test_messageHeader_generatesOurInstanceTagLazily(t *testing.T) {
	c := &Conversation{}

	_, err := otrV3{}.messageHeader(c, msgTypeDHCommit)

	assertEquals(t, err, nil)
	assertEquals(t, c.ourInstanceTag < minValidInstanceTag, false)

	previousInstanceTag := c.ourInstanceTag

	_, err = otrV3{}.messageHeader(c, msgTypeDHCommit)
	assertEquals(t, err, nil)
	assertEquals(t, c.ourInstanceTag, previousInstanceTag)
}
