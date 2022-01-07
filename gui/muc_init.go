package gui

func (u *gtkUI) initMUC() {
	initMUCSupportedErrors()
	initMUCTextsAndMessages()
	initMUCStyles(u, u.currentMUCColorSet())
}

func initMUCTextsAndMessages() {
	initMUCRoomConfigTexts()
}
