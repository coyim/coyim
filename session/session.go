package session

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
)

type connStatus int

// These constants represent the different connection states
const (
	DISCONNECTED connStatus = iota
	CONNECTING
	CONNECTED
)

// Session contains information about one specific connection
type Session struct {
	Conn             *xmpp.Conn
	ConnectionLogger io.Writer
	R                *roster.List
	ConnStatus       connStatus

	// conversations maps from a JID (without the resource) to an OTR
	// conversation. (Note that unencrypted conversations also pass through
	// OTR.)
	Conversations   map[string]*otr3.Conversation
	OtrEventHandler map[string]*event.OtrEventHandler

	PrivateKey     *otr3.PrivateKey
	Config         *config.Accounts
	CurrentAccount *config.Account

	// timeouts maps from Cookies (from outstanding requests) to the
	// absolute time when that request should timeout.
	Timeouts map[xmpp.Cookie]time.Time

	// lastActionTime is the time at which the user last entered a command,
	// or was last notified.
	LastActionTime      time.Time
	SessionEventHandler EventHandler

	timeoutTicker *time.Ticker
	rosterCookie  xmpp.Cookie

	subscribers struct {
		sync.RWMutex
		subs []chan<- interface{}
	}

	SaveConfiguration func()

	GroupDelimiter string

	xmppLogger io.Writer
}

// NewSession creates a new session from the given config
func NewSession(c *config.Accounts, cu *config.Account) *Session {
	s := &Session{
		Config:         c,
		CurrentAccount: cu,

		R:               roster.New(),
		Conversations:   make(map[string]*otr3.Conversation),
		OtrEventHandler: make(map[string]*event.OtrEventHandler),
		PrivateKey:      new(otr3.PrivateKey),
		LastActionTime:  time.Now(),

		xmppLogger: openLogFile(c.RawLogFile),
	}

	s.PrivateKey.Parse(cu.PrivateKey)

	return s
}

//TODO: error
func openLogFile(logFile string) io.Writer {
	if len(logFile) == 0 {
		return nil
	}

	log.Println("Logging XMPP messages to:", logFile)

	rawLog, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Println("Failed to open log file.", err)
		//return nil, errors.New("Failed to open raw log file: " + err.Error())
		return nil
	}

	return rawLog
}

func (s *Session) info(m string) {
	s.publishEvent(LogEvent{
		Level:   Info,
		Message: m,
	})
}

func (s *Session) warn(m string) {
	s.publishEvent(LogEvent{
		Level:   Warn,
		Message: m,
	})
}

func (s *Session) alert(m string) {
	s.publishEvent(LogEvent{
		Level:   Alert,
		Message: m,
	})
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

func either(l, r string) string {
	if l == "" {
		return r
	}
	return l
}

func (s *Session) receivedClientPresence(stanza *xmpp.ClientPresence) bool {
	switch stanza.Type {
	case "subscribe":
		s.R.SubscribeRequest(stanza.From, either(stanza.ID, "0000"), s.CurrentAccount.ID())
		s.publishPeerEvent(
			SubscriptionRequest,
			xmpp.RemoveResourceFromJid(stanza.From),
		)
	case "unavailable":
		if !s.R.PeerBecameUnavailable(stanza.From) {
			return true
		}

		s.publishEvent(PresenceEvent{
			Session:        s,
			ClientPresence: stanza,
			Gone:           true,
		})
	case "":
		if !s.R.PeerPresenceUpdate(stanza.From, stanza.Show, stanza.Status, s.CurrentAccount.ID()) {
			return true
		}

		s.publishEvent(PresenceEvent{
			Session:        s,
			ClientPresence: stanza,
			Gone:           false,
		})
	case "subscribed":
		s.R.Subscribed(stanza.From)
		s.publishPeerEvent(
			Subscribed,
			xmpp.RemoveResourceFromJid(stanza.From),
		)
	case "unsubscribe":
		s.R.Unsubscribed(stanza.From)
		s.publishPeerEvent(
			Unsubscribe,
			xmpp.RemoveResourceFromJid(stanza.From),
		)
	case "unsubscribed":
		// Ignore
	case "error":
		s.warn(fmt.Sprintf("Got a presence error from %s: %s\n", stanza.From, stanza.Error))
	default:
		s.info(fmt.Sprintf("unrecognized presence: %#v", stanza))
	}
	return true
}

func (s *Session) receivedClientIQ(stanza *xmpp.ClientIQ) bool {
	if stanza.Type == "get" || stanza.Type == "set" {
		reply, ignore := s.processIQ(stanza)
		if ignore {
			return true
		}

		if reply == nil {
			reply = xmpp.ErrorReply{
				Type:  "cancel",
				Error: xmpp.ErrorBadRequest{},
			}
		}

		if err := s.Conn.SendIQReply(stanza.From, "result", stanza.ID, reply); err != nil {
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

// WatchStanzas listens to XMPP stanzas and reacts on them
func (s *Session) WatchStanzas() {
	defer s.Close()

	stanzaChan := make(chan xmpp.Stanza)
	go s.readMessages(stanzaChan)
	for s.receiveStanza(stanzaChan) {
	}
}

func (s *Session) rosterReceived() {
	s.publish(RosterReceived)
}

func (s *Session) iqReceived(uid string) {
	s.publishPeerEvent(IQReceived, uid)
}

func (s *Session) receivedIQDiscoInfo() xmpp.DiscoveryReply {
	return xmpp.DiscoveryReply{
		Identities: []xmpp.DiscoveryIdentity{
			{
				Category: "client",
				Type:     "pc",
				Name:     s.CurrentAccount.Account,
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

func (s *Session) receivedIQRosterQuery(stanza *xmpp.ClientIQ) (ret interface{}, ignore bool) {
	// TODO: we should deal with "ask" attributes here

	if len(stanza.From) > 0 && !s.CurrentAccount.Is(stanza.From) {
		s.warn("Ignoring roster IQ from bad address: " + stanza.From)
		return nil, true
	}
	var rst xmpp.Roster
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&rst); err != nil || len(rst.Item) == 0 {
		s.warn("Failed to parse roster push IQ")
		return nil, false
	}

	for _, entry := range rst.Item {
		if entry.Subscription == "remove" {
			s.R.Remove(entry.Jid)
		} else if s.R.AddOrMerge(roster.PeerFrom(entry, s.CurrentAccount.ID())) {
			s.iqReceived(entry.Jid)
		}
	}

	return xmpp.EmptyReply{}, false
}

func (s *Session) processIQ(stanza *xmpp.ClientIQ) (ret interface{}, ignore bool) {
	buf := bytes.NewBuffer(stanza.Query)
	parser := xml.NewDecoder(buf)
	token, _ := parser.Token()
	isGet := stanza.Type == "get"
	if token == nil {
		return nil, false
	}
	startElem, ok := token.(xml.StartElement)
	if !ok {
		return nil, false
	}

	switch startElem.Name.Space + " " + startElem.Name.Local {
	case "http://jabber.org/protocol/disco#info query":
		if isGet {
			return s.receivedIQDiscoInfo(), false
		}
	case "jabber:iq:version query":
		if isGet {
			return s.receivedIQVersion(), false
		}
	case "jabber:iq:roster query":
		if !isGet {
			return s.receivedIQRosterQuery(stanza)
		}
	}
	s.info("Unknown IQ: " + startElem.Name.Space + " " + startElem.Name.Local)

	return nil, false
}

// HandleConfirmOrDeny is used to handle a users response to a subscription request
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
	s.info(fmt.Sprintf("New OTR session with %s established", from))

	s.publishPeerEvent(OTRNewKeys, from)
}

func (s *Session) otrEnded(uid string) {
	s.publishPeerEvent(OTREnded, uid)
}

func (s *Session) newConversation(peer string) *otr3.Conversation {
	conversation := &otr3.Conversation{}
	conversation.SetOurKey(s.PrivateKey)

	hadInstanceTag := s.CurrentAccount.InstanceTag != 0
	s.CurrentAccount.InstanceTag = conversation.InitializeInstanceTag(s.CurrentAccount.InstanceTag)
	if !hadInstanceTag {
		s.SaveConfiguration()
	}

	s.CurrentAccount.SetOTRPoliciesFor(peer, conversation)

	return conversation
}

// GetConversationWith returns the OTR conversation for the given peer
func (s *Session) GetConversationWith(peer string) *otr3.Conversation {
	if conversation, ok := s.Conversations[peer]; ok {
		return conversation
	}

	conversation := s.newConversation(peer)
	s.Conversations[peer] = conversation

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
	log.Printf("-> Stanza %#v\n", stanza)

	from := xmpp.RemoveResourceFromJid(stanza.From)

	if stanza.Type == "error" {
		s.alert("Error reported from " + from + ": " + stanza.Body)
		return
	}

	//TODO: Add a more general solution to XEP's
	if len(stanza.Body) == 0 && len(stanza.Extensions) > 0 {
		//Extension only stanza
		return
	}

	var err error
	var messageTime time.Time
	if stanza.Delay != nil && len(stanza.Delay.Stamp) > 0 {
		// An XEP-0203 Delayed Delivery <delay/> element exists for
		// this message, meaning that someone sent it while we were
		// offline. Let's show the timestamp for when the message was
		// sent, rather than time.Now().
		messageTime, err = time.Parse(time.RFC3339, stanza.Delay.Stamp)
		if err != nil {
			s.alert("Can not parse Delayed Delivery timestamp, using quoted string instead.")
		}
	} else {
		messageTime = time.Now()
	}

	s.receiveClientMessage(from, messageTime, stanza.Body)
}

func (s *Session) receiveClientMessage(from string, when time.Time, body string) {
	conversation := s.GetConversationWith(from)
	out, toSend, err := conversation.Receive([]byte(body))
	encrypted := conversation.IsEncrypted()

	if err != nil {
		s.alert("While processing message from " + from + ": " + err.Error())
	}

	for _, msg := range toSend {
		s.Conn.Send(from, string(msg))
	}

	eh, _ := s.OtrEventHandler[from]
	change := eh.ConsumeSecurityChange()
	switch change {
	case event.NewKeys:
		s.newOTRKeys(from, conversation)
	case event.ConversationEnded:
		s.otrEnded(from)

		// TODO: all this stuff is very CLI specific, we should move it out and create good interaction
		// for the gui

		// TODO: twstrike/otr3 does not allow sending messages after the channel has
		// been terminated, so this should not be a problem.
		// This is probably unsafe without a policy that _forces_ crypto to
		// _everyone_ by default and refuses plaintext. Users might not notice
		// their buddy has ended a session, which they have also ended, and they
		// might send a plain text message. So we should ensure they _want_ this
		// feature and have set it as an explicit preference.
		if s.CurrentAccount.OTRAutoTearDown {
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
		s.CurrentAccount.AuthorizeFingerprint(from, fpr)
		s.SaveConfiguration()
	case event.SMPFailed:
		s.alert(fmt.Sprintf("Authentication with %s failed", from))
	}

	if len(out) == 0 {
		return
	}

	s.messageReceived(from, when, encrypted, out)
}

func (s *Session) messageReceived(from string, timestamp time.Time, encrypted bool, message []byte) {
	s.publishEvent(MessageEvent{
		Session:   s,
		From:      from,
		When:      timestamp,
		Body:      message,
		Encrypted: encrypted,
	})

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

// AwaitVersionReply listens on the channel and waits for the version reply
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

// WatchTimeout watches timeouts
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

const defaultDelimiter = "::"

// WatchRosterEvents waits for roster events
func (s *Session) WatchRosterEvents() {
	defer s.Close()

	s.Conn.SignalPresence("")
	s.info("Fetching roster")

	delim, err := s.Conn.GetRosterDelimiter()
	if err != nil || delim == "" {
		delim = defaultDelimiter
	}
	s.GroupDelimiter = delim

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
				s.R.AddOrMerge(roster.PeerFrom(rr, s.CurrentAccount.ID()))
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

// IsDisconnected returns true if this account is disconnected and is not in the process of connecting
func (s *Session) IsDisconnected() bool {
	return s.ConnStatus == DISCONNECTED
}

// Connect connects to the server and starts the main threads
func (s *Session) Connect(password string, registerCallback xmpp.FormCallback) error {
	if !s.IsDisconnected() {
		return nil
	}

	s.ConnStatus = CONNECTING
	if s.ConnectionLogger == nil {
		s.ConnectionLogger = newLogger()
	}

	conf := s.CurrentAccount
	policy := config.ConnectionPolicy{
		RequireTor:       conf.RequireTor,
		UseHiddenService: true,

		Logger:     s.ConnectionLogger,
		XMPPLogger: s.xmppLogger,
	}

	conn, err := policy.Connect(password, conf)
	if err != nil {
		s.alert(err.Error())
		s.ConnStatus = DISCONNECTED
		s.publish(Disconnected)

		return err
	}

	s.Conn = conn
	s.ConnStatus = CONNECTED
	s.publish(Connected)

	go s.WatchTimeout()
	go s.WatchRosterEvents()
	go s.WatchStanzas()

	return nil
}

// EncryptAndSendTo encrypts and sends the message to the given peer
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

// SendMessageTo sends the given message directly to the peer
func (s *Session) SendMessageTo(peer string, message string) error {
	return s.Conn.Send(peer, message)
}

// StartEncryptedChatWith starts an encrypted chat with the given peer
func (s *Session) StartEncryptedChatWith(peer string) error {
	conversation := s.GetConversationWith(peer)
	return s.SendMessageTo(peer, string(conversation.QueryMessage()))
}

// TerminateConversationWith terminates the conversation with a specific peer
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

	//TODO: add wipe for conversation
	//conversation.Wipe()

	return nil
}

func (s *Session) terminateConversations() {
	for peer := range s.Conversations {
		s.TerminateConversationWith(peer)
	}
}

// Close terminates all outstanding OTR conversations and closes the connection to the server
func (s *Session) Close() {
	if s.IsDisconnected() {
		return
	}

	s.terminateConversations()

	//Stops all
	s.Conn.Cancel(s.rosterCookie)
	s.timeoutTicker.Stop()

	s.Conn.Close()

	s.R.Clear()
	s.rosterReceived()

	s.ConnStatus = DISCONNECTED
	s.publish(Disconnected)
}

// Ping does a Ping
func (s *Session) Ping() {
	if s.ConnStatus == DISCONNECTED {
		return
	}
	/* fmt.Println("Publish Ping") */
	s.publish(Ping)
}
