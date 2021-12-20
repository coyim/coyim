package gui

import (
	"fmt"
	"math"
)

const (
	colorNone                       = "none"
	colorTransparent                = "transparent"
	colorThemeBase                  = "@theme_base_color"
	colorThemeBackground            = "@theme_bg_color"
	colorThemeForeground            = "@theme_fg_color"
	colorThemeInsensitiveBackground = "@insensitive_bg_color"
)

var colorThemeInsensitiveForeground = colorFormat(rgb{131, 119, 119}, 1)

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
		warningForeground:                      colorFormat(rgbFrom(194, 23, 29), 1),
		warningBackground:                      colorFormat(rgbFrom(254, 202, 202), 1),
		errorForeground:                        colorFormat(rgbFrom(163, 7, 7), 1),
		someoneJoinedForeground:                colorFormat(rgbFrom(41, 115, 22), 1),
		someoneLeftForeground:                  colorFormat(rgbFrom(115, 22, 41), 1),
		timestampForeground:                    colorThemeInsensitiveForeground,
		nicknameForeground:                     colorFormat(rgbFrom(57, 91, 163), 1),
		subjectForeground:                      colorFormat(rgbFrom(0, 0, 128), 1),
		infoMessageForeground:                  colorFormat(rgbFrom(57, 91, 163), 1),
		messageForeground:                      colorFormat(rgbFrom(0, 0, 0), 1),
		configurationForeground:                colorFormat(rgbFrom(154, 4, 191), 1),
		roomMessagesBackground:                 colorThemeBase,
		roomMessagesBoxShadow:                  colorFormat(rgbFrom(0, 0, 0), 0.35),
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBackground,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBackground,
		roomOverlayBackground:                  colorFormat(rgbFrom(0, 0, 0), 0.5),
		roomOverlayContentForeground:           colorThemeForeground,
		roomOverlayContentBoxShadow:            colorFormat(rgbFrom(0, 0, 0), 0.5),
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogDecorationShadow:     colorFormat(rgbFrom(0, 0, 0), 0.15),
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		roomNotificationsBackground:            colorThemeBackground,
		rosterGroupBackground:                  colorFormat(rgbFrom(245, 245, 244), 1),
		rosterGroupForeground:                  colorFormat(rgbFrom(28, 25, 23), 1),
		rosterOccupantRoleForeground:           colorFormat(rgbFrom(168, 162, 158), 1),
		occupantStatusAvailableForeground:      colorFormat(rgbFrom(22, 101, 52), 1),
		occupantStatusAvailableBackground:      colorFormat(rgbFrom(240, 253, 244), 1),
		occupantStatusAvailableBorder:          colorFormat(rgbFrom(22, 163, 74), 1),
		occupantStatusNotAvailableForeground:   colorFormat(rgbFrom(30, 41, 59), 1),
		occupantStatusNotAvailableBackground:   colorFormat(rgbFrom(248, 250, 252), 1),
		occupantStatusNotAvailableBorder:       colorFormat(rgbFrom(71, 85, 105), 1),
		occupantStatusAwayForeground:           colorFormat(rgbFrom(154, 52, 18), 1),
		occupantStatusAwayBackground:           colorFormat(rgbFrom(255, 247, 237), 1),
		occupantStatusAwayBorder:               colorFormat(rgbFrom(234, 88, 12), 1),
		occupantStatusBusyForeground:           colorFormat(rgbFrom(159, 18, 57), 1),
		occupantStatusBusyBackground:           colorFormat(rgbFrom(255, 241, 242), 1),
		occupantStatusBusyBorder:               colorFormat(rgbFrom(190, 18, 60), 1),
		occupantStatusFreeForChatForeground:    colorFormat(rgbFrom(30, 64, 175), 1),
		occupantStatusFreeForChatBackground:    colorFormat(rgbFrom(239, 246, 255), 1),
		occupantStatusFreeForChatBorder:        colorFormat(rgbFrom(29, 78, 216), 1),
		occupantStatusExtendedAwayForeground:   colorFormat(rgbFrom(146, 64, 14), 1),
		occupantStatusExtendedAwayBackground:   colorFormat(rgbFrom(255, 251, 235), 1),
		occupantStatusExtendedAwayBorder:       colorFormat(rgbFrom(217, 119, 6), 1),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         colorFormat(rgbFrom(63, 98, 18), 1),
		infoBarTypeInfoBackgroundStop:          colorFormat(rgbFrom(77, 124, 15), 1),
		infoBarTypeInfoTitle:                   colorFormat(rgbFrom(236, 254, 255), 1),
		infoBarTypeInfoTime:                    colorFormat(rgbFrom(236, 254, 255), 0.5),
		infoBarTypeWarningBackgroundStart:      colorFormat(rgbFrom(195, 149, 7), 1),
		infoBarTypeWarningBackgroundStop:       colorFormat(rgbFrom(222, 173, 20), 1),
		infoBarTypeWarningTitle:                colorFormat(rgbFrom(255, 247, 237), 1),
		infoBarTypeWarningTime:                 colorFormat(rgbFrom(255, 247, 237), 0.5),
		infoBarTypeQuestionBackgroundStart:     colorFormat(rgbFrom(234, 88, 12), 1),
		infoBarTypeQuestionBackgroundStop:      colorFormat(rgbFrom(249, 115, 22), 1),
		infoBarTypeQuestionTitle:               colorFormat(rgbFrom(254, 252, 232), 1),
		infoBarTypeQuestionTime:                colorFormat(rgbFrom(254, 252, 232), 0.5),
		infoBarTypeErrorBackgroundStart:        colorFormat(rgbFrom(185, 28, 28), 1),
		infoBarTypeErrorBackgroundStop:         colorFormat(rgbFrom(203, 35, 35), 1),
		infoBarTypeErrorTitle:                  colorFormat(rgbFrom(255, 241, 242), 1),
		infoBarTypeErrorTime:                   colorFormat(rgbFrom(255, 241, 242), 0.5),
		infoBarTypeOtherBackgroundStart:        colorFormat(rgbFrom(7, 89, 133), 1),
		infoBarTypeOtherBackgroundStop:         colorFormat(rgbFrom(3, 105, 161), 1),
		infoBarTypeOtherTitle:                  colorFormat(rgbFrom(240, 253, 250), 1),
		infoBarTypeOtherTime:                   colorFormat(rgbFrom(240, 253, 250), 0.5),
		infoBarButtonBackground:                colorFormat(rgbFrom(0, 0, 0), 0.25),
		infoBarButtonForeground:                colorFormat(rgbFrom(255, 255, 255), 1),
		infoBarButtonHoverBackground:           colorFormat(rgbFrom(0, 0, 0), 0.35),
		infoBarButtonHoverForeground:           colorFormat(rgbFrom(255, 255, 255), 1),
		infoBarButtonActiveBackground:          colorFormat(rgbFrom(0, 0, 0), 0.45),
		infoBarButtonActiveForeground:          colorFormat(rgbFrom(255, 255, 255), 1),
		entryErrorBackground:                   colorFormat(rgbFrom(255, 245, 246), 1),
		entryErrorBorderShadow:                 colorFormat(rgbFrom(255, 127, 80), 1),
		entryErrorBorder:                       colorFormat(rgbFrom(228, 70, 53), 1),
		entryErrorLabel:                        colorFormat(rgbFrom(228, 70, 53), 1),
		occupantLostConnection:                 colorFormat(rgbFrom(115, 22, 41), 1),
		occupantRestablishConnection:           colorFormat(rgbFrom(41, 115, 22), 1),
	}
}

func (u *gtkUI) defaultMUCDarkColorSet() mucColorSet {
	return mucColorSet{
		warningForeground:                      colorFormat(rgbFrom(194, 23, 29), 1),
		warningBackground:                      colorFormat(rgbFrom(254, 202, 202), 1),
		errorForeground:                        colorFormat(rgbFrom(209, 104, 96), 1),
		someoneJoinedForeground:                colorFormat(rgbFrom(41, 115, 22), 1),
		someoneLeftForeground:                  colorFormat(rgbFrom(115, 22, 41), 1),
		timestampForeground:                    colorThemeInsensitiveForeground,
		nicknameForeground:                     colorFormat(rgbFrom(57, 91, 163), 1),
		subjectForeground:                      colorFormat(rgbFrom(0, 0, 128), 1),
		infoMessageForeground:                  colorFormat(rgbFrom(227, 66, 103), 1),
		messageForeground:                      colorFormat(rgbFrom(0, 0, 0), 1),
		configurationForeground:                colorFormat(rgbFrom(154, 4, 191), 1),
		roomMessagesBackground:                 colorThemeBase,
		roomMessagesBoxShadow:                  colorFormat(rgbFrom(0, 0, 0), 0.35),
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBase,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBase,
		roomOverlayBackground:                  colorFormat(rgbFrom(0, 0, 0), 0.5),
		roomOverlayContentForeground:           colorFormat(rgbFrom(51, 51, 51), 1),
		roomOverlayContentBoxShadow:            colorFormat(rgbFrom(0, 0, 0), 0.5),
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogDecorationShadow:     colorFormat(rgbFrom(0, 0, 0), 0.15),
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		roomNotificationsBackground:            colorThemeBackground,
		rosterGroupBackground:                  colorFormat(rgbFrom(28, 25, 23), 1),
		rosterGroupForeground:                  colorFormat(rgbFrom(250, 250, 249), 1),
		rosterOccupantRoleForeground:           colorFormat(rgbFrom(231, 229, 228), 1),
		occupantStatusAvailableForeground:      colorFormat(rgbFrom(22, 101, 52), 1),
		occupantStatusAvailableBackground:      colorFormat(rgbFrom(240, 253, 244), 1),
		occupantStatusAvailableBorder:          colorFormat(rgbFrom(22, 163, 74), 1),
		occupantStatusNotAvailableForeground:   colorFormat(rgbFrom(30, 41, 59), 1),
		occupantStatusNotAvailableBackground:   colorFormat(rgbFrom(248, 250, 252), 1),
		occupantStatusNotAvailableBorder:       colorFormat(rgbFrom(71, 85, 105), 1),
		occupantStatusAwayForeground:           colorFormat(rgbFrom(154, 52, 18), 1),
		occupantStatusAwayBackground:           colorFormat(rgbFrom(255, 247, 237), 1),
		occupantStatusAwayBorder:               colorFormat(rgbFrom(234, 88, 12), 1),
		occupantStatusBusyForeground:           colorFormat(rgbFrom(159, 18, 57), 1),
		occupantStatusBusyBackground:           colorFormat(rgbFrom(255, 241, 242), 1),
		occupantStatusBusyBorder:               colorFormat(rgbFrom(190, 18, 60), 1),
		occupantStatusFreeForChatForeground:    colorFormat(rgbFrom(30, 64, 175), 1),
		occupantStatusFreeForChatBackground:    colorFormat(rgbFrom(239, 246, 255), 1),
		occupantStatusFreeForChatBorder:        colorFormat(rgbFrom(29, 78, 216), 1),
		occupantStatusExtendedAwayForeground:   colorFormat(rgbFrom(146, 64, 14), 1),
		occupantStatusExtendedAwayBackground:   colorFormat(rgbFrom(255, 251, 235), 1),
		occupantStatusExtendedAwayBorder:       colorFormat(rgbFrom(217, 119, 6), 1),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         colorFormat(rgbFrom(63, 98, 18), 1),
		infoBarTypeInfoBackgroundStop:          colorFormat(rgbFrom(77, 124, 15), 1),
		infoBarTypeInfoTitle:                   colorFormat(rgbFrom(236, 254, 255), 1),
		infoBarTypeInfoTime:                    colorFormat(rgbFrom(236, 254, 255), 0.5),
		infoBarTypeWarningBackgroundStart:      colorFormat(rgbFrom(195, 149, 7), 1),
		infoBarTypeWarningBackgroundStop:       colorFormat(rgbFrom(222, 173, 20), 1),
		infoBarTypeWarningTitle:                colorFormat(rgbFrom(255, 247, 237), 1),
		infoBarTypeWarningTime:                 colorFormat(rgbFrom(255, 247, 237), 0.5),
		infoBarTypeQuestionBackgroundStart:     colorFormat(rgbFrom(234, 88, 12), 1),
		infoBarTypeQuestionBackgroundStop:      colorFormat(rgbFrom(249, 115, 22), 1),
		infoBarTypeQuestionTitle:               colorFormat(rgbFrom(254, 252, 232), 1),
		infoBarTypeQuestionTime:                colorFormat(rgbFrom(254, 252, 232), 0.5),
		infoBarTypeErrorBackgroundStart:        colorFormat(rgbFrom(185, 28, 28), 1),
		infoBarTypeErrorBackgroundStop:         colorFormat(rgbFrom(203, 35, 35), 1),
		infoBarTypeErrorTitle:                  colorFormat(rgbFrom(255, 241, 242), 1),
		infoBarTypeErrorTime:                   colorFormat(rgbFrom(255, 241, 242), 0.5),
		infoBarTypeOtherBackgroundStart:        colorFormat(rgbFrom(7, 89, 133), 1),
		infoBarTypeOtherBackgroundStop:         colorFormat(rgbFrom(3, 105, 161), 1),
		infoBarTypeOtherTitle:                  colorFormat(rgbFrom(240, 253, 250), 1),
		infoBarTypeOtherTime:                   colorFormat(rgbFrom(240, 253, 250), 0.5),
		infoBarButtonBackground:                colorFormat(rgbFrom(0, 0, 0), 0.25),
		infoBarButtonForeground:                colorFormat(rgbFrom(255, 255, 255), 1),
		infoBarButtonHoverBackground:           colorFormat(rgbFrom(0, 0, 0), 0.35),
		infoBarButtonHoverForeground:           colorFormat(rgbFrom(255, 255, 255), 1),
		infoBarButtonActiveBackground:          colorFormat(rgbFrom(0, 0, 0), 0.45),
		infoBarButtonActiveForeground:          colorFormat(rgbFrom(255, 255, 255), 1),
		entryErrorBackground:                   colorFormat(rgbFrom(255, 245, 246), 1),
		entryErrorBorderShadow:                 colorFormat(rgbFrom(255, 127, 80), 1),
		entryErrorBorder:                       colorFormat(rgbFrom(228, 70, 53), 1),
		entryErrorLabel:                        colorFormat(rgbFrom(228, 70, 53), 1),
		occupantLostConnection:                 colorFormat(rgbFrom(115, 22, 41), 1),
		occupantRestablishConnection:           colorFormat(rgbFrom(41, 115, 22), 1),
	}
}

type colorValue float64

type rgb struct {
	red   colorValue
	green colorValue
	blue  colorValue
}

func createColorValueFrom(v uint8) colorValue {
	return colorValue(float64(v) / 255)
}

func (v colorValue) toScaledValue() uint8 {
	return uint8(v * 255)
}

func rgbFrom(r, g, b uint8) rgb {
	return rgb{
		red:   createColorValueFrom(r),
		green: createColorValueFrom(g),
		blue:  createColorValueFrom(b),
	}
}

func rgbFromPercent(r, g, b float64) rgb {
	return rgb{
		red:   colorValue(r),
		green: colorValue(g),
		blue:  colorValue(b),
	}
}

type rgbaGetters interface {
	GetRed() float64
	GetGreen() float64
	GetBlue() float64
}

func rgbFromGetters(v rgbaGetters) rgb {
	return rgbFromPercent(v.GetRed(), v.GetGreen(), v.GetBlue())
}

func (r *rgb) toScaledColorValues() (uint8, uint8, uint8) {
	return r.red.toScaledValue(), r.green.toScaledValue(), r.blue.toScaledValue()
}

func colorFormat(c rgb, alpha float64) string {
	r, g, b := c.toScaledColorValues()
	if alpha == 1 {
		return fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
	}
	return fmt.Sprintf("rgba(%d, %d, %d, %f)", r, g, b, alpha)
}

const lightnessThreshold = 0.8

func (r *rgb) isDark() bool {
	return r.lightness() < lightnessThreshold
}

func (r *rgb) lightness() float64 {
	// We are using the formula found in https://en.wikipedia.org/wiki/HSL_and_HSV#From_RGB
	max := math.Max(math.Max(float64(r.red), float64(r.green)), float64(r.blue))
	min := math.Min(math.Min(float64(r.red), float64(r.green)), float64(r.blue))

	return (max + min) / 2
}
