package gui

import (
	"fmt"
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"

	"github.com/coyim/coyim/session/muc/data"
	. "gopkg.in/check.v1"
)

type MUCRoomConversationDisplaySuite struct{}

var _ = Suite(&JoinMUCRoomSuite{})

type mucRoomConversationDisplayMockGlib struct {
	glib_mock.Mock
}

func (*mucRoomConversationDisplayMockGlib) Local(vx string) string {
	return "[localized] " + vx
}

func (*mucRoomConversationDisplayMockGlib) Localf(vx string, args ...interface{}) string {
	return fmt.Sprintf("[localized] "+vx, args...)
}

func (*SignalsSuite) Test_mucRoomConversationDisplay_displayForAffiliationUpdate(c *C) {
	i18n.InitLocalization(&mucRoomConversationDisplayMockGlib{})

	none := newAffiliationFromString(data.AffiliationNone)
	outcast := newAffiliationFromString(data.AffiliationOutcast)
	member := newAffiliationFromString(data.AffiliationMember)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: "nick",
		New:      member,
		Previous: none,
		Actor:    "alex",
	})

	c.Assert(displayAffiliationUpdateMessage(d), Equals,
		"[localized] alex changed the position of nick to [localized] member")

	c.Assert(getDisplayForOccupantAffiliationUpdate(data.AffiliationUpdate{
		Nickname: "robin",
		New:      none,
		Previous: member,
		Actor:    "batman",
		Reason:   "I'm batman",
	}), Equals, "[localized] batman removed the [localized] member position from robin[localized]  because I'm batman")

	c.Assert(getDisplayForOccupantAffiliationUpdate(data.AffiliationUpdate{
		Nickname: "bob",
		New:      outcast,
		Previous: member,
		Actor:    "alice",
		Reason:   "he was rude",
	}), Equals, "[localized] alice banned bob from the room[localized]  because he was rude")

	c.Assert(getDisplayForOccupantAffiliationUpdate(data.AffiliationUpdate{
		Nickname: "nick",
		New:      none,
		Previous: outcast,
		Actor:    "jonathan",
	}), Equals, "[localized] jonathan removed the [localized] outcast position from nick")
}

func (*SignalsSuite) Test_mucRoomConversationDisplay_displayForAffiliationRemoved(c *C) {
	i18n.InitLocalization(&mucRoomConversationDisplayMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)
	none := newAffiliationFromString(data.AffiliationNone)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: "jonathan",
		New:      none,
		Previous: admin,
	})

	c.Assert(strings.Contains(d.displayForAffiliationRemoved(), "nick"), Equals, false)

	d.nickname = "alice"
	d.actor = ""
	c.Assert(strings.Contains(d.displayForAffiliationRemoved(), "administrator"), Equals, true)

	d.nickname = "alberto"
	d.actor = "me"
	c.Assert(strings.Contains(d.displayForAffiliationRemoved(), "me"), Equals, true)

	d.nickname = "nick"
	d.previousAffiliation = member
	d.actor = ""
	c.Assert(d.displayForAffiliationRemoved(), Equals,
		"[localized] The [localized] member position of nick was removed")

	d.nickname = "007"
	d.previousAffiliation = owner
	d.actor = "maria"
	c.Assert(d.displayForAffiliationRemoved(), Equals,
		"[localized] maria removed the [localized] owner position from 007")
}

func (*SignalsSuite) Test_mucRoomConversationDisplay_displayForAffiliationOutcast(c *C) {
	i18n.InitLocalization(&mucRoomConversationDisplayMockGlib{})

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: "nick",
		New:      newAffiliationFromString(data.AffiliationOutcast),
	})

	c.Assert(d.displayForAffiliationOutcast(), Equals,
		"[localized] nick was banned from the room")

	d.nickname = "jonathan"
	d.actor = "maria"
	c.Assert(d.displayForAffiliationOutcast(), Equals,
		"[localized] maria banned jonathan from the room")
}

func (*SignalsSuite) Test_mucRoomConversationDisplay_displayForAffiliationAdded(c *C) {
	i18n.InitLocalization(&mucRoomConversationDisplayMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)
	none := newAffiliationFromString(data.AffiliationNone)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: "nick",
		New:      member,
		Previous: none,
	})

	c.Assert(d.displayForAffiliationAdded(), Equals,
		"[localized] nick is now [localized] a member")

	d.nickname = "maria"
	d.newAffiliation = admin
	d.actor = "alberto"
	c.Assert(d.displayForAffiliationAdded(), Equals,
		"[localized] alberto changed the position of maria to [localized] administrator")

	d.nickname = "alice"
	d.newAffiliation = owner
	d.actor = "bob"
	c.Assert(d.displayForAffiliationAdded(), Equals,
		"[localized] bob changed the position of alice to [localized] owner")
}

func (*SignalsSuite) Test_mucRoomConversationDisplay_displayForAffiliationChanged(c *C) {
	i18n.InitLocalization(&mucRoomConversationDisplayMockGlib{})

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: "nick",
		New:      member,
		Previous: admin,
	})

	c.Assert(d.displayForAffiliationChanged(), Equals,
		"[localized] The position of nick was changed from [localized] administrator to [localized] member")

	d.nickname = "maria"
	d.newAffiliation = admin
	d.previousAffiliation = member
	d.actor = "juan"
	c.Assert(d.displayForAffiliationChanged(), Equals,
		"[localized] juan changed the position of maria from [localized] member to [localized] administrator")

	d.nickname = "alice"
	d.newAffiliation = owner
	d.previousAffiliation = member
	d.actor = "bob"
	c.Assert(d.displayForAffiliationChanged(), Equals,
		"[localized] bob changed the position of alice from [localized] member to [localized] owner")
}

func newAffiliationFromString(s string) data.Affiliation {
	a, err := data.AffiliationFromString(s)
	if err != nil {
		return nil
	}
	return a
}
