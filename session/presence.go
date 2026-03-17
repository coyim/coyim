package session

import (
	"github.com/coyim/coyim/xmpp/jid"
)

// AutoApprove will automatically approve an incoming subscription request for the given peer
func (s *session) AutoApprove(jid string) {
	s.autoApprovesLock.Lock()
	defer s.autoApprovesLock.Unlock()

	s.autoApproves[jid] = true
}

func (s *session) hasAndRemoveAutoApprove(jid string) bool {
	s.autoApprovesLock.Lock()
	defer s.autoApprovesLock.Unlock()

	has := s.autoApproves[jid]
	if has {
		delete(s.autoApproves, jid)
	}

	return has
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
