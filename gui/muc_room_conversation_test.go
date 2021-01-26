package gui

import (
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"

	"github.com/coyim/coyim/session/muc/data"
	. "gopkg.in/check.v1"
)

type MUCRoomConversationSuite struct{}

var _ = Suite(&JoinMUCRoomSuite{})

type mucRoomConversationMockGlib struct {
	glib_mock.Mock
}

func (*mucRoomConversationMockGlib) Local(vx string) string {
	return "[localized] " + vx
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationUpdate(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	none := newAffiliationFromString(data.AffiliationNone)
	outcast := newAffiliationFromString(data.AffiliationOutcast)
	member := newAffiliationFromString(data.AffiliationMember)

	outcasti18n := displayNameForAffiliation(outcast)
	memberi18n := displayNameForAffiliation(member)

	c.Assert(getDisplayForOccupantAffiliationUpdate("nick", member, none, "me", ""), Equals,
		i18n.Localf("%s updated the position of %s to %s", "me", "nick", memberi18n))

	c.Assert(getDisplayForOccupantAffiliationUpdate("nick", none, member, "me", "I wanted"), Equals,
		i18n.Localf("%s removed the %s position of %s%s", "me", memberi18n, "nick", i18n.Localf(" because %s", "I wanted")))

	c.Assert(getDisplayForOccupantAffiliationUpdate("nick", outcast, member, "me", "I wanted"), Equals,
		i18n.Localf("%s has banned %s in the room%s", "me", "nick", i18n.Localf(" because %s", "I wanted")))

	c.Assert(getDisplayForOccupantAffiliationUpdate("nick", none, outcast, "me", ""), Equals,
		i18n.Localf("%s removed the %s position of %s", "me", outcasti18n, "nick"))
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationRemoved(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	memberi18n := displayNameForAffiliation(member)
	admini18n := displayNameForAffiliation(admin)
	owneri18n := displayNameForAffiliation(owner)

	c.Assert(strings.Contains(getDisplayForOccupantAffiliationRemoved("jonathan", admin, ""), "nick"), Equals, false)
	c.Assert(strings.Contains(getDisplayForOccupantAffiliationRemoved("alice", admin, ""), admini18n), Equals, true)
	c.Assert(strings.Contains(getDisplayForOccupantAffiliationRemoved("alberto", admin, "me"), "me"), Equals, true)

	c.Assert(getDisplayForOccupantAffiliationRemoved("nick", member, ""), Equals,
		i18n.Localf("The %s position of %s was removed", memberi18n, "nick"))

	c.Assert(getDisplayForOccupantAffiliationRemoved("007", owner, "maria"), Equals,
		i18n.Localf("%s removed the %s position of %s", "maria", owneri18n, "007"))
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationOutcast(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	c.Assert(getDisplayForOccupantAffiliationOutcast("nick", ""), Equals,
		i18n.Localf("%s has been banned in the room", "nick"))

	c.Assert(getDisplayForOccupantAffiliationOutcast("jonathan", "maria"), Equals,
		i18n.Localf("%s has banned %s in the room", "maria", "jonathan"))
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationAdded(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	memberi18n := displayNameForAffiliation(member)
	admini18n := displayNameForAffiliation(admin)
	owneri18n := displayNameForAffiliation(owner)

	c.Assert(getDisplayForOccupantAffiliationAdded("nick", member, ""), Equals,
		i18n.Localf("The position of %s was updated to %s", "nick", memberi18n))

	c.Assert(getDisplayForOccupantAffiliationAdded("maria", admin, "me"), Equals,
		i18n.Localf("%s updated the position of %s to %s", "me", "maria", admini18n))

	c.Assert(getDisplayForOccupantAffiliationAdded("alice", owner, "bob"), Equals,
		i18n.Localf("%s updated the position of %s to %s", "bob", "alice", owneri18n))
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationChanged(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	memberi18n := displayNameForAffiliation(member)
	admini18n := displayNameForAffiliation(admin)
	owneri18n := displayNameForAffiliation(owner)

	c.Assert(getDisplayForOccupantAffiliationChanged("nick", member, admin, ""), Equals,
		i18n.Localf("The position of %s was updated from %s to %s", "nick", admini18n, memberi18n))

	c.Assert(getDisplayForOccupantAffiliationChanged("maria", admin, member, "me"), Equals,
		i18n.Localf("%s updated the position of %s from %s to %s", "me", "maria", memberi18n, admini18n))

	c.Assert(getDisplayForOccupantAffiliationChanged("alice", owner, member, "bob"), Equals,
		i18n.Localf("%s updated the position of %s from %s to %s", "bob", "alice", memberi18n, owneri18n))
}

func newAffiliationFromString(s string) data.Affiliation {
	a, err := data.AffiliationFromString(s)
	if err != nil {
		return newAffiliationFromString(data.AffiliationNone)
	}
	return a
}
