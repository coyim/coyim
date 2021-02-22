package gui

import (
	"github.com/coyim/coyim/session/muc/data"
	. "gopkg.in/check.v1"
)

type MUCNotificationMessagesSuite struct{}

var _ = Suite(&MUCNotificationMessagesSuite{})

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationNone(c *C) {
	initMUCRoomConversationDisplayI18n()

	au := data.AffiliationUpdate{
		Nickname: "batman",
		New:      newAffiliationFromString(data.AffiliationNone),
		Previous: newAffiliationFromString(data.AffiliationAdmin),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not an administrator anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not an administrator anymore because: batman lost his mind.")

	au.Previous = newAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not an owner anymore because: batman lost his mind.")

	au.Previous = newAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not a member anymore because: batman lost his mind.")

	au.Reason = ""
	au.Previous = newAffiliationFromString(data.AffiliationAdmin)
	au.Actor = newTestActor("robin", newAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not an administrator anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not an administrator anymore because: batman lost his mind.")

	au.Previous = newAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not an owner anymore because: batman lost his mind.")

	au.Previous = newAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not a member anymore because: batman lost his mind.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationOutcast(c *C) {
	initMUCRoomConversationDisplayI18n()

	au := data.AffiliationUpdate{
		Nickname: "alice",
		New:      newAffiliationFromString(data.AffiliationOutcast),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "alice was banned from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "alice was banned from the room because: she was rude.")

	au.Reason = ""
	au.Actor = newTestActor("bob", newAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The administrator bob banned alice from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The administrator bob banned alice from the room because: she was rude.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationAdded(c *C) {
	initMUCRoomConversationDisplayI18n()

	au := data.AffiliationUpdate{
		Nickname: "juanito",
		New:      newAffiliationFromString(data.AffiliationMember),
		Previous: newAffiliationFromString(data.AffiliationNone),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "juanito is now a member.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "juanito is now a member because: el es súper chévere.")

	au.Reason = ""
	au.Actor = newTestActor("pepito", newAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner pepito changed the position of juanito; juanito is now a member.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner pepito changed the position of juanito; juanito is now a member because: el es súper chévere.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationChanged(c *C) {
	initMUCRoomConversationDisplayI18n()

	au := data.AffiliationUpdate{
		Nickname: "thor",
		New:      newAffiliationFromString(data.AffiliationAdmin),
		Previous: newAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "The position of thor was changed from member to administrator.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The position of thor was changed from member to administrator because: he is the strongest avenger.")

	au.Reason = ""
	au.Actor = newTestActor("odin", newAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner odin changed the position of thor from member to administrator.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner odin changed the position of thor from member to administrator because: he is the strongest avenger.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_affiliationUpdate(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "chavo",
		New:      newAffiliationFromString(data.AffiliationAdmin),
		Previous: newAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getMUCNotificationMessageFrom(au), Equals, "The position of chavo was changed from member to administrator.")

	au.Previous = newAffiliationFromString(data.AffiliationNone)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "chavo is now an administrator.")

	au.New = newAffiliationFromString(data.AffiliationOwner)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "chavo is now an owner.")

	au.Previous = newAffiliationFromString(data.AffiliationOwner)
	au.New = newAffiliationFromString(data.AffiliationMember)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "The position of chavo was changed from owner to member.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleModerator(c *C) {
	ru := data.RoleUpdate{
		Nickname: "wanda",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleParticipant),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "The role of wanda was changed from participant to moderator.")

	ru.Reason = "vision wanted it"
	c.Assert(getRoleUpdateMessage(ru), Equals, "The role of wanda was changed from participant to moderator because: vision wanted it.")

	ru.Reason = ""
	ru.Actor = newTestActor("vision", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "The administrator vision changed the role of wanda from participant to moderator.")

	ru.Reason = "vision wanted it"
	c.Assert(getRoleUpdateMessage(ru), Equals, "The administrator vision changed the role of wanda from participant to moderator because: vision wanted it.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleParticipant(c *C) {
	ru := data.RoleUpdate{
		Nickname: "sancho",
		New:      newTestRoleFromString(data.RoleParticipant),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "The role of sancho was changed from moderator to participant.")

	ru.Reason = "los molinos son gigantes"
	c.Assert(getRoleUpdateMessage(ru), Equals, "The role of sancho was changed from moderator to participant because: los molinos son gigantes.")

	ru.Reason = ""
	ru.Actor = newTestActor("panza", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "The owner panza changed the role of sancho from moderator to participant.")

	ru.Reason = "los molinos son gigantes"
	c.Assert(getRoleUpdateMessage(ru), Equals, "The owner panza changed the role of sancho from moderator to participant because: los molinos son gigantes.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleVisitor(c *C) {
	ru := data.RoleUpdate{
		Nickname: "chapulin",
		New:      newTestRoleFromString(data.RoleVisitor),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "The role of chapulin was changed from moderator to visitor.")

	ru.Reason = "no contaban con mi astucia"
	c.Assert(getRoleUpdateMessage(ru), Equals, "The role of chapulin was changed from moderator to visitor because: no contaban con mi astucia.")

	ru.Reason = ""
	ru.Actor = newTestActor("chespirito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "The owner chespirito changed the role of chapulin from moderator to visitor.")

	ru.Reason = "no contaban con mi astucia"
	c.Assert(getRoleUpdateMessage(ru), Equals, "The owner chespirito changed the role of chapulin from moderator to visitor because: no contaban con mi astucia.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_roleUpdate(c *C) {
	ru := data.RoleUpdate{
		Nickname: "pablo",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "The role of pablo was changed from visitor to moderator.")

	ru.Previous = newTestRoleFromString(data.RoleParticipant)
	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "The role of pablo was changed from participant to moderator.")

	ru.New = newTestRoleFromString(data.RoleVisitor)
	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "The role of pablo was changed from participant to visitor.")
}

func newTestActor(nickname string, affiliation data.Affiliation, role data.Role) *data.Actor {
	return &data.Actor{
		Nickname:    nickname,
		Affiliation: affiliation,
		Role:        role,
	}
}

func newTestAffiliationFromString(s string) data.Affiliation {
	a, err := data.AffiliationFromString(s)
	if err != nil {
		return nil
	}
	return a
}

func newTestRoleFromString(s string) data.Role {
	r, err := data.RoleFromString(s)
	if err != nil {
		return nil
	}
	return r
}
