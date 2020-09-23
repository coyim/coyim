package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (s *session) JoinRoom(ident jid.Bare, nickname string) error {
	// TODO: The problem with this method is that it only _starts_ the process of joining the room
	// It would be good to have a method that takes responsibility for the whole flow
	to := ident.WithResource(jid.NewResource(nickname))
	err := s.conn.SendMUCPresence(to.String())
	if err != nil {
		s.log.WithFields(log.Fields{
			"room":     ident,
			"nickname": nickname,
		}).WithError(err).Error("An error occurred trying join the room")
		return err
	}
	return nil
}

func (s *session) HasRoom(rj jid.Bare, wantRoomInfo chan<- *muc.RoomListing) (<-chan bool, <-chan error) {
	// TODO: This should be refactored into smaller, more focused helper methods
	// TODO: unify logging in all the below, so everything logs the room name
	// TODO: rename these variables
	resultChannel := make(chan bool)
	errorChannel := make(chan error)
	go func() {
		r, err := s.Conn().EntityExists(rj.String())
		if err != nil {
			s.log.WithError(err).Error("An error occurred searching the entity on the server")
			errorChannel <- err
			return
		}
		if !r {
			resultChannel <- false
			return
		}
		idents, features, ok := s.Conn().DiscoveryFeaturesAndIdentities(rj.String())
		if !ok {
			err := errors.New("Something went wrong discovering the features and identities of the room")
			s.log.WithField("room", rj).WithError(err).Error("An error occurred trying to get the features and identities from the server")
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

		if wantRoomInfo != nil {
			s.GetRoom(rj, wantRoomInfo)
		}
	}()
	return resultChannel, errorChannel
}

// GetRoom will block, waiting to get the room information
func (s *session) GetRoom(rj jid.Bare, result chan<- *muc.RoomListing) {
	// TODO: make this method unnecessary by changing the GUI parts to not use it

	rl := muc.NewRoomListing()
	rl.Jid = rj
	// This is a little bit redundant since we already asked for this once
	// The right solution is to use the values from above, but that would be an extensive refactoring
	// so we will wait with that for now
	s.findOutMoreInformationAboutRoom(rl)
	result <- rl
}

func createRoomRecipient(room jid.Bare, nickname string) jid.Full {
	return jid.NewFull(room.Local(), room.Host(), jid.NewResource(nickname))
}

func (s *session) LeaveRoom(room jid.Bare, nickname string) (chan bool, chan error) {
	to := createRoomRecipient(room, nickname).String()

	// TODO: rename variables
	resultCh := make(chan bool)
	errorCh := make(chan error)
	go func() {
		err := s.conn.SendPresence(to, "unavailable", "", "")
		if err != nil {
			s.log.WithField("to", to).WithError(err).Error("error trying to leave room")
			errorCh <- err
			return
		}
		s.muc.roomManager.LeaveRoom(room)
		resultCh <- true
	}()

	return resultCh, errorCh
}
