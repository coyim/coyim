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

func (s *session) HasRoom(rj jid.Bare) (<-chan bool, <-chan error) {
	resultChannel := make(chan bool, 1)
	errorChannel := make(chan error)
	go func() {
		r, err := s.Conn().EntityExists(rj.String())
		if !r || err != nil {
			if err != nil {
				s.log.WithError(err).Warn("HasRoom() had an error")
				errorChannel <- err
				return
			}
			resultChannel <- false
			return
		}
		// Make sure the entity is a Room
		idents, features, ok := s.Conn().DiscoveryFeaturesAndIdentities(rj.String())
		if !ok {
			resultChannel <- false
			return
		}
		// Checking Identities
		ident, hasIdent := hasIdentity(idents, "conference", "text")
		if !hasIdent {
			resultChannel <- false
			return
		}
		// Checking Features
		if !hasFeatures(features, "http://jabber.org/protocol/muc") {
			resultChannel <- false
			return
		}
		// Checking Bare JID
		bares := fmt.Sprintf("%s@%s", ident, rj.Host())
		barerj, ok := jid.Parse(bares).(jid.Bare)
		if !ok || barerj != rj {
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
