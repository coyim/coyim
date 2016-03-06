package otr3

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const debugString = "?OTR!"
const debugPrefix = "[DEBUG] "

var standardErrorOutput io.Writer = os.Stderr

// SetDebug sets the debug mode for this conversation.
// If debug mode is enabled, calls to Send with a message equals to "?OTR!"
// will dump debug information about the current conversation state to stderr
func (c *Conversation) SetDebug(d bool) {
	c.debug = d
}

func (c *Conversation) otrOffer() string {
	switch c.whitespaceState {
	case whitespaceNotSent:
		return "NOT"
	case whitespaceSent:
		if c.msgState == encrypted {
			return "ACCEPTED"
		}
		return "SENT"
	case whitespaceRejected:
		return "REJECTED"
	default:
		return "INVALID"
	}
}

func (c *Conversation) dump(w *bufio.Writer) {
	w.WriteString("Context:\n\n")
	w.WriteString(fmt.Sprintf("  Our instance:   %08X\n", c.ourInstanceTag))
	w.WriteString(fmt.Sprintf("  Their instance: %08X\n\n", c.theirInstanceTag))
	w.WriteString(fmt.Sprintf("  Msgstate: %d (%s)\n\n", c.msgState, c.msgState.identityString()))
	w.WriteString(fmt.Sprintf("  Protocol version: %d\n", c.version.protocolVersion()))
	w.WriteString(fmt.Sprintf("  OTR offer: %s\n\n", c.otrOffer()))
	if c.ake == nil {
		w.WriteString("  Auth info: NULL\n")
	} else {
		c.dumpAKE(w)
	}
	w.WriteString("\n")

	c.dumpSMP(w)

	w.Flush()
}

// Will only be called if AKE is valid
func (c *Conversation) dumpAKE(w *bufio.Writer) {
	w.WriteString("  Auth info:\n")

	if c.ake != nil {
		w.WriteString(fmt.Sprintf("    State: %d (%s)\n", c.ake.state.identity(), c.ake.state.identityString()))
	}

	w.WriteString(fmt.Sprintf("    Our keyid:   %d\n", c.keys.ourKeyID))
	w.WriteString(fmt.Sprintf("    Their keyid: %d\n", c.keys.theirKeyID))
	w.WriteString(fmt.Sprintf("    Their fingerprint: %X\n", c.theirKey.Fingerprint()))
	w.WriteString(fmt.Sprintf("    Proto version = %d\n", c.version.protocolVersion()))
	w.Flush()
}

func (c *Conversation) dumpSMP(w *bufio.Writer) {
	w.WriteString("  SM state:\n")

	if c.smp.state != nil {
		w.WriteString(fmt.Sprintf("    Next expected: %d (%s)\n", c.smp.state.identity(), c.smp.state.identityString()))
	}

	receivedQ := 0
	if c.smp.question != nil {
		receivedQ = 1
	}
	w.WriteString(fmt.Sprintf("    Received_Q: %d\n", receivedQ))

	w.Flush()
}

func (smpStateExpect1) identity() int {
	return 0
}

func (smpStateExpect2) identity() int {
	return 2
}

func (smpStateExpect3) identity() int {
	return 3
}

func (smpStateExpect4) identity() int {
	return 4
}

func (smpStateWaitingForSecret) identity() int {
	return 1
}

func (smpStateExpect1) identityString() string {
	return "EXPECT1"
}

func (smpStateExpect2) identityString() string {
	return "EXPECT2"
}

func (smpStateExpect3) identityString() string {
	return "EXPECT3"
}

func (smpStateExpect4) identityString() string {
	return "EXPECT4"
}

func (smpStateWaitingForSecret) identityString() string {
	return "EXPECT1_WQ"
}

func (authStateNone) identity() int {
	return 0
}

func (authStateAwaitingDHKey) identity() int {
	return 1
}

func (authStateAwaitingRevealSig) identity() int {
	return 2
}

func (authStateAwaitingSig) identity() int {
	return 3
}

func (authStateNone) identityString() string {
	return "NONE"
}

func (authStateAwaitingDHKey) identityString() string {
	return "AWAITING_DHKEY"
}

func (authStateAwaitingRevealSig) identityString() string {
	return "AWAITING_REVEALSIG"
}

func (authStateAwaitingSig) identityString() string {
	return "AWAITING_SIG"
}

func (m msgState) identityString() string {
	switch m {
	case plainText:
		return "PLAINTEXT"
	case encrypted:
		return "ENCRYPTED"
	case finished:
		return "FINISHED"
	default:
		return "INVALID"
	}
}
