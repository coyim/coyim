package otr3

// MessagePlaintext contains regular plaintext to send or receive
type MessagePlaintext []byte

// A messageWithHeader is simply a full message with all content but not valid to send
type messageWithHeader []byte

// An encoded message have been encoded and is in the final form of an OTR message.
type encodedMessage []byte

// ValidMessage is a message that has gone through fragmentation and is valid to send through the IM client
// Some encodedMessage instances are validMessage instances, but this depends on the fragmentation size
type ValidMessage []byte

// Bytes will turn a slice of valid messages into a slice of byte slices
func Bytes(m []ValidMessage) [][]byte {
	ret := make([][]byte, len(m))
	//copy because we don't want to hold references to m's fragments
	for i, f := range m {
		ret[i] = make([]byte, len(f))
		copy(ret[i], []byte(f))
	}
	return ret
}

func compactMessagesWithHeader(msgs ...messageWithHeader) []messageWithHeader {
	var res []messageWithHeader
	for _, m := range msgs {
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}
