package muc

import (
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"

	. "gopkg.in/check.v1"

	"github.com/coyim/coyim/xmpp/jid"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func (s *MucSuite) Test_ByOccupantNick(c *C) {
	vals := []*Occupant{
		&Occupant{Nick: "Rafael"},
		&Occupant{Nick: "Ola"},
		&Occupant{Nick: "Cristian"},
		&Occupant{Nick: "Ola2"},
		&Occupant{Nick: "Reinaldo"},
	}

	sort.Sort(ByOccupantNick(vals))
	c.Assert(vals[0].Nick, Equals, "Cristian")
	c.Assert(vals[1].Nick, Equals, "Ola")
	c.Assert(vals[2].Nick, Equals, "Ola2")
	c.Assert(vals[3].Nick, Equals, "Rafael")
	c.Assert(vals[4].Nick, Equals, "Reinaldo")
}

func (s *MucSuite) Test_ByOccupantJid(c *C) {
	vals := []*Occupant{
		&Occupant{Nick: "Rafael"},
		&Occupant{Nick: "Ola"},
		&Occupant{Nick: "Cristian"},
		&Occupant{Nick: "Ola2"},
		&Occupant{Nick: "Reinaldo"},
	}

	sort.Sort(ByOccupantJid(vals))
	c.Assert(vals[0].Nick, Equals, "Cristian")
	c.Assert(vals[1].Nick, Equals, "Ola")
	c.Assert(vals[2].Nick, Equals, "Ola2")
	c.Assert(vals[3].Nick, Equals, "Rafael")
	c.Assert(vals[4].Nick, Equals, "Reinaldo")

	vals2 := []*Occupant{
		&Occupant{Nick: "Rafael", Jid: jid.R("foo@jabber.org/abc")},
		&Occupant{Nick: "Rafael2", Jid: jid.R("foo@jabber.im/abc")},
		&Occupant{Nick: "Ola", Jid: nil},
		&Occupant{Nick: "Cristian", Jid: nil},
	}

	sort.Sort(ByOccupantJid(vals2))
	c.Assert(vals2[0].Nick, Equals, "Cristian")
	c.Assert(vals2[1].Nick, Equals, "Ola")
	c.Assert(vals2[2].Nick, Equals, "Rafael2")
	c.Assert(vals2[3].Nick, Equals, "Rafael")
}
