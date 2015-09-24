//TODO change this package
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	coyconf "github.com/twstrike/coyim/config"
	coyui "github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
)

type Session struct {
	//TODO: This feels bad.
	//it only needs a reference to UI callbacks
	//maybe use the event handler?
	ui coyui.UI

	account string
	conn    *xmpp.Conn
	roster  []xmpp.RosterEntry

	// conversations maps from a JID (without the resource) to an OTR
	// conversation. (Note that unencrypted conversations also pass through
	// OTR.)
	conversations map[string]*otr3.Conversation
	eh            map[string]*eventHandler

	// knownStates maps from a JID (without the resource) to the last known
	// presence state of that contact. It's used to deduping presence
	// notifications.
	knownStates map[string]string
	privateKey  *otr3.PrivateKey
	config      *coyconf.Config

	// timeouts maps from Cookies (from outstanding requests) to the
	// absolute time when that request should timeout.
	timeouts map[xmpp.Cookie]time.Time
	// pendingRosterEdit, if non-nil, contains information about a pending
	// roster edit operation.
	pendingRosterEdit *coyui.RosterEdit
	// pendingRosterChan is the channel over which roster edit information
	// is received.
	pendingRosterChan chan *coyui.RosterEdit
	// pendingSubscribes maps JID with pending subscription requests to the
	// ID if the iq for the reply.
	pendingSubscribes map[string]string
	// lastActionTime is the time at which the user last entered a command,
	// or was last notified.
	lastActionTime time.Time
	sessionHandler sessionHandler

	timeoutTicker *time.Ticker
}

type sessionHandler interface {
	Info(string)
	Warn(string)
	Alert(string)
}

func (c *Session) info(m string) {
	c.sessionHandler.Info(m)
}

func (c *Session) warn(m string) {
	c.sessionHandler.Warn(m)
}

func (c *Session) alert(m string) {
	c.sessionHandler.Alert(m)
}

func (s *Session) readMessages(stanzaChan chan<- xmpp.Stanza) {
	defer close(stanzaChan)

	for {
		stanza, err := s.conn.Next()
		if err != nil {
			s.alert(err.Error())
			return
		}
		stanzaChan <- stanza
	}
}

func (s *Session) WatchStanzas() {
	defer s.Terminate()

	stanzaChan := make(chan xmpp.Stanza)
	go s.readMessages(stanzaChan)

StanzaLoop:
	for {
		select {
		case rawStanza, ok := <-stanzaChan:
			if !ok {
				s.warn("Exiting because channel to server closed")
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

				s.alert("Exiting in response to fatal error from server: " + text)
				break StanzaLoop
			case *xmpp.ClientMessage:
				s.processClientMessage(stanza)
			case *xmpp.ClientPresence:
				ignore, gone := s.processPresence(stanza)
				s.ui.ProcessPresence(stanza, ignore, gone)
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
					s.alert("Failed to send IQ message: " + err.Error())
				}
			default:
				s.info(fmt.Sprintf("%s %s", rawStanza.Name, rawStanza.Value))
			}
		}
	}
}

func (s *Session) rosterReceived() {
	s.ui.RosterReceived(s.roster)
}

func (s *Session) iqReceived(uid string) {
	s.ui.IQReceived(uid)
}

func (s *Session) processIQ(stanza *xmpp.ClientIQ) interface{} {
	buf := bytes.NewBuffer(stanza.Query)
	parser := xml.NewDecoder(buf)
	token, _ := parser.Token()
	if token == nil {
		return nil
	}
	startElem, ok := token.(xml.StartElement)
	if !ok {
		return nil
	}
	switch startElem.Name.Space + " " + startElem.Name.Local {
	case "http://jabber.org/protocol/disco#info query":
		return xmpp.DiscoveryReply{
			Identities: []xmpp.DiscoveryIdentity{
				{
					Category: "client",
					Type:     "pc",
					Name:     s.config.Account,
				},
			},
		}
	case "jabber:iq:version query":
		return xmpp.VersionReply{
			Name:    "testing",
			Version: "version",
			OS:      "none",
		}
	case "jabber:iq:roster query":
		if len(stanza.From) > 0 && stanza.From != s.account {
			s.warn("Ignoring roster IQ from bad address: " + stanza.From)
			return nil
		}
		var roster xmpp.Roster
		if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&roster); err != nil || len(roster.Item) == 0 {
			s.warn("Failed to parse roster push IQ")
			return nil
		}
		entry := roster.Item[0]

		if entry.Subscription == "remove" {
			for i, rosterEntry := range s.roster {
				if rosterEntry.Jid == entry.Jid {
					copy(s.roster[i:], s.roster[i+1:])
					s.roster = s.roster[:len(s.roster)-1]
				}
			}
			return xmpp.EmptyReply{}
		}

		found := false
		for i, rosterEntry := range s.roster {
			if rosterEntry.Jid == entry.Jid {
				s.roster[i] = entry
				found = true
				break
			}
		}

		if !found {
			s.roster = append(s.roster, entry)
			s.iqReceived(entry.Jid)
		}

		return xmpp.EmptyReply{}
	default:
		s.info("Unknown IQ: " + startElem.Name.Space + " " + startElem.Name.Local)
	}

	return nil
}

func (s *Session) handleConfirmOrDeny(jid string, isConfirm bool) {
	id, ok := s.pendingSubscribes[jid]
	if !ok {
		s.warn("No pending subscription from " + jid)
		return
	}
	delete(s.pendingSubscribes, id)
	typ := "unsubscribed"
	if isConfirm {
		typ = "subscribed"
	}
	if err := s.conn.SendPresence(jid, typ, id); err != nil {
		s.warn("Error sending presence stanza: " + err.Error())
	}
}

func (s *Session) newOTRKeys(from string, conversation *otr3.Conversation) {
	s.ui.NewOTRKeys(from, conversation)
}

func (s *Session) otrEnded(uid string) {
	s.ui.OTREnded(uid)
}

func (s *Session) getConversationWith(peer string) *otr3.Conversation {
	if conversation, ok := s.conversations[peer]; ok {
		return conversation
	}

	conversation := &otr3.Conversation{}
	conversation.SetOurKey(s.privateKey)

	//TODO: review this conf
	conversation.Policies.AllowV2()
	conversation.Policies.AllowV3()
	conversation.Policies.SendWhitespaceTag()
	conversation.Policies.WhitespaceStartAKE()
	// conversation.Policies.RequireEncryption()

	s.conversations[peer] = conversation

	//TODO: Why do we need a reference to the event handler in the session?
	eh, ok := s.eh[peer]
	if !ok {
		eh = new(eventHandler)
		conversation.SetSMPEventHandler(eh)
		conversation.SetErrorMessageHandler(eh)
		conversation.SetMessageEventHandler(eh)
		conversation.SetSecurityEventHandler(eh)
		s.eh[peer] = eh
	}

	return conversation
}

func (s *Session) processClientMessage(stanza *xmpp.ClientMessage) {
	from := xmpp.RemoveResourceFromJid(stanza.From)

	if stanza.Type == "error" {
		s.alert("Error reported from " + from + ": " + stanza.Body)
		return
	}

	conversation := s.getConversationWith(from)
	out, toSend, err := conversation.Receive([]byte(stanza.Body))
	encrypted := conversation.IsEncrypted()
	if err != nil {
		s.alert("While processing message from " + from + ": " + err.Error())
		s.conn.Send(stanza.From, ErrorPrefix+"Error processing message")
	}

	for _, msg := range toSend {
		s.conn.Send(stanza.From, string(msg))
	}

	eh, _ := s.eh[from]
	change := eh.consumeSecurityChange()
	switch change {
	case NewKeys:
		s.info(fmt.Sprintf("New OTR session with %s established", from))
		s.newOTRKeys(from, conversation)
	case ConversationEnded:
		s.otrEnded(from)

		// This is probably unsafe without a policy that _forces_ crypto to
		// _everyone_ by default and refuses plaintext. Users might not notice
		// their buddy has ended a session, which they have also ended, and they
		// might send a plain text message. So we should ensure they _want_ this
		// feature and have set it as an explicit preference.
		if s.config.OTRAutoTearDown {
			if s.conversations[from] == nil {
				s.alert(fmt.Sprintf("No secure session established; unable to automatically tear down OTR conversation with %s.", from))
				break
			} else {
				s.info(fmt.Sprintf("%s has ended the secure conversation.", from))
				msgs, err := conversation.End()
				if err != nil {
					//TODO: error handle
					panic("this should not happen")
				}
				for _, msg := range msgs {
					s.conn.Send(from, string(msg))
				}
				s.info(fmt.Sprintf("Secure session with %s has been automatically ended. Messages will be sent in the clear until another OTR session is established.", from))
			}
		} else {
			s.info(fmt.Sprintf("%s has ended the secure conversation. You should do likewise with /otr-end %s", from, from))
		}
	case SMPSecretNeeded:
		s.info(fmt.Sprintf("%s is attempting to authenticate. Please supply mutual shared secret with /otr-auth user secret", from))
		if question := eh.smpQuestion; len(question) > 0 {
			s.info(fmt.Sprintf("%s asks: %s", from, question))
		}
	case SMPComplete:
		s.info(fmt.Sprintf("Authentication with %s successful", from))
		fpr := conversation.GetTheirKey().DefaultFingerprint()
		if len(s.config.UserIdForFingerprint(fpr)) == 0 {
			s.config.KnownFingerprints = append(s.config.KnownFingerprints, coyconf.KnownFingerprint{Fingerprint: fpr, UserId: from})
		}
		s.config.Save()
	case SMPFailed:
		s.alert(fmt.Sprintf("Authentication with %s failed", from))
	}

	if len(out) == 0 {
		return
	}

	detectedOTRVersion := 0
	// We don't need to alert about tags encoded inside of messages that are
	// already encrypted with OTR
	whitespaceTagLength := len(coyui.OTRWhitespaceTagStart) + len(coyui.OTRWhiteSpaceTagV1)
	if !encrypted && len(out) >= whitespaceTagLength {
		whitespaceTag := out[len(out)-whitespaceTagLength:]
		if bytes.Equal(whitespaceTag[:len(coyui.OTRWhitespaceTagStart)], coyui.OTRWhitespaceTagStart) {
			if bytes.HasSuffix(whitespaceTag, coyui.OTRWhiteSpaceTagV1) {
				s.info(fmt.Sprintf("%s appears to support OTRv1. You should encourage them to upgrade their OTR client!", from))
				detectedOTRVersion = 1
			}
			if bytes.HasSuffix(whitespaceTag, coyui.OTRWhiteSpaceTagV2) {
				detectedOTRVersion = 2
			}
			if bytes.HasSuffix(whitespaceTag, coyui.OTRWhiteSpaceTagV3) {
				detectedOTRVersion = 3
			}
		}
	}

	if s.config.OTRAutoStartSession && detectedOTRVersion >= 2 {
		s.info(fmt.Sprintf("%s appears to support OTRv%d. We are attempting to start an OTR session with them.", from, detectedOTRVersion))
		s.conn.Send(from, QueryMessage)
	} else if s.config.OTRAutoStartSession && detectedOTRVersion == 1 {
		s.info(fmt.Sprintf("%s appears to support OTRv%d. You should encourage them to upgrade their OTR client!", from, detectedOTRVersion))
	}

	var timestamp string
	var messageTime time.Time
	if stanza.Delay != nil && len(stanza.Delay.Stamp) > 0 {
		// An XEP-0203 Delayed Delivery <delay/> element exists for
		// this message, meaning that someone sent it while we were
		// offline. Let's show the timestamp for when the message was
		// sent, rather than time.Now().
		messageTime, err = time.Parse(time.RFC3339, stanza.Delay.Stamp)
		if err != nil {
			s.alert("Can not parse Delayed Delivery timestamp, using quoted string instead.")
			timestamp = fmt.Sprintf("%q", stanza.Delay.Stamp)
		}
	} else {
		messageTime = time.Now()
	}
	if len(timestamp) == 0 {
		timestamp = messageTime.Format(time.Stamp)
	}

	s.messageReceived(from, timestamp, encrypted, out)
}

func (s *Session) messageReceived(from, timestamp string, encrypted bool, message []byte) {
	s.ui.MessageReceived(from, timestamp, encrypted, message)
	s.maybeNotify()
}

func (s *Session) maybeNotify() {
	now := time.Now()
	idleThreshold := s.config.IdleSecondsBeforeNotification
	if idleThreshold == 0 {
		idleThreshold = 60
	}
	notifyTime := s.lastActionTime.Add(time.Duration(idleThreshold) * time.Second)
	if now.Before(notifyTime) {
		return
	}

	s.lastActionTime = now
	if len(s.config.NotifyCommand) == 0 {
		return
	}

	cmd := exec.Command(s.config.NotifyCommand[0], s.config.NotifyCommand[1:]...)
	go func() {
		if err := cmd.Run(); err != nil {
			s.alert("Failed to run notify command: " + err.Error())
		}
	}()
}

func (s *Session) processPresence(stanza *xmpp.ClientPresence) (ignore, gone bool) {

	switch stanza.Type {
	case "subscribe":
		// This is a subscription request
		jid := xmpp.RemoveResourceFromJid(stanza.From)
		s.pendingSubscribes[jid] = stanza.Id
		ignore = true
		return
	case "unavailable":
		gone = true
	case "":
		break
	default:
		ignore = true
		return
	}

	from := xmpp.RemoveResourceFromJid(stanza.From)

	if gone {
		if _, ok := s.knownStates[from]; !ok {
			// They've gone, but we never knew they were online.
			ignore = true
			return
		}
		delete(s.knownStates, from)
	} else {
		if _, ok := s.knownStates[from]; !ok && coyui.IsAwayStatus(stanza.Show) {
			// Skip people who are initially away.
			ignore = true
			return
		}

		if lastState, ok := s.knownStates[from]; ok && lastState == stanza.Show {
			// No change. Ignore.
			ignore = true
			return
		}
		s.knownStates[from] = stanza.Show
	}

	return
}

func (s *Session) awaitVersionReply(ch <-chan xmpp.Stanza, user string) {
	stanza, ok := <-ch
	if !ok {
		s.warn("Version request to " + user + " timed out")
		return
	}
	reply, ok := stanza.Value.(*xmpp.ClientIQ)
	if !ok {
		s.warn("Version request to " + user + " resulted in bad reply type")
		return
	}

	if reply.Type == "error" {
		s.warn("Version request to " + user + " resulted in XMPP error")
		return
	} else if reply.Type != "result" {
		s.warn("Version request to " + user + " resulted in response with unknown type: " + reply.Type)
		return
	}

	buf := bytes.NewBuffer(reply.Query)
	var versionReply xmpp.VersionReply
	if err := xml.NewDecoder(buf).Decode(&versionReply); err != nil {
		s.warn("Failed to parse version reply from " + user + ": " + err.Error())
		return
	}

	s.info(fmt.Sprintf("Version reply from %s: %#v", user, versionReply))
}

func (s *Session) WatchTimeout() {
	s.timeoutTicker = time.NewTicker(1 * time.Second)

	for now := range s.timeoutTicker.C {
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

func (s *Session) WatchRosterEvents() {
	defer s.Terminate()

	s.info("Fetching roster")

	rosterReply, _, err := s.conn.RequestRoster()
	if err != nil {
		s.alert("Failed to request roster: " + err.Error())
		return
	}

	//TODO: not sure if this belongs here
	s.conn.SignalPresence("")

RosterLoop:
	for {
		select {
		case rosterStanza, ok := <-rosterReply:
			if !ok {
				s.alert("Failed to read roster: " + err.Error())
				break RosterLoop
			}

			if s.roster, err = xmpp.ParseRoster(rosterStanza); err != nil {
				s.alert("Failed to parse roster: " + err.Error())
				break RosterLoop
			}

			s.rosterReceived()

		case edit := <-s.pendingRosterChan:
			if !edit.IsComplete {
				//TODO: this is specific to CLI
				s.info("Please edit " + edit.FileName + " and run /rostereditdone when complete")
				s.pendingRosterEdit = edit
				continue
			}

			if s.processEditedRoster(edit) {
				s.pendingRosterEdit = nil
			} else {
				//TODO: this is specific to CLI
				s.alert("Please reedit file and run /rostereditdone again")
			}
		}
	}
}

func (s *Session) Terminate() {
	s.timeoutTicker.Stop()
	s.timeoutTicker = nil

	s.ui.Disconnected()
}

// editRoster runs in a goroutine and writes the roster to a file that the user
// can edit.
func (s *Session) editRoster(roster []xmpp.RosterEntry) {
	// In case the editor rewrites the file, we work inside a temp
	// directory.
	dir, err := ioutil.TempDir("" /* system default temp dir */, "xmpp-client")
	if err != nil {
		s.alert("Failed to create temp dir to edit roster: " + err.Error())
		return
	}

	mode, err := os.Stat(dir)
	if err != nil || mode.Mode()&os.ModePerm != 0700 {
		panic("broken system libraries gave us an insecure temp dir")
	}

	fileName := filepath.Join(dir, "roster")
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		s.alert("Failed to create temp file: " + err.Error())
		return
	}

	io.WriteString(f, `# Use this file to edit your roster.
# The file is tab deliminated and you need to preserve that. Otherwise you
# can delete lines to remove roster entries, add lines to subscribe (only
# the account is needed when adding a line) and change lines to change the
# corresponding entry.

# Once you are done, use the /rostereditdone command to process the result.

# Since there are multiple levels of unspecified character encoding, we give up
# and hex escape anything outside of printable ASCII in "\x01" form.

`)

	// Calculate the number of tabs which covers the longest escaped JID.
	maxLen := 0
	escapedJids := make([]string, len(roster))
	for i, item := range roster {
		escapedJids[i] = coyui.EscapeNonASCII(item.Jid)
		if l := len(escapedJids[i]); l > maxLen {
			maxLen = l
		}
	}
	tabs := (maxLen + 7) / 8

	for i, item := range s.roster {
		line := escapedJids[i]
		tabsUsed := len(escapedJids[i]) / 8

		if len(item.Name) > 0 || len(item.Group) > 0 {
			// We're going to put something else on the line to tab
			// across to the next column.
			for i := 0; i < tabs-tabsUsed; i++ {
				line += "\t"
			}
		}

		if len(item.Name) > 0 {
			line += "name:" + coyui.EscapeNonASCII(item.Name)
			if len(item.Group) > 0 {
				line += "\t"
			}
		}

		for j, group := range item.Group {
			if j > 0 {
				line += "\t"
			}
			line += "group:" + coyui.EscapeNonASCII(group)
		}
		line += "\n"
		io.WriteString(f, line)
	}
	f.Close()

	s.pendingRosterChan <- &coyui.RosterEdit{
		FileName: fileName,
		Roster:   roster,
	}
}

func (s *Session) loadEditedRoster(edit coyui.RosterEdit) {
	contents, err := ioutil.ReadFile(edit.FileName)
	if err != nil {
		s.alert("Failed to load edited roster: " + err.Error())
		return
	}
	os.Remove(edit.FileName)
	os.Remove(filepath.Dir(edit.FileName))

	edit.IsComplete = true
	edit.Contents = contents
	s.pendingRosterChan <- &edit
}

func (s *Session) processEditedRoster(edit *coyui.RosterEdit) bool {
	parsedRoster := make(map[string]xmpp.RosterEntry)
	lines := bytes.Split(edit.Contents, coyui.NewLine)
	tab := []byte{'\t'}

	// Parse roster entries from the file.
	for i, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := bytes.Split(line, tab)

		var entry xmpp.RosterEntry
		var err error

		if entry.Jid, err = coyui.UnescapeNonASCII(string(string(parts[0]))); err != nil {
			s.alert(fmt.Sprintf("Failed to parse JID on line %d: %s", i+1, err))
			return false
		}
		for _, part := range parts[1:] {
			if len(part) == 0 {
				continue
			}

			pos := bytes.IndexByte(part, ':')
			if pos == -1 {
				s.alert(fmt.Sprintf("Failed to find colon in item on line %d", i+1))
				return false
			}

			typ := string(part[:pos])
			value, err := coyui.UnescapeNonASCII(string(part[pos+1:]))
			if err != nil {
				s.alert(fmt.Sprintf("Failed to unescape item on line %d: %s", i+1, err))
				return false
			}

			switch typ {
			case "name":
				if len(entry.Name) > 0 {
					s.alert(fmt.Sprintf("Multiple names given for contact on line %d", i+1))
					return false
				}
				entry.Name = value
			case "group":
				if len(value) > 0 {
					entry.Group = append(entry.Group, value)
				}
			default:
				s.alert(fmt.Sprintf("Unknown item tag '%s' on line %d", typ, i+1))
				return false
			}
		}

		parsedRoster[entry.Jid] = entry
	}

	// Now diff them from the original roster
	var toDelete []string
	var toEdit []xmpp.RosterEntry
	var toAdd []xmpp.RosterEntry

	for _, entry := range edit.Roster {
		newEntry, ok := parsedRoster[entry.Jid]
		if !ok {
			toDelete = append(toDelete, entry.Jid)
			continue
		}
		if newEntry.Name != entry.Name || !setEqual(newEntry.Group, entry.Group) {
			toEdit = append(toEdit, newEntry)
		}
	}

NextAdd:
	for jid, newEntry := range parsedRoster {
		for _, entry := range edit.Roster {
			if entry.Jid == jid {
				continue NextAdd
			}
		}
		toAdd = append(toAdd, newEntry)
	}

	for _, jid := range toDelete {
		s.info("Deleting roster entry for " + jid)
		_, _, err := s.conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
			Item: xmpp.RosterRequestItem{
				Jid:          jid,
				Subscription: "remove",
			},
		})
		if err != nil {
			s.alert("Failed to remove roster entry: " + err.Error())
		}

		// Filter out any known fingerprints.
		newKnownFingerprints := make([]coyconf.KnownFingerprint, 0, len(s.config.KnownFingerprints))
		for _, fpr := range s.config.KnownFingerprints {
			if fpr.UserId == jid {
				continue
			}
			newKnownFingerprints = append(newKnownFingerprints, fpr)
		}
		s.config.KnownFingerprints = newKnownFingerprints
		s.config.Save()
	}

	for _, entry := range toEdit {
		s.info("Updating roster entry for " + entry.Jid)
		_, _, err := s.conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
			Item: xmpp.RosterRequestItem{
				Jid:   entry.Jid,
				Name:  entry.Name,
				Group: entry.Group,
			},
		})
		if err != nil {
			s.alert("Failed to update roster entry: " + err.Error())
		}
	}

	for _, entry := range toAdd {
		s.info("Adding roster entry for " + entry.Jid)
		_, _, err := s.conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
			Item: xmpp.RosterRequestItem{
				Jid:   entry.Jid,
				Name:  entry.Name,
				Group: entry.Group,
			},
		})
		if err != nil {
			s.alert("Failed to add roster entry: " + err.Error())
		}
	}

	return true
}

func setEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

EachValue:
	for _, v := range a {
		for _, v2 := range b {
			if v == v2 {
				continue EachValue
			}
		}
		return false
	}

	return true
}
