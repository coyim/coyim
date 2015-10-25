package cli

import (
	"fmt"

	"github.com/twstrike/coyim/session"
)

func (c *cliUI) observeSessionEvents() {
	//TODO: check for channel close
	for ev := range c.events {
		switch ev.EventType {
		case session.Connected:
			info(c.term, fmt.Sprintf("Your fingerprint is %x", ev.Session.PrivateKey.DefaultFingerprint()))

		case session.Disconnected:
			c.terminate <- true
		case session.RosterReceived:
			for _, entry := range ev.Session.R.ToSlice() {
				c.input.addUser(entry.Jid)
			}
		case session.IQReceived:
			c.input.addUser(ev.From)
		}
	}
}
