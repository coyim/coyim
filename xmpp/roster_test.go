package xmpp

import (
	"encoding/xml"

	"github.com/coyim/coyim/xmpp/data"

	. "gopkg.in/check.v1"
)

type RosterXMPPSuite struct{}

var _ = Suite(&RosterXMPPSuite{})

type testStanzaValue struct{}

func (s *RosterXMPPSuite) Test_ParseRoster_parsesRosterStanza(c *C) {
	rep := data.Stanza{
		Name: xml.Name{Local: "iq"},
		Value: &data.ClientIQ{
			Query: []byte("<query xmlns=\"jabber:iq:roster\"><item jid=\"alice@coy.im\" subscription=\"both\" name=\"Alice\"/><item jid=\"alice@coy.im/coyim\" subscription=\"both\" name=\"Alice using CoyIM\"/></query>"),
		},
	}

	roster, _ := data.ParseRoster(rep)

	c.Assert(len(roster), Equals, 2)
	c.Assert(roster[0].Jid, Equals, "alice@coy.im")
	c.Assert(roster[1].Jid, Equals, "alice@coy.im/coyim")
}

func (s *RosterXMPPSuite) Test_ParseRoster_failsIfItDoesntReceiveAClientIQ(c *C) {
	rep := data.Stanza{
		Name:  xml.Name{Local: "Foobarium"},
		Value: testStanzaValue{},
	}

	_, err := data.ParseRoster(rep)
	c.Assert(err.Error(), Equals, "xmpp: roster request resulted in tag of type Foobarium")
}

func (s *RosterXMPPSuite) Test_ParseRoster_failsIfTheRosterContentIsIncorrect(c *C) {
	rep := data.Stanza{
		Name: xml.Name{Local: "iq"},
		Value: &data.ClientIQ{
			Query: []byte("<foo></bar>"),
		},
	}

	_, err := data.ParseRoster(rep)
	c.Assert(err.Error(), Equals, "expected element type <query> but have <foo>")
}
