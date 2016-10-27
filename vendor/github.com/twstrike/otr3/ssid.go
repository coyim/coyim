package otr3

import "fmt"

// SecureSessionID returns the secure session ID as two formatted strings
// The index returned points to the string that should be highlighted
func (c *Conversation) SecureSessionID() (parts []string, highlightIndex int) {
	l := fmt.Sprintf("%0x", c.ssid[0:4])
	r := fmt.Sprintf("%0x", c.ssid[4:])

	ix := 1
	if c.sentRevealSig {
		ix = 0
	}

	return []string{l, r}, ix
}
