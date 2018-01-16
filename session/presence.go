package session

import (
	"github.com/coyim/coyim/xmpp/jid"
)

// AutoApprove will automatically approve an incoming subscription request for the given peer
func (s *session) AutoApprove(jid string) {
	s.autoApproves[jid] = true
}

// ApprovePresenceSubscription is used to request subscription approval
func (s *session) ApprovePresenceSubscription(jid jid.WithoutResource, id string) error {
	return s.conn.SendPresence(jid.String(), "subscribed", id, "")
}

// DenyPresenceSubscription is called to request subscription denial
func (s *session) DenyPresenceSubscription(jid jid.WithoutResource, id string) error {
	return s.conn.SendPresence(jid.String(), "unsubscribed", id, "")
}

// RequestPresenceSubscription is called to request presence subscription
func (s *session) RequestPresenceSubscription(jid jid.WithoutResource, msg string) error {
	return s.conn.SendPresence(jid.String(), "subscribe", "" /* generate id */, msg)
}
