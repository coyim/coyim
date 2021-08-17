package gui

import (
	"fmt"
)

const (
	colorNone                       = "none"
	colorTransparent                = "transparent"
	colorThemeBase                  = "@theme_base_color"
	colorThemeBackground            = "@theme_bg_color"
	colorThemeForeground            = "@theme_fg_color"
	colorThemeInsensitiveBackground = "@insensitive_bg_color"
	colorThemeInsensitiveForeground = "@insensitive_fg_color"
)

type mucColorSet struct {
	warningForeground                      string
	warningBackground                      string
	someoneJoinedForeground                string
	someoneLeftForeground                  string
	timestampForeground                    string
	nicknameForeground                     string
	subjectForeground                      string
	infoMessageForeground                  string
	messageForeground                      string
	errorForeground                        string
	configurationForeground                string
	roomMessagesBackground                 string
	roomMessagesBoxShadow                  string
	roomNameDisabledForeground             string
	roomSubjectForeground                  string
	roomOverlaySolidBackground             string
	roomOverlayContentSolidBackground      string
	roomOverlayContentBackground           string
	roomOverlayBackground                  string
	roomOverlayContentForeground           string
	roomOverlayContentBoxShadow            string
	roomWarningsDialogBackground           string
	roomWarningsDialogDecorationBackground string
	roomWarningsDialogDecorationShadow     string
	roomWarningsDialogHeaderBackground     string
	roomWarningsDialogContentBackground    string
	roomWarningsCurrentInfoForeground      string
	roomNotificationsBackground            string
	rosterGroupBackground                  string
	rosterGroupForeground                  string
	rosterOccupantRoleForeground           string
	occupantStatusAvailableForeground      string
	occupantStatusAvailableBackground      string
	occupantStatusAvailableBorder          string
	occupantStatusNotAvailableForeground   string
	occupantStatusNotAvailableBackground   string
	occupantStatusNotAvailableBorder       string
	occupantStatusAwayForeground           string
	occupantStatusAwayBackground           string
	occupantStatusAwayBorder               string
	occupantStatusBusyForeground           string
	occupantStatusBusyBackground           string
	occupantStatusBusyBorder               string
	occupantStatusFreeForChatForeground    string
	occupantStatusFreeForChatBackground    string
	occupantStatusFreeForChatBorder        string
	occupantStatusExtendedAwayForeground   string
	occupantStatusExtendedAwayBackground   string
	occupantStatusExtendedAwayBorder       string
	infoBarDefaultBorderColor              string
	infoBarTypeInfoBackgroundStart         string
	infoBarTypeInfoBackgroundStop          string
	infoBarTypeInfoTitle                   string
	infoBarTypeInfoTime                    string
	infoBarTypeWarningBackgroundStart      string
	infoBarTypeWarningBackgroundStop       string
	infoBarTypeWarningTitle                string
	infoBarTypeWarningTime                 string
	infoBarTypeQuestionBackgroundStart     string
	infoBarTypeQuestionBackgroundStop      string
	infoBarTypeQuestionTitle               string
	infoBarTypeQuestionTime                string
	infoBarTypeErrorBackgroundStart        string
	infoBarTypeErrorBackgroundStop         string
	infoBarTypeErrorTitle                  string
	infoBarTypeErrorTime                   string
	infoBarTypeOtherBackgroundStart        string
	infoBarTypeOtherBackgroundStop         string
	infoBarTypeOtherTitle                  string
	infoBarTypeOtherTime                   string
	infoBarButtonBackground                string
	infoBarButtonForeground                string
	infoBarButtonHoverBackground           string
	infoBarButtonHoverForeground           string
	infoBarButtonActiveBackground          string
	infoBarButtonActiveForeground          string
	entryErrorBackground                   string
	entryErrorBorderShadow                 string
	entryErrorBorder                       string
	entryErrorLabel                        string
	occupantLostConnection                 string
	occupantRestablishConnection           string
}

func (u *gtkUI) currentMUCColorSet() mucColorSet {
	if u.isDarkThemeVariant() {
		return u.defaultMUCDarkColorSet()
	}
	return u.defaultMUCLightColorSet()
}

func (u *gtkUI) defaultMUCLightColorSet() mucColorSet {
	return mucColorSet{
		warningForeground:                      colorFormat(rgb{194, 23, 29}, 1),
		warningBackground:                      colorFormat(rgb{254, 202, 202}, 1),
		errorForeground:                        colorFormat(rgb{163, 7, 7}, 1),
		someoneJoinedForeground:                colorFormat(rgb{41, 115, 22}, 1),
		someoneLeftForeground:                  colorFormat(rgb{115, 22, 41}, 1),
		timestampForeground:                    colorThemeInsensitiveForeground,
		nicknameForeground:                     colorFormat(rgb{57, 91, 163}, 1),
		subjectForeground:                      colorFormat(rgb{0, 0, 128}, 1),
		infoMessageForeground:                  colorFormat(rgb{57, 91, 163}, 1),
		messageForeground:                      colorFormat(rgb{0, 0, 0}, 1),
		configurationForeground:                colorFormat(rgb{154, 4, 191}, 1),
		roomMessagesBackground:                 colorThemeBase,
		roomMessagesBoxShadow:                  colorFormat(rgb{0, 0, 0}, 0.35),
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBackground,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBackground,
		roomOverlayBackground:                  colorFormat(rgb{0, 0, 0}, 0.5),
		roomOverlayContentForeground:           colorThemeForeground,
		roomOverlayContentBoxShadow:            colorFormat(rgb{0, 0, 0}, 0.5),
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogDecorationShadow:     colorFormat(rgb{0, 0, 0}, 0.15),
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		roomNotificationsBackground:            colorThemeBackground,
		rosterGroupBackground:                  colorFormat(rgb{245, 245, 244}, 1),
		rosterGroupForeground:                  colorFormat(rgb{28, 25, 23}, 1),
		rosterOccupantRoleForeground:           colorFormat(rgb{168, 162, 158}, 1),
		occupantStatusAvailableForeground:      colorFormat(rgb{22, 101, 52}, 1),
		occupantStatusAvailableBackground:      colorFormat(rgb{240, 253, 244}, 1),
		occupantStatusAvailableBorder:          colorFormat(rgb{22, 163, 74}, 1),
		occupantStatusNotAvailableForeground:   colorFormat(rgb{30, 41, 59}, 1),
		occupantStatusNotAvailableBackground:   colorFormat(rgb{248, 250, 252}, 1),
		occupantStatusNotAvailableBorder:       colorFormat(rgb{71, 85, 105}, 1),
		occupantStatusAwayForeground:           colorFormat(rgb{154, 52, 18}, 1),
		occupantStatusAwayBackground:           colorFormat(rgb{255, 247, 237}, 1),
		occupantStatusAwayBorder:               colorFormat(rgb{234, 88, 12}, 1),
		occupantStatusBusyForeground:           colorFormat(rgb{159, 18, 57}, 1),
		occupantStatusBusyBackground:           colorFormat(rgb{255, 241, 242}, 1),
		occupantStatusBusyBorder:               colorFormat(rgb{190, 18, 60}, 1),
		occupantStatusFreeForChatForeground:    colorFormat(rgb{30, 64, 175}, 1),
		occupantStatusFreeForChatBackground:    colorFormat(rgb{239, 246, 255}, 1),
		occupantStatusFreeForChatBorder:        colorFormat(rgb{29, 78, 216}, 1),
		occupantStatusExtendedAwayForeground:   colorFormat(rgb{146, 64, 14}, 1),
		occupantStatusExtendedAwayBackground:   colorFormat(rgb{255, 251, 235}, 1),
		occupantStatusExtendedAwayBorder:       colorFormat(rgb{217, 119, 6}, 1),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         colorFormat(rgb{63, 98, 18}, 1),
		infoBarTypeInfoBackgroundStop:          colorFormat(rgb{77, 124, 15}, 1),
		infoBarTypeInfoTitle:                   colorFormat(rgb{236, 254, 255}, 1),
		infoBarTypeInfoTime:                    colorFormat(rgb{236, 254, 255}, 0.5),
		infoBarTypeWarningBackgroundStart:      colorFormat(rgb{195, 149, 7}, 1),
		infoBarTypeWarningBackgroundStop:       colorFormat(rgb{222, 173, 20}, 1),
		infoBarTypeWarningTitle:                colorFormat(rgb{255, 247, 237}, 1),
		infoBarTypeWarningTime:                 colorFormat(rgb{255, 247, 237}, 0.5),
		infoBarTypeQuestionBackgroundStart:     colorFormat(rgb{234, 88, 12}, 1),
		infoBarTypeQuestionBackgroundStop:      colorFormat(rgb{249, 115, 22}, 1),
		infoBarTypeQuestionTitle:               colorFormat(rgb{254, 252, 232}, 1),
		infoBarTypeQuestionTime:                colorFormat(rgb{254, 252, 232}, 0.5),
		infoBarTypeErrorBackgroundStart:        colorFormat(rgb{185, 28, 28}, 1),
		infoBarTypeErrorBackgroundStop:         colorFormat(rgb{203, 35, 35}, 1),
		infoBarTypeErrorTitle:                  colorFormat(rgb{255, 241, 242}, 1),
		infoBarTypeErrorTime:                   colorFormat(rgb{255, 241, 242}, 0.5),
		infoBarTypeOtherBackgroundStart:        colorFormat(rgb{7, 89, 133}, 1),
		infoBarTypeOtherBackgroundStop:         colorFormat(rgb{3, 105, 161}, 1),
		infoBarTypeOtherTitle:                  colorFormat(rgb{240, 253, 250}, 1),
		infoBarTypeOtherTime:                   colorFormat(rgb{240, 253, 250}, 0.5),
		infoBarButtonBackground:                colorFormat(rgb{0, 0, 0}, 0.25),
		infoBarButtonForeground:                colorFormat(rgb{255, 255, 255}, 1),
		infoBarButtonHoverBackground:           colorFormat(rgb{0, 0, 0}, 0.35),
		infoBarButtonHoverForeground:           colorFormat(rgb{255, 255, 255}, 1),
		infoBarButtonActiveBackground:          colorFormat(rgb{0, 0, 0}, 0.45),
		infoBarButtonActiveForeground:          colorFormat(rgb{255, 255, 255}, 1),
		entryErrorBackground:                   colorFormat(rgb{255, 245, 246}, 1),
		entryErrorBorderShadow:                 colorFormat(rgb{255, 127, 80}, 1),
		entryErrorBorder:                       colorFormat(rgb{228, 70, 53}, 1),
		entryErrorLabel:                        colorFormat(rgb{228, 70, 53}, 1),
		occupantLostConnection:                 colorFormat(rgb{115, 22, 41}, 1),
		occupantRestablishConnection:           colorFormat(rgb{41, 115, 22}, 1),
	}
}

func (u *gtkUI) defaultMUCDarkColorSet() mucColorSet {
	return mucColorSet{
		warningForeground:                      colorFormat(rgb{194, 23, 29}, 1),
		warningBackground:                      colorFormat(rgb{254, 202, 202}, 1),
		errorForeground:                        colorFormat(rgb{209, 104, 96}, 1),
		someoneJoinedForeground:                colorFormat(rgb{41, 115, 22}, 1),
		someoneLeftForeground:                  colorFormat(rgb{115, 22, 41}, 1),
		timestampForeground:                    colorThemeInsensitiveForeground,
		nicknameForeground:                     colorFormat(rgb{57, 91, 163}, 1),
		subjectForeground:                      colorFormat(rgb{0, 0, 128}, 1),
		infoMessageForeground:                  colorFormat(rgb{227, 66, 103}, 1),
		messageForeground:                      colorFormat(rgb{0, 0, 0}, 1),
		configurationForeground:                colorFormat(rgb{154, 4, 191}, 1),
		roomMessagesBackground:                 colorThemeBase,
		roomMessagesBoxShadow:                  colorFormat(rgb{0, 0, 0}, 0.35),
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBase,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBase,
		roomOverlayBackground:                  colorFormat(rgb{0, 0, 0}, 0.5),
		roomOverlayContentForeground:           colorFormat(rgb{51, 51, 51}, 1),
		roomOverlayContentBoxShadow:            colorFormat(rgb{0, 0, 0}, 0.5),
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogDecorationShadow:     colorFormat(rgb{0, 0, 0}, 0.15),
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		roomNotificationsBackground:            colorThemeBackground,
		rosterGroupBackground:                  colorFormat(rgb{28, 25, 23}, 1),
		rosterGroupForeground:                  colorFormat(rgb{250, 250, 249}, 1),
		rosterOccupantRoleForeground:           colorFormat(rgb{231, 229, 228}, 1),
		occupantStatusAvailableForeground:      colorFormat(rgb{22, 101, 52}, 1),
		occupantStatusAvailableBackground:      colorFormat(rgb{240, 253, 244}, 1),
		occupantStatusAvailableBorder:          colorFormat(rgb{22, 163, 74}, 1),
		occupantStatusNotAvailableForeground:   colorFormat(rgb{30, 41, 59}, 1),
		occupantStatusNotAvailableBackground:   colorFormat(rgb{248, 250, 252}, 1),
		occupantStatusNotAvailableBorder:       colorFormat(rgb{71, 85, 105}, 1),
		occupantStatusAwayForeground:           colorFormat(rgb{154, 52, 18}, 1),
		occupantStatusAwayBackground:           colorFormat(rgb{255, 247, 237}, 1),
		occupantStatusAwayBorder:               colorFormat(rgb{234, 88, 12}, 1),
		occupantStatusBusyForeground:           colorFormat(rgb{159, 18, 57}, 1),
		occupantStatusBusyBackground:           colorFormat(rgb{255, 241, 242}, 1),
		occupantStatusBusyBorder:               colorFormat(rgb{190, 18, 60}, 1),
		occupantStatusFreeForChatForeground:    colorFormat(rgb{30, 64, 175}, 1),
		occupantStatusFreeForChatBackground:    colorFormat(rgb{239, 246, 255}, 1),
		occupantStatusFreeForChatBorder:        colorFormat(rgb{29, 78, 216}, 1),
		occupantStatusExtendedAwayForeground:   colorFormat(rgb{146, 64, 14}, 1),
		occupantStatusExtendedAwayBackground:   colorFormat(rgb{255, 251, 235}, 1),
		occupantStatusExtendedAwayBorder:       colorFormat(rgb{217, 119, 6}, 1),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         colorFormat(rgb{63, 98, 18}, 1),
		infoBarTypeInfoBackgroundStop:          colorFormat(rgb{77, 124, 15}, 1),
		infoBarTypeInfoTitle:                   colorFormat(rgb{236, 254, 255}, 1),
		infoBarTypeInfoTime:                    colorFormat(rgb{236, 254, 255}, 0.5),
		infoBarTypeWarningBackgroundStart:      colorFormat(rgb{195, 149, 7}, 1),
		infoBarTypeWarningBackgroundStop:       colorFormat(rgb{222, 173, 20}, 1),
		infoBarTypeWarningTitle:                colorFormat(rgb{255, 247, 237}, 1),
		infoBarTypeWarningTime:                 colorFormat(rgb{255, 247, 237}, 0.5),
		infoBarTypeQuestionBackgroundStart:     colorFormat(rgb{234, 88, 12}, 1),
		infoBarTypeQuestionBackgroundStop:      colorFormat(rgb{249, 115, 22}, 1),
		infoBarTypeQuestionTitle:               colorFormat(rgb{254, 252, 232}, 1),
		infoBarTypeQuestionTime:                colorFormat(rgb{254, 252, 232}, 0.5),
		infoBarTypeErrorBackgroundStart:        colorFormat(rgb{185, 28, 28}, 1),
		infoBarTypeErrorBackgroundStop:         colorFormat(rgb{203, 35, 35}, 1),
		infoBarTypeErrorTitle:                  colorFormat(rgb{255, 241, 242}, 1),
		infoBarTypeErrorTime:                   colorFormat(rgb{255, 241, 242}, 0.5),
		infoBarTypeOtherBackgroundStart:        colorFormat(rgb{7, 89, 133}, 1),
		infoBarTypeOtherBackgroundStop:         colorFormat(rgb{3, 105, 161}, 1),
		infoBarTypeOtherTitle:                  colorFormat(rgb{240, 253, 250}, 1),
		infoBarTypeOtherTime:                   colorFormat(rgb{240, 253, 250}, 0.5),
		infoBarButtonBackground:                colorFormat(rgb{0, 0, 0}, 0.25),
		infoBarButtonForeground:                colorFormat(rgb{255, 255, 255}, 1),
		infoBarButtonHoverBackground:           colorFormat(rgb{0, 0, 0}, 0.35),
		infoBarButtonHoverForeground:           colorFormat(rgb{255, 255, 255}, 1),
		infoBarButtonActiveBackground:          colorFormat(rgb{0, 0, 0}, 0.45),
		infoBarButtonActiveForeground:          colorFormat(rgb{255, 255, 255}, 1),
		entryErrorBackground:                   colorFormat(rgb{255, 245, 246}, 1),
		entryErrorBorderShadow:                 colorFormat(rgb{255, 127, 80}, 1),
		entryErrorBorder:                       colorFormat(rgb{228, 70, 53}, 1),
		entryErrorLabel:                        colorFormat(rgb{228, 70, 53}, 1),
		occupantLostConnection:                 colorFormat(rgb{115, 22, 41}, 1),
		occupantRestablishConnection:           colorFormat(rgb{41, 115, 22}, 1),
	}
}

type rgb struct {
	red   uint8
	green uint8
	blue  uint8
}

func colorFormat(c rgb, alpha float64) string {
	if alpha == 1 {
		return fmt.Sprintf("rgb(%d, %d, %d)", c.red, c.green, c.blue)
	}
	return fmt.Sprintf("rgba(%d, %d, %d, %f)", c.red, c.green, c.blue, alpha)
}
