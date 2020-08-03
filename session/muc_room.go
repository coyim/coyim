package session

import (
	"fmt"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) JoinRoom(rj jid.Bare, nickName string) {
	to := fmt.Sprintf("%s/%s", rj.String(), nickName)
	err := s.conn.SendMUCPresence(to)
	if err != nil {
		s.log.WithError(err).Warn("when trying to enter room")
	}
	return
}

func (s *session) HasRoom(rj jid.Bare) bool {
	_, e := s.Conn().QueryServiceInformation(rj.String())
	if e != nil {
		s.log.WithError(e).Debug("HasRoom() had an error")
		return false
	}
	return true
}

func (s *session) GetRoom(rj jid.Bare, rl *muc.RoomListing) {
	// TODO, checking should be called this function
	rl = muc.NewRoomListing()
	rl.Jid = rj
	go s.findOutMoreInformationAboutRoom(rl)
}
