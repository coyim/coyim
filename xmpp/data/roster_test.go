package data

import (
	. "gopkg.in/check.v1"
)

type RosterSuite struct{}

var _ = Suite(&RosterSuite{})

func (s *RosterSuite) Test_ParseRoster_failsIfStanzaNotIQ(c *C) {
	st := Stanza{
		Value: "foo",
	}
	res, e := ParseRoster(st)
	c.Assert(e, ErrorMatches, "xmpp: roster request resulted in tag of type.*")
	c.Assert(res, IsNil)
}

func (s *RosterSuite) Test_ParseRoster_failsOnBadXML(c *C) {
	st := Stanza{
		Value: &ClientIQ{
			Query: []byte("<badXML"),
		},
	}
	res, e := ParseRoster(st)
	c.Assert(e, ErrorMatches, "XML syntax error on.*")
	c.Assert(res, IsNil)
}

func (s *RosterSuite) Test_ParseRoster_succeedsInParsingRoster(c *C) {
	rosterXML := `
<query xmlns='jabber:iq:roster'>
  <item jid='nurse@example.com'/>
  <item jid='romeo@example.net'/>
  <item jid='foo@somewhere.com'/>
  <item jid='abc@example.org'/>
</query>
`

	st := Stanza{
		Value: &ClientIQ{
			Query: []byte(rosterXML),
		},
	}
	res, e := ParseRoster(st)
	c.Assert(e, IsNil)
	c.Assert(res, HasLen, 4)
	c.Assert(res[0].Jid, Equals, "abc@example.org")
	c.Assert(res[1].Jid, Equals, "foo@somewhere.com")
	c.Assert(res[2].Jid, Equals, "nurse@example.com")
	c.Assert(res[3].Jid, Equals, "romeo@example.net")
}
