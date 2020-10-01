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
		&Occupant{Nickname: "Rafael"},
		&Occupant{Nickname: "Ola"},
		&Occupant{Nickname: "Cristian"},
		&Occupant{Nickname: "Ola2"},
		&Occupant{Nickname: "Reinaldo"},
	}

	sort.Sort(ByOccupantNick(vals))
	c.Assert(vals[0].Nickname, Equals, "Cristian")
	c.Assert(vals[1].Nickname, Equals, "Ola")
	c.Assert(vals[2].Nickname, Equals, "Ola2")
	c.Assert(vals[3].Nickname, Equals, "Rafael")
	c.Assert(vals[4].Nickname, Equals, "Reinaldo")
}

func (s *MucSuite) Test_ByOccupantJid(c *C) {
	vals := []*Occupant{
		&Occupant{Nickname: "Rafael"},
		&Occupant{Nickname: "Ola"},
		&Occupant{Nickname: "Cristian"},
		&Occupant{Nickname: "Ola2"},
		&Occupant{Nickname: "Reinaldo"},
	}

	sort.Sort(ByOccupantJid(vals))
	c.Assert(vals[0].Nickname, Equals, "Cristian")
	c.Assert(vals[1].Nickname, Equals, "Ola")
	c.Assert(vals[2].Nickname, Equals, "Ola2")
	c.Assert(vals[3].Nickname, Equals, "Rafael")
	c.Assert(vals[4].Nickname, Equals, "Reinaldo")

	vals2 := []*Occupant{
		&Occupant{Nickname: "Rafael", RealJid: jid.ParseFull("foo@jabber.org/abc")},
		&Occupant{Nickname: "Rafael2", RealJid: jid.ParseFull("foo@jabber.im/abc")},
		&Occupant{Nickname: "Ola", RealJid: nil},
		&Occupant{Nickname: "Cristian", RealJid: nil},
	}

	sort.Sort(ByOccupantJid(vals2))
	c.Assert(vals2[0].Nickname, Equals, "Cristian")
	c.Assert(vals2[1].Nickname, Equals, "Ola")
	c.Assert(vals2[2].Nickname, Equals, "Rafael2")
	c.Assert(vals2[3].Nickname, Equals, "Rafael")
}
