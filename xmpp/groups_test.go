package xmpp

import (
	"bytes"
	"errors"
	"time"

	"github.com/coyim/coyim/xmpp/data"

	. "gopkg.in/check.v1"
)

type GroupsSuite struct{}

var _ = Suite(&GroupsSuite{})

func (s *GroupsSuite) Test_conn_RequestRosterDelimiter_works(c *C) {
	var out bytes.Buffer

	cn := &conn{
		out:       &out,
		inflights: make(map[data.Cookie]inflight),
	}

	ch, _, e := cn.RequestRosterDelimiter()
	c.Assert(e, IsNil)
	c.Assert(out.String(), Matches, ""+
		"\n"+
		"<iq type='get' id='.*?'>\n"+
		"  <query xmlns='jabber:iq:private'>\n"+
		"    <roster xmlns='roster:delimiter'/>\n"+
		"  </query>\n"+
		"</iq>\n")
	c.Assert(ch, Not(IsNil))
}

func (s *GroupsSuite) Test_conn_GetRosterDelimiter_failsIfRequestFails(c *C) {
	out := &mockConnIOReaderWriter{}
	out.err = errors.New("marker error")

	cn := &conn{
		out:       out,
		inflights: make(map[data.Cookie]inflight),
	}

	res, e := cn.GetRosterDelimiter()
	c.Assert(e, ErrorMatches, "marker error")
	c.Assert(res, Equals, "")
}

func (s *GroupsSuite) Test_conn_GetRosterDelimiter_works(c *C) {
	orgCreateInflight := createInflight
	defer func() {
		createInflight = orgCreateInflight
	}()

	mch := make(chan data.Stanza, 1)

	createInflight = func(c *conn, cookie data.Cookie, to string) (<-chan data.Stanza, data.Cookie, error) {
		return mch, 0, nil
	}

	out := &mockConnIOReaderWriter{}

	cc := &conn{
		out: out,
	}

	mch <- data.Stanza{
		Value: &data.ClientIQ{
			Query: []byte(`<query xmlns="jabber:iq:private">
  <roster xmlns="roster:delimiter">foobarium</roster>
</query>`),
		},
	}

	res, e := cc.GetRosterDelimiter()
	c.Assert(e, IsNil)
	c.Assert(res, Equals, "foobarium")
}

func (s *GroupsSuite) Test_conn_GetRosterDelimiter_failsIfNotIQ(c *C) {
	orgCreateInflight := createInflight
	defer func() {
		createInflight = orgCreateInflight
	}()

	mch := make(chan data.Stanza, 1)

	createInflight = func(c *conn, cookie data.Cookie, to string) (<-chan data.Stanza, data.Cookie, error) {
		return mch, 0, nil
	}

	out := &mockConnIOReaderWriter{}

	cc := &conn{
		out: out,
	}

	mch <- data.Stanza{
		Value: "foo",
	}

	res, e := cc.GetRosterDelimiter()
	c.Assert(e, IsNil)
	c.Assert(res, Equals, "")
}

func (s *GroupsSuite) Test_conn_GetRosterDelimiter_failsOnBadXML(c *C) {
	orgCreateInflight := createInflight
	defer func() {
		createInflight = orgCreateInflight
	}()

	mch := make(chan data.Stanza, 1)

	createInflight = func(c *conn, cookie data.Cookie, to string) (<-chan data.Stanza, data.Cookie, error) {
		return mch, 0, nil
	}

	out := &mockConnIOReaderWriter{}

	cc := &conn{
		out: out,
	}

	mch <- data.Stanza{
		Value: &data.ClientIQ{
			Query: []byte(`<query `),
		},
	}

	res, e := cc.GetRosterDelimiter()
	c.Assert(e, IsNil)
	c.Assert(res, Equals, "")
}

func (s *GroupsSuite) Test_conn_GetRosterDelimiter_timesOut(c *C) {
	orgRosterRequestTimeout := rosterRequestTimeout
	defer func() {
		rosterRequestTimeout = orgRosterRequestTimeout
	}()

	orgCreateInflight := createInflight
	defer func() {
		createInflight = orgCreateInflight
	}()

	mch := make(chan data.Stanza)

	createInflight = func(c *conn, cookie data.Cookie, to string) (<-chan data.Stanza, data.Cookie, error) {
		return mch, 0, nil
	}

	rosterRequestTimeout = 1 * time.Millisecond

	out := &mockConnIOReaderWriter{}

	cc := &conn{
		out: out,
	}

	res, e := cc.GetRosterDelimiter()
	c.Assert(e, IsNil)
	c.Assert(res, Equals, "")
}
