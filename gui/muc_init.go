package gui

func (u *gtkUI) initMUC() {
	initMUCSupportedErrors()
	initMUCTextsAndMessages()
	initMUCComponents()
	initMUCStyles(u.currentMUCColorSet())
}

func initMUCTextsAndMessages() {
	initMUCConfigUpdateMessages()
	initMUCRoomConfigTexts()
}

func initMUCComponents() {
	initMUCInfoBarComponent()
}
