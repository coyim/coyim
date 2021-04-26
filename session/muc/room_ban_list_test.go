package muc

import (
	"io/ioutil"

	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type MUCRoomBanListSuite struct{}

var _ = Suite(&MUCRoomBanListSuite{})

func (*MUCRoomBanListSuite) Test_RoomBanListItem_IsDifferentFrom(c *C) {
	itm := &RoomBanListItem{jid.Parse("bla"), &data.AdminAffiliation{}, "foo"}

	c.Assert(itm.IsDifferentFrom(&RoomBanListItem{jid.Parse("bla"), &data.AdminAffiliation{}, "foo"}), Equals, false)
	c.Assert(itm.IsDifferentFrom(&RoomBanListItem{jid.Parse("bla2"), &data.AdminAffiliation{}, "foo"}), Equals, true)
	c.Assert(itm.IsDifferentFrom(&RoomBanListItem{jid.Parse("bla"), &data.MemberAffiliation{}, "foo"}), Equals, true)
	c.Assert(itm.IsDifferentFrom(&RoomBanListItem{jid.Parse("bla"), &data.AdminAffiliation{}, "foo2"}), Equals, true)
}

func (*MUCRoomBanListSuite) Test_RoomBanList_Includes(c *C) {
	list := &RoomBanList{
		&RoomBanListItem{jid.Parse("bla"), &data.AdminAffiliation{}, "reason"},
		&RoomBanListItem{jid.Parse("foo"), &data.AdminAffiliation{}, ""},
		&RoomBanListItem{jid.Parse("org"), &data.AdminAffiliation{}, ""},
		&RoomBanListItem{jid.Parse("example.org"), &data.AdminAffiliation{}, ""},
		&RoomBanListItem{jid.Parse("id@domain.org"), &data.AdminAffiliation{}, ""},
	}

	c.Assert(list.IncludesJid("bla2"), Equals, false)
	c.Assert(list.IncludesJid("bla"), Equals, true)
	c.Assert(list.IncludesJid("org"), Equals, true)
	c.Assert(list.IncludesJid("example.org"), Equals, true)
	c.Assert(list.IncludesJid("id@domain.org"), Equals, true)
}
