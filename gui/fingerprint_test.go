package gui

import . "gopkg.in/check.v1"

type FingerprintXmppSuite struct{}

var _ = Suite(&FingerprintXmppSuite{})

func (s *FingerprintXmppSuite) Test_formatFingerprint(c *C) {
	testVal := []byte{0x5d, 0xfc, 0x9e, 0x41, 0x6b, 0xf7, 0x83, 0xea, 0x14, 0x90, 0xb8, 0x16, 0x9b, 0x86, 0x68, 0x21, 0xb5, 0x2e, 0xbb, 0xb7}

	res := formatFingerprint(testVal)

	c.Assert(res, Equals, "5DFC9E416B F783EA1490 B8169B8668 21B52EBBB7")
}
