package gui

import (
	"fmt"
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

func (*mucRoomConversationMockGlib) Localf(vx string, args ...interface{}) string {
	return fmt.Sprintf("[localized] "+vx, args...)
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationUpdate(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	none := newAffiliationFromString(data.AffiliationNone)
	outcast := newAffiliationFromString(data.AffiliationOutcast)
	member := newAffiliationFromString(data.AffiliationMember)

	affiliationUpdate := data.AffiliationUpdate{
		New:      member,
		Previous: none,
	}
	c.Assert(getDisplayForOccupantAffiliationUpdate("nick", affiliationUpdate, "alex", ""), Equals,
		"[localized] alex updated the position of nick to [localized] member")

	affiliationUpdate = data.AffiliationUpdate{
		New:      none,
		Previous: member,
	}
	c.Assert(getDisplayForOccupantAffiliationUpdate("robin", affiliationUpdate, "batman", "I'm batman"), Equals,
		"[localized] batman removed the [localized] member position of robin[localized]  because I'm batman")

	affiliationUpdate = data.AffiliationUpdate{
		New:      outcast,
		Previous: member,
	}
	c.Assert(getDisplayForOccupantAffiliationUpdate("bob", affiliationUpdate, "alice", "he was rude"), Equals,
		"[localized] alice has banned bob in the room[localized]  because he was rude")

	affiliationUpdate = data.AffiliationUpdate{
		New:      none,
		Previous: outcast,
	}
	c.Assert(getDisplayForOccupantAffiliationUpdate("nick", affiliationUpdate, "jonathan", ""), Equals,
		"[localized] jonathan removed the [localized] outcast position of nick")
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationRemoved(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	c.Assert(strings.Contains(getDisplayForOccupantAffiliationRemoved("jonathan", admin, ""), "nick"), Equals, false)
	c.Assert(strings.Contains(getDisplayForOccupantAffiliationRemoved("alice", admin, ""), "admin"), Equals, true)
	c.Assert(strings.Contains(getDisplayForOccupantAffiliationRemoved("alberto", admin, "me"), "me"), Equals, true)

	c.Assert(getDisplayForOccupantAffiliationRemoved("nick", member, ""), Equals,
		"[localized] The [localized] member position of nick was removed")

	c.Assert(getDisplayForOccupantAffiliationRemoved("007", owner, "maria"), Equals,
		"[localized] maria removed the [localized] owner position of 007")
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationOutcast(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	c.Assert(getDisplayForOccupantAffiliationOutcast("nick", ""), Equals,
		"[localized] nick has been banned in the room")

	c.Assert(getDisplayForOccupantAffiliationOutcast("jonathan", "maria"), Equals,
		"[localized] maria has banned jonathan in the room")
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationAdded(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	c.Assert(getDisplayForOccupantAffiliationAdded("nick", member, ""), Equals,
		"[localized] The position of nick was updated to [localized] member")

	c.Assert(getDisplayForOccupantAffiliationAdded("maria", admin, "alberto"), Equals,
		"[localized] alberto updated the position of maria to [localized] admin")

	c.Assert(getDisplayForOccupantAffiliationAdded("alice", owner, "bob"), Equals,
		"[localized] bob updated the position of alice to [localized] owner")
}

func (*SignalsSuite) Test_mucRoomConversation_getDisplayForOccupantAffiliationChanged(c *C) {
	i18n.InitLocalization(&mucRoomConversationMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	c.Assert(getDisplayForOccupantAffiliationChanged("nick", member, admin, ""), Equals,
		"[localized] The position of nick was updated from [localized] admin to [localized] member")

	c.Assert(getDisplayForOccupantAffiliationChanged("maria", admin, member, "juan"), Equals,
		"[localized] juan updated the position of maria from [localized] member to [localized] admin")

	c.Assert(getDisplayForOccupantAffiliationChanged("alice", owner, member, "bob"), Equals,
		"[localized] bob updated the position of alice from [localized] member to [localized] owner")
}

func newAffiliationFromString(s string) data.Affiliation {
	a, err := data.AffiliationFromString(s)
	if err != nil {
		return nil
	}
	return a
}
