package session

import (
	"errors"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (s *session) JoinRoom(roomID jid.Bare, nickname string) error {
	// TODO: The problem with this method is that it only _starts_ the process of joining the room
	// It would be good to have a method that takes responsibility for the whole flow
	to := roomID.WithResource(jid.NewResource(nickname))
	err := s.conn.SendMUCPresence(to.String())
	if err != nil {
		s.log.WithFields(log.Fields{
			"room":     roomID,
			"nickname": nickname,
		}).WithError(err).Error("An error occurred trying join the room")
		return err
	}
	return nil
}

type hasRoomContext struct {
	s               *session
	resultChannel   chan bool
	errorChannel    chan error
	roomID          jid.Bare
	foundIdentities []data.DiscoveryIdentity
	foundFeatures   []string
	log             coylog.Logger
}

func (rc *hasRoomContext) checkIfRoomExists(wantRoomInfo chan<- *muc.RoomListing) {
	steps := []func() (bool, error){
		rc.searchEntityOnServer,
		rc.discoverFeaturesAndIdentities,
		rc.hasRoomIdentity,
		rc.hasRoomFeature,
	}

	for _, f := range steps {
		ok, err := f()
		if err != nil {
			rc.s.log.WithError(err).Error("An error occurred while if the room exists on the server")
			rc.errorChannel <- err
			return
		}

		if !ok {
			rc.resultChannel <- false
			return
		}
	}

	rc.resultChannel <- true

	if wantRoomInfo != nil {
		rc.s.GetRoomListing(rc.roomID, wantRoomInfo)
	}
}

func (rc *hasRoomContext) searchEntityOnServer() (bool, error) {
	exists, err := rc.s.Conn().EntityExists(rc.roomID.String())
	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (rc *hasRoomContext) discoverFeaturesAndIdentities() (bool, error) {
	i, f, ok := rc.s.Conn().DiscoveryFeaturesAndIdentities(rc.roomID.String())
	if !ok {
		return false, errors.New("the room doesn't exists")
	}

	rc.foundIdentities = i
	rc.foundFeatures = f

	return true, nil
}

func (rc *hasRoomContext) hasRoomIdentity() (bool, error) {
	_, ok := hasIdentity(rc.foundIdentities, "conference", "text")
	if !ok {
		return false, errors.New("invalid room identity")
	}
	return true, nil
}

func (rc *hasRoomContext) hasRoomFeature() (bool, error) {
	if !hasFeatures(rc.foundFeatures, "http://jabber.org/protocol/muc") {
		return false, errors.New("invalid room feature")
	}
	return true, nil
}

func (s *session) HasRoom(roomID jid.Bare, wantRoomInfo chan<- *muc.RoomListing) (<-chan bool, <-chan error) {
	c := &hasRoomContext{
		s:             s,
		roomID:        roomID,
		resultChannel: make(chan bool),
		errorChannel:  make(chan error),
		log:           s.log.WithField("room", roomID),
	}

	go c.checkIfRoomExists(wantRoomInfo)

	return c.resultChannel, c.errorChannel
}

// LoadRoomInfo will load the given room configuration (RoomListing)
func (s *session) LoadRoomInfo(roomID jid.Bare) {
	s.muc.loadRoomInfo(roomID)
}

// GetRoomListing will block, waiting to get the room information
func (s *session) GetRoomListing(roomID jid.Bare, result chan<- *muc.RoomListing) {
	// TODO: make this method unnecessary by changing the GUI parts to not use it
	s.muc.getRoomListing(roomID, result)
}

func (s *session) LeaveRoom(room jid.Bare, nickname string) (chan bool, chan error) {
	to := createRoomRecipient(room, nickname).String()

	result := make(chan bool)
	errors := make(chan error)
	go func() {
		err := s.conn.SendPresence(to, "unavailable", "", "")
		if err != nil {
			s.log.WithField("to", to).WithError(err).Error("error trying to leave room")
			errors <- err
			return
		}
		s.muc.roomManager.LeaveRoom(room)
		result <- true
	}()

	return result, errors
}
