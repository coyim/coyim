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

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/i18n"
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

// Connector represents something that connect
type Connector interface {
	Connect()
}

// Session contains information about one specific connection
type Session struct {
	conn             *xmpp.Conn
	connectionLogger io.Writer
	r                *roster.List
	connStatus       connStatus

	otrEventHandler map[string]*event.OtrEventHandler

	privateKeys []otr3.PrivateKey

	//TODO: the session does not need all application config. Copy only what it needs to configure the session
	config *config.ApplicationConfig

	accountConfig *config.Account

	// timeouts maps from Cookies (from outstanding requests) to the
	// absolute time when that request should timeout.
	timeouts map[xmpp.Cookie]time.Time

	// LastActionTime is the time at which the user last entered a command,
	// or was last notified.
	lastActionTime      time.Time
	sessionEventHandler EventHandler

	// WantToBeOnline keeps track of whether a user has expressed a wish
	// to be online - if it's true, it will do more aggressive reconnecting
	wantToBeOnline bool

	subscribers struct {
		sync.RWMutex
		subs []chan<- interface{}
	}

	groupDelimiter string

	xmppLogger io.Writer

	connector Connector

	cmdManager  client.CommandManager
	convManager client.ConversationManager
}

// GetConfig returns the current account configuration
func (s *Session) GetConfig() *config.Account {
	return s.accountConfig
}

func parseFromConfig(cu *config.Account) []otr3.PrivateKey {
	var result []otr3.PrivateKey

	allKeys := cu.AllPrivateKeys()

	log.Printf("Loading %d configured keys", len(allKeys))
	for _, pp := range allKeys {
		_, ok, parsedKey := otr3.ParsePrivateKey(pp)
		if ok {
			result = append(result, parsedKey)
			log.Printf("Loaded key: %s", config.FormatFingerprint(parsedKey.PublicKey().Fingerprint()))
		}
	}

	return result
}

// NewSession creates a new session from the given config
func NewSession(c *config.ApplicationConfig, cu *config.Account) *Session {
	s := &Session{
		config:        c,
		accountConfig: cu,

		r:               roster.New(),
		otrEventHandler: make(map[string]*event.OtrEventHandler),
		lastActionTime:  time.Now(),

		timeouts: make(map[xmpp.Cookie]time.Time),

		xmppLogger: openLogFile(c.RawLogFile),
	}

	s.ReloadKeys()
	s.convManager = client.NewConversationManager(s, s)

	go observe(s)
	go checkReconnect(s)

	return s
}

// ReloadKeys will reload the keys from the configuration
func (s *Session) ReloadKeys() {
	s.privateKeys = parseFromConfig(s.accountConfig)
}

// Send will send the given message to the receiver given
func (s *Session) Send(to string, msg string) error {
	log.Printf("<- to=%v {%v}\n", to, msg)
	return s.conn.Send(to, msg)
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

func (s *Session) receivedStreamError(stanza *xmpp.StreamError) bool {
	s.alert("Exiting in response to fatal error from server: " + stanza.String())
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
		s.r.SubscribeRequest(stanza.From, either(stanza.ID, "0000"), s.GetConfig().ID())
		s.publishPeerEvent(
			SubscriptionRequest,
			xmpp.RemoveResourceFromJid(stanza.From),
		)
	case "unavailable":
		if !s.r.PeerBecameUnavailable(stanza.From) {
			return true
		}

		s.publishEvent(PresenceEvent{
			Session:        s,
			ClientPresence: stanza,
			Gone:           true,
		})
	case "":
		if !s.r.PeerPresenceUpdate(stanza.From, stanza.Show, stanza.Status, s.GetConfig().ID()) {
			return true
		}

		s.publishEvent(PresenceEvent{
			Session:        s,
			ClientPresence: stanza,
			Gone:           false,
		})
	case "subscribed":
		s.r.Subscribed(stanza.From)
		s.publishPeerEvent(
			Subscribed,
			xmpp.RemoveResourceFromJid(stanza.From),
		)
	case "unsubscribe":
		s.r.Unsubscribed(stanza.From)
		s.publishPeerEvent(
			Unsubscribe,
			xmpp.RemoveResourceFromJid(stanza.From),
		)
	case "unsubscribed":
		// Ignore
	case "error":
		s.warn(fmt.Sprintf("Got a presence error from %s: %#v\n", stanza.From, stanza.Error))
		s.r.LatestError(stanza.From, stanza.Error.Code, stanza.Error.Type, stanza.Error.Any.Space+" "+stanza.Error.Any.Local)
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

		if err := s.conn.SendIQReply(stanza.From, "result", stanza.ID, reply); err != nil {
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
			s.info(fmt.Sprintf("RECEIVED %s %s", rawStanza.Name, rawStanza.Value))
			return true
		}
	}
}

//TODO: differentiate errors from disconnect request
func (s *Session) watchStanzas() {
	defer s.connectionLost()

	stanzaChan := make(chan xmpp.Stanza)
	go s.readStanzasAndAlertOnErrors(stanzaChan)
	for s.receiveStanza(stanzaChan) {
	}
}

func (s *Session) readStanzasAndAlertOnErrors(stanzaChan chan xmpp.Stanza) {
	if err := s.conn.ReadStanzas(stanzaChan); err != nil {
		s.alert(fmt.Sprintf("error reading XMPP message: %s", err.Error()))
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
				Name:     s.GetConfig().Account,
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

func peerFrom(entry xmpp.RosterEntry, c *config.Account) *roster.Peer {
	belongsTo := c.ID()
	var nickname string
	p, ok := c.GetPeer(entry.Jid)
	if ok {
		nickname = p.Nickname
	}
	return roster.PeerFrom(entry, belongsTo, nickname)
}

func (s *Session) addOrMergeNewPeer(entry xmpp.RosterEntry, c *config.Account) bool {

	return s.r.AddOrMerge(peerFrom(entry, c))
}

func (s *Session) receivedIQRosterQuery(stanza *xmpp.ClientIQ) (ret interface{}, ignore bool) {
	// TODO: we should deal with "ask" attributes here

	if len(stanza.From) > 0 && !s.GetConfig().Is(stanza.From) {
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
			s.r.Remove(entry.Jid)
		} else if s.addOrMergeNewPeer(entry, s.GetConfig()) {
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
	id, ok := s.r.RemovePendingSubscribe(jid)
	if !ok {
		s.warn("No pending subscription from " + jid)
		return
	}

	var err error
	switch isConfirm {
	case true:
		err = s.ApprovePresenceSubscription(jid, id)
	default:
		err = s.DenyPresenceSubscription(jid, id)
	}

	if err != nil {
		s.warn("Error sending presence stanza: " + err.Error())
		return
	}

	if isConfirm {
		s.RequestPresenceSubscription(jid)
	}
}

func (s *Session) newOTRKeys(from string, conversation client.Conversation) {
	s.publishPeerEvent(OTRNewKeys, from)
}

func (s *Session) renewedOTRKeys(from string, conversation client.Conversation) {
	s.publishPeerEvent(OTRRenewedKeys, from)
}

func (s *Session) otrEnded(uid string) {
	s.publishPeerEvent(OTREnded, uid)
}

func (s *Session) listenToNotifications(c <-chan string, peer string) {
	for notification := range c {
		s.publishEvent(NotificationEvent{
			Session:      s,
			Peer:         peer,
			Notification: notification,
		})
	}
}

// NewConversation will create a new OTR conversation with the given peer
//TODO: why creating a conversation is coupled to the account config and the session
//TODO: does the creation of the OTR event handler need to be guarded with a lock?
func (s *Session) NewConversation(peer string) *otr3.Conversation {
	conversation := &otr3.Conversation{}
	conversation.SetOurKeys(s.privateKeys)
	conversation.SetFriendlyQueryMessage(i18n.Local("Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol."))

	instanceTag := conversation.InitializeInstanceTag(s.GetConfig().InstanceTag)

	if s.GetConfig().InstanceTag != instanceTag {
		s.cmdManager.ExecuteCmd(client.SaveInstanceTagCmd{
			Account:     s.GetConfig(),
			InstanceTag: instanceTag,
		})
	}

	s.GetConfig().SetOTRPoliciesFor(peer, conversation)

	eh, ok := s.otrEventHandler[peer]
	if !ok {
		eh = new(event.OtrEventHandler)
		eh.Account = s.GetConfig().Account
		eh.Peer = peer
		notificationsChan := make(chan string)
		eh.Notifications = notificationsChan
		go s.listenToNotifications(notificationsChan, peer)
		conversation.SetSMPEventHandler(eh)
		conversation.SetErrorMessageHandler(eh)
		conversation.SetMessageEventHandler(eh)
		conversation.SetSecurityEventHandler(eh)
		s.otrEventHandler[peer] = eh
	}

	return conversation
}

func (s *Session) processClientMessage(stanza *xmpp.ClientMessage) {
	log.Printf("-> Stanza %#v\n", stanza)

	from := xmpp.RemoveResourceFromJid(stanza.From)

	//TODO: investigate which errors are recoverable
	//https://xmpp.org/rfcs/rfc3920.html#stanzas-error
	if stanza.Type == "error" && stanza.Error != nil {
		s.alert(fmt.Sprintf("Error reported from %s: %#v", from, stanza.Error))
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
	conversation, _ := s.convManager.EnsureConversationWith(from)
	out, err := conversation.Receive(s, []byte(body))
	encrypted := conversation.IsEncrypted()

	if err != nil {
		s.alert("While processing message from " + from + ": " + err.Error())
	}

	eh, _ := s.otrEventHandler[from]
	change := eh.ConsumeSecurityChange()
	switch change {
	case event.NewKeys:
		s.newOTRKeys(from, conversation)
	case event.RenewedKeys:
		s.renewedOTRKeys(from, conversation)
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
		if s.GetConfig().OTRAutoTearDown {
			c, existing := s.convManager.GetConversationWith(from)
			if !existing {
				s.alert(fmt.Sprintf("No secure session established; unable to automatically tear down OTR conversation with %s.", from))
				break
			} else {
				s.info(fmt.Sprintf("%s has ended the secure conversation.", from))

				err := c.EndEncryptedChat(s)
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
		fpr := conversation.TheirFingerprint()
		s.cmdManager.ExecuteCmd(client.AuthorizeFingerprintCmd{
			Account:     s.GetConfig(),
			Peer:        from,
			Fingerprint: fpr,
		})
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

func (s *Session) watchTimeout() {
	tickInterval := time.Second

	for s.IsConnected() {
		now := <-time.After(tickInterval)
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
				log.Println("session: cookie", cookie, "has expired")
				s.conn.Cancel(cookie)
			} else {
				newTimeouts[cookie] = expiry
			}
		}

		s.timeouts = newTimeouts
	}
}

// Timeout set the timeout for an XMPP request
func (s *Session) Timeout(c xmpp.Cookie, t time.Time) {
	s.timeouts[c] = t
}

const defaultDelimiter = "::"

func (s *Session) watchRoster() {
	for s.requestRoster() {
		time.Sleep(time.Duration(3) * time.Minute)
	}
}

func (s *Session) requestRoster() bool {
	if !s.IsConnected() {
		return false
	}

	s.info("Fetching roster")

	delim, err := s.conn.GetRosterDelimiter()
	if err != nil || delim == "" {
		delim = defaultDelimiter
	}
	s.groupDelimiter = delim

	rosterReply, _, err := s.conn.RequestRoster()
	if err != nil {
		s.alert("Failed to request roster: " + err.Error())
		return true
	}

	rosterStanza, ok := <-rosterReply
	if !ok {
		//TODO: should we retry the request in such case?
		log.Println("session: roster request cancelled or timedout")
		return true
	}

	rst, err := xmpp.ParseRoster(rosterStanza)
	if err != nil {
		s.alert("Failed to parse roster: " + err.Error())
		return true
	}

	for _, rr := range rst {
		s.addOrMergeNewPeer(rr, s.GetConfig())
	}

	s.rosterReceived()
	s.info("Roster received")

	return true
}

// IsDisconnected returns true if this account is disconnected and is not in the process of connecting
func (s *Session) IsDisconnected() bool {
	return s.connStatus == DISCONNECTED
}

// IsConnected returns true if this account is connected and is not in the process of connecting
func (s *Session) IsConnected() bool {
	return s.connStatus == CONNECTED
}

func (s *Session) setStatus(status connStatus) {
	s.connStatus = status

	switch status {
	case CONNECTED:
		s.publish(Connected)
	case DISCONNECTED:
		s.publish(Disconnected)
	case CONNECTING:
		s.publish(Connecting)
	}
}

// Connect connects to the server and starts the main threads
func (s *Session) Connect(password string) error {
	if !s.IsDisconnected() {
		return nil
	}

	s.setStatus(CONNECTING)

	if s.connectionLogger == nil {
		s.connectionLogger = newLogger()
	}

	conf := s.GetConfig()
	policy := config.ConnectionPolicy{
		Logger:     s.connectionLogger,
		XMPPLogger: s.xmppLogger,
	}

	conn, err := policy.Connect(password, conf)
	if err != nil {
		s.alert(err.Error())
		s.setStatus(DISCONNECTED)

		return err
	}

	s.conn = conn
	s.setStatus(CONNECTED)

	s.conn.SignalPresence("")
	go s.watchRoster()
	go s.watchTimeout()
	go s.watchStanzas()

	return nil
}

// EncryptAndSendTo encrypts and sends the message to the given peer
func (s *Session) EncryptAndSendTo(peer string, message string) error {
	//TODO: review whether it should create a conversation
	conversation, _ := s.convManager.EnsureConversationWith(peer)
	return conversation.Send(s, []byte(message))
}

// SendMessageTo sends the given message directly to the peer
func (s *Session) SendMessageTo(peer string, message string) error {
	return s.conn.Send(peer, message)
}

func (s *Session) terminateConversations() {
	s.convManager.TerminateAll()
}

func (s *Session) connectionLost() {
	if s.IsDisconnected() {
		return
	}

	s.Close()
	s.publish(ConnectionLost)
}

// Close terminates all outstanding OTR conversations and closes the connection to the server
func (s *Session) Close() {
	if s.IsDisconnected() {
		return
	}

	s.setStatus(DISCONNECTED)

	s.terminateConversations()
	s.conn.Close()
}

// Ping does a Ping
func (s *Session) Ping() {
	if s.IsDisconnected() {
		return
	}
	/* fmt.Println("Publish Ping") */
	s.publish(Ping)
}

func (s *Session) CommandManager() client.CommandManager {
	return s.cmdManager
}

func (s *Session) SetCommandManager(c client.CommandManager) {
	s.cmdManager = c
}

func (s *Session) ConversationManager() client.ConversationManager {
	return s.convManager
}

func (s *Session) SetWantToBeOnline(val bool) {
	s.wantToBeOnline = val
}

func (s *Session) PrivateKeys() []otr3.PrivateKey {
	return s.privateKeys
}

func (s *Session) R() *roster.List {
	return s.r
}

func (s *Session) SetConnector(c Connector) {
	s.connector = c
}

func (s *Session) GroupDelimiter() string {
	return s.groupDelimiter
}

func (s *Session) Config() *config.ApplicationConfig {
	return s.config
}

func (s *Session) Conn() *xmpp.Conn {
	return s.conn
}

func (s *Session) SetSessionEventHandler(eh EventHandler) {
	s.sessionEventHandler = eh
}

func (s *Session) SetConnectionLogger(l io.Writer) {
	s.connectionLogger = l
}

func (s *Session) OtrEventHandler() map[string]*event.OtrEventHandler {
	return s.otrEventHandler
}

func (s *Session) SetLastActionTime(t time.Time) {
	s.lastActionTime = t
}
