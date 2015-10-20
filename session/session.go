package session

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
)

type connStatus int

const (
	DISCONNECTED connStatus = iota
	CONNECTING
	CONNECTED
)

type Session struct {
	Conn       *xmpp.Conn
	R          *roster.List
	ConnStatus connStatus

	// conversations maps from a JID (without the resource) to an OTR
	// conversation. (Note that unencrypted conversations also pass through
	// OTR.)
	Conversations   map[string]*otr3.Conversation
	OtrEventHandler map[string]*event.OtrEventHandler

	PrivateKey *otr3.PrivateKey
	Config     *config.Config

	// timeouts maps from Cookies (from outstanding requests) to the
	// absolute time when that request should timeout.
	Timeouts map[xmpp.Cookie]time.Time

	// lastActionTime is the time at which the user last entered a command,
	// or was last notified.
	LastActionTime      time.Time
	SessionEventHandler SessionEventHandler

	timeoutTicker *time.Ticker
	rosterCookie  xmpp.Cookie

	Account interface{}
}

func NewSession(c *config.Config) *Session {
	s := &Session{
		Config: c,

		R:               roster.New(),
		Conversations:   make(map[string]*otr3.Conversation),
		OtrEventHandler: make(map[string]*event.OtrEventHandler),
		PrivateKey:      new(otr3.PrivateKey),
		LastActionTime:  time.Now(),
	}

	s.PrivateKey.Parse(c.PrivateKey)

	return s
}

func (s *Session) info(m string) {
	s.SessionEventHandler.Info(m)
}

func (s *Session) warn(m string) {
	s.SessionEventHandler.Warn(m)
}

func (s *Session) alert(m string) {
	s.SessionEventHandler.Alert(m)
}

func (s *Session) readMessages(stanzaChan chan<- xmpp.Stanza) {
	defer close(stanzaChan)

	for {
		stanza, err := s.Conn.Next()
		if err != nil {
			s.alert(err.Error())
			return
		}

		stanzaChan <- stanza
	}
}

func (s *Session) receivedStreamError(stanza *xmpp.StreamError) bool {
	var text string

	if len(stanza.Text) > 0 {
		text = stanza.Text
	} else {
		text = fmt.Sprintf("%s", stanza.Any)
	}

	s.alert("Exiting in response to fatal error from server: " + text)
	return false
}

func (s *Session) receivedClientMessage(stanza *xmpp.ClientMessage) bool {
	s.processClientMessage(stanza)
	return true
}

func (s *Session) receivedClientPresence(stanza *xmpp.ClientPresence) bool {
	switch stanza.Type {
	case "subscribe":
		s.R.SubscribeRequest(stanza.From, stanza.Id)
		s.SessionEventHandler.SubscriptionRequest(s, xmpp.RemoveResourceFromJid(stanza.From))
	case "unavailable":
		if s.R.PeerBecameUnavailable(stanza.From) &&
			!s.Config.HideStatusUpdates {
			s.SessionEventHandler.ProcessPresence(stanza.From, stanza.To, stanza.Show, stanza.Status, true)
		}
	case "":
		if s.R.PeerPresenceUpdate(stanza.From, stanza.Show, stanza.Status) &&
			!s.Config.HideStatusUpdates {
			s.SessionEventHandler.ProcessPresence(stanza.From, stanza.To, stanza.Show, stanza.Status, false)
		}
	case "subscribed":
		s.R.Subscribed(stanza.From)
		s.SessionEventHandler.Subscribed(xmpp.RemoveResourceFromJid(stanza.To), xmpp.RemoveResourceFromJid(stanza.From))
	case "unsubscribe":
		s.R.Unsubscribed(stanza.From)
		s.SessionEventHandler.Unsubscribe(xmpp.RemoveResourceFromJid(stanza.To), xmpp.RemoveResourceFromJid(stanza.From))
	case "unsubscribed":
		// Ignore
	default:
		s.info(fmt.Sprintf("unrecognized presence: %#v", stanza))
	}
	return true
}

func (s *Session) receivedClientIQ(stanza *xmpp.ClientIQ) bool {
	if stanza.Type == "get" || stanza.Type == "set" {
		reply := s.processIQ(stanza)
		if reply == nil {
			reply = xmpp.ErrorReply{
				Type:  "cancel",
				Error: xmpp.ErrorBadRequest{},
			}
		}

		if err := s.Conn.SendIQReply(stanza.From, "result", stanza.Id, reply); err != nil {
			s.alert("Failed to send IQ message: " + err.Error())
		}
		return true
	}
	s.info(fmt.Sprintf("unrecognized iq: %#v", stanza))
	return true
}

func (s *Session) receiveStanza(stanzaChan chan xmpp.Stanza) bool {
	select {
	case rawStanza, ok := <-stanzaChan:
		if !ok {
			s.warn("Exiting because channel to server closed")
			return false
		}

		switch stanza := rawStanza.Value.(type) {
		case *xmpp.StreamError:
			return s.receivedStreamError(stanza)
		case *xmpp.ClientMessage:
			return s.receivedClientMessage(stanza)
		case *xmpp.ClientPresence:
			return s.receivedClientPresence(stanza)
		case *xmpp.ClientIQ:
			return s.receivedClientIQ(stanza)
		default:
			s.info(fmt.Sprintf("%s %s", rawStanza.Name, rawStanza.Value))
			return true
		}
	}
}

func (s *Session) WatchStanzas() {
	defer s.Close()

	stanzaChan := make(chan xmpp.Stanza)
	go s.readMessages(stanzaChan)
	for s.receiveStanza(stanzaChan) {
	}
}

func (s *Session) rosterReceived() {
	s.SessionEventHandler.RosterReceived(s)
}

func (s *Session) iqReceived(uid string) {
	s.SessionEventHandler.IQReceived(uid)
}

func (s *Session) receivedIQDiscoInfo() xmpp.DiscoveryReply {
	return xmpp.DiscoveryReply{
		Identities: []xmpp.DiscoveryIdentity{
			{
				Category: "client",
				Type:     "pc",
				Name:     s.Config.Account,
			},
		},
	}
}

func (s *Session) receivedIQVersion() xmpp.VersionReply {
	return xmpp.VersionReply{
		Name:    "testing",
		Version: "version",
		OS:      "none",
	}
}

func (s *Session) receivedIQRosterQuery(stanza *xmpp.ClientIQ) interface{} {
	// TODO: this code can only be hit by a iq get or iq set. Is iq get actually reasonable for this?
	// No, a get should likely not even arrive here
	// TODO: we should deal with "ask" attributes here

	if len(stanza.From) > 0 && xmpp.RemoveResourceFromJid(stanza.From) != s.Config.Account {
		s.warn("Ignoring roster IQ from bad address: " + stanza.From)
		return nil
	}
	var rst xmpp.Roster
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&rst); err != nil || len(rst.Item) == 0 {
		s.warn("Failed to parse roster push IQ")
		return nil
	}

	// TODO: this is incorrect - you can get more than one roster Item
	entry := rst.Item[0]

	if entry.Subscription == "remove" {
		s.R.Remove(entry.Jid)
		return xmpp.EmptyReply{}
	}

	if s.R.AddOrMerge(roster.PeerFrom(entry)) {
		s.iqReceived(entry.Jid)
	}

	return xmpp.EmptyReply{}
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
		return s.receivedIQDiscoInfo()
	case "jabber:iq:version query":
		return s.receivedIQVersion()
	case "jabber:iq:roster query":
		return s.receivedIQRosterQuery(stanza)
	default:
		s.info("Unknown IQ: " + startElem.Name.Space + " " + startElem.Name.Local)
	}

	return nil
}

func (s *Session) HandleConfirmOrDeny(jid string, isConfirm bool) {
	id, ok := s.R.RemovePendingSubscribe(jid)
	if !ok {
		s.warn("No pending subscription from " + jid)
		return
	}
	typ := "unsubscribed"
	if isConfirm {
		typ = "subscribed"
	}
	if err := s.Conn.SendPresence(jid, typ, id); err != nil {
		s.warn("Error sending presence stanza: " + err.Error())
	}
}

func (s *Session) newOTRKeys(from string, conversation *otr3.Conversation) {
	s.SessionEventHandler.NewOTRKeys(from, conversation)
}

func (s *Session) otrEnded(uid string) {
	s.SessionEventHandler.OTREnded(uid)
}

func (s *Session) newConversation(peer string) *otr3.Conversation {
	conversation := &otr3.Conversation{}
	conversation.SetOurKey(s.PrivateKey)

	//TODO: review this conf
	conversation.Policies.AllowV2()
	conversation.Policies.AllowV3()
	conversation.Policies.SendWhitespaceTag()
	conversation.Policies.WhitespaceStartAKE()
	// conversation.Policies.RequireEncryption()

	return conversation
}

func (s *Session) GetConversationWith(peer string) *otr3.Conversation {
	if conversation, ok := s.Conversations[peer]; ok {
		return conversation
	}

	conversation := s.newConversation(peer)
	s.Conversations[peer] = conversation

	//TODO: Why do we need a reference to the event handler in the session?
	eh, ok := s.OtrEventHandler[peer]
	if !ok {
		eh = new(event.OtrEventHandler)
		conversation.SetSMPEventHandler(eh)
		conversation.SetErrorMessageHandler(eh)
		conversation.SetMessageEventHandler(eh)
		conversation.SetSecurityEventHandler(eh)
		s.OtrEventHandler[peer] = eh
	}

	return conversation
}

func (s *Session) processClientMessage(stanza *xmpp.ClientMessage) {
	from := xmpp.RemoveResourceFromJid(stanza.From)

	if stanza.Type == "error" {
		s.alert("Error reported from " + from + ": " + stanza.Body)
		return
	}

	conversation := s.GetConversationWith(from)
	out, toSend, err := conversation.Receive([]byte(stanza.Body))
	encrypted := conversation.IsEncrypted()
	if err != nil {
		s.alert("While processing message from " + from + ": " + err.Error())
		s.Conn.Send(stanza.From, event.ErrorPrefix+"Error processing message")
	}

	for _, msg := range toSend {
		s.Conn.Send(stanza.From, string(msg))
	}

	//TODO: refactor
	//This consumes the security change from OtrEventHandler and trigger an event
	//on the SessionEventHandler. Why not having a single event handler?
	eh, _ := s.OtrEventHandler[from]
	change := eh.ConsumeSecurityChange()
	switch change {
	case event.NewKeys:
		s.info(fmt.Sprintf("New OTR session with %s established", from))
		s.newOTRKeys(from, conversation)
	case event.ConversationEnded:
		s.otrEnded(from)

		// TODO: twstrike/otr3 does not allow sending messages after the channel has
		// been terminated, so this should not be a problem.
		// This is probably unsafe without a policy that _forces_ crypto to
		// _everyone_ by default and refuses plaintext. Users might not notice
		// their buddy has ended a session, which they have also ended, and they
		// might send a plain text message. So we should ensure they _want_ this
		// feature and have set it as an explicit preference.
		if s.Config.OTRAutoTearDown {
			if s.Conversations[from] == nil {
				s.alert(fmt.Sprintf("No secure session established; unable to automatically tear down OTR conversation with %s.", from))
				break
			} else {
				s.info(fmt.Sprintf("%s has ended the secure conversation.", from))
				s.TerminateConversationWith(from)
				if err != nil {
					s.info(fmt.Sprintf("Unable to automatically tear down OTR conversation with %s: %s\n", from, err.Error()))
					break
				}

				s.info(fmt.Sprintf("Secure session with %s has been automatically ended. Messages will be sent in the clear until another OTR session is established.", from))
			}
		} else {
			s.info(fmt.Sprintf("%s has ended the secure conversation. You should do likewise with /otr-end %s", from, from))
		}
	case event.SMPSecretNeeded:
		s.info(fmt.Sprintf("%s is attempting to authenticate. Please supply mutual shared secret with /otr-auth user secret", from))
		if question := eh.SmpQuestion; len(question) > 0 {
			s.info(fmt.Sprintf("%s asks: %s", from, question))
		}
	case event.SMPComplete:
		s.info(fmt.Sprintf("Authentication with %s successful", from))
		fpr := conversation.GetTheirKey().DefaultFingerprint()
		if len(s.Config.UserIdForFingerprint(fpr)) == 0 {
			s.Config.AddFingerprint(fpr, from)
			s.Config.Save()
		}
	case event.SMPFailed:
		s.alert(fmt.Sprintf("Authentication with %s failed", from))
	}

	if len(out) == 0 {
		return
	}

	//TODO: remove because twstrike/otr3 already handles whitespace tags
	s.processWhitespaceTag(encrypted, out, from)

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
	s.SessionEventHandler.MessageReceived(s, from, timestamp, encrypted, message)
	s.maybeNotify()
}

func (s *Session) maybeNotify() {
	now := time.Now()
	idleThreshold := s.Config.IdleSecondsBeforeNotification
	if idleThreshold == 0 {
		idleThreshold = 60
	}
	notifyTime := s.LastActionTime.Add(time.Duration(idleThreshold) * time.Second)
	if now.Before(notifyTime) {
		return
	}

	s.LastActionTime = now
	if len(s.Config.NotifyCommand) == 0 {
		return
	}

	cmd := exec.Command(s.Config.NotifyCommand[0], s.Config.NotifyCommand[1:]...)
	go func() {
		if err := cmd.Run(); err != nil {
			s.alert("Failed to run notify command: " + err.Error())
		}
	}()
}

func isAwayStatus(status string) bool {
	switch status {
	case "xa", "away":
		return true
	}
	return false
}

func (s *Session) AwaitVersionReply(ch <-chan xmpp.Stanza, user string) {
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
	s.Timeouts = make(map[xmpp.Cookie]time.Time)
	s.timeoutTicker = time.NewTicker(1 * time.Second)

	for now := range s.timeoutTicker.C {
		haveExpired := false
		for _, expiry := range s.Timeouts {
			if now.After(expiry) {
				haveExpired = true
				break
			}
		}

		if !haveExpired {
			continue
		}

		newTimeouts := make(map[xmpp.Cookie]time.Time)
		for cookie, expiry := range s.Timeouts {
			if now.After(expiry) {
				s.Conn.Cancel(cookie)
			} else {
				newTimeouts[cookie] = expiry
			}
		}

		s.Timeouts = newTimeouts
	}
}

func (s *Session) WatchRosterEvents() {
	defer s.Close()

	//TODO: not sure if this belongs here
	s.Conn.SignalPresence("")
	s.info("Fetching roster")

	rosterReply, c, err := s.Conn.RequestRoster()
	if err != nil {
		s.alert("Failed to request roster: " + err.Error())
		return
	}

	s.rosterCookie = c

	for {
		select {
		case rosterStanza, ok := <-rosterReply:
			if !ok {
				s.alert("Failed to read roster: " + err.Error())
				return
			}

			rst, err := xmpp.ParseRoster(rosterStanza)

			if err != nil {
				s.alert("Failed to parse roster: " + err.Error())
				return
			}

			for _, rr := range rst {
				s.R.AddOrMerge(roster.PeerFrom(rr))
			}

			s.rosterReceived()
			s.info("Roster received")

			//TODO: this is CLI specific
			//case edit := <-s.PendingRosterChan:
			//	if !edit.IsComplete {
			//		//TODO: this is specific to CLI
			//		s.info("Please edit " + edit.FileName + " and run /rostereditdone when complete")
			//		s.PendingRosterEdit = edit
			//		continue
			//	}

			//	if s.processEditedRoster(edit) {
			//		s.PendingRosterEdit = nil
			//	} else {
			//		//TODO: this is specific to CLI
			//		s.alert("Please reedit file and run /rostereditdone again")
			//	}
		}
	}
}

func (s *Session) Connect(password string, registerCallback xmpp.FormCallback) error {
	if s.ConnStatus != DISCONNECTED {
		return nil
	}

	s.ConnStatus = CONNECTING

	conn, err := config.NewXMPPConn(s.Config, password, registerCallback, logger{})
	if err != nil {
		s.alert(err.Error())
		s.ConnStatus = DISCONNECTED
		return err
	}

	s.Conn = conn
	s.ConnStatus = CONNECTED

	go s.WatchTimeout()
	go s.WatchRosterEvents()
	go s.WatchStanzas()

	return nil
}

func (s *Session) EncryptAndSendTo(peer string, message string) error {
	conversation := s.GetConversationWith(peer)
	toSend, err := conversation.Send(otr3.ValidMessage(message))
	if err != nil {
		return err
	}

	for _, m := range toSend {
		err := s.SendMessageTo(peer, string(m))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) SendMessageTo(peer string, message string) error {
	return s.Conn.Send(peer, message)
}

func (s *Session) StartEncryptedChatWith(peer string) error {
	conversation := s.GetConversationWith(peer)
	return s.SendMessageTo(peer, string(conversation.QueryMessage()))
}

func (s *Session) TerminateConversationWith(peer string) error {
	//Do not use GetConversationWith because we dont want to create a new conversation just to destroy it
	c, ok := s.Conversations[peer]
	if !ok {
		return nil
	}

	msgs, err := c.End()
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		err := s.Conn.Send(peer, string(msg))
		if err != nil {
			return err
		}
	}

	//conversation.Wipe()
	delete(s.Conversations, peer)

	return nil
}

func (s *Session) terminateConversations() {
	for peer := range s.Conversations {
		//TODO: errors
		s.TerminateConversationWith(peer)
	}
}

func (s *Session) Close() {
	//TODO: what should be done it states == CONNECTING?
	if s.ConnStatus == DISCONNECTED {
		return
	}

	s.terminateConversations()

	//Stops all
	s.Conn.Cancel(s.rosterCookie)
	s.timeoutTicker.Stop()

	s.Conn.Close()
	s.ConnStatus = DISCONNECTED

	//TODO Should we hide all contacts when the account is disconnected?
	// It wont show a "please connect to view your roster" message
	s.R.Clear()
	s.SessionEventHandler.RosterReceived(s)

	s.SessionEventHandler.Disconnected()
}
