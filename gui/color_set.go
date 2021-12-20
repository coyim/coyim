package gui

type colorSet struct {
	rosterPeerBackground                      hexColor
	rosterPeerOfflineForeground               hexColor
	rosterPeerOnlineForeground                hexColor
	rosterGroupBackground                     hexColor
	rosterAccountOnlineBackground             hexColor
	rosterAccountOfflineBackground            hexColor
	conversationOutgoingUserForeground        hexColor
	conversationIncomingUserForeground        hexColor
	conversationOutgoingTextForeground        hexColor
	conversationIncomingTextForeground        hexColor
	conversationStatusTextForeground          hexColor
	conversationOutgoingDelayedUserForeground hexColor
	conversationOutgoingDelayedTextForeground hexColor

	// these two don't seem to be used anymore
	conversationLockTypingBackground   hexColor
	conversationUnlockTypingBackground hexColor

	timestampForeground hexColor
}

func (cm *hasColorManagement) currentColorSet() colorSet {
	if cm.isDarkThemeVariant() {
		return defaultDarkColorSet
	}
	return defaultLightColorSet
}

var defaultLightColorSet = colorSet{
	rosterPeerBackground:                      rgbFromHex("#ffffff"),
	rosterPeerOfflineForeground:               rgbFromHex("#aaaaaa"),
	rosterPeerOnlineForeground:                rgbFromHex("#000000"),
	rosterGroupBackground:                     rgbFromHex("#e9e7f3"),
	rosterAccountOnlineBackground:             rgbFromHex("#918caa"),
	rosterAccountOfflineBackground:            rgbFromHex("#d5d3de"),
	conversationOutgoingUserForeground:        rgbFromHex("#3465a4"),
	conversationIncomingUserForeground:        rgbFromHex("#a40000"),
	conversationOutgoingTextForeground:        rgbFromHex("#000000"),
	conversationIncomingTextForeground:        rgbFromHex("#000000"),
	conversationStatusTextForeground:          rgbFromHex("#aaaaaa"),
	conversationOutgoingDelayedUserForeground: rgbFromHex("#aaaaaa"),
	conversationOutgoingDelayedTextForeground: rgbFromHex("#aaaaaa"),
	conversationLockTypingBackground:          rgbFromHex("#e0e0e0"),
	conversationUnlockTypingBackground:        rgbFromHex("#f9f9f9"),
	timestampForeground:                       rgbFromHex("#aaaaaa"),
}

var defaultDarkColorSet = colorSet{
	rosterPeerBackground:                      rgbFromHex("#7f7f7f"),
	rosterPeerOfflineForeground:               rgbFromHex("#aaaaaa"),
	rosterPeerOnlineForeground:                rgbFromHex("#e5e5e5"),
	rosterGroupBackground:                     rgbFromHex("#b8b6bf"),
	rosterAccountOnlineBackground:             rgbFromHex("#d5d3de"),
	rosterAccountOfflineBackground:            rgbFromHex("#918caa"),
	conversationOutgoingUserForeground:        rgbFromHex("#3465a4"),
	conversationIncomingUserForeground:        rgbFromHex("#a40000"),
	conversationOutgoingTextForeground:        rgbFromHex("#7f7f7f"),
	conversationIncomingTextForeground:        rgbFromHex("#7f7f7f"),
	conversationStatusTextForeground:          rgbFromHex("#4e9a06"),
	conversationOutgoingDelayedUserForeground: rgbFromHex("#444444"),
	conversationOutgoingDelayedTextForeground: rgbFromHex("#444444"),
	conversationLockTypingBackground:          rgbFromHex("#e0e0e0"),
	conversationUnlockTypingBackground:        rgbFromHex("#252a2c"),
	timestampForeground:                       rgbFromHex("#444444"),
}
