package gui

import (
	"strings"

	"github.com/coyim/coyim/session/muc/data"
	. "gopkg.in/check.v1"
)

type MUCRoomConversationDisplayAffiliationsSuite struct{}

var _ = Suite(&MUCRoomConversationDisplayAffiliationsSuite{})

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationRemoved(c *C) {
	initMUCRoomConversationDisplayI18n()

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
		"The member position of nick was removed.")

	d.nickname = "007"
	d.previousAffiliation = owner
	d.actor = "maria"
	d.actorAffiliation = owner
	c.Assert(d.displayForAffiliationRemoved(), Equals,
		"The owner maria removed the owner position from 007.")
}

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationOutcast(c *C) {
	initMUCRoomConversationDisplayI18n()

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: "nick",
		New:      newAffiliationFromString(data.AffiliationOutcast),
	})

	c.Assert(d.displayForAffiliationOutcast(), Equals,
		"nick was banned from the room.")

	d.nickname = "jonathan"
	d.actor = "maria"
	c.Assert(d.displayForAffiliationOutcast(), Equals,
		"maria banned jonathan from the room.")
}

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationAdded(c *C) {
	initMUCRoomConversationDisplayI18n()

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
		"nick is now a member.")

	d.nickname = "maria"
	d.newAffiliation = admin
	d.actor = "alberto"
	d.actorAffiliation = admin
	c.Assert(d.displayForAffiliationAdded(), Equals,
		"The administrator alberto changed the position of maria to administrator.")

	d.nickname = "alice"
	d.newAffiliation = owner
	d.actor = "bob"
	d.actorAffiliation = owner
	c.Assert(d.displayForAffiliationAdded(), Equals,
		"The owner bob changed the position of alice to owner.")
}

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationChanged(c *C) {
	initMUCRoomConversationDisplayI18n()

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: "nick",
		New:      member,
		Previous: admin,
	})

	c.Assert(d.displayForAffiliationChanged(), Equals,
		"The position of nick was changed from administrator to member.")

	d.nickname = "maria"
	d.newAffiliation = admin
	d.previousAffiliation = member
	d.actor = "juan"
	d.actorAffiliation = member
	c.Assert(d.displayForAffiliationChanged(), Equals,
		"The member juan changed the position of maria from member to administrator.")

	d.nickname = "alice"
	d.newAffiliation = owner
	d.previousAffiliation = member
	d.actor = "bob"
	d.actorAffiliation = member
	c.Assert(d.displayForAffiliationChanged(), Equals,
		"The member bob changed the position of alice from member to owner.")
}

func newAffiliationFromString(s string) data.Affiliation {
	a, err := data.AffiliationFromString(s)
	if err != nil {
		return nil
	}
	return a
}
