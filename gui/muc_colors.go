package gui

import (
	"strconv"
	"strings"
)

type mucColorSet struct {
	warningForeground                 string
	warningBackground                 string
	someoneJoinedForeground           string
	someoneLeftForeground             string
	timestampForeground               string
	nicknameForeground                string
	subjectForeground                 string
	infoMessageForeground             string
	messageForeground                 string
	errorForeground                   string
	configurationForeground           string
	roomMessagesBackground            string
	roomOverlaySolidBackground        string
	roomOverlayContentSolidBackground string
	roomOverlayContentBackground      string
	roomOverlayBackground             string
	roomOverlayContentForeground      string
	roomNameDisabledForeground        string
	roomRosterStatusForeground        string
	roomSubjectForeground             string
	roomWarningForeground             string
	roomWarningBackground             string
	roomWarningBorder                 string
	entryErrorBackground              string
	entryErrorBorderShadow            string
	entryErrorBorder                  string
	entryErrorLabel                   string
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
		warningForeground:                 cs.warningForeground,
		warningBackground:                 cs.warningBackground,
		errorForeground:                   cs.errorForeground,
		someoneJoinedForeground:           "#297316",
		someoneLeftForeground:             "#731629",
		timestampForeground:               "#AAB7B8",
		nicknameForeground:                "#395BA3",
		subjectForeground:                 "#000080",
		infoMessageForeground:             "#395BA3",
		messageForeground:                 "#000000",
		configurationForeground:           "#9a04bf",
		roomMessagesBackground:            "#FFFFFF",
		roomOverlaySolidBackground:        "#FFFFFF",
		roomOverlayContentSolidBackground: "transparent",
		roomOverlayContentBackground:      "#FFFFFF",
		roomOverlayBackground:             "#000000",
		roomOverlayContentForeground:      "#333333",
		roomNameDisabledForeground:        "#A9A9A9",
		roomRosterStatusForeground:        "#666666",
		roomSubjectForeground:             "#666666",
		roomWarningForeground:             "#744210",
		roomWarningBackground:             "#FEFCBF",
		roomWarningBorder:                 "#D69E2E",
		entryErrorBackground:              "#FFF5F6",
		entryErrorBorderShadow:            "#FF7F50",
		entryErrorBorder:                  "#E44635",
		entryErrorLabel:                   "#E44635",
	}
}

func (u *gtkUI) defaultMUCDarkColorSet() mucColorSet {
	cs := u.defaultDarkColorSet()

	return mucColorSet{
		warningForeground:                 cs.warningForeground,
		warningBackground:                 cs.warningBackground,
		errorForeground:                   cs.errorForeground,
		someoneJoinedForeground:           "#297316",
		someoneLeftForeground:             "#731629",
		timestampForeground:               "#AAB7B8",
		nicknameForeground:                "#395BA3",
		subjectForeground:                 "#000080",
		infoMessageForeground:             "#E34267",
		messageForeground:                 "#000000",
		configurationForeground:           "#9a04bf",
		roomMessagesBackground:            "#FFFFFF",
		roomOverlaySolidBackground:        "#FFFFFF",
		roomOverlayContentSolidBackground: "transparent",
		roomOverlayContentBackground:      "#FFFFFF",
		roomOverlayBackground:             "#000000",
		roomOverlayContentForeground:      "#333333",
		roomNameDisabledForeground:        "#A9A9A9",
		roomRosterStatusForeground:        "#666666",
		roomSubjectForeground:             "#666666",
		roomWarningForeground:             "#744210",
		roomWarningBackground:             "#FEFCBF",
		roomWarningBorder:                 "#D69E2E",
		entryErrorBackground:              "#FFF5F6",
		entryErrorBorderShadow:            "#FF7F50",
		entryErrorBorder:                  "#E44635",
		entryErrorLabel:                   "#E44635",
	}
}

type rgb struct {
	red, green, blue uint8
}

func (cs *mucColorSet) hexClean(hex string) string {
	return strings.Replace(hex, "#", "", -1)
}

func (cs *mucColorSet) hexToRGB(hex string) (*rgb, error) {
	values, err := strconv.ParseInt(cs.hexClean(hex), 16, 32)
	if err != nil {
		return nil, err
	}

	return &rgb{uint8(values >> 16), uint8((values >> 8) & 0xFF), uint8(values & 0xFF)}, nil
}
