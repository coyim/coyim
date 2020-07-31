package session

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) JoinRoom(rj jid.Bare, nickName string) {
	err := s.conn.SendPresence(rj.String(), "", "", "")
	if err != nil {
		s.log.WithError(err).Warn("when trying to enter room")
	}
	return
}

func (s *session) GetRoom(rj jid.Bare) (*muc.RoomListing, error) {
	rl := muc.NewRoomListing()
	rl.Jid = rj
	err := s.findOutMoreInformationAboutRoom(rl)
	return rl, err
}
