package xmpp

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	mucSupport = "<x xmlns='http://jabber.org/protocol/muc'/>"
	mucNS      = "http://jabber.org/protocol/muc"
)

func (c *conn) GetChatContext() interfaces.LegacyOldDoNotUseChat {
	return &muc{
		conn:   c,
		events: make(chan interface{}),
	}
}

type muc struct {
	*conn
	events chan interface{}
}

func (m *muc) Events() chan interface{} {
	return m.events
}

//See: Section "6.2 Discovering the Features Supported by a MUC Service"
func (m *muc) CheckForSupport(entity string) bool {
	return m.HasSupportTo(entity, mucNS)
}

//See: Section "6.3 Discovering Rooms"
func (m *muc) QueryRooms(entity string) ([]data.DiscoveryItem, error) {
	query, err := m.QueryServiceItems(entity)
	if err != nil {
		return nil, err
	}

	return query.DiscoveryItems, nil
}

//See: Section "6.4 Querying for Room Information"
func (m *muc) LegacyOldDoNotUseQueryRoomInformation(room string) (*data.LegacyOldDoNotUseRoomInfo, error) {
	j := jid.Parse(room)
	if j == jid.Domain("") {
		return nil, errors.New("invalid room")
	}

	local := string(jid.MaybeLocal(j))

	//TODO: this error is useless when it says ("expected query, got error")
	//It should give us a xmpp error
	query, err := m.queryRoomInformation(&data.LegacyOldDoNotUseRoom{
		ID:      local,
		Service: string(j.Host()),
	})

	if err != nil {
		return nil, err
	}

	return parseRoomInformation(query), nil
}

func parseRoomInfoForm(forms []data.Form) data.LegacyOldDoNotUseRoomInfoForm {
	ret := data.LegacyOldDoNotUseRoomInfoForm{}
	_ = parseForms(&ret, forms)
	return ret
}

func parseRoomType(features []data.DiscoveryFeature) data.LegacyOldDoNotUseRoomType {
	ret := data.LegacyOldDoNotUseRoomType{}

	for _, f := range features {
		switch f.Var {
		case "muc_public":
			ret.Public = true
		case "muc_open":
			ret.Open = true
		case "muc_moderated":
			ret.Moderated = true
		case "muc_semianonymous":
			ret.SemiAnonymous = true
		case "muc_passwordprotected":
			ret.PasswordProtected = true
		case "muc_persistenc":
			ret.Persistent = true
		}
	}

	return ret
}

func parseRoomInformation(query *data.DiscoveryInfoQuery) *data.LegacyOldDoNotUseRoomInfo {
	return &data.LegacyOldDoNotUseRoomInfo{
		LegacyOldDoNotUseRoomInfoForm: parseRoomInfoForm(query.Forms[:]),
		LegacyOldDoNotUseRoomType:     parseRoomType(query.Features),
	}
}

func (m *muc) queryRoomInformation(room *data.LegacyOldDoNotUseRoom) (*data.DiscoveryInfoQuery, error) {
	return m.QueryServiceInformation(room.JID())
}

//See: Section "7.2.2 Basic MUC Protocol"
func (m *muc) LegacyOldDoNotUseEnterRoom(occupant *data.LegacyOldDoNotUseOccupant) error {
	//TODO: Implement section "7.2.1 Groupchat 1.0 Protocol"?
	return m.sendPresence(&data.ClientPresence{
		To:    occupant.JID(),
		Extra: mucSupport,
	})
}

//See: Section "7.14 Exiting a Room"
func (m *muc) LegacyOldDoNotUseLeaveRoom(occupant *data.LegacyOldDoNotUseOccupant) error {
	return m.sendPresence(&data.ClientPresence{
		To:    occupant.JID(),
		Type:  "unavailable",
		Extra: mucSupport,
	})
}

//See: Section "7.4 Sending a Message to All Occupants"
func (m *muc) LegacyOldDoNotUseSendChatMessage(msg string, to *data.LegacyOldDoNotUseRoom) error {
	//TODO: How to disable archive for chat messages?
	//TODO: Can we just use the same conn.Send() with a different type?
	_, err := fmt.Fprintf(m.out, "<message "+
		"to='%s' "+
		"from='%s' "+
		"type='groupchat'>"+
		"<body>%s</body>"+
		"</message>",
		xmlEscape(to.JID()), xmlEscape(m.conn.jid), xmlEscape(msg))
	return err
}

//See: Section "10.2 Subsequent Room Configuration"
func (m *muc) LegacyOldDoNotUseRequestRoomConfigForm(room *data.LegacyOldDoNotUseRoom) (*data.Form, error) {
	reply, _, err := m.SendIQ(room.JID(), "get", &data.LegacyOldDoNotUseRoomConfigurationQuery{})
	if err != nil {
		return nil, err
	}

	stanza, ok := <-reply
	if !ok {
		return nil, errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, errors.New("xmpp: failed to parse response")
	}

	r := &data.LegacyOldDoNotUseRoomConfigurationQuery{}
	err = xml.Unmarshal(iq.Query, r)
	return r.Form, err
}

func (m *muc) LegacyOldDoNotUseRoomConfigForm(room *data.LegacyOldDoNotUseRoom, formCallback data.FormCallback) error {
	form, err := m.LegacyOldDoNotUseRequestRoomConfigForm(room)
	if err != nil {
		return err
	}

	var datas []data.BobData
	roomConfig, err := processForm(form, datas, formCallback)
	if err != nil {
		return err
	}

	return m.LegacyOldDoNotUseUpdateRoomConfig(room, roomConfig)
}

//See: Section "10.2 Subsequent Room Configuration"
func (m *muc) LegacyOldDoNotUseUpdateRoomConfig(room *data.LegacyOldDoNotUseRoom, form *data.Form) error {
	_, _, err := m.SendIQ(room.JID(), "set", &data.LegacyOldDoNotUseRoomConfigurationQuery{
		Form: form,
	})

	return err
}
