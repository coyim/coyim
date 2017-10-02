package otr3

import "testing"

func Test_SecureSessionID_returnsTheSessionIDAsTwoFormattedStrings(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ssid = [8]byte{0x01, 0x02, 0xF3, 0x04, 0x00, 0x06, 0x07, 0x08}

	parts, _ := c.SecureSessionID()

	assertDeepEquals(t, parts[0], "0102f304")
	assertDeepEquals(t, parts[1], "00060708")
}

func Test_SecureSessionID_returnsTheIndexOfTheValueThatShouldBeHighlighted(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ssid = [8]byte{0x01, 0x02, 0xF3, 0x04, 0x00, 0x06, 0x07, 0x08}

	c.sentRevealSig = true
	_, f := c.SecureSessionID()
	assertEquals(t, f, 0)

	c.sentRevealSig = false
	_, f = c.SecureSessionID()
	assertEquals(t, f, 1)
}
