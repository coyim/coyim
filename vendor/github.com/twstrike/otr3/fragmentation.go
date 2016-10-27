package otr3

import "bytes"

var (
	fragmentSeparator      = []byte{','}
	fragmentItagsSeparator = []byte{'|'}
)

// fragmentationContext store the current fragmentation running. A fragmentationContext is zero-valid and can be immediately used without initialization.
// In order to follow the fragmentation rules, when the context needs to be reset, just create a new one - don't bother resetting variables
type fragmentationContext struct {
	frag                     []byte
	currentIndex, currentLen uint16
}

func min(l, r uint16) uint16 {
	if l < r {
		return l
	}
	return r
}

func fragmentStart(i, fraglen uint16) uint16 {
	return uint16(i * fraglen)
}

func fragmentEnd(i, fraglen, l uint16) uint16 {
	return uint16(min((i+1)*fraglen, l))
}

func fragmentData(data []byte, i int, fraglen, l uint16) []byte {
	return data[fragmentStart(uint16(i), fraglen):fragmentEnd(uint16(i), fraglen, l)]
}

// SetFragmentSize sets the maximum size for a message fragment.
// If specified, all messages produced by Receive and Send
// will be fragmented into messages of, at most, this number of bytes.
func (c *Conversation) SetFragmentSize(size uint16) {
	c.fragmentSize = size
}

func (c *Conversation) fragment(data encodedMessage, fraglen uint16) []ValidMessage {
	l := len(data)

	if l <= int(fraglen) || fraglen == 0 {
		return []ValidMessage{ValidMessage(data)}
	}

	fakeHeader := c.version.fragmentPrefix(1, 1, c.ourInstanceTag, c.theirInstanceTag)
	realFraglen := (fraglen - uint16(len(fakeHeader))) - 1

	if realFraglen <= 0 {
		return []ValidMessage{ValidMessage(data)}
	}

	numFragments := (l / int(realFraglen)) + 1
	ret := make([]ValidMessage, numFragments)
	for i := 0; i < numFragments; i++ {
		prefix := c.version.fragmentPrefix(i, numFragments, c.ourInstanceTag, c.theirInstanceTag)
		ret[i] = append(append(prefix, fragmentData(data, i, realFraglen, uint16(l))...), fragmentSeparator[0])
	}
	return ret
}

func fragmentsFinished(fctx fragmentationContext) bool {
	return fctx.currentIndex > 0 && fctx.currentIndex == fctx.currentLen
}

func parseFragment(data []byte) (resultData []byte, ix uint16, length uint16, ok bool) {
	parts := bytes.Split(data, fragmentSeparator)
	if len(parts) != 4 {
		return nil, 0, 0, false
	}
	var e1, e2 error
	ix, e1 = bytesToUint16(parts[0])
	length, e2 = bytesToUint16(parts[1])
	resultData = parts[2]
	ok = e1 == nil && e2 == nil
	return
}

func fragmentIsInvalid(ix, l uint16) bool {
	return ix == 0 || l == 0 || ix > l
}

func fragmentIsFirstMessage(ix, l uint16) bool {
	return ix == 1
}

func fragmentIsNextMessage(beforeCtx fragmentationContext, ix, l uint16) bool {
	return beforeCtx.currentIndex+1 == ix && beforeCtx.currentLen == l
}

func (ctx fragmentationContext) discardFragment() fragmentationContext {
	return ctx
}

func (ctx fragmentationContext) appendFragment(data []byte, ix, l uint16) fragmentationContext {
	return fragmentationContext{append(ctx.frag, data...), ix, l}
}

func restartFragment(data []byte, ix, l uint16) fragmentationContext {
	return fragmentationContext{makeCopy(data), ix, l}
}

func forgetFragment() fragmentationContext {
	return fragmentationContext{}
}

func (c *Conversation) receiveFragment(beforeCtx fragmentationContext, data ValidMessage) (fragmentationContext, error) {
	fragBody, ignore, ok1 := c.parseFragmentPrefix(data)
	resultData, ix, l, ok2 := parseFragment(fragBody)

	if ignore {
		c.messageEvent(MessageEventReceivedMessageForOtherInstance)
		return beforeCtx, nil
	}

	if !ok1 || !ok2 {
		return beforeCtx, newOtrError("invalid OTR fragment")
	}

	switch {
	case fragmentIsInvalid(ix, l):
		return beforeCtx.discardFragment(), nil
	case fragmentIsFirstMessage(ix, l):
		return restartFragment(resultData, ix, l), nil
	case fragmentIsNextMessage(beforeCtx, ix, l):
		return beforeCtx.appendFragment(resultData, ix, l), nil
	default:
		return forgetFragment(), nil
	}
}
