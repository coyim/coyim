package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrUpdateOccupantRequest represents an early error that happened during occupant update request
	ErrUpdateOccupantRequest = errors.New("invalid occupant update request")
	// ErrUpdateOccupantResponse represents an invalid response for an occupant update request
	ErrUpdateOccupantResponse = errors.New("invalid response for room configuration request")
	// ErrOwnerAffiliationRevokeConflict represents a revocation conflict that happens when the only
	// owner of the room tries to change their affiliation
	ErrOwnerAffiliationRevokeConflict = errors.New("owner affiliation revoke conflict")
)

func (s *session) UpdateOccupantAffiliation(roomID jid.Bare, occupantNickname string, occupantID jid.Full, affiliation data.Affiliation, reason string) (<-chan bool, <-chan error) {
	return s.muc.updateOccupantAffiliation(roomID, occupantNickname, occupantID, affiliation, reason)
}

func (m *mucManager) updateOccupantAffiliation(roomID jid.Bare, occupantNickname string, occupantID jid.Full, affiliation data.Affiliation, reason string) (<-chan bool, <-chan error) {
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
			ec <- ErrUpdateOccupantRequest
			return
		}

		err = validateIqResponse(reply)
		if err != nil {
			ec <- err
			return
		}

		rc <- true
	}()

	return rc, ec
}
