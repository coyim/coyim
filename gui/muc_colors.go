package gui

type mucColorSet struct {
	warningForeground       string
	someoneJoinedForeground string
	someoneLeftForeground   string
	timestampForeground     string
	nicknameForeground      string
	subjectForeground       string
	infoMessageForeground   string
	messageForeground       string
	errorForeground         string
	configurationForeground string
	white                   string
	gray300                 string
	gray500                 string
	brown500                string
	yellow200               string
	yellow600               string
}

func (u *gtkUI) currentMUCColorSet() mucColorSet {
	if u.isDarkThemeVariant() {
		return u.defaultMUCDarkColorSet()
	}
	return u.defaultMUCLightColorSet()
}

func (u *gtkUI) defaultMUCLightColorSet() mucColorSet {
	cs := u.defaultLightColorSet()

	return mucColorSet{
		warningForeground:       cs.warningForeground,
		errorForeground:         cs.errorForeground,
		someoneJoinedForeground: "#297316",
		someoneLeftForeground:   "#731629",
		timestampForeground:     "#AAB7B8",
		nicknameForeground:      "#395BA3",
		subjectForeground:       "#000080",
		infoMessageForeground:   "#395BA3",
		messageForeground:       "#000000",
		configurationForeground: "#9a04bf",
		white:                   "#FFFFFF",
		gray300:                 "#A9A9A9",
		gray500:                 "#666666",
		brown500:                "#744210",
		yellow200:               "#FEFCBF",
		yellow600:               "#D69E2E",
	}
}

func (u *gtkUI) defaultMUCDarkColorSet() mucColorSet {
	cs := u.defaultDarkColorSet()

	return mucColorSet{
		warningForeground:       cs.warningForeground,
		errorForeground:         cs.errorForeground,
		someoneJoinedForeground: "#297316",
		someoneLeftForeground:   "#731629",
		timestampForeground:     "#AAB7B8",
		nicknameForeground:      "#395BA3",
		subjectForeground:       "#000080",
		infoMessageForeground:   "#E34267",
		messageForeground:       "#000000",
		configurationForeground: "#9a04bf",
		white:                   "#FFFFFF",
		gray300:                 "#A9A9A9",
		gray500:                 "#666666",
		brown500:                "#744210",
		yellow200:               "#FEFCBF",
		yellow600:               "#D69E2E",
	}
}
