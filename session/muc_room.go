package session

import (
	"fmt"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) JoinRoom(rj jid.Bare, nickName string) {
	// TODO[OB]-MUC: Better to use a factory function here
	// TODO[OB]-MUC: You don't need to call String() when using the %s modifier
	to := fmt.Sprintf("%s/%s", rj.String(), nickName)
	err := s.conn.SendMUCPresence(to)
	if err != nil {
		// TODO[OB]-MUC: This error condition shouldn't be returned to someone?
		s.log.WithError(err).Warn("when trying to enter room")
	}
}

// TODO[OB]-MUC: Lots of unnecessary comments in this method

func (s *session) HasRoom(rj jid.Bare) (<-chan bool, <-chan error) {
	// TODO[OB]-MUC: Why is this one buffered?
	resultChannel := make(chan bool, 1)
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
		// TODO[OB]-MUC: Better to use a factory composition method here
		bares := fmt.Sprintf("%s@%s", ident, rj.Host())
		barerj, ok := jid.Parse(bares).(jid.Bare)
		if !ok || barerj != rj {
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
