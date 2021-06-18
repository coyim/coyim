package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
)

func initMUCRoomConfigTexts() {
	initMUCRoomConfigPagesTexts()
	initMUCRoomConfigFieldsTexts()
	initMUCRoomConfigOccupantFieldTexts()
}

type roomConfigPageTextInfo struct {
	displayTitle string
	displayIntro string
}

var roomConfigPagesTexts map[mucRoomConfigPageID]roomConfigPageTextInfo

func initMUCRoomConfigPagesTexts() {
	roomConfigPagesTexts = map[mucRoomConfigPageID]roomConfigPageTextInfo{
		roomConfigInformationPageIndex: {
			displayTitle: i18n.Local("Basic information"),
			displayIntro: i18n.Local("This section contains basic configuration options that you can " +
				"set for the room."),
		},
		roomConfigAccessPageIndex: {
			displayTitle: i18n.Local("Access"),
			displayIntro: i18n.Local("Here you can manage access to the room. If you specify a password " +
				"for the room, you should share it in a secure way. This will help to protect the " +
				"people in the room. Remember that room passwords do not make the room encrypted, so " +
				"people that control the location of the room might still have access to it, even " +
				"without providing a password."),
		},
		roomConfigPermissionsPageIndex: {
			displayTitle: i18n.Local("Permissions"),
			displayIntro: i18n.Local("Here you can change settings that impact who can do what inside " +
				"the room."),
		},
		roomConfigPositionsPageIndex: {
			displayTitle: i18n.Local("Positions"),
			displayIntro: i18n.Local("Here you can define who the owners, administrators and banned people are."),
		},
		roomConfigOthersPageIndex: {
			displayTitle: i18n.Local("Other settings"),
			displayIntro: i18n.Local("Here you can find other configuration options that might be useful " +
				"to you. Note that if archiving is enabled, all the discussions in the room might be logged " +
				"and potentially made publicly accessible."),
		},
		roomConfigSummaryPageIndex: {
			displayTitle: i18n.Local("Summary"),
		},
	}
}

func getRoomConfigPageTexts(pageID mucRoomConfigPageID) roomConfigPageTextInfo {
	if t, ok := roomConfigPagesTexts[pageID]; ok {
		return t
	}

	return roomConfigPageTextInfo{
		displayTitle: fmt.Sprintf("UnsupportedPage(%d)", pageID),
	}
}

func configPageDisplayTitle(pageID mucRoomConfigPageID) string {
	t := getRoomConfigPageTexts(pageID)
	return t.displayTitle
}

func configPageDisplayIntro(pageID mucRoomConfigPageID) string {
	t := getRoomConfigPageTexts(pageID)
	return t.displayIntro
}

type roomConfigFieldTextInfo struct {
	displayLabel       string
	summaryLabel       string
	displayDescription string
}

func newRoomConfigFieldTextInfo(label, description string) roomConfigFieldTextInfo {
	return roomConfigFieldTextInfo{
		displayLabel:       label,
		summaryLabel:       label,
		displayDescription: description,
	}
}

var roomConfigFieldsTexts map[muc.RoomConfigFieldType]roomConfigFieldTextInfo

func initMUCRoomConfigFieldsTexts() {
	roomConfigFieldsTexts = map[muc.RoomConfigFieldType]roomConfigFieldTextInfo{
		muc.RoomConfigFieldName: {
			displayLabel:       i18n.Local("Title"),
			summaryLabel:       i18n.Local("Title"),
			displayDescription: i18n.Local("The room title can be used to find the room in the public list."),
		},
		muc.RoomConfigFieldDescription: {
			displayLabel: i18n.Local("Description"),
			summaryLabel: i18n.Local("Description"),
			displayDescription: i18n.Local("The room description can be used to add more information " +
				"about the room, such as the purpose, the discussion topics, interests, etc."),
		},
		muc.RoomConfigFieldEnableLogging: {
			displayLabel: i18n.Local("Enable archiving of discussions"),
			summaryLabel: i18n.Local("Enable archiving of discussions"),
			displayDescription: i18n.Local("The conversation of this room will be stored in an " +
				"archive that could be accessed publicly. CoyIM users will be notified about this " +
				"when joining the room, however, other clients might not."),
		},
		muc.RoomConfigFieldLanguage: {
			displayLabel: i18n.Local("Primary language of discussion"),
			summaryLabel: i18n.Local("Primary language of discussion"),
			displayDescription: i18n.Local("This is the primary language in which conversations are " +
				"held. Changing this will not impact the language of the application."),
		},
		muc.RoomConfigFieldPubsub: {
			displayLabel: i18n.Local("XMPP URI of associated publish-subscribe node"),
			summaryLabel: i18n.Local("XMPP URI of associated publish-subscribe node"),
			displayDescription: i18n.Local("A chat room can have an associated place where publication " +
				"and subscription of certain information can happen. This is a technical setting, " +
				"which should be left empty unless you know what it means."),
		},
		muc.RoomConfigFieldCanChangeSubject: {
			displayLabel:       i18n.Local("Allow anyone to set the subject of the room"),
			summaryLabel:       i18n.Local("Occupants can change subject of the room"),
			displayDescription: i18n.Local("If not set, only moderators can modify it."),
		},
		muc.RoomConfigFieldAllowInvites: {
			displayLabel: i18n.Local("Allow members to invite others to the room"),
			summaryLabel: i18n.Local("Allow members to invite others"),
		},
		muc.RoomConfigFieldAllowPrivateMessages: {
			displayLabel: i18n.Local("People can send private messages"),
			summaryLabel: i18n.Local("People can send private messages"),
		},
		muc.RoomConfigFieldMaxOccupantsNumber: {
			displayLabel: i18n.Local("Maximum number of people in the room"),
			summaryLabel: i18n.Local("Maximum number of people in the room"),
		},
		muc.RoomConfigFieldIsPublic: {
			displayLabel: i18n.Local("Make this room public"),
			summaryLabel: i18n.Local("Include in public lists"),
			displayDescription: i18n.Local("A public room can be found by all users in any public " +
				"listing."),
		},
		muc.RoomConfigFieldIsPersistent: {
			displayLabel: i18n.Local("Make this room persistent"),
			summaryLabel: i18n.Local("Persistent"),
			displayDescription: i18n.Local("A persistent room won't be destroyed when the last " +
				"occupant leaves the room."),
		},
		muc.RoomConfigFieldPresenceBroadcast: {
			displayLabel: i18n.Local("Roles that will receive information about other people in the room"),
			summaryLabel: i18n.Local("Roles that will receive information about other people in the room"),
		},
		muc.RoomConfigFieldIsModerated: {
			displayLabel: i18n.Local("Make this room moderated"),
			summaryLabel: i18n.Local("Moderated"),
			displayDescription: i18n.Local("In a moderated room, visitors must be given permission " +
				"to speak."),
		},
		muc.RoomConfigFieldIsMembersOnly: {
			displayLabel: i18n.Local("Make this room members-only"),
			summaryLabel: i18n.Local("Only allow members to join"),
		},
		muc.RoomConfigFieldMembers: {
			displayLabel: i18n.Local("Members"),
		},
		muc.RoomConfigFieldPassword: {
			displayLabel: i18n.Local("Enter the room password"),
			summaryLabel: i18n.Local("Password"),
		},
		muc.RoomConfigFieldOwners: {
			displayLabel: i18n.Local("Owners"),
		},
		muc.RoomConfigFieldWhoIs: {
			displayLabel: i18n.Local("The account address of others in the room may be viewed by:"),
			summaryLabel: i18n.Local("Account address of occupants can be viewed by"),
		},
		muc.RoomConfigFieldMaxHistoryFetch: {
			displayLabel: i18n.Local("Maximum previous messages sent to people when joining the room"),
			summaryLabel: i18n.Local("Maximum previous messages sent to people when joining the room"),
		},
		muc.RoomConfigFieldAdmins: {
			displayLabel: i18n.Local("Administrators"),
		},
		muc.RoomConfigFieldAllowQueryUsers: {
			displayLabel: i18n.Local("A person can find out more information about other people in the room"),
			summaryLabel: i18n.Local("A person can find out more information about other people in the room"),
		},
		muc.RoomConfigFieldAllowVisitorStatus: {
			displayLabel: i18n.Local("Visitors to a room can change their custom status information"),
			summaryLabel: i18n.Local("Visitors to a room can change their custom status information"),
			displayDescription: i18n.Local("Since visitors are not allowed to speak in moderated rooms, " +
				"this setting would allow visitors to still communicate with others by means of their status text."),
		},
		muc.RoomConfigFieldAllowVisitorNickchange: {
			displayLabel: i18n.Local("Visitors can change their nickname"),
			summaryLabel: i18n.Local("Visitors can change their nickname"),
			displayDescription: i18n.Local("Just as with the status updates, this setting could allow " +
				"visitors to show information to others in the room by changing their nickname."),
		},
		muc.RoomConfigFieldAllowVoiceRequest: {
			displayLabel: i18n.Local("Visitors can ask for permission to speak"),
			summaryLabel: i18n.Local("Visitors can ask for permission to speak"),
		},
		muc.RoomConfigFieldAllowSubscription: {
			displayLabel: i18n.Local("Allow subscription"),
			summaryLabel: i18n.Local("Allow subscription"),
			displayDescription: i18n.Local("This is a technical setting which allows a new, experimental, version " +
				"of group chat to be used with this room. CoyIM does not support this, but if you turn on this setting, " +
				"other clients might utilize the setting when connecting to the room."),
		},
		muc.RoomConfigFieldMembersByDefault: {
			displayLabel: i18n.Local("When joining the room, people will be able to talk"),
			summaryLabel: i18n.Local("When joining the room, people will be able to talk"),
		},
		muc.RoomConfigFieldVoiceRequestMinInteval: {
			displayLabel: i18n.Local("Minimum interval between voice requests (in seconds)"),
			summaryLabel: i18n.Local("Minimum interval between voice requests (in seconds)"),
		},
		muc.RoomConfigAllowPrivateMessagesFromVisitors: {
			displayLabel: i18n.Local("Allow visitors to send private messages to"),
			summaryLabel: i18n.Local("Allow visitors to send private messages to"),
		},
		muc.RoomConfigPublicList: {
			displayLabel: i18n.Local("Make list of participants public"),
			summaryLabel: i18n.Local("Make list of participants public"),
		},
	}
}

type roomConfigOccupantFieldText struct {
	headerLabel       string
	descriptionLabel  string
	dialogTitle       string
	dialogDescription string
}

var roomConfigOccupantFieldTexts map[data.Affiliation]roomConfigOccupantFieldText

func initMUCRoomConfigOccupantFieldTexts() {
	roomConfigOccupantFieldTexts = map[data.Affiliation]roomConfigOccupantFieldText{
		&data.OwnerAffiliation{}: {
			headerLabel:       i18n.Local("Owners"),
			descriptionLabel:  i18n.Local("Owners will always be moderators in a room. They can give or take away any position or role and control any aspect of the room."),
			dialogTitle:       i18n.Local("Owners"),
			dialogDescription: i18n.Local("Here you can add one or more new owners to the room. You will have to use the account address of the user in order to make them an owner. This address can either be a simple one, such as user@example.org or a full one, such as user@example.org/abcdef."),
		},
		&data.AdminAffiliation{}: {
			headerLabel:       i18n.Local("Administrators"),
			descriptionLabel:  i18n.Local("A room administrator will automatically become a moderator when entering the room. They can't change the room configuration or destroy the room."),
			dialogTitle:       i18n.Local("Administrators"),
			dialogDescription: i18n.Local("Here you can add one or more new administrators to the room. You will have to use the account address of the user in order to make them an administrator. This address can either be a simple one, such as user@example.org or a full one, such as user@example.org/abcdef."),
		},
		&data.OutcastAffiliation{}: {
			headerLabel:       i18n.Local("Banned people"),
			descriptionLabel:  i18n.Local("A banned person will not be able to join the room."),
			dialogTitle:       i18n.Local("Ban people"),
			dialogDescription: i18n.Local("Here you can add one or more people who will not be able to join this room. You will have to use the account address of the user in order to ban them. This address can either be a simple one, such as user@example.org or a full one, such as user@example.org/abcdef."),
		},
	}
}
