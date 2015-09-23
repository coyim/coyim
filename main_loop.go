package main

import (
	"fmt"
	"time"

	coyui "github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
	"golang.org/x/crypto/ssh/terminal"
)

func stanzaLoop(ui coyui.UI, s *Session, stanzaChan <-chan xmpp.Stanza, done chan<- bool) {

StanzaLoop:
	for {
		select {
		case rawStanza, ok := <-stanzaChan:
			if !ok {
				ui.Warn("Exiting because channel to server closed")
				break StanzaLoop
			}
			switch stanza := rawStanza.Value.(type) {
			case *xmpp.StreamError:
				var text string
				if len(stanza.Text) > 0 {
					text = stanza.Text
				} else {
					text = fmt.Sprintf("%s", stanza.Any)
				}

				ui.Alert("Exiting in response to fatal error from server: " + text)
				break StanzaLoop
			case *xmpp.ClientMessage:
				s.processClientMessage(stanza)
			case *xmpp.ClientPresence:
				ui.ProcessPresence(stanza)
			case *xmpp.ClientIQ:
				if stanza.Type != "get" && stanza.Type != "set" {
					continue
				}
				reply := s.processIQ(stanza)
				if reply == nil {
					reply = xmpp.ErrorReply{
						Type:  "cancel",
						Error: xmpp.ErrorBadRequest{},
					}
				}
				if err := s.conn.SendIQReply(stanza.From, "result", stanza.Id, reply); err != nil {
					ui.Alert("Failed to send IQ message: " + err.Error())
				}
			default:
				ui.Info(fmt.Sprintf("%s %s", rawStanza.Name, rawStanza.Value))
			}
		}
	}

	done <- true
}

func rosterLoop(term *terminal.Terminal, s *Session, rosterReply <-chan xmpp.Stanza, done chan<- bool) {
RosterLoop:
	for {
		var err error

		select {
		case rosterStanza, ok := <-rosterReply:
			if !ok {
				//TODO was this error supposed to print the latest error regardless of where it came from?
				alert(term, "Failed to read roster: "+err.Error())
				break RosterLoop
			}
			if s.roster, err = xmpp.ParseRoster(rosterStanza); err != nil {
				alert(term, "Failed to parse roster: "+err.Error())
				break RosterLoop
			}

			s.rosterReceived()
			info(term, "Roster received")

		case edit := <-s.pendingRosterChan:
			if !edit.IsComplete {
				info(term, "Please edit "+edit.FileName+" and run /rostereditdone when complete")
				s.pendingRosterEdit = edit
				continue
			}
			if s.processEditedRoster(edit) {
				s.pendingRosterEdit = nil
			} else {
				alert(term, "Please reedit file and run /rostereditdone again")
			}

		}
	}

	done <- true
}

func timeoutLoop(s *Session, tick <-chan time.Time) {
	for now := range tick {
		haveExpired := false
		for _, expiry := range s.timeouts {
			if now.After(expiry) {
				haveExpired = true
				break
			}
		}
		if !haveExpired {
			continue
		}

		newTimeouts := make(map[xmpp.Cookie]time.Time)
		for cookie, expiry := range s.timeouts {
			if now.After(expiry) {
				s.conn.Cancel(cookie)
			} else {
				newTimeouts[cookie] = expiry
			}
		}

		s.timeouts = newTimeouts
	}
}
