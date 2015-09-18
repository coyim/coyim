package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
	"golang.org/x/crypto/ssh/terminal"
)

func mainLoop(term *terminal.Terminal, s *Session, config *Config, tick <-chan time.Time, stanzaChan, rosterReply <-chan xmpp.Stanza, commandChan <-chan interface{}) {
	var err error
MainLoop:

	for {
		select {
		case now := <-tick:
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

		case edit := <-s.pendingRosterChan:
			if !edit.isComplete {
				info(s.term, "Please edit "+edit.fileName+" and run /rostereditdone when complete")
				s.pendingRosterEdit = edit
				continue
			}
			if s.processEditedRoster(edit) {
				s.pendingRosterEdit = nil
			} else {
				alert(s.term, "Please reedit file and run /rostereditdone again")
			}

		case rosterStanza, ok := <-rosterReply:
			if !ok {
				alert(s.term, "Failed to read roster: "+err.Error())
				return
			}
			if s.roster, err = xmpp.ParseRoster(rosterStanza); err != nil {
				alert(s.term, "Failed to parse roster: "+err.Error())
				return
			}
			for _, entry := range s.roster {
				s.input.AddUser(entry.Jid)
			}
			info(s.term, "Roster received")

		case cmd, ok := <-commandChan:
			if !ok {
				warn(term, "Exiting because command channel closed")
				break MainLoop
			}
			s.lastActionTime = time.Now()
			switch cmd := cmd.(type) {
			case quitCommand:
				for to, conversation := range s.conversations {
					msgs, err := conversation.End()
					if err != nil {
						//TODO: error handle
						panic("this should not happen")
					}
					for _, msg := range msgs {
						s.conn.Send(to, string(msg))
					}
				}
				break MainLoop
			case versionCommand:
				replyChan, cookie, err := s.conn.SendIQ(cmd.User, "get", xmpp.VersionQuery{})
				if err != nil {
					alert(s.term, "Error sending version request: "+err.Error())
					continue
				}
				s.timeouts[cookie] = time.Now().Add(5 * time.Second)
				go s.awaitVersionReply(replyChan, cmd.User)
			case rosterCommand:
				info(s.term, "Current roster:")
				maxLen := 0
				for _, item := range s.roster {
					if maxLen < len(item.Jid) {
						maxLen = len(item.Jid)
					}
				}

				for _, item := range s.roster {
					state, ok := s.knownStates[item.Jid]

					line := ""
					if ok {
						line += "[*] "
					} else if cmd.OnlineOnly {
						continue
					} else {
						line += "[ ] "
					}

					line += item.Jid
					numSpaces := 1 + (maxLen - len(item.Jid))
					for i := 0; i < numSpaces; i++ {
						line += " "
					}
					line += item.Subscription + "\t" + item.Name
					if ok {
						line += "\t" + state
					}
					info(s.term, line)
				}
			case rosterEditCommand:
				if s.pendingRosterEdit != nil {
					warn(s.term, "Aborting previous roster edit")
					s.pendingRosterEdit = nil
				}
				rosterCopy := make([]xmpp.RosterEntry, len(s.roster))
				copy(rosterCopy, s.roster)
				go s.editRoster(rosterCopy)
			case rosterEditDoneCommand:
				if s.pendingRosterEdit == nil {
					warn(s.term, "No roster edit in progress. Use /rosteredit to start one")
					continue
				}
				go s.loadEditedRoster(*s.pendingRosterEdit)
			case toggleStatusUpdatesCommand:
				s.config.HideStatusUpdates = !s.config.HideStatusUpdates
				s.config.Save()
				// Tell the user the current state of the statuses
				if s.config.HideStatusUpdates {
					info(s.term, "Status updates disabled")
				} else {
					info(s.term, "Status updates enabled")
				}
			case confirmCommand:
				s.handleConfirmOrDeny(cmd.User, true /* confirm */)
			case denyCommand:
				s.handleConfirmOrDeny(cmd.User, false /* deny */)
			case addCommand:
				s.conn.SendPresence(cmd.User, "subscribe", "" /* generate id */)
			case msgCommand:
				conversation, ok := s.conversations[cmd.to]
				isEncrypted := ok && conversation.IsEncrypted()
				if cmd.setPromptIsEncrypted != nil {
					cmd.setPromptIsEncrypted <- isEncrypted
				}
				if !isEncrypted && config.ShouldEncryptTo(cmd.to) {
					warn(s.term, fmt.Sprintf("Did not send: no encryption established with %s", cmd.to))
					continue
				}
				var msgs [][]byte
				message := []byte(cmd.msg)
				// Automatically tag all outgoing plaintext
				// messages with a whitespace tag that
				// indicates that we support OTR.
				if config.OTRAutoAppendTag &&
					!bytes.Contains(message, []byte("?OTR")) &&
					(!ok || !conversation.IsEncrypted()) {
					message = append(message, OTRWhitespaceTag...)
				}
				if ok {
					var err error
					validMsgs, err := conversation.Send(message)
					msgs = otr3.Bytes(validMsgs)
					if err != nil {
						alert(s.term, err.Error())
						break
					}
				} else {
					msgs = [][]byte{[]byte(message)}
				}
				for _, message := range msgs {
					s.conn.Send(cmd.to, string(message))
				}
			case otrCommand:
				s.conn.Send(string(cmd.User), QueryMessage)
			case otrInfoCommand:
				info(term, fmt.Sprintf("Your OTR fingerprint is %x", s.privateKey.DefaultFingerprint()))
				for to, conversation := range s.conversations {
					if conversation.IsEncrypted() {
						info(s.term, fmt.Sprintf("Secure session with %s underway:", to))
						printConversationInfo(s, to, conversation)
					}
				}
			case endOTRCommand:
				to := string(cmd.User)
				conversation, ok := s.conversations[to]
				if !ok {
					alert(s.term, "No secure session established")
					break
				}
				msgs, err := conversation.End()
				if err != nil {
					//TODO: error handle
					panic("this should not happen")
				}
				for _, msg := range msgs {
					s.conn.Send(to, string(msg))
				}
				s.input.SetPromptForTarget(cmd.User, false)
				warn(s.term, "OTR conversation ended with "+cmd.User)
			case authQACommand:
				to := string(cmd.User)
				conversation, ok := s.conversations[to]
				if !ok {
					alert(s.term, "Can't authenticate without a secure conversation established")
					break
				}
				var ret []otr3.ValidMessage
				if s.eh[to].waitingForSecret {
					s.eh[to].waitingForSecret = false
					ret, err = conversation.ProvideAuthenticationSecret([]byte(cmd.Secret))
				} else {
					ret, err = conversation.StartAuthenticate(cmd.Question, []byte(cmd.Secret))
				}
				msgs := otr3.Bytes(ret)
				if err != nil {
					alert(s.term, "Error while starting authentication with "+to+": "+err.Error())
				}
				for _, msg := range msgs {
					s.conn.Send(to, string(msg))
				}
			case authOobCommand:
				fpr, err := hex.DecodeString(cmd.Fingerprint)
				if err != nil {
					alert(s.term, fmt.Sprintf("Invalid fingerprint %s - not authenticated", cmd.Fingerprint))
					break
				}
				existing := s.config.UserIdForFingerprint(fpr)
				if len(existing) != 0 {
					alert(s.term, fmt.Sprintf("Fingerprint %s already belongs to %s", cmd.Fingerprint, existing))
					break
				}
				s.config.KnownFingerprints = append(s.config.KnownFingerprints, KnownFingerprint{fingerprint: fpr, UserId: cmd.User})
				s.config.Save()
				info(s.term, fmt.Sprintf("Saved manually verified fingerprint %s for %s", cmd.Fingerprint, cmd.User))
			case awayCommand:
				s.conn.SignalPresence("away")
			case chatCommand:
				s.conn.SignalPresence("chat")
			case dndCommand:
				s.conn.SignalPresence("dnd")
			case xaCommand:
				s.conn.SignalPresence("xa")
			case onlineCommand:
				s.conn.SignalPresence("")
			}
		case rawStanza, ok := <-stanzaChan:
			if !ok {
				warn(term, "Exiting because channel to server closed")
				break MainLoop
			}
			switch stanza := rawStanza.Value.(type) {
			case *xmpp.ClientMessage:
				s.processClientMessage(stanza)
			case *xmpp.ClientPresence:
				s.processPresence(stanza)
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
					alert(term, "Failed to send IQ message: "+err.Error())
				}
			case *xmpp.StreamError:
				var text string
				if len(stanza.Text) > 0 {
					text = stanza.Text
				} else {
					text = fmt.Sprintf("%s", stanza.Any)
				}
				alert(term, "Exiting in response to fatal error from server: "+text)
				break MainLoop
			default:
				info(term, fmt.Sprintf("%s %s", rawStanza.Name, rawStanza.Value))
			}
		}
	}
}
