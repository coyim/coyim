package session

import "github.com/coyim/coyim/xmpp/data"

// AutoApprove will automatically approve an incoming subscription request for the given peer
func (s *session) AutoApprove(jid string) {
	s.autoApproves[jid] = true
}

// ApprovePresenceSubscription is used to request subscription approval
func (s *session) ApprovePresenceSubscription(jid data.JIDWithoutResource, id string) error {
	return s.conn.SendPresence(jid.Representation(), "subscribed", id, "")
}

// DenyPresenceSubscription is called to request subscription denial
func (s *session) DenyPresenceSubscription(jid data.JIDWithoutResource, id string) error {
	return s.conn.SendPresence(jid.Representation(), "unsubscribed", id, "")
}

// RequestPresenceSubscription is called to request presence subscription
func (s *session) RequestPresenceSubscription(jid data.JIDWithoutResource, msg string) error {
	return s.conn.SendPresence(jid.Representation(), "subscribe", "" /* generate id */, msg)
}
