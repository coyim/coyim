package session

import (
	"encoding/xml"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/mock"
	. "gopkg.in/check.v1"
)

type RosterSuite struct{}

var _ = Suite(&RosterSuite{})

type sendIQXmppConnMock struct {
	*mock.Conn

	sendIQ func(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error)
}

func (m *sendIQXmppConnMock) SendIQ(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
	return m.sendIQ(v1, v2, v3)
}

func (s *RosterSuite) Test_session_RemoveContact(c *C) {
	conn := &sendIQXmppConnMock{}

	var v1s []string
	var v2s []string
	var v3s []interface{}

	conn.sendIQ = func(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
		v1s = append(v1s, v1)
		v2s = append(v2s, v2)
		v3s = append(v3s, v3)
		return nil, 0, nil
	}

	sess := &session{conn: conn}

	sess.RemoveContact("someone@somehere.com")

	c.Assert(v1s, DeepEquals, []string{""})
	c.Assert(v2s, DeepEquals, []string{"set"})
	c.Assert(v3s, DeepEquals, []interface{}{
		data.RosterRequest{XMLName: xml.Name{}, Item: data.RosterRequestItem{Jid: "someone@somehere.com", Subscription: "remove"}},
	})
}
