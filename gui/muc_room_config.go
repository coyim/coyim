package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
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

var roomConfigFriendlyMessages map[data.RoomConfigType]func(data.RoomDiscoInfo) string

func initMUCConfigUpdateMessages() {
	roomConfigFriendlyMessages = map[data.RoomConfigType]func(data.RoomDiscoInfo) string{
		data.RoomConfigSupportsVoiceRequests:     roomConfigSupportsVoiceRequests,
		data.RoomConfigAllowsRegistration:        roomConfigAllowsRegistration,
		data.RoomConfigPersistent:                roomConfigPersistent,
		data.RoomConfigModerated:                 roomConfigModerated,
		data.RoomConfigPasswordProtected:         roomConfigPasswordProtected,
		data.RoomConfigPublic:                    roomConfigPublic,
		data.RoomConfigLanguage:                  roomConfigLanguage,
		data.RoomConfigOccupantsCanChangeSubject: roomConfigOccupantsCanChangeSubject,
		data.RoomConfigTitle:                     roomConfigTitle,
		data.RoomConfigDescription:               roomConfigDescription,
		data.RoomConfigMembersCanInvite:          roomConfigMembersCanInvite,
		data.RoomConfigAllowPrivateMessages:      roomConfigAllowPrivateMessages,
		data.RoomConfigMaxHistoryFetch:           roomConfigMaxHistoryFetch,
	}
}

func getRoomConfigUpdatedFriendlyMessages(changes roomConfigChangedTypes, discoInfo data.RoomDiscoInfo) []string {
	messages := []string{}

	for k, f := range roomConfigFriendlyMessages {
		if changes.contains(k) {
			messages = append(messages, f(discoInfo))
		}
	}

	return messages
}

func roomConfigSupportsVoiceRequests(di data.RoomDiscoInfo) string {
	if di.SupportsVoiceRequests {
		return i18n.Local("The room's occupants with no \"voice\" can now " +
			"request permission to speak.")
	}
	return i18n.Local("This room's doesn't support \"voice\" request " +
		"now, which means that visitors can't ask for permission to speak.")
}

func roomConfigAllowsRegistration(di data.RoomDiscoInfo) string {
	if di.AllowsRegistration {
		return i18n.Local("This room supports user registration.")
	}
	return i18n.Local("This room doesn't support user registration.")
}

func roomConfigPersistent(di data.RoomDiscoInfo) string {
	if di.Persistent {
		return i18n.Local("This room is now persistent.")
	}
	return i18n.Local("This room is not persistent anymore.")
}

func roomConfigModerated(di data.RoomDiscoInfo) string {
	if di.Moderated {
		return i18n.Local("Only participants and moderators can now send messages in this room.")
	}
	return i18n.Local("Everyone can now send messages in this room.")
}

func roomConfigPasswordProtected(di data.RoomDiscoInfo) string {
	if di.PasswordProtected {
		return i18n.Local("This room is now protected by a password.")
	}
	return i18n.Local("This room is not protected by a password.")
}

func roomConfigPublic(di data.RoomDiscoInfo) string {
	if di.Public {
		return i18n.Local("This room is publicly listed.")
	}
	return i18n.Local("This room is not publicly listed anymore.")
}

func roomConfigLanguage(di data.RoomDiscoInfo) string {
	return i18n.Localf("The language of this room was changed to %s", getLanguage(di.Language))
}

func getLanguage(languageCode string) string {
	languageTag, _ := language.Parse(languageCode)
	l := display.Self.Name(languageTag)
	if l == "" {
		return languageCode
	}
	return l
}

func roomConfigOccupantsCanChangeSubject(config data.RoomDiscoInfo) string {
	if config.OccupantsCanChangeSubject {
		return i18n.Local("Occupants can now change the subject of this room.")
	}
	return i18n.Local("Occupants cannot change the subject of this room.")
}

func roomConfigTitle(config data.RoomDiscoInfo) string {
	return i18n.Localf("Title was changed to \"%s\".", config.Title)
}

func roomConfigDescription(config data.RoomDiscoInfo) string {
	return i18n.Localf("Description was changed to \"%s\".", config.Description)
}

func roomConfigAllowPrivateMessages(config data.RoomDiscoInfo) string {
	if config.AllowPrivateMessages == data.RoleNone {
		return i18n.Local("No one in this room can send private messages now.")
	}
	return i18n.Localf("Only the \"%s\" can send private messages "+
		"to the room's occupants.", rolePluralName(config.AllowPrivateMessages))
}

func roomConfigMembersCanInvite(config data.RoomDiscoInfo) string {
	if config.MembersCanInvite {
		return i18n.Local("Members can now invite other users to join.")
	}
	return i18n.Local("Members cannot invite other users to join anymore.")
}

func roomConfigMaxHistoryFetch(config data.RoomDiscoInfo) string {
	return i18n.Localf("The room's max history length was changed to %d", config.MaxHistoryFetch)
}
