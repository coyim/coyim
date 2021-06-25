package gui

import (
	"fmt"
	"strconv"
	"strings"
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
	roomNameDisabledForeground             string
	roomSubjectForeground                  string
	roomOverlaySolidBackground             string
	roomOverlayContentSolidBackground      string
	roomOverlayContentBackground           string
	roomOverlayBackground                  string
	roomOverlayContentForeground           string
	roomWarningsDialogBackground           string
	roomWarningsDialogDecorationBackground string
	roomWarningsDialogHeaderBackground     string
	roomWarningsDialogContentBackground    string
	roomWarningsCurrentInfoForeground      string
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
	infoBarTypeWarningBackgroundStart      string
	infoBarTypeWarningBackgroundStop       string
	infoBarTypeWarningTitle                string
	infoBarTypeQuestionBackgroundStart     string
	infoBarTypeQuestionBackgroundStop      string
	infoBarTypeQuestionTitle               string
	infoBarTypeErrorBackgroundStart        string
	infoBarTypeErrorBackgroundStop         string
	infoBarTypeErrorTitle                  string
	infoBarTypeOtherBackgroundStart        string
	infoBarTypeOtherBackgroundStop         string
	infoBarTypeOtherTitle                  string
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
		warningForeground:                      colorFormat(cs.warningForeground, 1),
		warningBackground:                      colorFormat(cs.warningBackground, 1),
		someoneJoinedForeground:                colorFormat("297316", 1),
		someoneLeftForeground:                  colorFormat("731629", 1),
		timestampForeground:                    colorFormat("AAB7B8", 1),
		nicknameForeground:                     colorFormat("395BA3", 1),
		subjectForeground:                      colorFormat("000080", 1),
		infoMessageForeground:                  colorFormat("395BA3", 1),
		messageForeground:                      colorFormat("000000", 1),
		errorForeground:                        colorFormat(cs.errorForeground, 1),
		configurationForeground:                colorFormat("9A04BF", 1),
		roomMessagesBackground:                 colorThemeBase,
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBackground,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBackground,
		roomOverlayBackground:                  colorFormat("000000", 1),
		roomOverlayContentForeground:           colorThemeForeground,
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		rosterGroupBackground:                  colorFormat("F5F5F4", 1),
		rosterGroupForeground:                  colorFormat("1C1917", 1),
		rosterOccupantRoleForeground:           colorFormat("A8A29E", 1),
		occupantStatusAvailableForeground:      colorFormat("166534", 1),
		occupantStatusAvailableBackground:      colorFormat("F0FDF4", 1),
		occupantStatusAvailableBorder:          colorFormat("16A34A", 1),
		occupantStatusNotAvailableForeground:   colorFormat("1E293B", 1),
		occupantStatusNotAvailableBackground:   colorFormat("F8FAFC", 1),
		occupantStatusNotAvailableBorder:       colorFormat("475569", 1),
		occupantStatusAwayForeground:           colorFormat("9A3412", 1),
		occupantStatusAwayBackground:           colorFormat("FFF7ED", 1),
		occupantStatusAwayBorder:               colorFormat("EA580C", 1),
		occupantStatusBusyForeground:           colorFormat("9F1239", 1),
		occupantStatusBusyBackground:           colorFormat("FFF1F2", 1),
		occupantStatusBusyBorder:               colorFormat("BE123C", 1),
		occupantStatusFreeForChatForeground:    colorFormat("1E40AF", 1),
		occupantStatusFreeForChatBackground:    colorFormat("EFF6FF", 1),
		occupantStatusFreeForChatBorder:        colorFormat("1D4ED8", 1),
		occupantStatusExtendedAwayForeground:   colorFormat("92400E", 1),
		occupantStatusExtendedAwayBackground:   colorFormat("FFFBEB", 1),
		occupantStatusExtendedAwayBorder:       colorFormat("D97706", 1),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         colorFormat("3F6212", 1),
		infoBarTypeInfoBackgroundStop:          colorFormat("4D7C0F", 1),
		infoBarTypeInfoTitle:                   colorFormat("ECFEFF", 1),
		infoBarTypeWarningBackgroundStart:      colorFormat("C39507", 1),
		infoBarTypeWarningBackgroundStop:       colorFormat("DEAD14", 1),
		infoBarTypeWarningTitle:                colorFormat("FFF7ED", 1),
		infoBarTypeQuestionBackgroundStart:     colorFormat("EA580C", 1),
		infoBarTypeQuestionBackgroundStop:      colorFormat("F97316", 1),
		infoBarTypeQuestionTitle:               colorFormat("FEFCE8", 1),
		infoBarTypeErrorBackgroundStart:        colorFormat("B91C1C", 1),
		infoBarTypeErrorBackgroundStop:         colorFormat("CB2323", 1),
		infoBarTypeErrorTitle:                  colorFormat("FFF1F2", 1),
		infoBarTypeOtherBackgroundStart:        colorFormat("075985", 1),
		infoBarTypeOtherBackgroundStop:         colorFormat("0369A1", 1),
		infoBarTypeOtherTitle:                  colorFormat("F0FDFA", 1),
		infoBarButtonBackground:                colorFormat("000000", 0.25),
		infoBarButtonForeground:                colorFormat("FFFFFF", 1),
		infoBarButtonHoverBackground:           colorFormat("000000", 0.35),
		infoBarButtonHoverForeground:           colorFormat("FFFFFF", 1),
		infoBarButtonActiveBackground:          colorFormat("000000", 0.45),
		infoBarButtonActiveForeground:          colorFormat("FFFFFF", 1),
		entryErrorBackground:                   colorFormat("FFF5F6", 1),
		entryErrorBorderShadow:                 colorFormat("FF7F50", 1),
		entryErrorBorder:                       colorFormat("E44635", 1),
		entryErrorLabel:                        colorFormat("E44635", 1),
	}
}

func (u *gtkUI) defaultMUCDarkColorSet() mucColorSet {
	cs := u.defaultDarkColorSet()

	return mucColorSet{
		warningForeground:                      cs.warningForeground,
		warningBackground:                      cs.warningBackground,
		errorForeground:                        cs.errorForeground,
		someoneJoinedForeground:                colorFormat("297316", 1),
		someoneLeftForeground:                  colorFormat("731629", 1),
		timestampForeground:                    colorThemeInsensitiveForeground,
		nicknameForeground:                     colorFormat("395BA3", 1),
		subjectForeground:                      colorFormat("000080", 1),
		infoMessageForeground:                  colorFormat("E34267", 1),
		messageForeground:                      colorFormat("000000", 1),
		configurationForeground:                colorFormat("9a04bf", 1),
		roomMessagesBackground:                 colorThemeBase,
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBase,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBase,
		roomOverlayBackground:                  colorFormat("000000", 1),
		roomOverlayContentForeground:           colorFormat("333333", 1),
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		rosterGroupBackground:                  colorFormat("1C1917", 1),
		rosterGroupForeground:                  colorFormat("FAFAF9", 1),
		rosterOccupantRoleForeground:           colorFormat("E7E5E4", 1),
		occupantStatusAvailableForeground:      colorFormat("166534", 1),
		occupantStatusAvailableBackground:      colorFormat("F0FDF4", 1),
		occupantStatusAvailableBorder:          colorFormat("16A34A", 1),
		occupantStatusNotAvailableForeground:   colorFormat("1E293B", 1),
		occupantStatusNotAvailableBackground:   colorFormat("F8FAFC", 1),
		occupantStatusNotAvailableBorder:       colorFormat("475569", 1),
		occupantStatusAwayForeground:           colorFormat("9A3412", 1),
		occupantStatusAwayBackground:           colorFormat("FFF7ED", 1),
		occupantStatusAwayBorder:               colorFormat("EA580C", 1),
		occupantStatusBusyForeground:           colorFormat("9F1239", 1),
		occupantStatusBusyBackground:           colorFormat("FFF1F2", 1),
		occupantStatusBusyBorder:               colorFormat("BE123C", 1),
		occupantStatusFreeForChatForeground:    colorFormat("1E40AF", 1),
		occupantStatusFreeForChatBackground:    colorFormat("EFF6FF", 1),
		occupantStatusFreeForChatBorder:        colorFormat("1D4ED8", 1),
		occupantStatusExtendedAwayForeground:   colorFormat("92400E", 1),
		occupantStatusExtendedAwayBackground:   colorFormat("FFFBEB", 1),
		occupantStatusExtendedAwayBorder:       colorFormat("D97706", 1),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         colorFormat("3F6212", 1),
		infoBarTypeInfoBackgroundStop:          colorFormat("4D7C0F", 1),
		infoBarTypeInfoTitle:                   colorFormat("ECFEFF", 1),
		infoBarTypeWarningBackgroundStart:      colorFormat("C39507", 1),
		infoBarTypeWarningBackgroundStop:       colorFormat("DEAD14", 1),
		infoBarTypeWarningTitle:                colorFormat("FFF7ED", 1),
		infoBarTypeQuestionBackgroundStart:     colorFormat("EA580C", 1),
		infoBarTypeQuestionBackgroundStop:      colorFormat("F97316", 1),
		infoBarTypeQuestionTitle:               colorFormat("FEFCE8", 1),
		infoBarTypeErrorBackgroundStart:        colorFormat("B91C1C", 1),
		infoBarTypeErrorBackgroundStop:         colorFormat("CB2323", 1),
		infoBarTypeErrorTitle:                  colorFormat("FFF1F2", 1),
		infoBarTypeOtherBackgroundStart:        colorFormat("075985", 1),
		infoBarTypeOtherBackgroundStop:         colorFormat("0369A1", 1),
		infoBarTypeOtherTitle:                  colorFormat("F0FDFA", 1),
		infoBarButtonBackground:                colorFormat("000000", 0.25),
		infoBarButtonForeground:                colorFormat("FFFFFF", 1),
		infoBarButtonHoverBackground:           colorFormat("000000", 0.35),
		infoBarButtonHoverForeground:           colorFormat("FFFFFF", 1),
		infoBarButtonActiveBackground:          colorFormat("000000", 0.45),
		infoBarButtonActiveForeground:          colorFormat("FFFFFF", 1),
		entryErrorBackground:                   colorFormat("FFF5F6", 1),
		entryErrorBorderShadow:                 colorFormat("FF7F50", 1),
		entryErrorBorder:                       colorFormat("E44635", 1),
		entryErrorLabel:                        colorFormat("E44635", 1),
	}
}

type rgb struct {
	red   uint8
	green uint8
	blue  uint8
}

func (c *rgb) format() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.red, c.green, c.blue)
}

func (c *rgb) formatWithAlpha(alpha float64) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %f)", c.red, c.green, c.blue, alpha)
}

var colorFallback = &rgb{0, 0, 0}

func colorFormat(color string, alpha float64) string {
	if alpha == 1 {
		return colorHexAddPrefix(color)
	}

	if c, err := colorHexToRGB(color); err == nil {
		return c.formatWithAlpha(alpha)
	}

	return colorFallback.formatWithAlpha(alpha)
}

func colorHexAddPrefix(color string) string {
	if strings.HasPrefix(color, "#") {
		return color
	}
	return fmt.Sprintf("#%s", color)
}

func colorHexToRGB(hex string) (*rgb, error) {
	values, err := strconv.ParseInt(strings.Replace(hex, "#", "", -1), 16, 32)
	if err != nil {
		return nil, err
	}

	return &rgb{
		red:   uint8(values >> 16),
		green: uint8((values >> 8) & 0xFF),
		blue:  uint8(values & 0xFF),
	}, nil
}
