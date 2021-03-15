package session

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/otrclient"
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
	conn xi.Conn
	log  coylog.Logger
	r    *roster.List

	muc *mucManager

	connStatus     connStatus
	connStatusLock sync.RWMutex

	privateKeys []otr3.PrivateKey

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

	cmdManager  otrclient.CommandManager
	convManager otrclient.ConversationManager

	dialerFactory func(tls.Verifier, tls.Factory) xi.Dialer

	autoApproves map[string]bool

	nicknames []string

	pendingEvents     int
	pendingEventsLock sync.Mutex
	eventsReachedZero chan bool

	resource string

	doneBadStanza chan<- bool
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

	acc := cu.Account
	l := log.WithField("account", acc)
	l.WithField("numKeys", len(allKeys)).Info("Loading configured keys")
	for _, pp := range allKeys {
		_, ok, parsedKey := otr3.ParsePrivateKey(pp)
		if ok {
			result = append(result, parsedKey)
			l.WithField("key", config.FormatFingerprint(parsedKey.PublicKey().Fingerprint())).Info("Loaded key")
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
func Factory(c *config.ApplicationConfig, cu *config.Account, df func(tls.Verifier, tls.Factory) xi.Dialer) access.Session {
	// Make xmppLogger go to in memory STRING and/or the log file

	inMemoryLog, xmppLogger := CreateXMPPLogger(c.RawLogFile)

	sessionLog := log.WithFields(log.Fields{
		"account": cu.Account,
	})

	s := &session{
		config:        c,
		accountConfig: cu,

		r:              roster.New(),
		lastActionTime: time.Now(),

		timeouts: make(map[data.Cookie]time.Time),

		autoApproves: make(map[string]bool),

		inMemoryLog:   inMemoryLog,
		xmppLogger:    xmppLogger,
		log:           sessionLog.WithField("component", "session"),
		dialerFactory: df,
	}

	s.ReloadKeys()
	s.convManager = otrclient.NewConversationManager(s.newConversation, s, cu.Account, s.onOtrEventHandlerCreate, sessionLog.WithField("component", "otr"))
	s.muc = newMUCManager(s.log, s.Conn, s.publishEvent)

	go observe(s)
	go checkReconnect(s)

	return s
}

// ReloadKeys will reload the keys from the configuration
func (s *session) ReloadKeys() {
	s.privateKeys = parseFromConfig(s.accountConfig)
}

// Send will send the given message to the receiver given
func (s *session) Send(peer jid.Any, msg string, otr bool) error {
	conn, ok := s.connection()
	if ok {
		s.log.WithFields(log.Fields{
			"to":      peer,
			"sentMsg": msg,
		}).Debug("Send()")
		return conn.Send(peer.String(), msg, otr)
	}
	return &access.OfflineError{Msg: i18n.Local("Couldn't send message since we are not connected")}
}

func (s *session) receivedStreamError(stanza *data.StreamError) bool {
	s.log.WithField("stanza", stanza.String()).Error("Exiting in response to fatal error from server")
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
	s.log.WithField("stanza", fmt.Sprintf("%#v", stanza)).Debug("receivedClientMessage()")

	if stanza.Body == "" && len(stanza.Extensions) > 0 {
		s.processExtensions(stanza)
		return true
	}

	peer := jid.Parse(stanza.From)

	if stanza.Encryption != nil {
		s.processEncryption(peer, stanza.Encryption)
	}

	// TODO: it feels iffy that we have error and groupchat special handled here
	// But not checking on the "message" type.
	switch stanza.Type {
	case "groupchat":
		s.muc.receiveClientMessage(stanza)
		return true
	case "error":
		// TODO: investigate which errors are NOT recoverable, and return false
		// to close the connection
		// https://xmpp.org/rfcs/rfc3920.html#stanzas-error

		if isMUCError(stanza.Error) {
			s.muc.handleMUCErrorMessage(jid.ParseBare(stanza.From), stanza.Error)
		} else {
			s.log.WithField("error", stanza.Error).WithField("peer", peer.NoResource()).Error("Error reported from peer")
		}

		return true
	}

	messageTime := retrieveMessageTime(stanza)
	s.receiveClientMessage(peer, messageTime, stanza.Body)

	return true
}

func (s *session) receivedClientPresence(stanza *data.ClientPresence) bool {
	jj := jid.Parse(stanza.From)
	jjnr := jj.NoResource()

	switch stanza.Type {
	case "subscribe":
		jjr := jjnr.String()
		if s.autoApproves[jjr] {
			delete(s.autoApproves, jjr)
			_ = s.ApprovePresenceSubscription(jjnr, stanza.ID)
		} else {
			s.r.SubscribeRequest(jjnr, either(stanza.ID, "0000"), s.GetConfig().ID())
			s.publishPeerEvent(
				events.SubscriptionRequest,
				jjnr,
			)
		}
	case "unavailable":
		if isMUCUserPresence(stanza) {
			s.muc.handlePresence(stanza)
			return true
		}

		if !s.r.PeerBecameUnavailable(jj) {
			return true
		}

		s.publishEvent(events.Presence{
			ClientPresence: stanza,
			Gone:           true,
		})
	case "":
		if jj.NoResource().String() == jj.String() {
			// This happens if a malfunctioning client/server is
			// sending presence information without a resource.
			// This is likely a bug
			s.log.WithFields(log.Fields{
				"from":   stanza.From,
				"stanza": stanza,
			}).Warn("Got a presence without resource in 'from' - this is likely an error")
			return true
		}

		if isMUCUserPresence(stanza) {
			s.muc.handlePresence(stanza)
			return true
		}

		if !s.r.PeerPresenceUpdate(jj.(jid.WithResource), stanza.Show, stanza.Status, s.GetConfig().ID()) {
			return true
		}

		//TODO: If to == "" this is our own presence confirmation.
		//"From" is how we are identified (will be JID/"some-id")
		//Same thing happens for the group-chat, but in this case it tell us also what are our affiliations and roles.
		//Thats why I'm worried about handling this as a regular peer presence - which is not.
		s.publishEvent(events.Presence{
			ClientPresence: stanza,
			Gone:           false,
		})
	case "subscribed":
		s.r.Subscribed(jjnr)
		s.publishPeerEvent(
			events.Subscribed,
			jjnr,
		)
	case "unsubscribe":
		s.r.Unsubscribed(jjnr)
		s.publishPeerEvent(
			events.Unsubscribe,
			jjnr,
		)
	case "unsubscribed":
		// Ignore
	case "error":
		s.log.WithFields(log.Fields{
			"from":  stanza.From,
			"error": stanza.Error,
		}).Error("Got a presence error")

		s.r.LatestError(jjnr, stanza.Error.Code, stanza.Error.Type, stanza.Error.Condition.XMLName.Space+" "+stanza.Error.Condition.XMLName.Local)

		if isMUCErrorPresence(stanza.Error) {
			s.muc.handleMUCErrorPresence(stanza.From, stanza.Error)
		}

	default:
		s.log.WithFields(log.Fields{
			"stanza": stanza,
			"from":   stanza.From,
		}).Warn("Unrecognized presence")
	}

	return true
}

func (s *session) receivedClientIQ(stanza *data.ClientIQ) bool {
	if stanza.Type == "get" || stanza.Type == "set" {
		reply, iqtype, ignore := s.processIQ(stanza)
		if ignore {
			return true
		}

		if reply == nil {
			reply = data.ErrorReply{
				Type:  "cancel",
				Error: data.ErrorBadRequest{},
			}
			if iqtype == "" {
				iqtype = "error"
			}
		}

		if iqtype == "" {
			iqtype = "result"
		}

		s.sendIQReply(stanza, iqtype, reply)
		return true
	}
	s.log.WithField("stanza", stanza).Info("unrecognized iq")
	return true
}

func (s *session) receiveStanza(stanzaChan chan data.Stanza) bool {
	rawStanza, ok := <-stanzaChan
	if !ok {
		return false
	}

	result := false

	switch stanza := rawStanza.Value.(type) {
	case *data.StreamError:
		result = s.receivedStreamError(stanza)
	case *data.ClientMessage:
		result = s.receivedClientMessage(stanza)
	case *data.ClientPresence:
		result = s.receivedClientPresence(stanza)
	case *data.ClientIQ:
		result = s.receivedClientIQ(stanza)
	default:
		s.log.WithFields(log.Fields{
			"name":  rawStanza.Name,
			"value": rawStanza.Value,
		}).Info("unhandled stanza")
		result = true
	}

	return result
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
		s.log.WithError(err).Error("error reading XMPP message")
	}
	if s.doneBadStanza != nil {
		s.doneBadStanza <- true
	}
}

func (s *session) rosterUpdated() {
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
		s.log.WithField("from", stanza.From).Warn("Ignoring roster IQ from bad address")
		return nil, "", true
	}
	var rst data.Roster
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&rst); err != nil || len(rst.Item) == 0 {
		s.log.Warn("Failed to parse roster push IQ")
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
		s.log.WithField("from", jid.String()).Warn("No pending subscription")
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
		s.log.WithError(err).Warn("Error sending presence stanza")
		return
	}

	if isConfirm {
		_ = s.RequestPresenceSubscription(jid, "")
	}
}

func (s *session) receiveClientMessage(peer jid.Any, when time.Time, body string) {
	msg := []byte(body)
	conversation, _ := s.convManager.EnsureConversationWith(peer, msg)
	out, err := conversation.Receive(msg)
	encrypted := conversation.IsEncrypted()

	if err != nil {
		s.log.WithError(err).WithField("peer", peer.String()).Error("While processing message from peer")
	}

	eh := conversation.EventHandler()
	change := eh.ConsumeSecurityChange()
	switch change {
	case otrclient.NewKeys:
		s.newOTRKeys(peer.(jid.WithResource), conversation)
	case otrclient.RenewedKeys:
		s.renewedOTRKeys(peer.(jid.WithResource), conversation)
	case otrclient.ConversationEnded:
		s.otrEnded(peer.(jid.WithResource))
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
			s.log.WithField("peer", peer).Info("Peer has ended the secure conversation.")
			err := conversation.EndEncryptedChat()
			if err != nil {
				s.log.WithError(err).WithField("peer", peer).Info("Unable to automatically tear down OTR conversation with peer")
				break
			}

			s.log.WithField("peer", peer).Info("Secure session with peer has been automatically ended. Messages will be sent in the clear until another OTR session is established.")
		} else {
			s.log.WithField("peer", peer).Info("Peer has ended the secure conversation. You should do likewise")
		}
	case otrclient.SMPSecretNeeded:
		s.publishSMPEvent(events.SecretNeeded, peer.(jid.WithResource), eh.SmpQuestion)
	case otrclient.SMPComplete:
		s.publishSMPEvent(events.Success, peer.(jid.WithResource), "")
		s.cmdManager.ExecuteCmd(otrclient.AuthorizeFingerprintCmd{
			Account:     s.GetConfig(),
			Session:     s,
			Peer:        peer.NoResource(),
			Fingerprint: conversation.TheirFingerprint(),
			Tag:         "SMP",
		})
	case otrclient.SMPFailed:
		s.publishSMPEvent(events.Failure, peer.(jid.WithResource), "")
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

	/* #nosec G204 */
	cmd := exec.Command(s.config.NotifyCommand[0], s.config.NotifyCommand[1:]...)
	go func() {
		if err := cmd.Run(); err != nil {
			s.log.WithError(err).Error("Failed to run notify command")
		}
	}()
}

// AwaitVersionReply listens on the channel and waits for the version reply
func (s *session) AwaitVersionReply(ch <-chan data.Stanza, user string) {
	stanza, ok := <-ch
	if !ok {
		s.log.WithField("user", user).Warn("Version request timed out")
		return
	}
	reply, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		s.log.WithField("user", user).Warn("Version request resulted in bad reply type")
		return
	}

	if reply.Type == "error" {
		s.log.WithField("user", user).Warn("Version request resulted in XMPP error")
		return
	} else if reply.Type != "result" {
		s.log.WithField("user", user).WithField("type", reply.Type).Warn("Version request resulted in response with unknown type")
		return
	}

	buf := bytes.NewBuffer(reply.Query)
	var versionReply data.VersionReply
	if err := xml.NewDecoder(buf).Decode(&versionReply); err != nil {
		s.log.WithField("user", user).WithError(err).Warn("Failed to parse version reply")
		return
	}

	s.log.WithField("user", user).WithField("version", versionReply).Info("Version reply from peer")
}

var tickInterval = time.Second

func (s *session) watchTimeout() {
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
				s.log.WithField("cookie", cookie).Debug("session: cookie has expired")
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

	s.log.Info("Fetching VCard")

	vcardReply, _, err := conn.RequestVCard()
	if err != nil {
		s.log.WithError(err).Error("Failed to request vcard")
		return
	}

	vcardStanza, ok := <-vcardReply
	if !ok {
		s.log.Debug("session: vcard request cancelled or timed out")
		return
	}

	vc, err := data.ParseVCard(vcardStanza)
	if err != nil {
		s.log.WithError(err).Error("Failed to parse vcard")
		return
	}

	s.nicknames = []string{vc.Nickname, vc.FullName}
}

func (s *session) DisplayName() string {
	return either(either(s.accountConfig.Nickname, firstNonEmpty(s.nicknames...)), s.accountConfig.Account)
}

func (s *session) requestRoster() bool {
	conn, ok := s.connection()
	if !ok {
		return false
	}

	s.log.Info("Fetching roster")

	delim, err := conn.GetRosterDelimiter()
	if err != nil || delim == "" {
		delim = defaultDelimiter
	}
	s.groupDelimiter = delim

	rosterReply, _, err := conn.RequestRoster()
	if err != nil {
		s.log.WithError(err).Error("Failed to request roster")
		return true
	}

	rosterStanza, ok := <-rosterReply
	if !ok {
		//TODO: should we retry the request in such case?
		s.log.Debug("session: roster request cancelled or timed out")
		return true
	}

	rst, err := data.ParseRoster(rosterStanza)
	if err != nil {
		s.log.WithError(err).Error("Failed to parse roster")
		return true
	}

	for _, rr := range rst {
		s.addOrMergeNewPeer(rr, s.GetConfig())
	}

	s.log.Info("Roster received")
	s.rosterUpdated()

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
		Log:           s.log,
		XMPPLogger:    s.xmppLogger,
		DialerFactory: s.dialerFactory,
	}

	resource := ""
	if s.wantToBeOnline {
		resource = s.resource
	}

	s.log.WithFields(log.Fields{
		"resource":       resource,
		"wantToBeOnline": s.wantToBeOnline,
	}).Debug("Connect()")

	conn, err := policy.Connect(password, resource, conf, verifier)
	if err != nil {
		s.log.WithError(err).Error("failed to connect")

		s.setStatus(DISCONNECTED)

		return err
	}

	if s.getConnStatus() == CONNECTING {
		s.conn = conn
		s.setStatus(CONNECTED)
		s.resource = s.conn.GetJIDResource()

		_ = conn.SignalPresence("")
		go s.watchRoster()
		go func() {
			if s.conn.ServerHasFeature("vcard-temp") {
				s.getVCard()
			}
		}()
		go s.watchTimeout()
		go s.watchStanzas()
	} else {
		if s.conn != nil {
			_ = s.conn.Close()
		}
	}

	return nil
}

// EncryptAndSendTo encrypts and sends the message to the given peer
func (s *session) EncryptAndSendTo(peer jid.Any, message string) (trace int, delayed bool, err error) {
	if s.IsConnected() {
		c, _ := s.convManager.EnsureConversationWith(peer, nil)
		trace, err = c.Send([]byte(message))
		delayed = c.EventHandler().ConsumeDelayedState(trace)
		return
	}
	return 0, false, &access.OfflineError{Msg: i18n.Local("Couldn't send message since we are not connected")}
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
		_ = conn.Close()
		s.conn = nil
	} else {
		s.setStatus(DISCONNECTED)
	}
}

func (s *session) CommandManager() otrclient.CommandManager {
	return s.cmdManager
}

func (s *session) SetCommandManager(c otrclient.CommandManager) {
	s.cmdManager = c
}

func (s *session) ConversationManager() otrclient.ConversationManager {
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

func (s *session) SetLastActionTime(t time.Time) {
	s.lastActionTime = t
}

// SendPing is called to checks if account's connection still alive
func (s *session) SendPing() {
	reply, _, err := s.conn.SendPing()
	if err != nil {
		s.log.WithError(err).Warn("Failure to ping server")
		return
	}

	pingTimeout := 10 * time.Second

	go func() {
		select {
		case <-time.After(pingTimeout):
			s.log.Info("Ping timeout. Disconnecting...")
			s.setStatus(DISCONNECTED)
		case stanza := <-reply:
			iq, ok := stanza.Value.(*data.ClientIQ)
			if !ok {
				return
			}
			if iq.Type == "error" {
				s.log.Warn("Server does not support Ping")
				return
			}
		}
	}()
}

// SendMUCMessage is used to send a group chat message
func (s *session) SendMUCMessage(to, from, body string) error {
	if !s.IsConnected() {
		return errors.New("session is not connected")
	}

	return s.muc.sendMessage(to, from, body)
}
