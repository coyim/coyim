package gui

import (
	"strings"

	"github.com/coyim/coyim/session/muc/data"
	. "gopkg.in/check.v1"
)

type MUCRoomConversationDisplayAffiliationsSuite struct{}

var _ = Suite(&MUCRoomConversationDisplayAffiliationsSuite{})

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationUpdate(c *C) {
	initMUCRoomConversationDisplayI18n()

	none := newAffiliationFromString(data.AffiliationNone)
	outcast := newAffiliationFromString(data.AffiliationOutcast)
	member := newAffiliationFromString(data.AffiliationMember)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		New:      member,
		Previous: none,
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "nick",
			Actor: &data.OccupantUpdateActor{
				Nickname: "alex",
			},
		},
	})

	c.Assert(displayAffiliationUpdateMessage(d), Equals,
		"alex changed the position of nick to member")

	c.Assert(getDisplayForOccupantAffiliationUpdate(data.AffiliationUpdate{
		New:      none,
		Previous: member,
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "robin",
			Reason:   "I'm batman",
			Actor: &data.OccupantUpdateActor{
				Nickname:    "batman",
				Affiliation: member,
			},
		},
	}), Equals, "The member batman removed the member position from robin because: I'm batman")

	c.Assert(getDisplayForOccupantAffiliationUpdate(data.AffiliationUpdate{
		New:      outcast,
		Previous: member,
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "bob",
			Reason:   "he was rude",
			Actor: &data.OccupantUpdateActor{
				Nickname:    "alice",
				Affiliation: member,
			},
		},
	}), Equals, "The member alice banned bob from the room because: he was rude")

	c.Assert(getDisplayForOccupantAffiliationUpdate(data.AffiliationUpdate{
		New:      none,
		Previous: outcast,
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "nick",
			Actor: &data.OccupantUpdateActor{
				Nickname:    "jonathan",
				Affiliation: outcast,
			},
		},
	}), Equals, "The outcast jonathan removed the outcast position from nick")
}

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationRemoved(c *C) {
	initMUCRoomConversationDisplayI18n()

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)
	none := newAffiliationFromString(data.AffiliationNone)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		New:      none,
		Previous: admin,
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "jonathan",
		},
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
		"The member position of nick was removed")

	d.nickname = "007"
	d.previousAffiliation = owner
	d.actor = "maria"
	d.actorAffiliation = owner
	c.Assert(d.displayForAffiliationRemoved(), Equals,
		"The owner maria removed the owner position from 007")
}

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationOutcast(c *C) {
	initMUCRoomConversationDisplayI18n()

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		New: newAffiliationFromString(data.AffiliationOutcast),
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "nick",
		},
	})

	c.Assert(d.displayForAffiliationOutcast(), Equals,
		"nick was banned from the room")

	d.nickname = "jonathan"
	d.actor = "maria"
	c.Assert(d.displayForAffiliationOutcast(), Equals,
		"maria banned jonathan from the room")
}

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationAdded(c *C) {
	initMUCRoomConversationDisplayI18n()

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)
	none := newAffiliationFromString(data.AffiliationNone)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		New:      member,
		Previous: none,
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "nick",
		},
	})

	c.Assert(d.displayForAffiliationAdded(), Equals,
		"nick is now a member")

	d.nickname = "maria"
	d.newAffiliation = admin
	d.actor = "alberto"
	d.actorAffiliation = admin
	c.Assert(d.displayForAffiliationAdded(), Equals,
		"The administrator alberto changed the position of maria to administrator")

	d.nickname = "alice"
	d.newAffiliation = owner
	d.actor = "bob"
	d.actorAffiliation = owner
	c.Assert(d.displayForAffiliationAdded(), Equals,
		"The owner bob changed the position of alice to owner")
}

func (*MUCRoomConversationDisplayAffiliationsSuite) Test_mucRoomConversationDisplay_displayForAffiliationChanged(c *C) {
	initMUCRoomConversationDisplayI18n()

	member := newAffiliationFromString(data.AffiliationMember)
	admin := newAffiliationFromString(data.AffiliationAdmin)
	owner := newAffiliationFromString(data.AffiliationOwner)

	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		New:      member,
		Previous: admin,
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Nickname: "nick",
		},
	})

	c.Assert(d.displayForAffiliationChanged(), Equals,
		"The position of nick was changed from administrator to member")

	d.nickname = "maria"
	d.newAffiliation = admin
	d.previousAffiliation = member
	d.actor = "juan"
	d.actorAffiliation = member
	c.Assert(d.displayForAffiliationChanged(), Equals,
		"The member juan changed the position of maria from member to administrator")

	d.nickname = "alice"
	d.newAffiliation = owner
	d.previousAffiliation = member
	d.actor = "bob"
	d.actorAffiliation = member
	c.Assert(d.displayForAffiliationChanged(), Equals,
		"The member bob changed the position of alice from member to owner")
}

func newAffiliationFromString(s string) data.Affiliation {
	a, err := data.AffiliationFromString(s)
	if err != nil {
		return nil
	}
	return a
}
