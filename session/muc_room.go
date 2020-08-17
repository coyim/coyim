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
}

func (s *session) HasRoom(rj jid.Bare) <-chan bool {
	result := make(chan bool, 1)
	go func() {
		r, err := s.Conn().EntityExists(rj.String())
		if !r {
			if err != nil {
				s.log.WithError(err).Debug("HasRoom() had an error")
			}
			result <- false
			return
		}
		// Make sure the entity is a Room
		idents, features, ok := s.Conn().DiscoveryFeaturesAndIdentities(rj.String())
		// Checking Identities
		identName, hasIdent := s.hasIdentity(idents, "conference", "text")
		if !ok || !hasIdent {
			result <- false
			return
		}
		// Checking Features
		if !s.hasFeatures(features, "http://jabber.org/protocol/muc") {
			result <- false
			return
		}
		// Checking Identity Name with the room specified
		if identName != rj.Local().String() {
			result <- false
			return
		}
		result <- true
	}()
	return result
}

func (s *session) GetRoom(rj jid.Bare, rl *muc.RoomListing) {
	// TODO, check it out the best way to do this function to get all
	// the information of the room from the server
	rl = muc.NewRoomListing()
	rl.Jid = rj
	go s.findOutMoreInformationAboutRoom(rl)
}
