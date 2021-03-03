package xmpp

import (
	"bytes"

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
