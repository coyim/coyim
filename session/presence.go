package session

func (s *Session) ApprovePresenceSubscription(jid, id string) error {
	return s.Conn.SendPresence(jid, "subscribed", id)
}

func (s *Session) DenyPresenceSubscription(jid, id string) error {
	return s.Conn.SendPresence(jid, "unsubscribed", id)
}

func (s *Session) RequestPresenceSubscription(jid string) error {
	return s.Conn.SendPresence(jid, "subscribe", "" /* generate id */)
}
