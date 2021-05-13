package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
)

func initMUCRoomConfigTexts() {
	initMUCRoomConfigPagesTexts()
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
