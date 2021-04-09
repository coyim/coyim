package gui

import (
	"strconv"
	"strings"
)

type mucColorSet struct {
	warningForeground       string
	warningBackground       string
	someoneJoinedForeground string
	someoneLeftForeground   string
	timestampForeground     string
	nicknameForeground      string
	subjectForeground       string
	infoMessageForeground   string
	messageForeground       string
	errorForeground         string
	configurationForeground string
	roomMessagesBackground  string

	// Room specifics
	roomNameDisabledForeground string

	// Room overlay
	roomOverlaySolidBackground        string
	roomOverlayContentSolidBackground string
	roomOverlayContentBackground      string
	roomOverlayBackground             string
	roomOverlayContentForeground      string

	// Room roster
	rosterGroupBackground        string
	rosterGroupForeground        string
	rosterOccupantRoleForeground string

	// Occupant statuses colors
	occupantStatusAvailableForeground    string
	occupantStatusAvailableBackground    string
	occupantStatusAvailableBorder        string
	occupantStatusNotAvailableForeground string
	occupantStatusNotAvailableBackground string
	occupantStatusNotAvailableBorder     string
	occupantStatusAwayForeground         string
	occupantStatusAwayBackground         string
	occupantStatusAwayBorder             string
	occupantStatusBusyForeground         string
	occupantStatusBusyBackground         string
	occupantStatusBusyBorder             string
	occupantStatusFreeForChatForeground  string
	occupantStatusFreeForChatBackground  string
	occupantStatusFreeForChatBorder      string
	occupantStatusExtendedAwayForeground string
	occupantStatusExtendedAwayBackground string
	occupantStatusExtendedAwayBorder     string

	roomSubjectForeground  string
	roomWarningForeground  string
	roomWarningBackground  string
	roomWarningBorder      string
	entryErrorBackground   string
	entryErrorBorderShadow string
	entryErrorBorder       string
	entryErrorLabel        string
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
		warningBackground:       cs.warningBackground,
		errorForeground:         cs.errorForeground,
		someoneJoinedForeground: "#297316",
		someoneLeftForeground:   "#731629",
		timestampForeground:     "#AAB7B8",
		nicknameForeground:      "#395BA3",
		subjectForeground:       "#000080",
		infoMessageForeground:   "#395BA3",
		messageForeground:       "#000000",
		configurationForeground: "#9a04bf",
		roomMessagesBackground:  "#FFFFFF",

		roomNameDisabledForeground: "@insensitive_fg_color",

		roomOverlaySolidBackground:        "@theme_bg_color",
		roomOverlayContentSolidBackground: "transparent",
		roomOverlayContentBackground:      "@theme_bg_color",
		roomOverlayBackground:             "#000000",
		roomOverlayContentForeground:      "@theme_fg_color",

		rosterGroupBackground:        "#F5F5F4",
		rosterGroupForeground:        "#1C1917",
		rosterOccupantRoleForeground: "#A8A29E",

		occupantStatusAvailableForeground:    "#166534",
		occupantStatusAvailableBackground:    "#F0FDF4",
		occupantStatusAvailableBorder:        "#16A34A",
		occupantStatusNotAvailableForeground: "#1E293B",
		occupantStatusNotAvailableBackground: "#F8FAFC",
		occupantStatusNotAvailableBorder:     "#475569",
		occupantStatusAwayForeground:         "#9A3412",
		occupantStatusAwayBackground:         "#FFF7ED",
		occupantStatusAwayBorder:             "#EA580C",
		occupantStatusBusyForeground:         "#9F1239",
		occupantStatusBusyBackground:         "#FFF1F2",
		occupantStatusBusyBorder:             "#BE123C",
		occupantStatusFreeForChatForeground:  "#1E40AF",
		occupantStatusFreeForChatBackground:  "#EFF6FF",
		occupantStatusFreeForChatBorder:      "#1D4ED8",
		occupantStatusExtendedAwayForeground: "#92400E",
		occupantStatusExtendedAwayBackground: "#FFFBEB",
		occupantStatusExtendedAwayBorder:     "#D97706",

		roomSubjectForeground:  "#666666",
		roomWarningForeground:  "#744210",
		roomWarningBackground:  "#FEFCBF",
		roomWarningBorder:      "#D69E2E",
		entryErrorBackground:   "#FFF5F6",
		entryErrorBorderShadow: "#FF7F50",
		entryErrorBorder:       "#E44635",
		entryErrorLabel:        "#E44635",
	}
}

func (u *gtkUI) defaultMUCDarkColorSet() mucColorSet {
	cs := u.defaultDarkColorSet()

	return mucColorSet{
		warningForeground:       cs.warningForeground,
		warningBackground:       cs.warningBackground,
		errorForeground:         cs.errorForeground,
		someoneJoinedForeground: "#297316",
		someoneLeftForeground:   "#731629",
		timestampForeground:     "@insensitive_fg_color",
		nicknameForeground:      "#395BA3",
		subjectForeground:       "#000080",
		infoMessageForeground:   "#E34267",
		messageForeground:       "#000000",
		configurationForeground: "#9a04bf",
		roomMessagesBackground:  "@theme_base_color",

		roomNameDisabledForeground: "#A9A9A9",

		roomOverlaySolidBackground:        "@theme_base_color",
		roomOverlayContentSolidBackground: "transparent",
		roomOverlayContentBackground:      "@theme_base_color",
		roomOverlayBackground:             "#000000",
		roomOverlayContentForeground:      "#333333",

		rosterGroupBackground:        "#1C1917",
		rosterGroupForeground:        "#FAFAF9",
		rosterOccupantRoleForeground: "#E7E5E4",

		occupantStatusAvailableForeground:    "#166534",
		occupantStatusAvailableBackground:    "#F0FDF4",
		occupantStatusAvailableBorder:        "#16A34A",
		occupantStatusNotAvailableForeground: "#1E293B",
		occupantStatusNotAvailableBackground: "#F8FAFC",
		occupantStatusNotAvailableBorder:     "#475569",
		occupantStatusAwayForeground:         "#9A3412",
		occupantStatusAwayBackground:         "#FFF7ED",
		occupantStatusAwayBorder:             "#EA580C",
		occupantStatusBusyForeground:         "#9F1239",
		occupantStatusBusyBackground:         "#FFF1F2",
		occupantStatusBusyBorder:             "#BE123C",
		occupantStatusFreeForChatForeground:  "#1E40AF",
		occupantStatusFreeForChatBackground:  "#EFF6FF",
		occupantStatusFreeForChatBorder:      "#1D4ED8",
		occupantStatusExtendedAwayForeground: "#92400E",
		occupantStatusExtendedAwayBackground: "#FFFBEB",
		occupantStatusExtendedAwayBorder:     "#D97706",

		roomSubjectForeground:  "#666666",
		roomWarningForeground:  "#744210",
		roomWarningBackground:  "#FEFCBF",
		roomWarningBorder:      "#D69E2E",
		entryErrorBackground:   "#FFF5F6",
		entryErrorBorderShadow: "#FF7F50",
		entryErrorBorder:       "#E44635",
		entryErrorLabel:        "#E44635",
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
