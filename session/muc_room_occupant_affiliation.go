package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrUpdateOccupantAffiliation represents an early error that happened during occupant affiliation update request
	ErrUpdateOccupantAffiliation = errors.New("invalid occupant update affiliation request")
	// ErrUpdateOccupantAffiliationResponse represents an invalid response for an occupant affiliation update request
	ErrUpdateOccupantAffiliationResponse = errors.New("invalid response for room configuration request")
)

func (s *session) UpdateOccupantAffiliation(roomID jid.Bare, occupantID jid.Full, affiliation data.Affiliation, reason string) (<-chan bool, <-chan error) {
	return s.muc.updateOccupantAffiliation(roomID, occupantID, affiliation, reason)
}

func (m *mucManager) updateOccupantAffiliation(roomID jid.Bare, occupantID jid.Full, affiliation data.Affiliation, reason string) (<-chan bool, <-chan error) {
	l := m.log.WithFields(log.Fields{
		"room":        roomID,
		"occupant":    occupantID,
		"affiliation": affiliation.Name(),
	})

	rc := make(chan bool)
	ec := make(chan error)

	go func() {
		reply, _, err := m.conn().SendIQ(roomID.String(), "set", &xmppData.MUCAdmin{
			Item: &xmppData.MUCItem{
				Affiliation: affiliation.Name(),
				Jid:         occupantID.String(),
				Reason:      reason,
			},
		})

		if err != nil {
			l.WithError(err).Error("An error occurred when updating the occupant affiliation")
			ec <- ErrUpdateOccupantAffiliation
			return
		}

		err = validateIqResponse(reply)
		if err != nil {
			l.WithError(err).Error("An error occurred when trying to read the response from the room configuration rollback request")
			ec <- ErrUpdateOccupantAffiliationResponse
			return
		}

		rc <- true
	}()

	return rc, ec
}
