package session

// ApprovePresenceSubscription is used to request subscription approval
func (s *Session) ApprovePresenceSubscription(jid, id string) error {
	return s.conn.SendPresence(jid, "subscribed", id)
}

// DenyPresenceSubscription is called to request subscription denial
func (s *Session) DenyPresenceSubscription(jid, id string) error {
	return s.conn.SendPresence(jid, "unsubscribed", id)
}

// RequestPresenceSubscription is called to request presence subscription
func (s *Session) RequestPresenceSubscription(jid string) error {
	return s.conn.SendPresence(jid, "subscribe", "" /* generate id */)
}
