package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
)

func initMUCRoomConfigTexts() {
	initMUCRoomConfigPagesTexts()
	initMUCRoomConfigFieldsTexts()
}

type roomConfigPageTextInfo struct {
	displayTitle string
	displayIntro string
}

var roomConfigPagesTexts map[string]roomConfigPageTextInfo

func initMUCRoomConfigPagesTexts() {
	roomConfigPagesTexts = map[string]roomConfigPageTextInfo{
		pageConfigInfo: {
			displayTitle: i18n.Local("Basic information"),
			displayIntro: i18n.Local("This section contains basic configuration options that you can " +
				"set for the room."),
		},
		pageConfigAccess: {
			displayTitle: i18n.Local("Access"),
			displayIntro: i18n.Local("Here you can manage access to the room. If you specify a password " +
				"for the room, you should share it in a secure way. This will help to protect the " +
				"people in the room. Remember that room passwords do not make the room encrypted, so " +
				"people that control the location of the room might still have access to it, even " +
				"without providing a password."),
		},
		pageConfigPermissions: {
			displayTitle: i18n.Local("Permissions"),
			displayIntro: i18n.Local("Here you can change settings that impact who can do what inside " +
				"the room."),
		},
		pageConfigOccupants: {
			displayTitle: i18n.Local("Occupants"),
			displayIntro: i18n.Local("Here you can define who the owners and administrators are. " +
				"Owners will always be moderators in a room. They can give or take away any position " +
				"or role and control any aspect of the room. A room administrator will automatically " +
				"become a moderator when entering the room. They can grant or revoke memberships, and " +
				"also ban or unban people from a room. An administrator can't change the room configuration " +
				"or destroy the room."),
		},
		pageConfigOthers: {
			displayTitle: i18n.Local("Other settings"),
			displayIntro: i18n.Local("Here you can find other configuration options that might be useful " +
				"to you. Note that if archiving is enabled, all the discussions in the room might be logged " +
				"and potentially made publicly accessible."),
		},
		pageConfigSummary: {
			displayTitle: i18n.Local("Summary"),
		},
	}
}

func getRoomConfigPageTexts(pageID string) roomConfigPageTextInfo {
	if t, ok := roomConfigPagesTexts[pageID]; ok {
		return t
	}

	return roomConfigPageTextInfo{
		displayTitle: fmt.Sprintf("UnsupportedPage(%s)", pageID),
	}
}

func configPageDisplayTitle(pageID string) string {
	t := getRoomConfigPageTexts(pageID)
	return t.displayTitle
}

type roomConfigFieldTextInfo struct {
	displayLabel       string
	displayDescription string
}

var roomConfigFieldsTexts map[muc.RoomConfigFieldType]roomConfigFieldTextInfo

func initMUCRoomConfigFieldsTexts() {
	roomConfigFieldsTexts = map[muc.RoomConfigFieldType]roomConfigFieldTextInfo{
		muc.RoomConfigFieldName: {
			displayLabel:       i18n.Local("Title"),
			displayDescription: i18n.Local("The room title can be used to find the room in the public list."),
		},
		muc.RoomConfigFieldDescription: {
			displayLabel: i18n.Local("Description"),
			displayDescription: i18n.Local("The room description can be used to add more information " +
				"about the room, such as the purpose, the discussion topics, interests, etc."),
		},
		muc.RoomConfigFieldEnableLogging: {
			displayLabel: i18n.Local("Enable archiving of discussions"),
			displayDescription: i18n.Local("The conversation of this room will be stored in an " +
				"archive that could be accessed publicly. CoyIM users will be notified about this " +
				"when enter in the room, other client might not."),
		},
		muc.RoomConfigFieldLanguage: {
			displayLabel: i18n.Local("Primary language of discussion"),
			displayDescription: i18n.Local("This is the primary language in which conversations are " +
				"held. Changing this will not impact the language of the application."),
		},
		muc.RoomConfigFieldPubsub: {
			displayLabel: i18n.Local("XMPP URI of associated publish-subscribe node"),
			displayDescription: i18n.Local("A chat room can have an associated place where publication " +
				"and subscription of certain information can happen. This is a technical setting, " +
				"which should be left empty unless you know what it means."),
		},
		muc.RoomConfigFieldCanChangeSubject: {
			displayLabel:       i18n.Local("Allow anyone to set the room's subject"),
			displayDescription: i18n.Local("If not set, only moderators can modify it."),
		},
		muc.RoomConfigFieldAllowInvites: {
			displayLabel: i18n.Local("Allow members to invite others to the room"),
		},
		muc.RoomConfigFieldAllowPrivateMessages: {
			displayLabel: i18n.Local("Private messages to others in the room can be sent by:"),
		},
		muc.RoomConfigFieldMaxOccupantsNumber: {
			displayLabel: i18n.Local("Maximum number of people in the room"),
		},
		muc.RoomConfigFieldIsPublic: {
			displayLabel: i18n.Local("Make this room public"),
			displayDescription: i18n.Local("A public room can be found by all users in any public " +
				"listing."),
		},
		muc.RoomConfigFieldIsPersistent: {
			displayLabel: i18n.Local("Make this room persistent"),
			displayDescription: i18n.Local("A persistent room won't be destroyed when the last " +
				"occupant leaves the room."),
		},
		muc.RoomConfigFieldPresenceBroadcast: {
			displayLabel: i18n.Local("What roles will receive information about other people in the room:"),
		},
		muc.RoomConfigFieldIsModerated: {
			displayLabel: i18n.Local("Make this room moderated"),
			displayDescription: i18n.Local("In a moderated room, visitors must be given permission " +
				"to speak."),
		},
		muc.RoomConfigFieldIsMembersOnly: {
			displayLabel: i18n.Local("Make this room members-only"),
		},
		muc.RoomConfigFieldMembers: {
			displayLabel: i18n.Local("Members"),
		},
		muc.RoomConfigFieldIsPasswordProtected: {
			displayLabel: i18n.Local("Make this room password protected"),
		},
		muc.RoomConfigFieldPassword: {
			displayLabel: i18n.Local("Enter the room password"),
		},
		muc.RoomConfigFieldOwners: {
			displayLabel: i18n.Local("Owners"),
		},
		muc.RoomConfigFieldWhoIs: {
			displayLabel: i18n.Local("The account address of others in the room may be viewed by:"),
		},
		muc.RoomConfigFieldMaxHistoryFetch: {
			displayLabel: i18n.Local("Maximum previous messages sent to people when joining the room"),
		},
		muc.RoomConfigFieldAdmins: {
			displayLabel: i18n.Local("Administrators"),
		},
	}
}
