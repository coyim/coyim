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

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/otr_client"
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/ui"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
)

type connStatus int

// These constants represent the different connection states
const (
	DISCONNECTED connStatus = iota
	CONNECTING
	CONNECTED
)

type session struct {
	conn             xi.Conn
	connectionLogger io.Writer
	r                *roster.List

	connStatus     connStatus
	connStatusLock sync.RWMutex

	otrEventHandlers *otr_client.EventHandlers

	privateKeys []otr3.PrivateKey

	//TODO: the session does not need all application config. Copy only what it needs to configure the session
	config *config.ApplicationConfig

	accountConfig *config.Account

	// timeouts maps from Cookies (from outstanding requests) to the
	// absolute time when that request should timeout.
	timeouts map[data.Cookie]time.Time

	// LastActionTime is the time at which the user last entered a command,
	// or was last notified.
	lastActionTime      time.Time
	sessionEventHandler access.EventHandler

	// WantToBeOnline keeps track of whether a user has expressed a wish
	// to be online - if it's true, it will do more aggressive reconnecting
	wantToBeOnline bool

	subscribers struct {
		sync.RWMutex
		subs []chan<- interface{}
	}

	groupDelimiter string

	inMemoryLog *bytes.Buffer
	xmppLogger  io.Writer

	connector access.Connector

	cmdManager  otr_client.CommandManager
	convManager otr_client.ConversationManager

	dialerFactory func(tls.Verifier) xi.Dialer

	autoApproves map[string]bool

	nicknames []string
}

// GetInMemoryLog returns the in memory log or nil
func (s *session) GetInMemoryLog() *bytes.Buffer {
	return s.inMemoryLog
}

// GetConfig returns the current account configuration
func (s *session) GetConfig() *config.Account {
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

// CreateXMPPLogger creates a XMPP log.
func CreateXMPPLogger(rawLog string) (*bytes.Buffer, io.Writer) {
	log := openLogFile(rawLog)

	var inMemory *bytes.Buffer
	if *config.DebugFlag {
		inMemory = new(bytes.Buffer)

		if log != nil {
			log = io.MultiWriter(log, inMemory)
		} else {
			log = inMemory
		}
	}

	return inMemory, log
}

// Factory creates a new session from the given config
func Factory(c *config.ApplicationConfig, cu *config.Account, df func(tls.Verifier) xi.Dialer) access.Session {
	// Make xmppLogger go to in memory STRING and/or the log file

	inMemoryLog, xmppLogger := CreateXMPPLogger(c.RawLogFile)

	s := &session{
		config:        c,
		accountConfig: cu,

		r:              roster.New(),
		lastActionTime: time.Now(),

		timeouts: make(map[data.Cookie]time.Time),

		autoApproves: make(map[string]bool),

		inMemoryLog:      inMemoryLog,
		xmppLogger:       xmppLogger,
		connectionLogger: LogToDebugLog(),
		dialerFactory:    df,
	}

	s.otrEventHandlers = otr_client.NewEventHandlers(cu.Account, s.onOtrEventHandlerCreate)

	s.ReloadKeys()
	s.convManager = otr_client.NewConversationManager(s, s)

	go observe(s)
	go checkReconnect(s)

	return s
}

func (s *session) onOtrEventHandlerCreate(peer jid.Any, eh *otr_client.EventHandler, nots chan string, dels chan int) {
	go s.listenToNotifications(nots, peer)
	go s.listenToDelayedMessageDelivery(dels, peer)
}

// ReloadKeys will reload the keys from the configuration
func (s *session) ReloadKeys() {
	s.privateKeys = parseFromConfig(s.accountConfig)
}

// Send will send the given message to the receiver given
func (s *session) Send(to jid.WithoutResource, resource jid.Resource, msg string) error {
	conn, ok := s.connection()
	if ok {
		log.Printf("<- to=%v {%v}\n", to.MaybeWithResource(resource), msg)
		return conn.Send(to.MaybeWithResource(resource).String(), msg)
	}
	return &access.OfflineError{Msg: i18n.Local("Couldn't send message since we are not connected")}
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

func (s *session) info(m string) {
	s.publishEvent(events.Log{
		Level:   events.Info,
		Message: m,
	})
}

func (s *session) Warn(m string) {
	s.warn(m)
}

func (s *session) Info(m string) {
	s.info(m)
}

func (s *session) warn(m string) {
	s.publishEvent(events.Log{
		Level:   events.Warn,
		Message: m,
	})
}

func (s *session) alert(m string) {
	s.publishEvent(events.Log{
		Level:   events.Alert,
		Message: m,
	})
}

func (s *session) receivedStreamError(stanza *data.StreamError) bool {
	s.alert("Exiting in response to fatal error from server: " + stanza.String())
	return false
}

func retrieveMessageTime(stanza *data.ClientMessage) time.Time {
	if stanza.Delay != nil && len(stanza.Delay.Stamp) > 0 {
		// An XEP-0203 Delayed Delivery <delay/> element exists for
		// this message, meaning that someone sent it while we were
		// offline. Let's show the timestamp for when the message was
		// sent, rather than time.Now().
		messageTime, err := time.Parse(time.RFC3339, stanza.Delay.Stamp)
		if err != nil {
			//TODO: use quoted string instead of timstamp.
			//s.alert("Can not parse Delayed Delivery timestamp, using quoted string instead.")
		} else {
			return messageTime
		}
	}

	return time.Now()
}

func (s *session) receivedClientMessage(stanza *data.ClientMessage) bool {
	log.Printf("-> Stanza %#v\n", stanza)

	if len(stanza.Body) == 0 && len(stanza.Extensions) > 0 {
		s.processExtensions(stanza)
		return true
	}

	peer := jid.Parse(stanza.From)

	// TODO: it feels iffy that we have error and groupchat special handled here
	// But not checking on the "message" type.
	switch stanza.Type {
	case "error":
		//TODO: investigate which errors are NOT recoverable, and return false
		//to close the connection
		//https://xmpp.org/rfcs/rfc3920.html#stanzas-error
		if stanza.Error != nil {
			s.alert(fmt.Sprintf("Error reported from %s: %#v", peer.NoResource(), stanza.Error))
			return true
		}
	case "groupchat":
		s.publishEvent(events.ChatMessage{
			// TODO: unclear if the following line is correct
			From:          jid.R(stanza.From),
			When:          retrieveMessageTime(stanza),
			Body:          stanza.Body,
			ClientMessage: stanza,
		})
		return true
	}

	messageTime := retrieveMessageTime(stanza)
	s.receiveClientMessage(peer, messageTime, stanza.Body)

	return true
}

func either(l, r string) string {
	if l == "" {
		return r
	}
	return l
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func (s *session) receivedClientPresence(stanza *data.ClientPresence) bool {
	//MUC is interested in every presence, so we publish regardless.
	//It is sad that not every presence stanza triggers a presence event.
	s.publishEvent(events.ChatPresence{
		Session:        s,
		ClientPresence: stanza,
	})

	switch stanza.Type {
	case "subscribe":
		jj := jid.NR(stanza.From)
		jjr := jj.String()
		if s.autoApproves[jjr] {
			delete(s.autoApproves, jjr)
			s.ApprovePresenceSubscription(jj, stanza.ID)
		} else {
			s.r.SubscribeRequest(jj, either(stanza.ID, "0000"), s.GetConfig().ID())
			s.publishPeerEvent(
				events.SubscriptionRequest,
				jj,
			)
		}
	case "unavailable":
		if !s.r.PeerBecameUnavailable(jid.Parse(stanza.From)) {
			return true
		}

		s.publishEvent(events.Presence{
			Session:        s,
			ClientPresence: stanza,
			Gone:           true,
		})
	case "":
		if !s.r.PeerPresenceUpdate(jid.R(stanza.From), stanza.Show, stanza.Status, s.GetConfig().ID()) {
			return true
		}

		//TODO: If to == "" this is our own presence confirmation.
		//"From" is how we are identified (will be JID/"some-id")
		//Same thing happens for the group-chat, but in this case it tell us also what are our affiliations and roles.
		//Thats why I'm worried about handling this as a regular peer presence - which is not.

		s.publishEvent(events.Presence{
			Session:        s,
			ClientPresence: stanza,
			Gone:           false,
		})
	case "subscribed":
		jj := jid.NR(stanza.From)
		s.r.Subscribed(jj)
		s.publishPeerEvent(
			events.Subscribed,
			jj,
		)
	case "unsubscribe":
		jj := jid.NR(stanza.From)
		s.r.Unsubscribed(jj)
		s.publishPeerEvent(
			events.Unsubscribe,
			jj,
		)
	case "unsubscribed":
		// Ignore
	case "error":
		s.warn(fmt.Sprintf("Got a presence error from %s: %#v\n", stanza.From, stanza.Error))
		s.r.LatestError(jid.NR(stanza.From), stanza.Error.Code, stanza.Error.Type, stanza.Error.Condition.XMLName.Space+" "+stanza.Error.Condition.XMLName.Local)
	default:
		s.info(fmt.Sprintf("unrecognized presence: %#v", stanza))
	}

	return true
}

func (s *session) receivedClientIQ(stanza *data.ClientIQ) bool {
	if stanza.Type == "get" || stanza.Type == "set" {
		reply, iqtype, ignore := s.processIQ(stanza)
		if ignore {
			return true
		}

		if iqtype == "" {
			iqtype = "result"
		}

		if reply == nil {
			reply = data.ErrorReply{
				Type:  "cancel",
				Error: data.ErrorBadRequest{},
			}
		}

		s.sendIQReply(stanza, iqtype, reply)
		return true
	}
	s.info(fmt.Sprintf("unrecognized iq: %#v", stanza))
	return true
}

func (s *session) receiveStanza(stanzaChan chan data.Stanza) bool {
	select {
	case rawStanza, ok := <-stanzaChan:
		if !ok {
			return false
		}

		switch stanza := rawStanza.Value.(type) {
		case *data.StreamError:
			return s.receivedStreamError(stanza)
		case *data.ClientMessage:
			return s.receivedClientMessage(stanza)
		case *data.ClientPresence:
			return s.receivedClientPresence(stanza)
		case *data.ClientIQ:
			return s.receivedClientIQ(stanza)
		default:
			s.info(fmt.Sprintf("unhandled stanza: %s %s", rawStanza.Name, rawStanza.Value))
			return true
		}
	}
}

//TODO: differentiate errors from disconnect request
func (s *session) watchStanzas() {
	defer s.connectionLost()

	stanzaChan := make(chan data.Stanza)
	go s.readStanzasAndAlertOnErrors(stanzaChan)
	for s.receiveStanza(stanzaChan) {
	}
}

func (s *session) readStanzasAndAlertOnErrors(stanzaChan chan data.Stanza) {
	if err := s.conn.ReadStanzas(stanzaChan); err != nil {
		s.alert(fmt.Sprintf("error reading XMPP message: %s", err.Error()))
	}
}

func (s *session) rosterReceived() {
	s.info("Roster received")
	s.publish(events.RosterReceived)
}

func (s *session) iqReceived(peer jid.Any) {
	s.publishPeerEvent(events.IQReceived, peer)
}

func (s *session) receivedIQVersion() data.VersionReply {
	return data.VersionReply{
		Name:    "testing",
		Version: "version",
		OS:      "none",
	}
}

func peerFrom(entry data.RosterEntry, c *config.Account) *roster.Peer {
	belongsTo := c.ID()

	var nickname string
	var groups []string
	if p, ok := c.GetPeer(entry.Jid); ok {
		nickname = p.Nickname
		groups = p.Groups
	}

	return roster.PeerFrom(entry, belongsTo, nickname, groups)
}

func (s *session) addOrMergeNewPeer(entry data.RosterEntry, c *config.Account) bool {
	return s.r.AddOrMerge(peerFrom(entry, c))
}

func (s *session) receivedIQRosterQuery(stanza *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	// TODO: we should deal with "ask" attributes here

	if len(stanza.From) > 0 && !s.GetConfig().Is(stanza.From) {
		s.warn("Ignoring roster IQ from bad address: " + stanza.From)
		return nil, "", true
	}
	var rst data.Roster
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&rst); err != nil || len(rst.Item) == 0 {
		s.warn("Failed to parse roster push IQ")
		return nil, "", false
	}

	for _, entry := range rst.Item {
		jj := jid.Parse(entry.Jid)
		if entry.Subscription == "remove" {
			s.r.Remove(jj.NoResource())
		} else if s.addOrMergeNewPeer(entry, s.GetConfig()) {
			s.iqReceived(jj)
		}
	}

	return data.EmptyReply{}, "", false
}

// HandleConfirmOrDeny is used to handle a users response to a subscription request
func (s *session) HandleConfirmOrDeny(jid jid.WithoutResource, isConfirm bool) {
	id, ok := s.r.RemovePendingSubscribe(jid)
	if !ok {
		s.warn("No pending subscription from " + jid.String())
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
		s.RequestPresenceSubscription(jid, "")
	}
}

func (s *session) newOTRKeys(peer jid.WithoutResource, conversation otr_client.Conversation) {
	s.publishPeerEvent(events.OTRNewKeys, peer)
}

func (s *session) renewedOTRKeys(peer jid.WithoutResource, conversation otr_client.Conversation) {
	s.publishPeerEvent(events.OTRRenewedKeys, peer)
}

func (s *session) otrEnded(peer jid.WithoutResource) {
	s.publishPeerEvent(events.OTREnded, peer)
}

func (s *session) listenToNotifications(c <-chan string, peer jid.Any) {
	for notification := range c {
		s.publishEvent(events.Notification{
			Session:      s,
			Peer:         peer,
			Notification: notification,
		})
	}
}

func (s *session) listenToDelayedMessageDelivery(c <-chan int, peer jid.Any) {
	for t := range c {
		s.publishEvent(events.DelayedMessageSent{
			Session: s,
			Peer:    peer,
			Tracer:  t,
		})
	}
}

// NewConversation will create a new OTR conversation with the given peer
//TODO: why creating a conversation is coupled to the account config and the session
//TODO: does the creation of the OTR event handler need to be guarded with a lock?
//TODO: Why starting a conversation requires being able to translate a message?
//TODO: we are not taking a proper JID type here. This is to stop leakage into other packages for now, but should be fixed.
//This also assumes it's useful to send friendly message to another person in
//the same language configured on your app.
func (s *session) NewConversation(peer jid.Any) *otr3.Conversation {
	conversation := &otr3.Conversation{}
	conversation.SetOurKeys(s.privateKeys)
	conversation.SetFriendlyQueryMessage(i18n.Local("Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol."))

	instanceTag := conversation.InitializeInstanceTag(s.GetConfig().InstanceTag)

	if s.GetConfig().InstanceTag != instanceTag {
		s.cmdManager.ExecuteCmd(otr_client.SaveInstanceTagCmd{
			Account:     s.GetConfig(),
			InstanceTag: instanceTag,
		})
	}

	s.GetConfig().SetOTRPoliciesFor(peer.String(), conversation)
	s.otrEventHandlers.EnsureExists(peer, conversation)

	return conversation
}

// ManuallyEndEncryptedChat allows a user to end the encrypted chat from this side
func (s *session) ManuallyEndEncryptedChat(peer jid.WithoutResource, resource jid.Resource) error {
	c, ok := s.ConversationManager().GetConversationWith(peer, resource)
	if !ok {
		return fmt.Errorf("couldn't find conversation with %s / %s", peer, string(resource))
	}

	defer s.otrEventHandlers.Get(peer).ConsumeSecurityChange()
	return c.EndEncryptedChat(s, resource)
}

func (s *session) receiveClientMessage(peer jid.Any, when time.Time, body string) {
	from := peer.NoResource()
	fr := from.String()
	resource := jid.Resource("")
	if px, ok := peer.(jid.WithResource); ok {
		resource = px.Resource()
	}

	// TODO: do we want to have different conversation instances for different resources?
	conversation, _ := s.convManager.EnsureConversationWith(from, resource)
	out, err := conversation.Receive(s, resource, []byte(body))
	encrypted := conversation.IsEncrypted()

	if err != nil {
		s.alert("While processing message from " + fr + ": " + err.Error())
	}

	eh := s.otrEventHandlers.Get(from)
	change := eh.ConsumeSecurityChange()
	switch change {
	case otr_client.NewKeys:
		s.newOTRKeys(from, conversation)
	case otr_client.RenewedKeys:
		s.renewedOTRKeys(from, conversation)
	case otr_client.ConversationEnded:
		s.otrEnded(from)

		// TODO: all this stuff is very CLI specific, we should move it out and create good interaction
		// for the gui

		// TODO: coyim/otr3 does not allow sending messages after the channel has
		// been terminated, so this should not be a problem.
		// This is probably unsafe without a policy that _forces_ crypto to
		// _everyone_ by default and refuses plaintext. Users might not notice
		// their buddy has ended a session, which they have also ended, and they
		// might send a plain text message. So we should ensure they _want_ this
		// feature and have set it as an explicit preference.
		if s.GetConfig().OTRAutoTearDown {
			c, existing := s.convManager.GetConversationWith(from, resource)
			if !existing {
				s.alert(fmt.Sprintf("No secure session established; unable to automatically tear down OTR conversation with %s.", fr))
				break
			} else {
				s.info(fmt.Sprintf("%s has ended the secure conversation.", fr))

				err := c.EndEncryptedChat(s, resource)
				if err != nil {
					s.info(fmt.Sprintf("Unable to automatically tear down OTR conversation with %s: %s\n", fr, err.Error()))
					break
				}

				s.info(fmt.Sprintf("Secure session with %s has been automatically ended. Messages will be sent in the clear until another OTR session is established.", fr))
			}
		} else {
			s.info(fmt.Sprintf("%s has ended the secure conversation. You should do likewise with /otr-end %s", fr, fr))
		}
	case otr_client.SMPSecretNeeded:
		s.publishSMPEvent(events.SecretNeeded, peer, eh.SmpQuestion)
	case otr_client.SMPComplete:
		s.publishSMPEvent(events.Success, peer, "")
		s.cmdManager.ExecuteCmd(otr_client.AuthorizeFingerprintCmd{
			Account:     s.GetConfig(),
			Session:     s,
			Peer:        from,
			Fingerprint: conversation.TheirFingerprint(),
		})
	case otr_client.SMPFailed:
		s.publishSMPEvent(events.Failure, peer, "")
	}

	if len(out) == 0 {
		return
	}

	if encrypted {
		out = ui.UnescapeNewlineTags(out)
	}

	s.messageReceived(peer, when, encrypted, out)
}

func (s *session) messageReceived(peer jid.Any, timestamp time.Time, encrypted bool, message []byte) {
	s.publishEvent(events.Message{
		Session:   s,
		From:      peer,
		When:      timestamp,
		Body:      message,
		Encrypted: encrypted,
	})

	s.maybeNotify()
}

func (s *session) maybeNotify() {
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
func (s *session) AwaitVersionReply(ch <-chan data.Stanza, user string) {
	stanza, ok := <-ch
	if !ok {
		s.warn("Version request to " + user + " timed out")
		return
	}
	reply, ok := stanza.Value.(*data.ClientIQ)
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
	var versionReply data.VersionReply
	if err := xml.NewDecoder(buf).Decode(&versionReply); err != nil {
		s.warn("Failed to parse version reply from " + user + ": " + err.Error())
		return
	}

	s.info(fmt.Sprintf("Version reply from %s: %#v", user, versionReply))
}

func (s *session) watchTimeout() {
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

		newTimeouts := make(map[data.Cookie]time.Time)
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
func (s *session) Timeout(c data.Cookie, t time.Time) {
	s.timeouts[c] = t
}

const defaultDelimiter = "::"

func (s *session) watchRoster() {
	for s.requestRoster() {
		time.Sleep(time.Duration(3) * time.Minute)
	}
}

func (s *session) getVCard() {
	conn, ok := s.connection()
	if !ok {
		return
	}

	s.info("Fetching VCard")

	vcardReply, _, err := conn.RequestVCard()
	if err != nil {
		s.alert("Failed to request vcard: " + err.Error())
		return
	}

	vcardStanza, ok := <-vcardReply
	if !ok {
		log.Println("session: vcard request cancelled or timedout")
		return
	}

	vc, err := data.ParseVCard(vcardStanza)
	if err != nil {
		s.alert("Failed to parse vcard: " + err.Error())
		return
	}

	s.nicknames = []string{vc.Nickname, vc.FullName}

	return
}

func (s *session) DisplayName() string {
	return either(either(s.accountConfig.Nickname, firstNonEmpty(s.nicknames...)), s.accountConfig.Account)
}

func (s *session) requestRoster() bool {
	conn, ok := s.connection()
	if !ok {
		return false
	}

	s.info("Fetching roster")

	delim, err := conn.GetRosterDelimiter()
	if err != nil || delim == "" {
		delim = defaultDelimiter
	}
	s.groupDelimiter = delim

	rosterReply, _, err := conn.RequestRoster()
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

	rst, err := data.ParseRoster(rosterStanza)
	if err != nil {
		s.alert("Failed to parse roster: " + err.Error())
		return true
	}

	for _, rr := range rst {
		s.addOrMergeNewPeer(rr, s.GetConfig())
	}

	s.rosterReceived()

	return true
}

func (s *session) getConnStatus() connStatus {
	s.connStatusLock.RLock()
	defer s.connStatusLock.RUnlock()
	return s.connStatus
}

func (s *session) setConnStatus(v connStatus) {
	s.connStatusLock.Lock()
	defer s.connStatusLock.Unlock()
	s.connStatus = v
}

// IsDisconnected returns true if this account is disconnected and is not in the process of connecting
func (s *session) IsDisconnected() bool {
	return s.getConnStatus() == DISCONNECTED
}

// IsConnected returns true if this account is connected and is not in the process of connecting
func (s *session) IsConnected() bool {
	return s.getConnStatus() == CONNECTED
}

func (s *session) connection() (xi.Conn, bool) {
	return s.conn, s.getConnStatus() == CONNECTED
}

func (s *session) setStatus(status connStatus) {
	s.setConnStatus(status)

	switch status {
	case CONNECTED:
		s.publish(events.Connected)
	case DISCONNECTED:
		s.publish(events.Disconnected)
	case CONNECTING:
		s.publish(events.Connecting)
	}
}

// Connect connects to the server and starts the main threads
func (s *session) Connect(password string, verifier tls.Verifier) error {
	if !s.IsDisconnected() {
		return nil
	}

	s.setStatus(CONNECTING)

	conf := s.GetConfig()
	policy := config.ConnectionPolicy{
		Logger:        s.connectionLogger,
		XMPPLogger:    s.xmppLogger,
		DialerFactory: s.dialerFactory,
	}

	conn, err := policy.Connect(password, conf, verifier)
	if err != nil {
		s.setStatus(DISCONNECTED)

		return err
	}

	if s.getConnStatus() == CONNECTING {
		s.conn = conn
		s.setStatus(CONNECTED)

		conn.SignalPresence("")
		go s.watchRoster()
		go s.getVCard()
		go s.watchTimeout()
		go s.watchStanzas()
	} else {
		if s.conn != nil {
			s.conn.Close()
		}
	}

	return nil
}

// EncryptAndSendTo encrypts and sends the message to the given peer
func (s *session) EncryptAndSendTo(peer jid.WithoutResource, resource jid.Resource, message string) (trace int, delayed bool, err error) {
	//TODO: review whether it should create a conversation
	if s.IsConnected() {
		conversation, _ := s.convManager.EnsureConversationWith(peer, resource)
		trace, err = conversation.Send(s, resource, []byte(message))
		delayed = s.otrEventHandlers.Get(peer).ConsumeDelayedState(trace)
		return
	}
	return 0, false, &access.OfflineError{Msg: i18n.Local("Couldn't send message since we are not connected")}
}

func (s *session) terminateConversations() {
	s.convManager.TerminateAll()
}

func (s *session) connectionLost() {
	if s.IsDisconnected() {
		return
	}

	s.Close()
	s.publish(events.ConnectionLost)
}

// Close terminates all outstanding OTR conversations and closes the connection to the server
func (s *session) Close() {
	if s.IsDisconnected() {
		return
	}

	conn := s.conn
	if conn != nil {
		if !s.wantToBeOnline {
			s.terminateConversations()
		}
		s.setStatus(DISCONNECTED)
		conn.Close()
		s.conn = nil
	} else {
		s.setStatus(DISCONNECTED)
	}
}

func (s *session) CommandManager() otr_client.CommandManager {
	return s.cmdManager
}

func (s *session) SetCommandManager(c otr_client.CommandManager) {
	s.cmdManager = c
}

func (s *session) ConversationManager() otr_client.ConversationManager {
	return s.convManager
}

func (s *session) SetWantToBeOnline(val bool) {
	s.wantToBeOnline = val
}

func (s *session) PrivateKeys() []otr3.PrivateKey {
	return s.privateKeys
}

func (s *session) R() *roster.List {
	return s.r
}

func (s *session) SetConnector(c access.Connector) {
	s.connector = c
}

func (s *session) GroupDelimiter() string {
	return s.groupDelimiter
}

func (s *session) Config() *config.ApplicationConfig {
	return s.config
}

func (s *session) Conn() xi.Conn {
	return s.conn
}

func (s *session) SetSessionEventHandler(eh access.EventHandler) {
	s.sessionEventHandler = eh
}

func (s *session) SetConnectionLogger(l io.Writer) {
	s.connectionLogger = l
}

func (s *session) OtrEventHandlers() *otr_client.EventHandlers {
	return s.otrEventHandlers
}

func (s *session) SetLastActionTime(t time.Time) {
	s.lastActionTime = t
}

// SendPing is called to checks if account's connection still alive
func (s *session) SendPing() {
	reply, _, err := s.conn.SendPing()
	if err != nil {
		s.warn(fmt.Sprintf("Failure to ping server: %#v\n", err))
		return
	}

	pingTimeout := 10 * time.Second

	go func() {
		select {
		case <-time.After(pingTimeout):
			s.info("Ping timeout. Disconnecting...")
			s.setStatus(DISCONNECTED)
		case stanza, _ := <-reply:
			iq, ok := stanza.Value.(*data.ClientIQ)
			if !ok {
				return
			}
			if iq.Type == "error" {
				s.warn("Server does not support Ping")
				return
			}
		}
	}()
}

// StartSMP begins the SMP interactions for a conversation
func (s *session) StartSMP(peer jid.WithoutResource, resource jid.Resource, question, answer string) {
	// TODO: probably not possible.
	if peer == nil {
		s.alert("error: tried to start SMP with a nameless peer")
		return
	}
	conv, ok := s.convManager.GetConversationWith(peer, resource)
	if !ok {
		s.alert("error: tried to start SMP when a conversation does not exist")
		return
	}
	if err := conv.StartAuthenticate(s, resource, question, []byte(answer)); err != nil {
		s.alert("error: cannot start SMP: " + err.Error())
	}
}

// FinishSMP takes a user's SMP answer and finishes the protocol
func (s *session) FinishSMP(peer jid.WithoutResource, resource jid.Resource, answer string) {
	conv, ok := s.convManager.GetConversationWith(peer, resource)
	if !ok {
		s.alert("error: tried to finish SMP when a conversation does not exist")
		return
	}
	if err := conv.ProvideAuthenticationSecret(s, resource, []byte(answer)); err != nil {
		s.alert("error: cannot provide an authentication secret for SMP: " + err.Error())
	}
}

// AbortSMP will abort the current SMP interaction for a conversation
func (s *session) AbortSMP(peer jid.WithoutResource, resource jid.Resource) {
	conv, ok := s.convManager.GetConversationWith(peer, resource)
	if !ok {
		s.alert("error: tried to abort SMP when a conversation does not exist")
		return
	}
	if err := conv.AbortAuthentication(s, resource); err != nil {
		s.alert("error: cannot abort SMP: " + err.Error())
	}
}

// TODO: we also need a way to deal with when the TLV is received.

func (s *session) CreateSymmetricKeyFor(peer jid.Any) []byte {
	r := peer.PotentialResource()
	conv, ok := s.convManager.GetConversationWith(peer.NoResource(), r)
	if !ok {
		return nil
	}

	key, err := conv.CreateExtraSymmetricKey(s, r)
	if err != nil {
		return nil
	}

	return key
}
