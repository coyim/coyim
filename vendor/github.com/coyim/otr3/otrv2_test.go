package otr3

import "testing"

func Test_otrv2_parseFragmentPrefix_returnsNotOKIfDataIsTooShort(t *testing.T) {
	_, ignore, ok := otrV2{}.parseFragmentPrefix(nil, []byte{0x00})
	assertFalse(t, ignore)
	assertFalse(t, ok)
}

func Test_otrv2_parseMessageHeader_returnsErrorIfTheMessageIsTooShort(t *testing.T) {
	_, _, err := otrV2{}.parseMessageHeader(nil, []byte{0x00})
	assertEquals(t, err, errInvalidOTRMessage)
}
