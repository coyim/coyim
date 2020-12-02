package gui

import (
	"strconv"
	"strings"
)

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
	black                   string
	light                   string
	dark                    string
	gray300                 string
	gray500                 string
	brown500                string
	yellow200               string
	yellow600               string
	red200                  string
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
		black:                   "#000000",
		light:                   "#FFFFFF",
		dark:                    "#333333",
		gray300:                 "#A9A9A9",
		gray500:                 "#666666",
		brown500:                "#744210",
		yellow200:               "#FEFCBF",
		yellow600:               "#D69E2E",
		red200:                  "#FEE2E2",
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
		black:                   "#000000",
		light:                   "#111111",
		dark:                    "#000000",
		gray300:                 "#A9A9A9",
		gray500:                 "#666666",
		brown500:                "#744210",
		yellow200:               "#FEFCBF",
		yellow600:               "#D69E2E",
		red200:                  "#FEE2E2",
	}
}

type rgb struct {
	red, green, blue uint8
}

func (cs *mucColorSet) hexClean(hex string) string {
	return strings.ReplaceAll(hex, "#", "")
}

func (cs *mucColorSet) hexToRGB(hex string) (*rgb, error) {
	values, err := strconv.ParseInt(cs.hexClean(hex), 16, 32)
	if err != nil {
		return nil, err
	}

	return &rgb{uint8(values >> 16), uint8((values >> 8) & 0xFF), uint8(values & 0xFF)}, nil
}
