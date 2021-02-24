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
		New:      newTestAffiliationFromString(data.AffiliationNone),
		Previous: newTestAffiliationFromString(data.AffiliationAdmin),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not an administrator anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not an administrator anymore because: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not an owner anymore because: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "batman is not a member anymore because: batman lost his mind.")

	au.Reason = ""
	au.Previous = newTestAffiliationFromString(data.AffiliationAdmin)
	au.Actor = newTestActor("robin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not an administrator anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not an administrator anymore because: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not an owner anymore because: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner robin changed the position of batman; batman is not a member anymore because: batman lost his mind.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationOutcast(c *C) {
	initMUCRoomConversationDisplayI18n()

	au := data.AffiliationUpdate{
		Nickname: "alice",
		New:      newTestAffiliationFromString(data.AffiliationOutcast),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "alice was banned from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "alice was banned from the room because: she was rude.")

	au.Reason = ""
	au.Actor = newTestActor("bob", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The administrator bob banned alice from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The administrator bob banned alice from the room because: she was rude.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationAdded(c *C) {
	initMUCRoomConversationDisplayI18n()

	au := data.AffiliationUpdate{
		Nickname: "juanito",
		New:      newTestAffiliationFromString(data.AffiliationMember),
		Previous: newTestAffiliationFromString(data.AffiliationNone),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "juanito is now a member.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "juanito is now a member because: el es súper chévere.")

	au.Reason = ""
	au.Actor = newTestActor("pepito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner pepito changed the position of juanito; juanito is now a member.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner pepito changed the position of juanito; juanito is now a member because: el es súper chévere.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationChanged(c *C) {
	initMUCRoomConversationDisplayI18n()

	au := data.AffiliationUpdate{
		Nickname: "thor",
		New:      newTestAffiliationFromString(data.AffiliationAdmin),
		Previous: newTestAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "The position of thor was changed from member to administrator.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The position of thor was changed from member to administrator because: he is the strongest avenger.")

	au.Reason = ""
	au.Actor = newTestActor("odin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner odin changed the position of thor from member to administrator.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "The owner odin changed the position of thor from member to administrator because: he is the strongest avenger.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_affiliationUpdate(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "chavo",
		New:      newTestAffiliationFromString(data.AffiliationAdmin),
		Previous: newTestAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getMUCNotificationMessageFrom(au), Equals, "The position of chavo was changed from member to administrator.")

	au.Previous = newTestAffiliationFromString(data.AffiliationNone)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "chavo is now an administrator.")

	au.New = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "chavo is now an owner.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	au.New = newTestAffiliationFromString(data.AffiliationMember)
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

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_noCommonCases(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "superman",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationNone),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The position and the role of superman were changed.")

	aru.Reason = "foo"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The position and the role of superman were changed because: foo.")

	aru.Reason = ""
	aru.Actor = newTestActor("louis", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The owner louis changed the position of superman. As a result, their role was also changed.")

	aru.Reason = "foo"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The owner louis changed the position of superman. As a result, their role was also changed. The reason given was: foo.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationRemoved(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "007",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "007 is not an administrator anymore. As a result, their role was changed from moderator to visitor.")

	aru.Reason = "he is an assassin"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "007 is not an administrator anymore. As a result, their role was changed from moderator to visitor. The reason given was: he is an assassin.")

	aru.Reason = ""
	aru.Actor = newTestActor("the enemy", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The owner the enemy changed the position of 007; 007 is not an administrator anymore. As a result, their role was changed from moderator to visitor.")

	aru.Reason = "bla"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The owner the enemy changed the position of 007; 007 is not an administrator anymore. As a result, their role was changed from moderator to visitor. The reason given was: bla.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationAdded(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "alice",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationNone),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The position of alice was changed to administrator. As a result, their role was changed from visitor to moderator.")

	aru.Reason = "she is lost in the world of wonders"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The position of alice was changed to administrator. As a result, their role was changed from visitor to moderator. The reason given was: she is lost in the world of wonders.")

	aru.Reason = ""
	aru.Actor = newTestActor("rabbit", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The owner rabbit changed the position of alice to administrator. As a result, their role was changed from visitor to moderator.")

	aru.Reason = "she is lost in the world of wonders"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "The owner rabbit changed the position of alice to administrator. As a result, their role was changed from visitor to moderator. The reason given was: she is lost in the world of wonders.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_affiliationRoleUpdate(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "chavo",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getMUCNotificationMessageFrom(aru), Equals, "chavo is not an administrator anymore. As a result, their role was changed from moderator to visitor.")

	aru.NewAffiliation = newTestAffiliationFromString(data.AffiliationAdmin)
	aru.PreviousAffiliation = newTestAffiliationFromString(data.AffiliationNone)
	aru.NewRole = newTestRoleFromString(data.RoleModerator)
	aru.PreviousRole = newTestRoleFromString(data.RoleVisitor)
	c.Assert(getMUCNotificationMessageFrom(aru), Equals, "The position of chavo was changed to administrator. As a result, their role was changed from visitor to moderator.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationRemoved(c *C) {
	initMUCRoomConversationDisplayI18n()

	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationNone),
			Previous: newTestAffiliationFromString(data.AffiliationAdmin),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are not an administrator anymore.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are not an administrator anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are not an owner anymore. The reason given was: you are funny.")

	sau.Reason = ""
	sau.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are not a member anymore.")

	sau.Actor = newTestActor("robin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner robin changed your position; you are not a member anymore.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner robin changed your position; you are not a member anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner robin changed your position; you are not an owner anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner robin changed your position; you are not a member anymore. The reason given was: you are funny.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationAdded(c *C) {
	initMUCRoomConversationDisplayI18n()

	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationAdmin),
			Previous: newTestAffiliationFromString(data.AffiliationNone),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are now an administrator.")

	sau.Reason = "estás encopetao"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are now an administrator. The reason given was: estás encopetao.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are now a member.")

	sau.Reason = "you dance very well"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are now a member. The reason given was: you dance very well.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are now an owner.")

	sau.Reason = "the day is cool"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You are now an owner. The reason given was: the day is cool.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationAdmin)
	sau.Actor = newTestActor("paco", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The administrator paco changed your position; you are now an administrator.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The administrator paco changed your position; you are now an administrator. The reason given was: you are funny.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationChanged(c *C) {
	initMUCRoomConversationDisplayI18n()

	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationAdmin),
			Previous: newTestAffiliationFromString(data.AffiliationMember),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "Your position was changed from member to administrator.")

	sau.Reason = "you are loco"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "Your position was changed from member to administrator because: you are loco.")

	sau.Reason = ""
	sau.Actor = newTestActor("chapulin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner chapulin changed your position from member to administrator.")

	sau.Reason = "you are locote"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner chapulin changed your position from member to administrator because: you are locote.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateSuccessMessage(c *C) {
	initMUCRoomConversationDisplayI18n()

	nickname := "Juan"
	owner := newTestAffiliationFromString(data.AffiliationOwner)
	admin := newTestAffiliationFromString(data.AffiliationAdmin)
	member := newTestAffiliationFromString(data.AffiliationMember)
	none := newTestAffiliationFromString(data.AffiliationNone)

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, none), Equals,
		"Juan no longer has a position.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, member, none), Equals,
		"Juan is not a member anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, admin, none), Equals,
		"Juan is not an administrator anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, owner, none), Equals,
		"Juan is not an owner anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, member), Equals,
		"The position of Juan was updated to member.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, admin), Equals,
		"The position of Juan was updated to administrator.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, owner), Equals,
		"The position of Juan was updated to owner.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationOutcast(c *C) {
	initMUCRoomConversationDisplayI18n()

	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New: newTestAffiliationFromString(data.AffiliationOutcast),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You has been banned from the room.")

	sau.Reason = "it's so cold"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "You has been banned from the room. The reason given was: it's so cold.")

	sau.Reason = ""
	sau.Actor = newTestActor("calvin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner calvin banned you from the room.")

	sau.Reason = "it isn't cool"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "The owner calvin banned you from the room. The reason given was: it isn't cool.")
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
