package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	mucData "github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrRoomAffiliationsUpdate occurred when the room affiliation list couldn't be updated
	ErrRoomAffiliationsUpdate = errors.New("list affiliation couldn't be updated")
)

// UpdateOccupantAffiliations modifies the affiliations of a list of occupants in the given room
func (s *session) UpdateOccupantAffiliations(roomID jid.Bare, occupantAffiliations []*muc.RoomOccupantItem) (<-chan bool, <-chan error) {
	rc := make(chan bool)
	ec := make(chan error)

	go func() {
		ok, err := s.muc.updateOccupantAffiliations(roomID, occupantAffiliations)
		if err != nil {
			s.muc.log.WithField("room", roomID).WithError(err).Error("The occupant affiliations couldn't be updated")
			ec <- err
		} else {
			rc <- ok
		}
	}()

	return rc, ec
}

func (m *mucManager) updateOccupantAffiliations(roomID jid.Bare, occupantAffiliations []*muc.RoomOccupantItem) (bool, error) {
	// All xmpp servers should accept the modification of more than one occupant's affiliation
	// in a single request but the current version of Prosody server doesn't support that.
	// So, this is a temporary iteration that was implemented in order to
	// configure more than one occupant affiliation on Prosody servers.
	for _, i := range occupantAffiliations {
		err := m.modifyOccupantAffiliation(roomID, i)
		if err != nil {
			m.log.WithFields(log.Fields{
				"occupant":    i.Jid.String(),
				"affiliation": i.Affiliation.Name(),
				"reason":      i.Reason,
			}).WithError(err).Error("The occupant affiliation couldn't be updated")
			return false, err
		}
	}
	return true, nil
}

func (m *mucManager) modifyOccupantAffiliation(roomID jid.Bare, occupantAffiliations *muc.RoomOccupantItem) error {
	stanza, _, err := m.conn().SendIQ(roomID.String(), "set", newRoomOccupantAffiliationQuery(occupantAffiliations))
	if err != nil {
		return err
	}

	iq := <-stanza

	reply, ok := iq.Value.(*xmppData.ClientIQ)
	if !ok || reply.Type != "result" {
		return ErrRoomAffiliationsUpdate
	}

	return nil
}

func newRoomOccupantAffiliationQuery(item *muc.RoomOccupantItem) xmppData.MUCAdmin {
	return xmppData.MUCAdmin{
		Items: []xmppData.MUCItem{
			{
				Jid:         item.Jid.String(),
				Affiliation: item.Affiliation.Name(),
				Reason:      item.Reason,
			},
		},
	}
}

// GetRoomOccupantsByAffiliation retrieves an occupant list based on a given affiliation
func (s *session) GetRoomOccupantsByAffiliation(roomID jid.Bare, a mucData.Affiliation) (<-chan []*muc.RoomOccupantItem, <-chan error) {
	occupantItems := make(chan []*muc.RoomOccupantItem)
	errorChannel := make(chan error)

	rc, ec := s.muc.requestRoomOccupantsByAffiliation(roomID, a)
	go func() {
		select {
		case items := <-rc:
			occ, err := parseMUCItemsToRoomOccupantItems(items)
			if err != nil {
				s.muc.log.WithError(err).WithField("affiliation", a).Error("cannot parse MUCItems to OccupantItems")
				errorChannel <- err
				return
			}
			occupantItems <- occ
		case err := <-ec:
			s.muc.log.WithError(err).WithField("affiliation", a).Error("cannot retrieve occupants")
			errorChannel <- err
		}
	}()

	return occupantItems, errorChannel
}
