package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
)

func (c *cliUI) handleSessionEvent(ev session.Event) {
	switch ev.Type {
	case session.Connected:
		info(c.term, fmt.Sprintf("Your fingerprint is %x", ev.Session.PrivateKey.DefaultFingerprint()))

	case session.Disconnected:
		c.terminate <- true
	case session.RosterReceived:
		for _, entry := range ev.Session.R.ToSlice() {
			c.input.addUser(entry.Jid)
		}
	}
}

func (c *cliUI) handlePeerEvent(ev session.PeerEvent) {
	switch ev.Type {
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

func (c *cliUI) handlePresenceEvent(ev session.PresenceEvent) {
	if ev.Session.CurrentAccount.HideStatusUpdates {
		return
	}

	from := xmpp.RemoveResourceFromJid(ev.From)

	var line []byte
	line = append(line, []byte(fmt.Sprintf("   (%s) ", time.Now().Format(time.Kitchen)))...)
	line = append(line, c.term.Escape.Magenta...)
	line = append(line, []byte(from)...)
	line = append(line, ':')
	line = append(line, c.term.Escape.Reset...)
	line = append(line, ' ')

	if ev.Gone {
		line = append(line, []byte("offline")...)
	} else if len(ev.Show) > 0 {
		line = append(line, []byte(ev.Show)...)
	} else {
		line = append(line, []byte("online")...)
	}
	line = append(line, ' ')
	line = append(line, []byte(ev.Status)...)
	line = append(line, '\n')
	c.term.Write(line)
}

func (c *cliUI) handleMessageEvent(ev session.MessageEvent) {
	var line []byte
	if ev.Encrypted {
		line = append(line, c.term.Escape.Green...)
	} else {
		line = append(line, c.term.Escape.Red...)
	}

	t := fmt.Sprintf("(%s) %s: ", ev.When.Format(time.Stamp), ev.From)
	line = append(line, []byte(t)...)
	line = append(line, c.term.Escape.Reset...)
	line = appendTerminalEscaped(line, ui.StripHTML(ev.Body))
	line = append(line, '\n')
	if c.session.Config.Bell {
		line = append(line, '\a')
	}

	c.term.Write(line)
}

func (c *cliUI) observeSessionEvents() {
	//TODO: check for channel close
	for ev := range c.events {
		switch t := ev.(type) {
		case session.Event:
			c.handleSessionEvent(t)
		case session.PeerEvent:
			c.handlePeerEvent(t)
		case session.PresenceEvent:
			c.handlePresenceEvent(t)
		case session.MessageEvent:
			c.handleMessageEvent(t)
		default:
			log.Printf("unsupported event %#v\n", t)
		}
	}
}
