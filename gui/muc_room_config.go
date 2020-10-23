package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

type roomConfigChangedTypes []data.RoomConfigType

func (c roomConfigChangedTypes) contains(k data.RoomConfigType) bool {
	for _, kk := range c {
		if kk == k {
			return true
		}
	}
	return false
}

var roomConfigFriendlyMessages map[data.RoomConfigType]func(data.RoomConfig) string

func initMUCConfigUpdateMessages() {
	roomConfigFriendlyMessages = map[data.RoomConfigType]func(data.RoomConfig) string{
		data.RoomConfigSupportsVoiceRequests:     roomConfigSupportsVoiceRequests,
		data.RoomConfigAllowsRegistration:        roomConfigAllowsRegistration,
		data.RoomConfigPersistent:                roomConfigPersistent,
		data.RoomConfigModerated:                 roomConfigModerated,
		data.RoomConfigOpen:                      roomConfigOpen,
		data.RoomConfigPasswordProtected:         roomConfigPasswordProtected,
		data.RoomConfigPublic:                    roomConfigPublic,
		data.RoomConfigLanguage:                  roomConfigLanguage,
		data.RoomConfigOccupantsCanChangeSubject: roomConfigOccupantsCanChangeSubject,
		data.RoomConfigTitle:                     roomConfigTitle,
		data.RoomConfigDescription:               roomConfigDescription,
		data.RoomConfigOccupants:                 roomConfigOccupants,
		data.RoomConfigOccupantsCanInvite:        roomConfigOccupantsCanInvite,
		data.RoomConfigAllowPrivateMessages:      roomConfigAllowPrivateMessages,
		data.RoomConfigLogged:                    roomConfigLogged,
	}
}

func getRoomConfigUpdatedFriendlyMessages(changes roomConfigChangedTypes, config data.RoomConfig) []string {
	messages := []string{}

	for k, f := range roomConfigFriendlyMessages {
		if changes.contains(k) {
			messages = append(messages, f(config))
		}
	}

	return messages
}

func roomConfigSupportsVoiceRequests(config data.RoomConfig) string {
	if config.SupportsVoiceRequests {
		return i18n.Local("The room's occupants with no \"voice\" can now " +
			"request permission to speak.")
	}
	return i18n.Local("This room's doesn't support \"voice\" request " +
		"now, which means that visitors can't ask for permission to speak.")
}

func roomConfigAllowsRegistration(config data.RoomConfig) string {
	if config.AllowsRegistration {
		return i18n.Local("This room is not open anymore, which means " +
			"a user can't enter without being on the member list.")
	}
	return i18n.Local("The room is open now, which means any user " +
		"can enter without being on the member list.")
}

func roomConfigPersistent(config data.RoomConfig) string {
	if config.Persistent {
		return i18n.Local("The room is now persistent, which means it will " +
			"not be destroyed if the last occupant exits.")
	}
	return i18n.Local("The room isn't persistent, which means it will be " +
		"destroyed when the last occupant exits.")
}

func roomConfigModerated(config data.RoomConfig) string {
	if config.Moderated {
		return i18n.Local("The room now only allows sending messages " +
			"to all occupants to those with \"voice\".")
	}
	return i18n.Local("The room is non-moderated now, which means " +
		"everyone can send messages to all occupants.")
}

func roomConfigOpen(config data.RoomConfig) string {
	if config.Open {
		return i18n.Local("Anyone (non-banned members) is allowed " +
			"now to enter this room without being on the member list.")
	}
	return i18n.Local("Only registered members are allowed now to enter " +
		"this room.")
}

func roomConfigPasswordProtected(config data.RoomConfig) string {
	if config.PasswordProtected {
		return i18n.Local("This room is now protected by a password, " +
			"which means a user can't enter without first providing the correct password.")
	}
	return i18n.Local("This room is not protected by password (unsecure).")
}

func roomConfigPublic(config data.RoomConfig) string {
	if config.Public {
		return i18n.Local("Anyone can find this room publicly now.")
	}
	return i18n.Local("This room is hidden now, which means no one " +
		"can find it through normal means such as searching and service discovery.")
}

func roomConfigLanguage(config data.RoomConfig) string {
	return i18n.Localf("Room's language of discussion was changed "+
		"to \"%s\".", config.Language)
}

func roomConfigOccupantsCanChangeSubject(config data.RoomConfig) string {
	if config.OccupantsCanChangeSubject {
		return i18n.Local("Occupants can now change this room's subject.")
	}
	return i18n.Local("Occupants can't change this room's subject.")
}

func roomConfigTitle(config data.RoomConfig) string {
	return i18n.Localf("Room's title was changed to \"%s\".", config.Title)
}

func roomConfigDescription(config data.RoomConfig) string {
	return i18n.Localf("Room's description was changed to \"%s\".", config.Description)
}

func roomConfigAllowPrivateMessages(config data.RoomConfig) string {
	if config.AllowPrivateMessages == data.RoleNone {
		return i18n.Local("No one in this room can send private messages now.")
	}
	return i18n.Localf("Only the \"%s\" can send private messages "+
		"to the room's occupants.", rolePluralName(config.AllowPrivateMessages))
}

func roomConfigOccupants(config data.RoomConfig) string {
	return i18n.Localf("The number of occupants in this room was "+
		"limited to \"%d\".", config.Occupants)
}

func roomConfigOccupantsCanInvite(config data.RoomConfig) string {
	if config.OccupantsCanInvite {
		return i18n.Local("Room's cccupants can now invite other users " +
			"to join it.")
	}
	return i18n.Local("Occupants are now unable to invite other users " +
		"to join this room.")
}

func roomConfigLogged(config data.RoomConfig) string {
	if config.Logged {
		return i18n.Local("The discussion in this room is logged " +
			"to a public archive now.")
	}
	return i18n.Local("The room's discussion is private now, " +
		"the service is not logging any message.")
}
