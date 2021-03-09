package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/session/muc/data"
)

func getAffiliationUpdateSuccessMessage(nickname string, previousAffiliation, affiliation data.Affiliation) string {
	if affiliation.IsNone() {
		// This is impossible to happen but we need to cover all cases.
		if previousAffiliation.IsNone() {
			return i18n.Localf("%[1]s no longer has a position.", nickname)
		}
		return getAffiliationRemovedSuccessMessage(nickname, previousAffiliation)
	}
	return getAffiliationChangedSuccessMessage(nickname, affiliation)
}

func getAffiliationRemovedSuccessMessage(nickname string, previousAffiliation data.Affiliation) string {
	switch {
	case previousAffiliation.IsOwner():
		return i18n.Localf("%[1]s is not an owner anymore.", nickname)
	case previousAffiliation.IsAdmin():
		return i18n.Localf("%[1]s is not an administrator anymore.", nickname)
	case previousAffiliation.IsMember():
		return i18n.Localf("%[1]s is not a member anymore.", nickname)
	default:
		return i18n.Localf("%[1]s is not banned anymore.", nickname)
	}

}

func getAffiliationChangedSuccessMessage(nickname string, affiliation data.Affiliation) string {
	switch {
	case affiliation.IsOwner():
		return i18n.Localf("The position of %[1]s was changed to owner.", nickname)
	case affiliation.IsAdmin():
		return i18n.Localf("The position of %[1]s was changed to administrator.", nickname)
	case affiliation.IsMember():
		return i18n.Localf("The position of %[1]s was changed to member.", nickname)
	default:
		return i18n.Localf("The position of %[1]s was changed.", nickname)
	}
}

// getRoleUpdateSuccessMessage returns a friendly notification message for the role update process
// This function receives the following params:
// nickname - The nickname of the occupant to whom the role was changed
// previousRole - The previous role of the occupant
// newRole - The new role of the occupant
func getRoleUpdateSuccessMessage(nickname string, previousRole, newRole data.Role) string {
	if newRole.IsDifferentFrom(previousRole) {
		switch {
		case newRole.IsNone():
			return i18n.Localf("%s is not %s anymore.", nickname, displayNameForRoleWithPreposition(previousRole))

		case previousRole.IsNone():
			return i18n.Localf("The role of %s was changed to %s.", nickname, displayNameForRole(newRole))

		default:
			return i18n.Localf("The role of %s was changed from %s to %s.", nickname, displayNameForRole(previousRole), displayNameForRole(newRole))
		}
	}

	// This is impossible to happen but we need to cover all cases.
	return i18n.Localf("The role of %s wasn't changed.", nickname)
}

type updateFailureMessages struct {
	notificationMessage string
	errorDialogTitle    string
	errorDialogHeader   string
	errorDialogMessage  string
}

func getAffiliationUpdateFailureMessage(nickname string, newAffiliation data.Affiliation, err error) *updateFailureMessages {
	return &updateFailureMessages{
		notificationMessage: i18n.Localf("The position of %[1]s couldn't be changed", nickname),
		errorDialogTitle:    i18n.Local("Changing the position failed"),
		errorDialogHeader:   i18n.Localf("The position of %[1]s couldn't be changed", nickname),
		errorDialogMessage:  getAffiliationFailureErrorMessage(nickname, newAffiliation, err),
	}
}

func getAffiliationFailureErrorMessage(nickname string, newAffiliation data.Affiliation, err error) string {
	if err == session.ErrRemoveOwnerAffiliation {
		return i18n.Local("You can't change your own position because you are the only owner for this room. Every room must have at least one owner.")
	}
	return getUpdateAffiliationFailureErrorMessage(nickname, newAffiliation)
}

func getUpdateAffiliationFailureErrorMessage(nickname string, newAffiliation data.Affiliation) string {
	switch {
	case newAffiliation.IsOwner():
		return i18n.Localf("An error occurred trying to change the position of %[1]s to owner.", nickname)
	case newAffiliation.IsAdmin():
		return i18n.Localf("An error occurred trying to change the position of %[1]s to administrator.", nickname)
	case newAffiliation.IsMember():
		return i18n.Localf("An error occurred trying to change the position of %[1]s to member.", nickname)
	default:
		return i18n.Localf("An error occurred trying to change the position of %[1]s.", nickname)
	}
}

func getRoleUpdateFailureMessage(nickname string, newRole data.Role) *updateFailureMessages {
	if newRole.IsNone() {
		return getRoleRemoveFailureMessage(nickname, nil, nil)
	}

	return &updateFailureMessages{
		notificationMessage: i18n.Localf("The role of %s couldn't be changed", nickname),
		errorDialogTitle:    i18n.Local("Changing the role failed"),
		errorDialogHeader:   i18n.Localf("The role of %s couldn't be changed", nickname),
		errorDialogMessage:  i18n.Localf("An error occurred trying to change the role of %s to %s.", nickname, displayNameForRole(newRole)),
	}
}

func getRoleRemoveFailureMessage(nickname string, actorAffiliation data.Affiliation, err error) *updateFailureMessages {
	return &updateFailureMessages{
		notificationMessage: i18n.Localf("%[1]s couldn't be expelled.", nickname),
		errorDialogTitle:    i18n.Local("Expelling process failed"),
		errorDialogHeader:   i18n.Localf("%[1]s couldn't be expelled", nickname),
		errorDialogMessage:  getRoleRemoveFailureMessageBasedOnError(nickname, actorAffiliation, err),
	}
}

func getRoleRemoveFailureMessageBasedOnError(nickname string, actorAffiliation data.Affiliation, err error) string {
	switch err {
	case session.ErrNotAllowedKickOccupant:
		return getRoleRemoveFailureMessageWithActor(nickname, actorAffiliation)
	default:
		return i18n.Localf("An error occurred expelling to %[1]s.", nickname)
	}
}

func getRoleRemoveFailureMessageWithActor(nickname string, actorAffiliation data.Affiliation) string {
	if actorAffiliation != nil {
		switch {
		case actorAffiliation.IsOwner():
			return i18n.Localf("As an owner you don't have permissions to expel %[1]s.", nickname)
		case actorAffiliation.IsAdmin():
			return i18n.Localf("As an administrator you don't have permissions to expel %[1]s.", nickname)
		case actorAffiliation.IsMember():
			return i18n.Localf("As a member you don't have permissions to expel %[1]s.", nickname)
		}
	}

	return i18n.Localf("You don't have permissions to expel %[1]s.", nickname)
}

func getMUCNotificationMessageFrom(d interface{}) string {
	switch t := d.(type) {
	case data.AffiliationUpdate:
		return getAffiliationUpdateMessage(t)
	case data.SelfAffiliationUpdate:
		return getSelfAffiliationUpdateMessage(t)
	case data.RoleUpdate:
		return getRoleUpdateMessage(t)
	case data.AffiliationRoleUpdate:
		return getAffiliationRoleUpdateMessage(t)
	default:
		return ""
	}
}

func getAffiliationUpdateMessage(affiliationUpdate data.AffiliationUpdate) (m string) {
	switch {
	case affiliationUpdate.New.IsNone():
		m = getAffiliationRemovedMessage(affiliationUpdate)
	case affiliationUpdate.New.IsBanned():
		m = getAffiliationBannedMessage(affiliationUpdate)
	case affiliationUpdate.Previous.IsNone():
		m = getAffiliationAddedMessage(affiliationUpdate)
	default:
		m = getAffiliationChangedMessage(affiliationUpdate)
	}

	return appendReasonToMessage(m, affiliationUpdate.Reason)
}

func getAffiliationRemovedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationRemovedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationRemovedMessageWithActor(affiliationUpdate)
}

func getAffiliationRemovedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner():
		return i18n.Localf("%[1]s is not an owner anymore.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("%[1]s is not an administrator anymore.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember():
		return i18n.Localf("%[1]s is not a member anymore.", affiliationUpdate.Nickname)
	default:
		return i18n.Localf("%[1]s is not banned anymore.", affiliationUpdate.Nickname)
	}
}

func getAffiliationRemovedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor.Affiliation.IsOwner() {
		return getAffiliationRemovedMessageWithOwnerActor(affiliationUpdate)
	}
	return getAffiliationRemovedMessageWithAdminActor(affiliationUpdate)
}
func getAffiliationRemovedMessageWithOwnerActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[3]s is not an owner anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[3]s is not an administrator anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[3]s is not a member anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	default:
		return i18n.Localf("The owner %[1]s changed the position of %[2]s.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}

}

func getAffiliationRemovedMessageWithAdminActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[3]s is not an owner anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[3]s is not an administrator anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[3]s is not a member anymore.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	default:
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
}

func getAffiliationBannedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationBannedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationBannedMessageWithActor(affiliationUpdate)
}

func getAffiliationBannedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	return i18n.Localf("%[1]s was banned from the room.", affiliationUpdate.Nickname)
}

func getAffiliationBannedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Actor.Affiliation.IsOwner():
		return i18n.Localf("The owner %[1]s banned %[2]s from the room.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Actor.Affiliation.IsAdmin():
		return i18n.Localf("The administrator %[1]s banned %[2]s from the room.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	default:
		return i18n.Localf("%s banned %s from the room.", affiliationUpdate.Actor.Nickname, affiliationUpdate.Nickname)
	}
}

func getAffiliationAddedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationAddedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationAddedMessageWithActor(affiliationUpdate)
}

func getAffiliationAddedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsOwner():
		return i18n.Localf("%[1]s is now an owner.", affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsAdmin():
		return i18n.Localf("%[1]s is now an administrator.", affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsMember():
		return i18n.Localf("%[1]s is now a member.", affiliationUpdate.Nickname)
	default:
		return i18n.Localf("%[1]s is now banned.", affiliationUpdate.Nickname)

	}
}

func getAffiliationAddedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor.Affiliation.IsOwner() {
		return getAffiliationAddedMessageWithOwnerActor(affiliationUpdate)
	}
	return getAffiliationAddedMessageWithAdminActor(affiliationUpdate)
}

func getAffiliationAddedMessageWithOwnerActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[3]s is now an owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[3]s is now an administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[3]s is now a member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	default:
		return i18n.Localf("The owner %[1]s changed the position of %[2]s; %[3]s is now banned",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	}
}

func getAffiliationAddedMessageWithAdminActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[3]s is now an owner",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[3]s is now an administrator",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[3]s is now a member",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	default:
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s; %[3]s is now banned",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname)
	}
}

func getAffiliationChangedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		return getAffiliationChangedMessageWithoutActor(affiliationUpdate)
	}
	return getAffiliationChangedMessageWithActor(affiliationUpdate)
}

func getAffiliationChangedMessageWithoutActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The position of %[1]s was changed from owner to administrator.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The position of %[1]s was changed from owner to member.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The position of %[1]s was changed from administrator to member.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The position of %[1]s was changed from administrator to owner.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The position of %[1]s was changed from member to administrator.", affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The position of %[1]s was changed from member to owner.", affiliationUpdate.Nickname)
	default:
		return i18n.Localf("The position of %[1]s was changed.", affiliationUpdate.Nickname)
	}
}

func getAffiliationChangedMessageWithActor(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor.Affiliation.IsOwner() {
		return getAffiliationChangedMessageWithOwnerActor(affiliationUpdate)
	}
	return getAffiliationChangedMessageWithAdminActor(affiliationUpdate)
}

func getAffiliationChangedMessageWithOwnerActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from owner to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from owner to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from administrator to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from administrator to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from member to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed the position of %[2]s from member to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	default:
		return i18n.Localf("The owner %[1]s changed the position of %[2]s.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
}

func getAffiliationChangedMessageWithAdminActor(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from owner to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsOwner() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from owner to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from administrator to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsAdmin() && affiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from administrator to member.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from member to owner.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	case affiliationUpdate.Previous.IsMember() && affiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s from member to administrator.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	default:
		return i18n.Localf("The administrator %[1]s changed the position of %[2]s.",
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}
}

func getRoleUpdateMessage(roleUpdate data.RoleUpdate) string {
	if roleUpdate.New.IsNone() {
		return getRoleRemovedMessage(roleUpdate)
	}
	return getRoleChangedMessage(roleUpdate)
}

func getRoleRemovedMessage(roleUpdate data.RoleUpdate) string {
	m := getRoleRemovedMessageWithoutActor(roleUpdate)
	if roleUpdate.Actor != nil {
		m = getRoleRemovedMessageWithActor(roleUpdate)
	}

	return appendReasonToMessage(m, roleUpdate.Reason)
}

func getRoleRemovedMessageWithoutActor(roleUpdate data.RoleUpdate) string {
	return i18n.Localf("%s was temporarily removed from the room.", roleUpdate.Nickname)
}

func getRoleRemovedMessageWithActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Actor.Affiliation.IsOwner():
		return i18n.Localf("The owner %s temporarily removed %s from the room.", roleUpdate.Actor.Nickname, roleUpdate.Nickname)
	case roleUpdate.Actor.Affiliation.IsAdmin():
		return i18n.Localf("The administrator %s temporarily removed %s from the room.", roleUpdate.Actor.Nickname, roleUpdate.Nickname)
	default:
		return i18n.Localf("%s temporarily removed %s from the room.", roleUpdate.Actor.Nickname, roleUpdate.Nickname)
	}
}

func getRoleChangedMessage(roleUpdate data.RoleUpdate) string {
	m := getRoleChangedMessageWithoutActor(roleUpdate)
	if roleUpdate.Actor != nil {
		m = getRoleChangedMessageWithActor(roleUpdate)
	}

	return appendReasonToMessage(m, roleUpdate.Reason)
}

func getRoleChangedMessageWithoutActor(roleUpdate data.RoleUpdate) string {
	return i18n.Localf("The role of %s was changed from %s to %s.",
		roleUpdate.Nickname,
		displayNameForRole(roleUpdate.Previous),
		displayNameForRole(roleUpdate.New))
}

func getRoleChangedMessageWithActor(roleUpdate data.RoleUpdate) string {
	switch {
	case roleUpdate.Actor.Affiliation.IsOwner():
		return i18n.Localf("The owner %s changed the role of %s from %s to %s.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname,
			displayNameForRole(roleUpdate.Previous),
			displayNameForRole(roleUpdate.New))
	case roleUpdate.Actor.Affiliation.IsAdmin():
		return i18n.Localf("The administrator %s changed the role of %s from %s to %s.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname,
			displayNameForRole(roleUpdate.Previous),
			displayNameForRole(roleUpdate.New))
	default:
		return i18n.Localf("%s changed the role of %s from %s to %s.",
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname,
			displayNameForRole(roleUpdate.Previous),
			displayNameForRole(roleUpdate.New))
	}
}

func getSelfRoleUpdateMessage(selfRoleUpdate data.RoleUpdate) string {
	m := getSelfRoleUpdateMessageWithoutActor(selfRoleUpdate)
	if selfRoleUpdate.Actor != nil {
		m = getSelfRoleUpdateMessageWithActor(selfRoleUpdate)
	}

	return appendReasonToMessage(m, selfRoleUpdate.Reason)
}

func getSelfRoleUpdateMessageWithoutActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Local("Your role was changed from participant to visitor.")
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Local("Your role was changed from participant to moderator.")
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Local("Your role was changed from visitor to participant.")
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Local("Your role was changed from visitor to moderator.")
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Local("Your role was changed from moderator to visitor.")
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Local("Your role was changed from moderator to participant.")
	default:
		return i18n.Local("Your role was changed.")
	}
}

func getSelfRoleUpdateMessageWithActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Actor.Affiliation.IsOwner():
		return getSelfRoleUpdateMessageForOwnerActor(selfRoleUpdate)
	case selfRoleUpdate.Actor.Affiliation.IsAdmin():
		return getSelfRoleUpdateMessageForAdminActor(selfRoleUpdate)
	default:
		return getSelfRoleUpdateMessageForActor(selfRoleUpdate)
	}
}

func getSelfRoleUpdateMessageForOwnerActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The owner %s changed your role from participant to visitor.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The owner %s changed your role from participant to moderator.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The owner %s changed your role from visitor to participant.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The owner %s changed your role from visitor to moderator.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The owner %s changed your role from moderator to visitor.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The owner %s changed your role from moderator to participant.", selfRoleUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The owner %s changed your role.", selfRoleUpdate.Actor.Nickname)
	}
}

func getSelfRoleUpdateMessageForAdminActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The administrator %s changed your role from participant to visitor.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The administrator %s changed your role from participant to moderator.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The administrator %s changed your role from visitor to participant.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("The administrator %s changed your role from visitor to moderator.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("The administrator %s changed your role from moderator to visitor.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("The administrator %s changed your role from moderator to participant.", selfRoleUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The administrator %s changed your role.", selfRoleUpdate.Actor.Nickname)
	}
}

func getSelfRoleUpdateMessageForActor(selfRoleUpdate data.RoleUpdate) string {
	switch {
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("%s changed your role from participant to visitor.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsParticipant() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("%s changed your role from participant to moderator.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("%s changed your role from visitor to participant.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsVisitor() && selfRoleUpdate.New.IsModerator():
		return i18n.Localf("%s changed your role from visitor to moderator.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsVisitor():
		return i18n.Localf("%s changed your role from moderator to visitor.", selfRoleUpdate.Actor.Nickname)
	case selfRoleUpdate.Previous.IsModerator() && selfRoleUpdate.New.IsParticipant():
		return i18n.Localf("%s changed your role from moderator to participant.", selfRoleUpdate.Actor.Nickname)
	default:
		return i18n.Localf("%s changed your role.", selfRoleUpdate.Actor.Nickname)
	}
}

func getAffiliationRoleUpdateMessage(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	affiliationMessage := getAffiliationRoleUpdateMessageForAffiliation(affiliationRoleUpdate)
	roleMessage := getAffiliationRoleUpdateMessageForRole(affiliationRoleUpdate)

	return appendReasonToMessage(
		i18n.Localf("%s %s", affiliationMessage, roleMessage),
		affiliationRoleUpdate.Reason,
	)
}

func getAffiliationRoleUpdateMessageForAffiliation(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsNone():
		return getAffiliationRoleUpateForAffiliationRemoved(affiliationRoleUpdate)
	case affiliationRoleUpdate.PreviousAffiliation.IsNone():
		return getAffiliationRoleUpdateForAffiliationAdded(affiliationRoleUpdate)
	default:
		return getAffiliationRoleUpdateForAffiliationUpdated(affiliationRoleUpdate)
	}
}

func getAffiliationRoleUpdateMessageForRole(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, their role was changed from visitor to participant.")
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, their role was changed from visitor to moderator.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, their role was changed from participant to visitor.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, their role was changed from participant to moderator.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, their role was changed from moderator to visitor.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, their role was changed from moderator to participant.")
	default:
		return i18n.Local("As a result, their role was also changed.")
	}
}

func getSelfAffiliationRoleUpdateMessageForRole(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, your role was changed from visitor to participant.")
	case affiliationRoleUpdate.PreviousRole.IsVisitor() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, your role was changed from visitor to moderator.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, your role was changed from participant to visitor.")
	case affiliationRoleUpdate.PreviousRole.IsParticipant() && affiliationRoleUpdate.NewRole.IsModerator():
		return i18n.Local("As a result, your role was changed from participant to moderator.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsVisitor():
		return i18n.Local("As a result, your role was changed from moderator to visitor.")
	case affiliationRoleUpdate.PreviousRole.IsModerator() && affiliationRoleUpdate.NewRole.IsParticipant():
		return i18n.Local("As a result, your role was changed from moderator to participant.")
	default:
		return i18n.Local("As a result, your role was also changed.")
	}
}

func getAffiliationRoleUpateForAffiliationRemoved(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		return getAffiliationRoleUpateForAffiliationRemovedWithoutActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpateForAffiliationRemovedWithActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpateForAffiliationRemovedWithoutActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("%s is not an owner anymore.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("%s is not an administrator anymore.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("%s is not a member anymore.", affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("%s is not banned anymore.", affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpateForAffiliationRemovedWithActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getAffiliationRoleUpateForAffiliationRemovedWithOwnerActor(affiliationRoleUpdate)
	case affiliationRoleUpdate.Actor.Affiliation.IsAdmin():
		return getAffiliationRoleUpateForAffiliationRemovedWithAdminActor(affiliationRoleUpdate)
	default:
		return getAffiliationRoleUpateForAffiliationRemovedForActor(affiliationRoleUpdate)
	}
}

func getAffiliationRoleUpateForAffiliationRemovedWithOwnerActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("The owner %s changed the position of %s; %s is not an owner anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("The owner %s changed the position of %s; %s is not an administrator anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("The owner %s changed the position of %s; %s is not a member anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The owner %s changed the position of %s; %s is not banned anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpateForAffiliationRemovedWithAdminActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("The administrator %s changed the position of %s; %s is not an owner anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("The administrator %s changed the position of %s; %s is not an administrator anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("The administrator %s changed the position of %s; %s is not a member anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The administrator %s changed the position of %s; %s is not banned anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpateForAffiliationRemovedForActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Localf("%s changed the position of %s; %s is not an owner anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Localf("%s changed the position of %s; %s is not an administrator anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember():
		return i18n.Localf("%s changed the position of %s; %s is not a member anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("%s changed the position of %s; %s is not banned anymore.",
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationAdded(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		return getAffiliationRoleUpdateForAffiliationAddedWithoutActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpdateForAffiliationAddedWithActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpdateForAffiliationAddedWithoutActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The position of %s was changed to owner.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The position of %s was changed to administrator.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The position of %s was changed to member.", affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The position of %s was changed.", affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationAddedWithActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getAffiliationRoleUpdateForAffiliationAddedWithOwnerActor(affiliationRoleUpdate)
	case affiliationRoleUpdate.Actor.Affiliation.IsAdmin():
		return getAffiliationRoleUpdateForAffiliationAddedWithAdminActor(affiliationRoleUpdate)
	default:
		return getAffiliationRoleUpdateForAffiliationAddedForActor(affiliationRoleUpdate)
	}
}

func getAffiliationRoleUpdateForAffiliationAddedWithOwnerActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner %s changed the position of %s to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner %s changed the position of %s to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner %s changed the position of %s to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The owner %s changed the position of %s.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationAddedWithAdminActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator %s changed the position of %s to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator %s changed the position of %s to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator %s changed the position of %s to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The administrator %s changed the position of %s.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationAddedForActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("%s changed the position of %s to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("%s changed the position of %s to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("%s changed the position of %s to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("%s changed the position of %s.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationUpdated(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		return getAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(affiliationRoleUpdate)
	}
	return getAffiliationRoleUpdateForAffiliationUpdatedWithActor(affiliationRoleUpdate)
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The position of %s was changed from owner to member.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The position of %s was changed from administrator to member.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The position of %s was changed from owner to administrator.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The position of %s was changed from member to administrator.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The position of %s was changed from administrator to owner.", affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The position of %s was changed from member to owner.", affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The position of %s was changed.", affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(affiliationRoleUpdate)
	case affiliationRoleUpdate.Actor.Affiliation.IsAdmin():
		return getAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(affiliationRoleUpdate)
	default:
		return getAffiliationRoleUpdateForAffiliationUpdatedForActor(affiliationRoleUpdate)
	}
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner %s changed the position of %s from member to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner %s changed the position of %s from member to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner %s changed the position of %s from administrator to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner %s changed the position of %s from administrator to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner %s changed the position of %s from owner to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner %s changed the position of %s from owner to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The owner %s changed the position of %s.", affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator %s changed the position of %s from member to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator %s changed the position of %s from member to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator %s changed the position of %s from administrator to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator %s changed the position of %s from administrator to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The administrator %s changed the position of %s from owner to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator %s changed the position of %s from owner to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("The administrator %s changed the position of %s.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	}
}

func getAffiliationRoleUpdateForAffiliationUpdatedForActor(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("%s changed the position of %s from member to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsMember() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("%s changed the position of %s from member to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("%s changed the position of %s from administrator to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsAdmin() && affiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("%s changed the position of %s from administrator to owner.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("%s changed the position of %s from owner to member.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	case affiliationRoleUpdate.PreviousAffiliation.IsOwner() && affiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("%s changed the position of %s from owner to administrator.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	default:
		return i18n.Localf("%s changed the position of %s.", affiliationRoleUpdate.Actor.Nickname, affiliationRoleUpdate.Nickname)
	}
}

func getSelfAffiliationRoleUpdateMessage(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	affiliationMessage := getSelfAffiliationRoleUpdateMessageForAffiliation(selfAffiliationRoleUpdate)
	roleMessage := getSelfAffiliationRoleUpdateMessageForRole(selfAffiliationRoleUpdate)

	return appendReasonToMessage(
		i18n.Localf("%s %s", affiliationMessage, roleMessage),
		selfAffiliationRoleUpdate.Reason,
	)
}

func getSelfAffiliationRoleUpdateMessageForAffiliation(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.NewAffiliation.IsNone():
		return getSelfAffiliationRoleUpateForAffiliationRemoved(selfAffiliationRoleUpdate)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsNone():
		return getSelfAffiliationRoleUpdateForAffiliationAdded(selfAffiliationRoleUpdate)
	default:
		return getSelfAffiliationRoleUpdateForAffiliationUpdated(selfAffiliationRoleUpdate)
	}
}

func getSelfAffiliationRoleUpateForAffiliationRemoved(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if selfAffiliationRoleUpdate.Actor == nil {
		return getSelfAffiliationRoleUpateForAffiliationRemovedWithoutActor(selfAffiliationRoleUpdate)
	}

	return getSelfAffiliationRoleUpateForAffiliationRemovedWithActor(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpateForAffiliationRemovedWithoutActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner():
		return i18n.Local("You are not an owner anymore.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		return i18n.Local("You are not an administrator anymore.")
	default:
		return i18n.Local("You are not a member anymore.")

	}
}

func getSelfAffiliationRoleUpateForAffiliationRemovedWithActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) (m string) {
	switch {
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner():
		m = i18n.Localf("The owner %s changed your position; you are not an owner anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		m = i18n.Localf("The owner %s changed your position; you are not an administrator anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.PreviousAffiliation.IsMember():
		m = i18n.Localf("The owner %s changed your position; you are not a member anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner():
		m = i18n.Localf("The administrator %s changed your position; you are not an owner anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin():
		m = i18n.Localf("The administrator %s changed your position; you are not an administrator anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.PreviousAffiliation.IsMember():
		m = i18n.Localf("The administrator %s changed your position; you are not a member anymore.",
			selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return m
}

func getSelfAffiliationRoleUpdateForAffiliationAdded(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if selfAffiliationRoleUpdate.Actor == nil {
		return getSelfAffiliationRoleUpdateForAffiliationAddedWithoutActor(selfAffiliationRoleUpdate)
	}
	return getSelfAffiliationRoleUpdateForAffiliationAddedWithActor(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpdateForAffiliationAddedWithoutActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("Your position was changed to owner.")
	case selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("Your position was changed to administrator.")
	default:
		return i18n.Localf("Your position was changed to member.")

	}

}

func getSelfAffiliationRoleUpdateForAffiliationAddedWithActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The owner %s changed your position to owner.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The owner %s changed your position to administrator.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		return i18n.Localf("The owner %s changed your position to member.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		return i18n.Localf("The administrator %s changed your position to owner.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		return i18n.Localf("The administrator %s changed your position to administrator.", selfAffiliationRoleUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The administrator %s changed your position to member.", selfAffiliationRoleUpdate.Actor.Nickname)
	}
}

func getSelfAffiliationRoleUpdateForAffiliationUpdated(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if selfAffiliationRoleUpdate.Actor == nil {
		return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(selfAffiliationRoleUpdate)
	}

	return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithActor(selfAffiliationRoleUpdate)
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithoutActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) (m string) {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		m = i18n.Localf("Your position was changed from owner to administrator.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		m = i18n.Localf("Your position was changed from owner to member.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		m = i18n.Localf("Your position was changed from administrator to owner.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		m = i18n.Localf("Your position was changed from administrator to member.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		m = i18n.Localf("Your position was changed from member to owner.")
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		m = i18n.Localf("Your position was changed from member to administrator.")
	}
	return m
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) (m string) {
	switch {
	case selfAffiliationRoleUpdate.Actor.Affiliation.IsOwner():
		return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(selfAffiliationRoleUpdate)
	default: // administrator
		return getSelfAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(selfAffiliationRoleUpdate)
	}
}

func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithOwnerActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) (m string) {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		m = i18n.Localf("The owner %s changed your position from owner to administrator.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		m = i18n.Localf("The owner %s changed your position from owner to member.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		m = i18n.Localf("The owner %s changed your position from administrator to owner.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		m = i18n.Localf("The owner %s changed your position from administrator to member.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		m = i18n.Localf("The owner %s changed your position from member to owner.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		m = i18n.Localf("The owner %s changed your position from member to administrator.", selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return m
}
func getSelfAffiliationRoleUpdateForAffiliationUpdatedWithAdminActor(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) (m string) {
	switch {
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		m = i18n.Localf("The administrator %s changed your position from owner to administrator.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsOwner() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		m = i18n.Localf("The administrator %s changed your position from owner to member.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		m = i18n.Localf("The administrator %s changed your position from administrator to owner.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsAdmin() && selfAffiliationRoleUpdate.NewAffiliation.IsMember():
		m = i18n.Localf("The administrator %s changed your position from administrator to member.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsOwner():
		m = i18n.Localf("The administrator %s changed your position from member to owner.", selfAffiliationRoleUpdate.Actor.Nickname)
	case selfAffiliationRoleUpdate.PreviousAffiliation.IsMember() && selfAffiliationRoleUpdate.NewAffiliation.IsAdmin():
		m = i18n.Localf("The administrator %s changed your position from member to administrator.", selfAffiliationRoleUpdate.Actor.Nickname)
	}
	return m
}

func getSelfAffiliationUpdateMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) (m string) {
	switch {
	case selfAffiliationUpdate.New.IsNone():
		m = getSelfAffiliationRemovedMessage(selfAffiliationUpdate)
	case selfAffiliationUpdate.New.IsBanned():
		m = getSelfAffiliationBannedMessage(selfAffiliationUpdate)
	case selfAffiliationUpdate.Previous.IsNone():
		m = getSelfAffiliationAddedMessage(selfAffiliationUpdate)
	default:
		m = getSelfAffiliationChangedMessage(selfAffiliationUpdate)
	}

	return appendReasonToMessage(m, selfAffiliationUpdate.Reason)
}

func getSelfAffiliationRemovedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return getSelfAffiliationRemovedMessageWithoutActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationRemovedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationRemovedMessageWithoutActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner():
		return i18n.Local("You are not an owner anymore.")
	case selfAffiliationUpdate.Previous.IsAdmin():
		return i18n.Local("You are not an administrator anymore.")
	case selfAffiliationUpdate.Previous.IsMember():
		return i18n.Local("You are not a member anymore.")
	default:
		return i18n.Local("You are not banned anymore.")
	}
}

func getSelfAffiliationRemovedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return getSelfAffiliationRemovedMessageWithOwnerActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationRemovedMessageWithAdminActor(selfAffiliationUpdate)
}

func getSelfAffiliationRemovedMessageWithOwnerActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The owner %[1]s changed your position; you are not an owner anymore.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The owner %[1]s changed your position; you are not an administrator anymore.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember():
		return i18n.Localf("The owner %[1]s changed your position; you are not a member anymore.", selfAffiliationUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The owner %[1]s changed your position; you are not banned anymore.", selfAffiliationUpdate.Actor.Nickname)
	}
}

func getSelfAffiliationRemovedMessageWithAdminActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner():
		return i18n.Localf("The administrator %[1]s changed your position; you are not an owner anymore.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed your position; you are not an administrator anymore.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember():
		return i18n.Localf("The administrator %[1]s changed your position; you are not a member anymore.", selfAffiliationUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The administrator %[1]s changed your position; you are not banned anymore.", selfAffiliationUpdate.Actor.Nickname)
	}
}

func getSelfAffiliationAddedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return getSelfAffiliationAddedMessageWithoutActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationAddedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationAddedMessageWithoutActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.New.IsOwner():
		return i18n.Local("You are now an owner.")
	case selfAffiliationUpdate.New.IsAdmin():
		return i18n.Local("You are now an administrator.")
	case selfAffiliationUpdate.New.IsMember():
		return i18n.Local("You are now a member.")
	default:
		return i18n.Local("You are now banned.")
	}
}

func getSelfAffiliationAddedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return getSelfAffiliationAddedMessageWithOwnerActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationAddedMessageWithAdminActor(selfAffiliationUpdate)
}

func getSelfAffiliationAddedMessageWithOwnerActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed your position; you are now an owner.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed your position; you are now an administrator.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed your position; you are now an member.", selfAffiliationUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The owner %[1]s changed your position; you are now banned.", selfAffiliationUpdate.Actor.Nickname)
	}
}

func getSelfAffiliationAddedMessageWithAdminActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed your position; you are now an owner.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed your position; you are now an administrator.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed your position; you are now an member.", selfAffiliationUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The administrator %[1]s changed your position; you are now banned.", selfAffiliationUpdate.Actor.Nickname)
	}

}

func getSelfAffiliationChangedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return getSelfAffiliationChangedMessageWithoutActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationChangedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationChangedMessageWithoutActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Local("Your position was changed from owner to administrator.")
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsMember():
		return i18n.Local("Your position was changed from owner to member.")
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Local("Your position was changed from administrator to owner.")
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsMember():
		return i18n.Local("Your position was changed from administrator to member.")
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Local("Your position was changed from member to administrator.")
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Local("Your position was changed from member to owner.")
	default:
		return i18n.Local("Your position was changed.")
	}
}

func getSelfAffiliationChangedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return getSelfAffiliationChangedMessageWithOwnerActor(selfAffiliationUpdate)
	}
	return getSelfAffiliationChangedMessageWithAdminActor(selfAffiliationUpdate)
}

func getSelfAffiliationChangedMessageWithOwnerActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed your position from owner to administrator.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed your position from owner to member.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed your position from administrator to owner.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The owner %[1]s changed your position from administrator to member.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The owner %[1]s changed your position from member to administrator.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The owner %[1]s changed your position from member to owner.", selfAffiliationUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The owner %[1]s changed your position.", selfAffiliationUpdate.Actor.Nickname)
	}
}

func getSelfAffiliationChangedMessageWithAdminActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	switch {
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed your position from owner to administrator.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsOwner() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed your position from owner to member.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed your position from administrator to owner.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsAdmin() && selfAffiliationUpdate.New.IsMember():
		return i18n.Localf("The administrator %[1]s changed your position from administrator to member.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsAdmin():
		return i18n.Localf("The administrator %[1]s changed your position from member to administrator.", selfAffiliationUpdate.Actor.Nickname)
	case selfAffiliationUpdate.Previous.IsMember() && selfAffiliationUpdate.New.IsOwner():
		return i18n.Localf("The administrator %[1]s changed your position from member to owner.", selfAffiliationUpdate.Actor.Nickname)
	default:
		return i18n.Localf("The administrator %[1]s changed your position.", selfAffiliationUpdate.Actor.Nickname)
	}
}

func getSelfAffiliationBannedMessage(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor == nil {
		return i18n.Localf("You has been banned from the room.")
	}
	return getSelfAffiliationBannedMessageWithActor(selfAffiliationUpdate)
}

func getSelfAffiliationBannedMessageWithActor(selfAffiliationUpdate data.SelfAffiliationUpdate) string {
	if selfAffiliationUpdate.Actor.Affiliation.IsOwner() {
		return i18n.Localf("The owner %[1]s banned you from the room.", selfAffiliationUpdate.Actor.Nickname)
	}
	return i18n.Localf("The administrator %[1]s banned you from the room.", selfAffiliationUpdate.Actor.Nickname)
}

func appendReasonToMessage(message, reason string) string {
	if reason != "" {
		return i18n.Localf("%s The reason given was: %s.", message, reason)
	}
	return message
}
