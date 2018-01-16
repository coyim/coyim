package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/ui"
	"github.com/coyim/coyim/xmpp/jid"
)

func (c *cliUI) handleSessionEvent(ev events.Event) {
	switch ev.Type {
	case events.Connected:
		for _, pk := range ev.Session.PrivateKeys() {
			info(c.term, c.termControl, fmt.Sprintf("Your fingerprint is %x", pk.PublicKey().Fingerprint()))
		}
	case events.Disconnected:
		c.terminate <- true
	case events.RosterReceived:
		for _, entry := range ev.Session.R().ToSlice() {
			c.input.addUser(entry.Jid)
		}
	}
}

func (c *cliUI) handlePeerEvent(ev events.Peer) {
	switch ev.Type {
	case events.IQReceived:
		c.input.addUser(ev.From.NoResource())
	case events.OTREnded:
		c.input.SetPromptForTarget(ev.From, false)
	case events.OTRNewKeys, events.OTRRenewedKeys:
		uid := ev.From
		info(c.term, c.termControl, fmt.Sprintf("New OTR session with %s established", uid.String()))
		//TODO: review whether it should create conversations
		conversation, _ := ev.Session.ConversationManager().EnsureConversationWith(uid.NoResource(), jid.Resource(""))

		c.input.SetPromptForTarget(uid, true)
		c.printConversationInfo(uid, conversation)
	case events.SubscriptionRequest:
		msg := fmt.Sprintf("%[1]s wishes to see when you're online. Use '/confirm %[1]s' to confirm (or likewise with /deny to decline)", ev.From)

		info(c.term, c.termControl, msg)
		c.input.addUser(ev.From.NoResource())
	}
}

func (c *cliUI) handlePresenceEvent(ev events.Presence) {
	if ev.Session.GetConfig().HideStatusUpdates {
		return
	}

	from := jid.Parse(ev.From).NoResource().String()

	var line []byte
	line = append(line, []byte(fmt.Sprintf("   (%s) ", time.Now().Format(time.Kitchen)))...)
	line = append(line, c.termControl.Escape(c.term).Magenta...)
	line = append(line, []byte(from)...)
	line = append(line, ':')
	line = append(line, c.termControl.Escape(c.term).Reset...)
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
	line = append(line, '\r', '\n')
	c.term.Write(line)
}

func (c *cliUI) handleMessageEvent(ev events.Message) {
	var line []byte
	if ev.Encrypted {
		line = append(line, c.termControl.Escape(c.term).Green...)
	} else {
		line = append(line, c.termControl.Escape(c.term).Red...)
	}

	t := fmt.Sprintf("(%s) %s: ", ev.When.Format(time.Stamp), ev.From)
	line = append(line, []byte(t)...)
	line = append(line, c.termControl.Escape(c.term).Reset...)
	line = appendTerminalEscaped(line, ui.StripHTML(ev.Body))
	line = append(line, '\n')
	if c.session.Config().Bell {
		line = append(line, '\a')
	}

	c.term.Write(line)
}

func (c *cliUI) handleLogEvent(ev events.Log) {
	switch ev.Level {
	case events.Info:
		info(c.term, c.termControl, ev.Message)
	case events.Warn:
		warn(c.term, c.termControl, ev.Message)
	case events.Alert:
		alert(c.term, c.termControl, ev.Message)
	}
}

func (c *cliUI) handleSMPEvent(ev events.SMP) {
	switch ev.Type {
	case events.SecretNeeded:
		info(c.term, c.termControl, fmt.Sprintf("%s is attempting to authenticate. Please supply mutual shared secret with /otr-auth user secret", ev.From))
		if question := ev.Body; len(question) > 0 {
			info(c.term, c.termControl, fmt.Sprintf("%s asks: %s", ev.From, question))
		}
	case events.Success:
		info(c.term, c.termControl, fmt.Sprintf("Authentication with %s successful", ev.From))
	case events.Failure:
		alert(c.term, c.termControl, fmt.Sprintf("Authentication with %s failed", ev.From))
	}
}

func (c *cliUI) observeSessionEvents() {
	for ev := range c.events {
		switch t := ev.(type) {
		case events.Event:
			c.handleSessionEvent(t)
		case events.Peer:
			c.handlePeerEvent(t)
		case events.Presence:
			c.handlePresenceEvent(t)
		case events.Message:
			c.handleMessageEvent(t)
		case events.SMP:
			c.handleSMPEvent(t)
		default:
			log.Printf("unsupported event %#v\n", t)
		}
	}
}
