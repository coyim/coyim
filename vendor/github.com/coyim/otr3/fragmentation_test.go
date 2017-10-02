package otr3

import (
	"crypto/rand"
	"testing"
)

const defaultInstanceTag = 0x00000100

func Test_isFragmented_returnsFalseForAShortValue(t *testing.T) {
	ctx := newConversation(otrV2{}, rand.Reader)
	assertEquals(t, ctx.version.isFragmented([]byte("")), false)
}

func Test_isFragmented_returnsFalseForALongValue(t *testing.T) {
	ctx := newConversation(otrV2{}, rand.Reader)
	assertEquals(t, ctx.version.isFragmented([]byte("?OTR:BLA")), false)
}

func Test_isFragmented_returnsFalseForAFragmentedV3MessageWhenRunningV2(t *testing.T) {
	ctx := newConversation(otrV2{}, rand.Reader)
	assertEquals(t, ctx.version.isFragmented([]byte("?OTR|BLA")), false)
}

func Test_isFragmented_returnsTrueForAFragmentedV3MessageWhenRunningV3(t *testing.T) {
	ctx := newConversation(otrV3{}, rand.Reader)
	assertEquals(t, ctx.version.isFragmented([]byte("?OTR|BLA")), true)
}

func Test_isFragmented_returnsTrueForAFragmentedV2MessageWhenRunningV2(t *testing.T) {
	ctx := newConversation(otrV2{}, rand.Reader)
	assertEquals(t, ctx.version.isFragmented([]byte("?OTR,BLA")), true)
}

func Test_isFragmented_returnsTrueForAFragmentedV2MessageWhenRunningV3(t *testing.T) {
	ctx := newConversation(otrV3{}, rand.Reader)
	assertEquals(t, ctx.version.isFragmented([]byte("?OTR,BLA")), true)
}

func Test_fragment_returnsNoChangeForASmallerPackage(t *testing.T) {
	ctx := newConversation(otrV3{}, rand.Reader)
	ctx.ourInstanceTag = defaultInstanceTag
	ctx.theirInstanceTag = defaultInstanceTag

	data := []byte("one two three")

	assertDeepEquals(t, ctx.fragment(data, 13), []ValidMessage{data})
}

func Test_fragment_returnsFragmentsForNeededFragmentation(t *testing.T) {
	ctx := newConversation(otrV3{}, rand.Reader)
	ctx.ourInstanceTag = defaultInstanceTag
	ctx.theirInstanceTag = defaultInstanceTag + 2

	data := []byte("one one one two two two three three three")
	assertDeepEquals(t, ctx.fragment(data, 40), []ValidMessage{
		[]byte("?OTR|00000100|00000102,00001,00011,one ,"),
		[]byte("?OTR|00000100|00000102,00002,00011,one ,"),
		[]byte("?OTR|00000100|00000102,00003,00011,one ,"),
		[]byte("?OTR|00000100|00000102,00004,00011,two ,"),
		[]byte("?OTR|00000100|00000102,00005,00011,two ,"),
		[]byte("?OTR|00000100|00000102,00006,00011,two ,"),
		[]byte("?OTR|00000100|00000102,00007,00011,thre,"),
		[]byte("?OTR|00000100|00000102,00008,00011,e th,"),
		[]byte("?OTR|00000100|00000102,00009,00011,ree ,"),
		[]byte("?OTR|00000100|00000102,00010,00011,thre,"),
		[]byte("?OTR|00000100|00000102,00011,00011,e,"),
	})
}

func Test_fragment_returnsFragmentsForNeededFragmentationForV2(t *testing.T) {
	ctx := newConversation(otrV2{}, rand.Reader)
	ctx.ourInstanceTag = defaultInstanceTag
	ctx.theirInstanceTag = defaultInstanceTag + 1

	data := []byte("one one one two two two three three three")

	res := ctx.fragment(data, 22)
	assertDeepEquals(t, res, []ValidMessage{
		[]byte("?OTR,00001,00011,one ,"),
		[]byte("?OTR,00002,00011,one ,"),
		[]byte("?OTR,00003,00011,one ,"),
		[]byte("?OTR,00004,00011,two ,"),
		[]byte("?OTR,00005,00011,two ,"),
		[]byte("?OTR,00006,00011,two ,"),
		[]byte("?OTR,00007,00011,thre,"),
		[]byte("?OTR,00008,00011,e th,"),
		[]byte("?OTR,00009,00011,ree ,"),
		[]byte("?OTR,00010,00011,thre,"),
		[]byte("?OTR,00011,00011,e,"),
	})
}

func Test_receiveFragment_returnsANewFragmentationContextForANewMessage(t *testing.T) {
	c := newConversation(otrV2{}, rand.Reader)
	data := []byte("?OTR,00001,00004,one ,")

	fctx, e := c.receiveFragment(fragmentationContext{}, data)

	assertDeepEquals(t, fctx.frag, []byte("one "))
	assertDeepEquals(t, e, nil)
	assertEquals(t, fctx.currentIndex, uint16(1))
	assertEquals(t, fctx.currentLen, uint16(4))
}

func Test_receiveFragment_returnsANewFragmentationContextForANewV3Message(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourInstanceTag = 0x102
	c.theirInstanceTag = 0x100
	data := []byte("?OTR|00000100|00000102,00001,00004,one ,")

	fctx, e := c.receiveFragment(fragmentationContext{}, data)

	assertDeepEquals(t, fctx.frag, []byte("one "))
	assertDeepEquals(t, e, nil)
	assertEquals(t, fctx.currentIndex, uint16(1))
	assertEquals(t, fctx.currentLen, uint16(4))
}

func Test_receiveFragment_returnsTheExistingContextIfTheInstanceTagsDoesNotMatch(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourInstanceTag = 0x103
	c.theirInstanceTag = 0x104

	existingContext := fragmentationContext{frag: []byte("shouldn't change")}

	fctx, _ := c.receiveFragment(existingContext, []byte("?OTR|00000204|00000103,00001,00004,one ,"))
	assertDeepEquals(t, fctx, existingContext)

	fctx, _ = c.receiveFragment(existingContext, []byte("?OTR|00000104|00000203,00001,00004,one ,"))
	assertDeepEquals(t, fctx, existingContext)
}

func Test_receiveFragment_signalsMessageEventIfInstanceTagsDoesNotMatch(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourInstanceTag = 0x103
	c.theirInstanceTag = 0x104

	existingContext := fragmentationContext{frag: []byte("shouldn't change")}

	c.expectMessageEvent(t, func() {
		c.receiveFragment(existingContext, []byte("?OTR|00000204|00000103,00001,00004,one ,"))
	}, MessageEventReceivedMessageForOtherInstance, nil, nil)
}

func Test_receiveFragment_sendsAnErrorMessageAboutMalformedIfHandlerExists(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourInstanceTag = 0x103
	c.theirInstanceTag = 0x0A

	existingContext := fragmentationContext{frag: []byte("shouldn't change")}

	c.errorMessageHandler = dynamicErrorMessageHandler{
		func(error ErrorCode) []byte {
			if error == ErrorCodeMessageMalformed {
				return []byte("black happened")
			}
			return []byte("white happened")
		}}

	c.receiveFragment(existingContext, []byte("?OTR|0000000A|00000103,00001,00004,one ,"))
	ts, _ := c.withInjections(nil, nil)
	assertDeepEquals(t, string(ts[0]), "?OTR Error: black happened")
}

func Test_receiveFragment_signalsMalformedMessageIfTheirInstanceTagIsBelowTheLimit(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourInstanceTag = 0x103
	c.theirInstanceTag = 0x0A

	existingContext := fragmentationContext{frag: []byte("shouldn't change")}

	c.expectMessageEvent(t, func() {
		c.receiveFragment(existingContext, []byte("?OTR|0000000A|00000103,00001,00004,one ,"))
	}, MessageEventReceivedMessageMalformed, nil, nil)
}

func Test_receiveFragment_returnsTheSameContextIfMessageNumberIsZero(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	data := []byte("?OTR,00000,00004,one ,")
	fctx, _ := c.receiveFragment(fragmentationContext{}, data)
	assertDeepEquals(t, fctx, fragmentationContext{})
}

func Test_receiveFragment_returnsTheSameContextIfMessageCountIsZero(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	data := []byte("?OTR,00001,00000,one ,")
	fctx, _ := c.receiveFragment(fragmentationContext{}, data)
	assertDeepEquals(t, fctx, fragmentationContext{})
}

func Test_receiveFragment_returnsTheSameContextIfMessageNumberIsAboveMessageCount(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	data := []byte("?OTR,00005,00004,one ,")
	fctx, _ := c.receiveFragment(fragmentationContext{}, data)
	assertDeepEquals(t, fctx, fragmentationContext{})
}

func Test_receiveFragment_returnsTheNextContextIfMessageNumberIsOneMoreThanThePreviousOne(t *testing.T) {
	c := newConversation(otrV2{}, rand.Reader)
	data := []byte("?OTR,00003,00004, one,")
	fctx, _ := c.receiveFragment(fragmentationContext{[]byte("blarg one two"), 2, 4}, data)
	assertDeepEquals(t, fctx, fragmentationContext{[]byte("blarg one two one"), 3, 4})
}

func Test_receiveFragment_resetsTheContextIfTheMessageCountIsNotTheSame(t *testing.T) {
	c := newConversation(otrV2{}, rand.Reader)
	data := []byte("?OTR,00003,00005, one,")
	fctx, _ := c.receiveFragment(fragmentationContext{[]byte("blarg one two"), 2, 4}, data)
	assertDeepEquals(t, fctx, fragmentationContext{})
}

func Test_receiveFragment_resetsTheContextIfTheMessageNumberIsNotExactlyOnePlus(t *testing.T) {
	c := newConversation(otrV2{}, rand.Reader)
	data := []byte("?OTR,00004,00005, one,")
	fctx, _ := c.receiveFragment(fragmentationContext{[]byte("blarg one two"), 2, 5}, data)
	assertDeepEquals(t, fctx, fragmentationContext{})
}

func Test_fragmentFinished_isFalseIfThereAreNoFragments(t *testing.T) {
	assertDeepEquals(t, fragmentsFinished(fragmentationContext{[]byte{}, 0, 0}), false)
}

func Test_fragmentFinished_isFalseIfTheNumberOfFragmentsIsNotTheSame(t *testing.T) {
	assertDeepEquals(t, fragmentsFinished(fragmentationContext{[]byte{}, 1, 2}), false)
}

func Test_fragmentFinished_isFalseIfTheNumberOfFragmentsIsNotTheSameWhereTheNumberIsHigher(t *testing.T) {
	assertDeepEquals(t, fragmentsFinished(fragmentationContext{[]byte{}, 3, 2}), false)
}

func Test_fragmentFinished_isTrueIfTheNumberIsTheSameAsTheCount(t *testing.T) {
	assertDeepEquals(t, fragmentsFinished(fragmentationContext{[]byte{}, 3, 3}), true)
}

func Test_parseFragment_returnsNotOKIfThereAreNotEnoughParts(t *testing.T) {
	_, _, _, ok := parseFragment([]byte{0x2C, 0x2C})
	assertDeepEquals(t, ok, false)
}

func Test_parseFragment_returnsNotOKIfThereAreTooManyParts(t *testing.T) {
	_, _, _, ok := parseFragment([]byte{0x2C, 0x2C, 0x2C, 0x2C})
	assertDeepEquals(t, ok, false)
}

func Test_parseFragment_returnsNotOKIfTheIndexIsNotAValidUint(t *testing.T) {
	_, _, _, ok := parseFragment([]byte{0x30, 0x30, 0x30, 0x30, 0x29, 0x2C, 0x30, 0x30, 0x30, 0x30, 0x31, 0x2C, 0x01, 0x2C})
	assertDeepEquals(t, ok, false)
}

func Test_parseFragment_returnsNotOKIfTheLengthIsNotAValidUint(t *testing.T) {
	_, _, _, ok := parseFragment([]byte{0x30, 0x30, 0x30, 0x30, 0x31, 0x2C, 0x30, 0x30, 0x30, 0x30, 0x29, 0x2C, 0x01, 0x2C})
	assertDeepEquals(t, ok, false)
}

func Test_parseFragment_returnsOKIfThereAreExactlyTheRightAmountOfParts(t *testing.T) {
	_, _, _, ok := parseFragment([]byte{0x30, 0x30, 0x30, 0x30, 0x31, 0x2C, 0x30, 0x30, 0x30, 0x30, 0x31, 0x2C, 0x01, 0x2C})
	assertDeepEquals(t, ok, true)
}

func Test_receiveFragment_returnsErrorIfTheFragmentIsNotCorrect(t *testing.T) {
	c := newConversation(otrV2{}, rand.Reader)
	_, e := c.receiveFragment(fragmentationContext{}, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x30, 0x30, 0x30, 0x29, 0x2C, 0x30, 0x30, 0x30, 0x30, 0x31, 0x2C, 0x01, 0x2C})
	assertDeepEquals(t, e, newOtrError("invalid OTR fragment"))
}

func Test_parseFragmentPrefix_resolveVersion2IfNotDefined(t *testing.T) {
	fragment := []byte("?OTR,00001,00004,?OTR:AAICAAAAxJh7YMX8vCry1O+3ewL88,")

	c := &Conversation{Policies: policies(allowV2)}
	c.parseFragmentPrefix(fragment)

	assertEquals(t, c.version, otrV2{})

}

func Test_parseFragmentPrefix_rejectsVersion2IfNotAllowedByThePolicy(t *testing.T) {
	fragment := []byte("?OTR,00001,00004,?OTR:AAICAAAAxJh7YMX8vCry1O+3ewL88,")

	c := &Conversation{Policies: policies(allowV3)}
	_, ignore, ok := c.parseFragmentPrefix(fragment)

	assertEquals(t, ok, false)
	assertEquals(t, ignore, true)
	assertEquals(t, c.version, nil)

}

func Test_parseFragmentPrefix_resolveVersion3IfNotDefined(t *testing.T) {
	fragment := []byte("?OTR|5a73a599|27e31597,00001,00003,?OTR:AAMDJ+MVmSfjF,")

	c := &Conversation{Policies: policies(allowV3)}
	c.parseFragmentPrefix(fragment)

	assertEquals(t, c.version, otrV3{})
}

func Test_parseFragmentPrefix_rejectsVersion3IfNotAllowedByThePolicy(t *testing.T) {
	fragment := []byte("?OTR|5a73a599|27e31597,00001,00003,?OTR:AAMDJ+MVmSfjF,")

	c := &Conversation{Policies: policies(allowV2)}
	_, ignore, ok := c.parseFragmentPrefix(fragment)

	assertEquals(t, ok, false)
	assertEquals(t, ignore, true)
	assertEquals(t, c.version, nil)
}
