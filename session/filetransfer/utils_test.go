package filetransfer

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/mock"
	. "gopkg.in/check.v1"
)

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

type sendIQMock struct {
	*mock.Conn

	sendIQ func(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error)
}

type hasConn struct {
	c xi.Conn
}

func (h *hasConn) Conn() xi.Conn {
	return h.c
}

func (m *sendIQMock) SendIQ(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
	return m.sendIQ(v1, v2, v3)
}

func (s *UtilsSuite) Test_basicIQ_failsWhenSendingIQ(c *C) {
	m := &sendIQMock{
		sendIQ: func(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error) {
			return nil, 0, errors.New("marker interface")
		},
	}

	e := basicIQ(&hasConn{m}, "foo@bar", "result", nil, nil, func(*data.ClientIQ) {})

	c.Assert(e, ErrorMatches, "marker interface")
}

func (s *UtilsSuite) Test_basicIQ_failsOnClosedReturnChannel(c *C) {
	m := &sendIQMock{
		sendIQ: func(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error) {
			ch := make(chan data.Stanza)
			close(ch)
			return ch, 0, nil
		},
	}

	e := basicIQ(&hasConn{m}, "foo@bar", "result", nil, nil, func(*data.ClientIQ) {})

	c.Assert(e, Equals, errChannelClosed)
}

func (s *UtilsSuite) Test_basicIQ_succeeds(c *C) {
	m := &sendIQMock{
		sendIQ: func(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error) {
			ch := make(chan data.Stanza, 1)
			ch <- data.Stanza{
				Value: &data.ClientIQ{
					Type: "result",
				},
			}
			return ch, 0, nil
		},
	}

	successCalled := false
	e := basicIQ(&hasConn{m}, "foo@bar", "result", nil, nil, func(*data.ClientIQ) {
		successCalled = true
	})

	c.Assert(e, IsNil)
	c.Assert(successCalled, Equals, true)
}

func (s *UtilsSuite) Test_basicIQ_succeedsAndUnpacks(c *C) {
	m := &sendIQMock{
		sendIQ: func(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error) {
			ch := make(chan data.Stanza, 1)
			ch <- data.Stanza{
				Value: &data.ClientIQ{
					Type:  "result",
					Query: []byte(`<delay xmlns="urn:xmpp:delay" from="hello"></delay>`),
				},
			}
			return ch, 0, nil
		},
	}

	var unpackData data.Delay

	successCalled := false
	e := basicIQ(&hasConn{m}, "foo@bar", "result", nil, &unpackData, func(*data.ClientIQ) {
		successCalled = true
	})

	c.Assert(e, IsNil)
	c.Assert(successCalled, Equals, true)
	c.Assert(unpackData.From, Equals, "hello")
}
