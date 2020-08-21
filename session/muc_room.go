package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) JoinRoom(rj jid.Bare, nickName string) error {
	to := rj.WithResource(jid.NewResource(nickName))
	err := s.conn.SendMUCPresence(to.String())
	if err != nil {
		s.log.WithError(err).Warn("when trying to enter room")
		return err
	}
	return nil
}

func (s *session) HasRoom(rj jid.Bare) (<-chan bool, <-chan error) {
	resultChannel := make(chan bool)
	errorChannel := make(chan error)
	go func() {
		r, err := s.Conn().EntityExists(rj.String())
		if err != nil {
			s.log.WithError(err).Error("HasRoom() had an error")
			errorChannel <- err
			return
		}
		if !r {
			resultChannel <- false
			return
		}
		// Make sure the entity is a Room
		idents, features, ok := s.Conn().DiscoveryFeaturesAndIdentities(rj.String())
		if !ok {
			err := errors.New("An error ocurred trying to get the features and identities from the server")
			s.log.WithError(err).Error("DiscoveryFeaturesAndIdentities() had an error")
			errorChannel <- err
			return
		}
		_, hasIdent := hasIdentity(idents, "conference", "text")
		if !hasIdent {
			resultChannel <- false
			return
		}
		if !hasFeatures(features, "http://jabber.org/protocol/muc") {
			resultChannel <- false
			return
		}
		resultChannel <- true
	}()
	return resultChannel, errorChannel
}

func (s *session) GetRoom(rj jid.Bare, rl *muc.RoomListing) {
	// TODO, check it out the best way to do this function to get all
	// the information of the room from the server
	rl = muc.NewRoomListing()
	rl.Jid = rj
	go s.findOutMoreInformationAboutRoom(rl)
}
