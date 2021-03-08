package session

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/mock"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type IQSuite struct{}

var _ = Suite(&IQSuite{})

type sendIQReplyXMPPConnMock struct {
	*mock.Conn

	sendIQReply func(string, string, string, interface{}) error
}

func (m *sendIQReplyXMPPConnMock) SendIQReply(v1, v2, v3 string, v4 interface{}) error {
	return m.sendIQReply(v1, v2, v3, v4)
}

func (s *IQSuite) Test_session_SendIQError_works(c *C) {
	conn := &sendIQReplyXMPPConnMock{}

	var v1s []string
	var v2s []string
	var v3s []string
	var v4s []interface{}

	conn.sendIQReply = func(v1, v2, v3 string, v4 interface{}) error {
		v1s = append(v1s, v1)
		v2s = append(v2s, v2)
		v3s = append(v3s, v3)
		v4s = append(v4s, v4)
		return nil
	}
	sess := &session{conn: conn}

	sess.SendIQError(&data.ClientIQ{
		From: "some@one.org",
		ID:   "1234",
	}, "something")

	c.Assert(v1s, DeepEquals, []string{"some@one.org"})
	c.Assert(v2s, DeepEquals, []string{"error"})
	c.Assert(v3s, DeepEquals, []string{"1234"})
	c.Assert(v4s, DeepEquals, []interface{}{"something"})
}

func (s *IQSuite) Test_session_SendIQResult_works(c *C) {
	conn := &sendIQReplyXMPPConnMock{}

	var v1s []string
	var v2s []string
	var v3s []string
	var v4s []interface{}

	conn.sendIQReply = func(v1, v2, v3 string, v4 interface{}) error {
		v1s = append(v1s, v1)
		v2s = append(v2s, v2)
		v3s = append(v3s, v3)
		v4s = append(v4s, v4)
		return nil
	}
	sess := &session{conn: conn}

	sess.SendIQResult(&data.ClientIQ{
		From: "some@one.org",
		ID:   "1234",
	}, "something")

	c.Assert(v1s, DeepEquals, []string{"some@one.org"})
	c.Assert(v2s, DeepEquals, []string{"result"})
	c.Assert(v3s, DeepEquals, []string{"1234"})
	c.Assert(v4s, DeepEquals, []interface{}{"something"})
}

func (s *IQSuite) Test_session_SendIQResult_logsOnFailure(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	conn := &sendIQReplyXMPPConnMock{}

	conn.sendIQReply = func(string, string, string, interface{}) error {
		return errors.New("marker error")
	}
	sess := &session{conn: conn, log: l}

	sess.SendIQResult(&data.ClientIQ{
		From: "some@one.org",
		ID:   "1234",
	}, "something")

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Failed to send IQ message")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "marker error")
}

func (s *IQSuite) Test_session_processIQ_fails(c *C) {
	sess := &session{}

	ret, iqtype, ignore := sess.processIQ(&data.ClientIQ{Query: []byte("")})

	c.Assert(ret, IsNil)
	c.Assert(iqtype, Equals, "")
	c.Assert(ignore, Equals, false)
}
