package session

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) JoinRoom(rj jid.Bare, nickName string) error {
	to := jid.NewFull(rj.Local(), rj.Host(), jid.NewResource(nickName))
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

		// TODO[OB]-MUC: It reads a bit weirdly mixing up these two results
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
			// TODO[OB]-MUC: Is this really correct?
			resultChannel <- false
			return
		}
		ident, hasIdent := hasIdentity(idents, "conference", "text")
		if !hasIdent {
			resultChannel <- false
			return
		}
		if !hasFeatures(features, "http://jabber.org/protocol/muc") {
			resultChannel <- false
			return
		}
		barerj := jid.NewBare(jid.NewLocal(ident), rj.Host())
		if barerj != rj {
			// TODO[OB]-MUC: I feel like this mixes up two conerns
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
