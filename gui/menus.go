package gui

type isSensitive interface {
	SetSensitive(bool)
}

var globalSensitives = []isSensitive{}

func addItemsThatShouldToggleOnGlobalMenuStatus(s isSensitive) {
	globalSensitives = append(globalSensitives, s)
}

func (u *gtkUI) updateGlobalMenuStatus() {
	haveOnline := false
	for _, a := range u.accounts {
		if a.connected() {
			haveOnline = true
		}
	}

	for _, gs := range globalSensitives {
		gs.SetSensitive(haveOnline)
	}
}
