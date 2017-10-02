package otr3

import "testing"

func Test_FragmentedMessage_canBeConvertedToSliceOfByteSlices(t *testing.T) {
	fragmented := []ValidMessage{
		ValidMessage{0x01, 0x02},
		ValidMessage{0x03, 0x04},
		ValidMessage{0x05, 0x06},
	}

	assertDeepEquals(t, Bytes(fragmented), [][]byte{
		[]byte{0x01, 0x02},
		[]byte{0x03, 0x04},
		[]byte{0x05, 0x06},
	})
}
