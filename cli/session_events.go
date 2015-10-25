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
		case session.OTREnded:
			c.input.SetPromptForTarget(ev.From, false)
		case session.OTRNewKeys:
			uid := ev.From
			c.input.SetPromptForTarget(uid, true)
			c.printConversationInfo(uid, ev.Session.GetConversationWith(uid))
		case session.SubscriptionRequest:
			msg := fmt.Sprintf("%[1]s wishes to see when you're online. Use '/confirm %[1]s' to confirm (or likewise with /deny to decline)", ev.From)

			info(c.term, msg)
			c.input.addUser(ev.From)
		}
	}
}
