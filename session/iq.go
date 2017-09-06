package session

import "github.com/twstrike/coyim/xmpp/data"

func (s *session) sendIQError(stanza *data.ClientIQ, reply interface{}) {
	s.sendIQReply(stanza, "error", reply)
}

func (s *session) sendIQResult(stanza *data.ClientIQ, reply interface{}) {
	s.sendIQReply(stanza, "result", reply)
}

func (s *session) sendIQReply(stanza *data.ClientIQ, tp string, reply interface{}) {
	if err := s.conn.SendIQReply(stanza.From, tp, stanza.ID, reply); err != nil {
		s.alert("Failed to send IQ message: " + err.Error())
	}
}
