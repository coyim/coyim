package gui

import (
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/session/muc/data"
	. "gopkg.in/check.v1"
)

type MUCNotificationMessagesSuite struct{}

var _ = Suite(&MUCNotificationMessagesSuite{})

func (s *MUCNotificationMessagesSuite) SetUpSuite(c *C) {
	initMUCi18n()
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationNone(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "batman",
		New:      newTestAffiliationFromString(data.AffiliationNone),
		Previous: newTestAffiliationFromString(data.AffiliationAdmin),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] $nickname{batman} is not $affiliation{an administrator} anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] $nickname{batman} is not $affiliation{an administrator} anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] $nickname{batman} is not $affiliation{an owner} anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] $nickname{batman} is not $affiliation{a member} anymore. The reason given was: batman lost his mind.")

	au.Reason = ""
	au.Previous = newTestAffiliationFromString(data.AffiliationAdmin)
	au.Actor = newTestActor("robin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The owner $nickname{robin} changed the position of $nickname{batman}; $nickname{batman} is not $affiliation{an administrator} anymore.")

	au.Reason = "batman lost his mind"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner $nickname{robin} changed the position of $nickname{batman}; $nickname{batman} is not $affiliation{an administrator} anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner $nickname{robin} changed the position of $nickname{batman}; $nickname{batman} is not $affiliation{an owner} anymore. The reason given was: batman lost his mind.")

	au.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner $nickname{robin} changed the position of $nickname{batman}; $nickname{batman} is not $affiliation{a member} anymore. The reason given was: batman lost his mind.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationOutcast(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "alice",
		New:      newTestAffiliationFromString(data.AffiliationOutcast),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] $nickname{alice} was banned from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] $nickname{alice} was banned from the room. The reason given was: she was rude.")

	au.Reason = ""
	au.Actor = newTestActor("bob", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The administrator $nickname{bob} banned $nickname{alice} from the room.")

	au.Reason = "she was rude"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The administrator $nickname{bob} banned $nickname{alice} from the room. The reason given was: she was rude.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationAdded(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "juanito",
		New:      newTestAffiliationFromString(data.AffiliationMember),
		Previous: newTestAffiliationFromString(data.AffiliationNone),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] $nickname{juanito} is now $affiliation{a member}.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] $nickname{juanito} is now $affiliation{a member}. The reason given was: el es súper chévere.")

	au.Reason = ""
	au.Actor = newTestActor("pepito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The owner $nickname{pepito} changed the position of $nickname{juanito}; $nickname{juanito} is now $affiliation{a member}.")

	au.Reason = "el es súper chévere"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner $nickname{pepito} changed the position of $nickname{juanito}; $nickname{juanito} is now $affiliation{a member}. The reason given was: el es súper chévere.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateMessage_affiliationChanged(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "thor",
		New:      newTestAffiliationFromString(data.AffiliationAdmin),
		Previous: newTestAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The position of $nickname{thor} was changed from $affiliation{member} to $affiliation{administrator}.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The position of $nickname{thor} was changed from $affiliation{member} to $affiliation{administrator}. The reason given was: he is the strongest avenger.")

	au.Reason = ""
	au.Actor = newTestActor("odin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] The owner $nickname{odin} changed the position of $nickname{thor} from $affiliation{member} to $affiliation{administrator}.")

	au.Reason = "he is the strongest avenger"
	c.Assert(getAffiliationUpdateMessage(au), Equals, "[localized] [localized] The owner $nickname{odin} changed the position of $nickname{thor} from $affiliation{member} to $affiliation{administrator}. The reason given was: he is the strongest avenger.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_affiliationUpdate(c *C) {
	au := data.AffiliationUpdate{
		Nickname: "chavo",
		New:      newTestAffiliationFromString(data.AffiliationAdmin),
		Previous: newTestAffiliationFromString(data.AffiliationMember),
	}

	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] The position of $nickname{chavo} was changed from $affiliation{member} to $affiliation{administrator}.")

	au.Previous = newTestAffiliationFromString(data.AffiliationNone)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] $nickname{chavo} is now $affiliation{an administrator}.")

	au.New = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] $nickname{chavo} is now $affiliation{an owner}.")

	au.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	au.New = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getMUCNotificationMessageFrom(au), Equals, "[localized] The position of $nickname{chavo} was changed from $affiliation{owner} to $affiliation{member}.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleModerator(c *C) {
	ru := data.RoleUpdate{
		Nickname: "wanda",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleParticipant),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The role of wanda was changed from participant to moderator.")

	ru.Reason = "vision wanted it"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The role of wanda was changed from participant to moderator. The reason given was: vision wanted it.")

	ru.Reason = ""
	ru.Actor = newTestActor("vision", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The administrator $nickname{vision} changed the role of $nickname{wanda} from $role{participant} to $role{moderator}.")

	ru.Reason = "vision wanted it"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The administrator $nickname{vision} changed the role of $nickname{wanda} from $role{participant} to $role{moderator}. The reason given was: vision wanted it.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleParticipant(c *C) {
	ru := data.RoleUpdate{
		Nickname: "sancho",
		New:      newTestRoleFromString(data.RoleParticipant),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The role of sancho was changed from moderator to participant.")

	ru.Reason = "los molinos son gigantes"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The role of sancho was changed from moderator to participant. The reason given was: los molinos son gigantes.")

	ru.Reason = ""
	ru.Actor = newTestActor("panza", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The owner $nickname{panza} changed the role of $nickname{sancho} from $role{moderator} to $role{participant}.")

	ru.Reason = "los molinos son gigantes"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The owner $nickname{panza} changed the role of $nickname{sancho} from $role{moderator} to $role{participant}. The reason given was: los molinos son gigantes.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleVisitor(c *C) {
	ru := data.RoleUpdate{
		Nickname: "chapulin",
		New:      newTestRoleFromString(data.RoleVisitor),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The role of chapulin was changed from moderator to visitor.")

	ru.Reason = "no contaban con mi astucia"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The role of chapulin was changed from moderator to visitor. The reason given was: no contaban con mi astucia.")

	ru.Reason = ""
	ru.Actor = newTestActor("chespirito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The owner $nickname{chespirito} changed the role of $nickname{chapulin} from $role{moderator} to $role{visitor}.")

	ru.Reason = "no contaban con mi astucia"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The owner $nickname{chespirito} changed the role of $nickname{chapulin} from $role{moderator} to $role{visitor}. The reason given was: no contaban con mi astucia.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateMessage_roleNone(c *C) {
	ru := data.RoleUpdate{
		Nickname: "alberto",
		New:      newTestRoleFromString(data.RoleNone),
		Previous: newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] $nickname{alberto} was expelled from the room.")

	ru.Reason = "bla"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] $nickname{alberto} was expelled from the room. The reason given was: bla.")

	ru.Reason = ""
	ru.Actor = newTestActor("foo", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] The owner $nickname{foo} expelled $nickname{alberto} from the room.")

	ru.Reason = "bla"
	c.Assert(getRoleUpdateMessage(ru), Equals, "[localized] [localized] The owner $nickname{foo} expelled $nickname{alberto} from the room. The reason given was: bla.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfRoleUpdateMessage_roleModerator(c *C) {
	sru := data.SelfRoleUpdate{}
	sru.RoleUpdate = data.RoleUpdate{
		Nickname: "wanda",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleParticipant),
	}

	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] Your role was changed from $role{participant} to $role{moderator}.")

	sru.Reason = "vision wanted it"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] Your role was changed from $role{participant} to $role{moderator}. The reason given was: vision wanted it.")

	sru.Reason = ""
	sru.Actor = newTestActor("vision", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] The administrator $nickname{vision} changed your role from $role{participant} to $role{moderator}.")

	sru.Reason = "vision wanted it"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] The administrator $nickname{vision} changed your role from $role{participant} to $role{moderator}. The reason given was: vision wanted it.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfRoleUpdateMessage_roleParticipant(c *C) {
	sru := data.SelfRoleUpdate{}
	sru.RoleUpdate = data.RoleUpdate{
		Nickname: "sancho",
		New:      newTestRoleFromString(data.RoleParticipant),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] Your role was changed from $role{moderator} to $role{participant}.")

	sru.Reason = "los molinos son gigantes"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] Your role was changed from $role{moderator} to $role{participant}. The reason given was: los molinos son gigantes.")

	sru.Reason = ""
	sru.Actor = newTestActor("panza", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] The owner $nickname{panza} changed your role from $role{moderator} to $role{participant}.")

	sru.Reason = "los molinos son gigantes"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] The owner $nickname{panza} changed your role from $role{moderator} to $role{participant}. The reason given was: los molinos son gigantes.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfRoleUpdateMessage_roleVisitor(c *C) {
	sru := data.SelfRoleUpdate{}
	sru.RoleUpdate = data.RoleUpdate{
		Nickname: "chapulin",
		New:      newTestRoleFromString(data.RoleVisitor),
		Previous: newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] Your role was changed from $role{moderator} to $role{visitor}.")

	sru.Reason = "no contaban con mi astucia"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] Your role was changed from $role{moderator} to $role{visitor}. The reason given was: no contaban con mi astucia.")

	sru.Reason = ""
	sru.Actor = newTestActor("chespirito", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] The owner $nickname{chespirito} changed your role from $role{moderator} to $role{visitor}.")

	sru.Reason = "no contaban con mi astucia"
	c.Assert(getSelfRoleUpdateMessage(sru), Equals, "[localized] [localized] The owner $nickname{chespirito} changed your role from $role{moderator} to $role{visitor}. The reason given was: no contaban con mi astucia.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_roleUpdate(c *C) {
	ru := data.RoleUpdate{
		Nickname: "pablo",
		New:      newTestRoleFromString(data.RoleModerator),
		Previous: newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "[localized] The role of pablo was changed from visitor to moderator.")

	ru.Previous = newTestRoleFromString(data.RoleParticipant)
	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "[localized] The role of pablo was changed from participant to moderator.")

	ru.New = newTestRoleFromString(data.RoleVisitor)
	c.Assert(getMUCNotificationMessageFrom(ru), Equals, "[localized] The role of pablo was changed from participant to visitor.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationRemoved(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "007",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] $nickname{007} is not $affiliation{an administrator} anymore. [localized] As a result, their role was changed from $role{moderator} to $role{visitor}.")

	aru.Reason = "he is an assassin"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] $nickname{007} is not $affiliation{an administrator} anymore. [localized] As a result, their role was changed from $role{moderator} to $role{visitor}. The reason given was: he is an assassin.")

	aru.Reason = ""
	aru.Actor = newTestActor("the enemy", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The owner $nickname{the enemy} changed the position of $nickname{007}; $nickname{007} is not $affiliation{an administrator} anymore. [localized] As a result, their role was changed from $role{moderator} to $role{visitor}.")

	aru.Reason = "bla"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The owner $nickname{the enemy} changed the position of $nickname{007}; $nickname{007} is not $affiliation{an administrator} anymore. [localized] As a result, their role was changed from $role{moderator} to $role{visitor}. The reason given was: bla.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationAdded(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "alice",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationNone),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The position of $nickname{alice} was changed to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}.")

	aru.Reason = "she is lost in the world of wonders"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The position of $nickname{alice} was changed to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}. The reason given was: she is lost in the world of wonders.")

	aru.Reason = ""
	aru.Actor = newTestActor("rabbit", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The owner $nickname{rabbit} changed the position of $nickname{alice} to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}.")

	aru.Reason = "she is lost in the world of wonders"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The owner $nickname{rabbit} changed the position of $nickname{alice} to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}. The reason given was: she is lost in the world of wonders.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationRoleUpdateMessage_affiliationUpdated(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "Pegassus",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationOwner),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The position of $nickname{Pegassus} was changed from $affiliation{owner} to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}.")

	aru.Reason = "he is a silver warrior"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The position of $nickname{Pegassus} was changed from $affiliation{owner} to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}. The reason given was: he is a silver warrior.")

	aru.Reason = ""
	aru.Actor = newTestActor("Ikki", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] The owner $nickname{Ikki} changed the position of $nickname{Pegassus} from $affiliation{owner} to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}.")

	aru.Reason = "he has the phoenix flame"
	c.Assert(getAffiliationRoleUpdateMessage(aru), Equals, "[localized] [localized] [localized] The owner $nickname{Ikki} changed the position of $nickname{Pegassus} from $affiliation{owner} to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}. The reason given was: he has the phoenix flame.")
}

func (*MUCNotificationMessagesSuite) Test_getMUCNotificationMessageFrom_affiliationRoleUpdate(c *C) {
	aru := data.AffiliationRoleUpdate{
		Nickname:            "chavo",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getMUCNotificationMessageFrom(aru), Equals, "[localized] [localized] $nickname{chavo} is not $affiliation{an administrator} anymore. [localized] As a result, their role was changed from $role{moderator} to $role{visitor}.")

	aru.NewAffiliation = newTestAffiliationFromString(data.AffiliationAdmin)
	aru.PreviousAffiliation = newTestAffiliationFromString(data.AffiliationNone)
	aru.NewRole = newTestRoleFromString(data.RoleModerator)
	aru.PreviousRole = newTestRoleFromString(data.RoleVisitor)
	c.Assert(getMUCNotificationMessageFrom(aru), Equals, "[localized] [localized] The position of $nickname{chavo} was changed to $affiliation{administrator}. [localized] As a result, their role was changed from $role{visitor} to $role{moderator}.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationRemoved(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationNone),
			Previous: newTestAffiliationFromString(data.AffiliationAdmin),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are not $affiliation{an administrator} anymore.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are not $affiliation{an administrator} anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are not $affiliation{an owner} anymore. The reason given was: you are funny.")

	sau.Reason = ""
	sau.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are not $affiliation{a member} anymore.")

	sau.Actor = newTestActor("robin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The owner $nickname{robin} changed your position; you are not $affiliation{a member} anymore.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner $nickname{robin} changed your position; you are not $affiliation{a member} anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner $nickname{robin} changed your position; you are not $affiliation{an owner} anymore. The reason given was: you are funny.")

	sau.Previous = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner $nickname{robin} changed your position; you are not $affiliation{a member} anymore. The reason given was: you are funny.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationAdded(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationAdmin),
			Previous: newTestAffiliationFromString(data.AffiliationNone),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are now $affiliation{an administrator}.")

	sau.Reason = "estás encopetao"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are now $affiliation{an administrator}. The reason given was: estás encopetao.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationMember)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are now $affiliation{a member}.")

	sau.Reason = "you dance very well"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are now $affiliation{a member}. The reason given was: you dance very well.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationOwner)
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You are now $affiliation{an owner}.")

	sau.Reason = "the day is cool"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You are now $affiliation{an owner}. The reason given was: the day is cool.")

	sau.Reason = ""
	sau.AffiliationUpdate.New = newTestAffiliationFromString(data.AffiliationAdmin)
	sau.Actor = newTestActor("paco", newTestAffiliationFromString(data.AffiliationAdmin), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The administrator $nickname{paco} changed your position; you are now $affiliation{an administrator}.")

	sau.Reason = "you are funny"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The administrator $nickname{paco} changed your position; you are now $affiliation{an administrator}. The reason given was: you are funny.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationChanged(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New:      newTestAffiliationFromString(data.AffiliationAdmin),
			Previous: newTestAffiliationFromString(data.AffiliationMember),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] Your position was changed from $affiliation{member} to $affiliation{administrator}.")

	sau.Reason = "you are loco"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] Your position was changed from $affiliation{member} to $affiliation{administrator}. The reason given was: you are loco.")

	sau.Reason = ""
	sau.Actor = newTestActor("chapulin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The owner $nickname{chapulin} changed your position from $affiliation{member} to $affiliation{administrator}.")

	sau.Reason = "you are locote"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner $nickname{chapulin} changed your position from $affiliation{member} to $affiliation{administrator}. The reason given was: you are locote.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateSuccessMessage(c *C) {
	nickname := "Juan"
	owner := newTestAffiliationFromString(data.AffiliationOwner)
	admin := newTestAffiliationFromString(data.AffiliationAdmin)
	member := newTestAffiliationFromString(data.AffiliationMember)
	outcast := newTestAffiliationFromString(data.AffiliationOutcast)
	none := newTestAffiliationFromString(data.AffiliationNone)

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, member, none), Equals,
		"[localized] $nickname{Juan} is not $affiliation{a member} anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, admin, none), Equals,
		"[localized] $nickname{Juan} is not $affiliation{an administrator} anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, owner, none), Equals,
		"[localized] $nickname{Juan} is not $affiliation{an owner} anymore.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, member), Equals,
		"[localized] The position of $nickname{Juan} was changed to $affiliation{member}.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, admin), Equals,
		"[localized] The position of $nickname{Juan} was changed to $affiliation{administrator}.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, owner), Equals,
		"[localized] The position of $nickname{Juan} was changed to $affiliation{owner}.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, none, outcast), Equals,
		"[localized] $nickname{Juan} has been banned from the room.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, member, outcast), Equals,
		"[localized] $nickname{Juan} has been banned from the room.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, admin, outcast), Equals,
		"[localized] $nickname{Juan} has been banned from the room.")

	c.Assert(getAffiliationUpdateSuccessMessage(nickname, owner, outcast), Equals,
		"[localized] $nickname{Juan} has been banned from the room.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateSuccessMessage(c *C) {
	moderator := newTestRoleFromString(data.RoleModerator)
	participant := newTestRoleFromString(data.RoleParticipant)
	visitor := newTestRoleFromString(data.RoleVisitor)
	none := newTestRoleFromString(data.RoleNone)

	c.Assert(getRoleUpdateSuccessMessage("Maria", moderator, none), Equals, "[localized] $nickname{Maria} was expelled from the room.")
	c.Assert(getRoleUpdateSuccessMessage("Carlos", participant, none), Equals, "[localized] $nickname{Carlos} was expelled from the room.")
	c.Assert(getRoleUpdateSuccessMessage("Mauricio", visitor, none), Equals, "[localized] $nickname{Mauricio} was expelled from the room.")

	c.Assert(getRoleUpdateSuccessMessage("Jose", none, moderator), Equals, "[localized] The role of $nickname{Jose} was changed to $role{moderator}.")
	c.Assert(getRoleUpdateSuccessMessage("Alberto", none, participant), Equals, "[localized] The role of $nickname{Alberto} was changed to $role{participant}.")
	c.Assert(getRoleUpdateSuccessMessage("Juan", none, visitor), Equals, "[localized] The role of $nickname{Juan} was changed to $role{visitor}.")

	c.Assert(getRoleUpdateSuccessMessage("Alberto", moderator, participant), Equals, "[localized] The role of $nickname{Alberto} was changed from $role{moderator} to $role{participant}.")
	c.Assert(getRoleUpdateSuccessMessage("Alberto", moderator, visitor), Equals, "[localized] The role of $nickname{Alberto} was changed from $role{moderator} to $role{visitor}.")
	c.Assert(getRoleUpdateSuccessMessage("Alberto", participant, moderator), Equals, "[localized] The role of $nickname{Alberto} was changed from $role{participant} to $role{moderator}.")
	c.Assert(getRoleUpdateSuccessMessage("Carlos", participant, visitor), Equals, "[localized] The role of $nickname{Carlos} was changed from $role{participant} to $role{visitor}.")
	c.Assert(getRoleUpdateSuccessMessage("Carlos", visitor, participant), Equals, "[localized] The role of $nickname{Carlos} was changed from $role{visitor} to $role{participant}.")
	c.Assert(getRoleUpdateSuccessMessage("Juan", visitor, moderator), Equals, "[localized] The role of $nickname{Juan} was changed from $role{visitor} to $role{moderator}.")
}

func (*MUCNotificationMessagesSuite) Test_getAffiliationUpdateFailureMessage(c *C) {
	owner := newTestAffiliationFromString(data.AffiliationOwner)
	admin := newTestAffiliationFromString(data.AffiliationAdmin)
	member := newTestAffiliationFromString(data.AffiliationMember)
	outcast := newTestAffiliationFromString(data.AffiliationOutcast)
	none := newTestAffiliationFromString(data.AffiliationNone)

	messages := getAffiliationUpdateFailureMessage("Luisa", owner, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{Luisa} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of Luisa couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of Luisa to owner.")

	messages = getAffiliationUpdateFailureMessage("Marco", admin, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{Marco} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of Marco couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of Marco to administrator.")

	messages = getAffiliationUpdateFailureMessage("Pedro", member, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{Pedro} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of Pedro couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of Pedro to member.")

	messages = getAffiliationUpdateFailureMessage("Luisa", outcast, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{Luisa} couldn't be banned.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Banning failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] Luisa couldn't be banned")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to ban Luisa.")

	messages = getAffiliationUpdateFailureMessage("José", none, nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] The position of $nickname{José} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the position failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The position of José couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the position of José.")

}

func (*MUCNotificationMessagesSuite) Test_getRoleUpdateFailureMessage(c *C) {
	moderator := newTestRoleFromString(data.RoleModerator)
	participant := newTestRoleFromString(data.RoleParticipant)
	visitor := newTestRoleFromString(data.RoleVisitor)
	none := newTestRoleFromString(data.RoleNone)

	messages := getRoleUpdateFailureMessage("Mauricio", moderator)
	c.Assert(messages.notificationMessage, Equals, "[localized] The role of $nickname{Mauricio} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the role failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The role of Mauricio couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the role of Mauricio to moderator.")

	messages = getRoleUpdateFailureMessage("Juan", participant)
	c.Assert(messages.notificationMessage, Equals, "[localized] The role of $nickname{Juan} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the role failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The role of Juan couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the role of Juan to participant.")

	messages = getRoleUpdateFailureMessage("Pepe", visitor)
	c.Assert(messages.notificationMessage, Equals, "[localized] The role of $nickname{Pepe} couldn't be changed.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Changing the role failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] The role of Pepe couldn't be changed")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred trying to change the role of Pepe to visitor.")

	messages = getRoleUpdateFailureMessage("Juana", none)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{Juana} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] Juana couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred expelling Juana.")
}

func (*MUCNotificationMessagesSuite) Test_getRoleRemoveFailureMessage(c *C) {
	messages := getRoleRemoveFailureMessage("foo", newTestAffiliationFromString(data.AffiliationOwner), nil)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{foo} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] foo couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] An error occurred expelling foo.")

	messages = getRoleRemoveFailureMessage("nil", nil, session.ErrNotAllowedKickOccupant)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{nil} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] nil couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] You don't have permissions to expel nil.")

	messages = getRoleRemoveFailureMessage("bla", newTestAffiliationFromString(data.AffiliationAdmin), session.ErrNotAllowedKickOccupant)
	c.Assert(messages.notificationMessage, Equals, "[localized] $nickname{bla} couldn't be expelled.")
	c.Assert(messages.errorDialogTitle, Equals, "[localized] Expelling failed")
	c.Assert(messages.errorDialogHeader, Equals, "[localized] bla couldn't be expelled")
	c.Assert(messages.errorDialogMessage, Equals, "[localized] As an administrator you don't have permissions to expel bla.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationUpdateMessage_affiliationOutcast(c *C) {
	sau := data.SelfAffiliationUpdate{
		AffiliationUpdate: data.AffiliationUpdate{
			New: newTestAffiliationFromString(data.AffiliationOutcast),
		},
	}

	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] You have been banned from the room.")

	sau.Reason = "it's so cold"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] You have been banned from the room. The reason given was: it's so cold.")

	sau.Reason = ""
	sau.Actor = newTestActor("calvin", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] The owner $nickname{calvin} banned you from the room.")

	sau.Reason = "it isn't cool"
	c.Assert(getSelfAffiliationUpdateMessage(sau), Equals, "[localized] [localized] The owner $nickname{calvin} banned you from the room. The reason given was: it isn't cool.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationRoleUpdateMessage_affiliationRemoved(c *C) {
	saru := data.SelfAffiliationRoleUpdate{}
	saru.AffiliationRoleUpdate = data.AffiliationRoleUpdate{
		Nickname:            "007",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationNone),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationAdmin),
		NewRole:             newTestRoleFromString(data.RoleVisitor),
		PreviousRole:        newTestRoleFromString(data.RoleModerator),
	}

	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] You are not $affiliation{an administrator} anymore. [localized] As a result, your role was changed from $role{moderator} to $role{visitor}.")

	saru.Reason = "he is an assassin"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] You are not $affiliation{an administrator} anymore. [localized] As a result, your role was changed from $role{moderator} to $role{visitor}. The reason given was: he is an assassin.")

	saru.Reason = ""
	saru.Actor = newTestActor("the enemy", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] The owner $nickname{the enemy} changed your position; you are not $affiliation{an administrator} anymore. [localized] As a result, your role was changed from $role{moderator} to $role{visitor}.")

	saru.Reason = "bla"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] The owner $nickname{the enemy} changed your position; you are not $affiliation{an administrator} anymore. [localized] As a result, your role was changed from $role{moderator} to $role{visitor}. The reason given was: bla.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationRoleUpdateMessage_affiliationAdded(c *C) {
	saru := data.SelfAffiliationRoleUpdate{}
	saru.AffiliationRoleUpdate = data.AffiliationRoleUpdate{
		Nickname:            "alice",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationNone),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] Your position was changed to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}.")

	saru.Reason = "she is lost in the world of wonders"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] Your position was changed to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}. The reason given was: she is lost in the world of wonders.")

	saru.Reason = ""
	saru.Actor = newTestActor("rabbit", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] The owner $nickname{rabbit} changed your position to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}.")

	saru.Reason = "she is lost in the world of wonders"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] The owner $nickname{rabbit} changed your position to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}. The reason given was: she is lost in the world of wonders.")
}

func (*MUCNotificationMessagesSuite) Test_getSelfAffiliationRoleUpdateMessage_affiliationUpdated(c *C) {
	saru := data.SelfAffiliationRoleUpdate{}
	saru.AffiliationRoleUpdate = data.AffiliationRoleUpdate{
		Nickname:            "goku",
		NewAffiliation:      newTestAffiliationFromString(data.AffiliationAdmin),
		PreviousAffiliation: newTestAffiliationFromString(data.AffiliationMember),
		NewRole:             newTestRoleFromString(data.RoleModerator),
		PreviousRole:        newTestRoleFromString(data.RoleVisitor),
	}

	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] Your position was changed from $affiliation{member} to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}.")

	saru.Reason = "you are a powerfull Saiyajin"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] Your position was changed from $affiliation{member} to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}. The reason given was: you are a powerfull Saiyajin.")

	saru.Reason = ""
	saru.Actor = newTestActor("vegeta", newTestAffiliationFromString(data.AffiliationOwner), newTestRoleFromString(data.RoleModerator))
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] The owner $nickname{vegeta} changed your position from $affiliation{member} to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}.")

	saru.Reason = "he is the prince of the Saiyajins"
	c.Assert(getSelfAffiliationRoleUpdateMessage(saru), Equals, "[localized] [localized] [localized] The owner $nickname{vegeta} changed your position from $affiliation{member} to $affiliation{administrator}. [localized] As a result, your role was changed from $role{visitor} to $role{moderator}. The reason given was: he is the prince of the Saiyajins.")
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
		panic("Error produced trying to get an affiliation from a string")
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
