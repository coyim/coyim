package xmpp

import (
	"encoding/xml"

	. "gopkg.in/check.v1"
)

type RosterXmppSuite struct{}

var _ = Suite(&RosterXmppSuite{})

type testStanzaValue struct{}

func (s *RosterXmppSuite) Test_ParseRoster_failsIfItDoesntReceiveAClientIQ(c *C) {
	rep := Stanza{
		Name:  xml.Name{Local: "Foobarium"},
		Value: testStanzaValue{},
	}

	_, err := ParseRoster(rep)
	c.Assert(err.Error(), Equals, "xmpp: roster request resulted in tag of type Foobarium")
}

func (s *RosterXmppSuite) Test_ParseRoster_failsIfTheRosterContentIsIncorrect(c *C) {
	rep := Stanza{
		Name: xml.Name{Local: "iq"},
		Value: &ClientIQ{
			Query: []byte("<foo></bar>"),
		},
	}

	_, err := ParseRoster(rep)
	c.Assert(err.Error(), Equals, "expected element type <query> but have <foo>")
}
