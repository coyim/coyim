package session

// AutoApprove will automatically approve an incoming subscription request for the given peer
func (s *session) AutoApprove(jid string) {
	s.autoApproves[jid] = true
}

// ApprovePresenceSubscription is used to request subscription approval
func (s *session) ApprovePresenceSubscription(jid, id string) error {
	return s.conn.SendPresence(jid, "subscribed", id, "")
}

// DenyPresenceSubscription is called to request subscription denial
func (s *session) DenyPresenceSubscription(jid, id string) error {
	return s.conn.SendPresence(jid, "unsubscribed", id, "")
}

// RequestPresenceSubscription is called to request presence subscription
func (s *session) RequestPresenceSubscription(jid, msg string) error {
	return s.conn.SendPresence(jid, "subscribe", "" /* generate id */, msg)
}
