package session

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/session/muc"
	mucData "github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrRoomConfigSubmit represents an early error that happened during room configuration submit
	ErrRoomConfigSubmit = errors.New("invalid room configuration submit request")
	// ErrRoomConfigSubmitResponse represents an invalid response for a room configuration request
	ErrRoomConfigSubmitResponse = errors.New("invalid response for room configuration request")
	// ErrRoomConfigCancel represents an early error that happened during a room configuration cancel request
	ErrRoomConfigCancel = errors.New("invalid room configuration cancel request")
	// ErrRoomConfigCancelResponse represents an invalid response for a room configuration cancel request
	ErrRoomConfigCancelResponse = errors.New("invalid response for room configuration cancel request")
)

const mucRequestGetRoomConfigForm mucRequestType = "getRoomConfigForm"

var affiliationsToRequest = []mucData.Affiliation{
	&mucData.OwnerAffiliation{},
	&mucData.AdminAffiliation{},
	&mucData.OutcastAffiliation{},
}

func (s *session) GetRoomConfigurationForm(roomID jid.Bare) (<-chan *muc.RoomConfigForm, <-chan error) {
	fc := make(chan *muc.RoomConfigForm)

	req := s.muc.newMUCRoomRequest(roomID, mucRequestGetRoomConfigForm, func(response []byte) error {
		cf := &data.MUCRoomConfiguration{}
		if err := xml.Unmarshal(response, cf); err != nil {
			return err
		}
		rcf := muc.NewRoomConfigForm(cf.Form)
		s.muc.addOccupantsIntoRoomConfigForm(roomID, rcf)
		fc <- rcf
		return nil
	})

	go req.get(data.MUCRoomConfiguration{})

	return fc, req.errorChannel
}

func (m *mucManager) addOccupantsIntoRoomConfigForm(roomID jid.Bare, rcf *muc.RoomConfigForm) {
	for _, a := range affiliationsToRequest {
		rc, err := m.requestRoomOccupantsByAffiliation(roomID, a)
		select {
		case items := <-rc:
			occ, err := parseMUCItemsToRoomOccupantItems(items)
			if err != nil {
				m.log.WithError(err).WithField("affiliation", a).Error("cannot parse MUCItems to OccupantItems")
				continue
			}
			setOccupantList(a, rcf, occ)
		case e := <-err:
			m.log.WithError(e).WithField("affiliation", a).Error("cannot retrieve occupants")
			continue
		}
	}
}

func setOccupantList(a mucData.Affiliation, rcf *muc.RoomConfigForm, occ muc.RoomOccupantItemList) {
	switch {
	case a.IsOwner():
		rcf.SetOwnerList(occ)
	case a.IsAdmin():
		rcf.SetAdminList(occ)
	case a.IsBanned():
		rcf.SetBanList(occ)
	}
}

func parseMUCItemsToRoomOccupantItems(items []xmppData.MUCItem) ([]*muc.RoomOccupantItem, error) {
	list := []*muc.RoomOccupantItem{}
	for _, itm := range items {
		affiliation, err := mucData.AffiliationFromString(itm.Affiliation)
		if err != nil {
			return nil, err
		}

		list = append(list, &muc.RoomOccupantItem{
			Jid:         jid.Parse(itm.Jid),
			Affiliation: affiliation,
			Reason:      itm.Reason,
		})
	}

	return list, nil
}

const mucRequestRoomOccupantsByAffiliation mucRequestType = "requestRoomOccupantsByAffiliation"

func (m *mucManager) requestRoomOccupantsByAffiliation(roomID jid.Bare, a mucData.Affiliation) (<-chan []xmppData.MUCItem, <-chan error) {
	oc := make(chan []xmppData.MUCItem)
	req := m.newMUCRoomRequest(roomID, mucRequestRoomOccupantsByAffiliation, func(response []byte) error {
		var list xmppData.MUCAdmin
		if err := xml.Unmarshal(response, &list); err != nil {
			m.log.WithError(err).Error("failed to unmarshall the response of room's occupants list")
			return err
		}

		oc <- list.Items
		return nil
	})

	go req.get(newRoomOccupantsRequestQueryByAffiliation(a))

	return oc, req.errorChannel
}

func newRoomOccupantsRequestQueryByAffiliation(affiliation mucData.Affiliation) xmppData.MUCAdmin {
	return xmppData.MUCAdmin{
		Items: []xmppData.MUCItem{
			{Affiliation: affiliation.Name()},
		},
	}
}

func (s *session) SubmitRoomConfigurationForm(roomID jid.Bare, form *muc.RoomConfigForm) (<-chan bool, <-chan muc.SubmitFormError) {
	log := log.WithFields(log.Fields{
		"room":  roomID,
		"where": "SubmitRoomConfigurationForm",
	})

	sc := make(chan bool)
	ec := make(chan muc.SubmitFormError)

	go func() {
		reply, _, err := s.conn.SendIQ(roomID.String(), "set", data.MUCRoomConfiguration{
			Form: form.GetFormData(),
		})

		if err != nil {
			log.WithError(err).Error("An error occurred when trying to send the information query to save the room configuration")
			ec <- muc.NewSubmitFormError(ErrRoomConfigSubmit)
			return
		}

		err = validateIqResponse(reply)
		if err != nil {
			log.WithError(ErrInformationQueryResponse).Error("An error occurred when trying to read the response from the room configuration request")
			ec <- muc.NewSubmitFormError(ErrRoomConfigSubmitResponse)
			return
		}

		sc <- true

	}()

	return sc, ec
}

func (s *session) CancelRoomConfiguration(roomID jid.Bare) <-chan error {
	log := log.WithFields(log.Fields{
		"room":  roomID,
		"where": "CancelRoomConfiguration",
	})

	ec := make(chan error)

	go func() {
		reply, _, err := s.conn.SendIQ(roomID.String(), "set", data.MUCRoomConfiguration{
			Form: &data.Form{
				Type: "cancel",
			},
		})

		if err != nil {
			log.WithError(err).Error("An error occurred when trying to send the request to rollback the room configuration")
			ec <- ErrRoomConfigCancel
			return
		}

		err = validateIqResponse(reply)
		if err != nil {
			log.WithError(ErrInformationQueryResponse).Error("An error occurred when trying to read the response from the room configuration rollback request")
			ec <- ErrRoomConfigCancelResponse
			return
		}

		close(ec)
	}()

	return ec
}

func validateIqResponse(reply <-chan data.Stanza) error {
	stanza, ok := <-reply
	if !ok {
		return ErrInformationQueryResponse
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return ErrInformationQueryResponse
	}

	if iq.Type == "error" {
		if iq.Error.MUCConflict != nil {
			return ErrOwnerAffiliationRevokeConflict
		}

		if iq.Error.MUCNotAllowed != nil {
			return ErrNotAllowedKickOccupant
		}

		return ErrUnexpectedResponse
	}

	return nil
}
