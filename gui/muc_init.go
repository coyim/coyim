package gui

func (u *gtkUI) initMUC() {
	initMUCSupportedErrors()
	initMUCTextsAndMessages()
	initMUCStyles(u.currentMUCColorSet())
}

func initMUCTextsAndMessages() {
	initMUCConfigUpdateMessages()
	initMUCRoomConfigTexts()
	initMUCRoomConversationTexts()
}
