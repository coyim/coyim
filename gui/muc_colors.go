package gui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	colorNone                       = cssColorReferenceFrom("none")
	colorTransparent                = cssColorReferenceFrom("transparent")
	colorThemeBase                  = cssColorReferenceFrom("@theme_base_color")
	colorThemeBackground            = cssColorReferenceFrom("@theme_bg_color")
	colorThemeForeground            = cssColorReferenceFrom("@theme_fg_color")
	colorThemeInsensitiveBackground = cssColorReferenceFrom("@insensitive_bg_color")
)

var colorThemeInsensitiveForeground = rgbFrom(131, 119, 119)

type mucColorSet struct {
	warningForeground                      cssColor
	warningBackground                      cssColor
	someoneJoinedForeground                cssColor
	someoneLeftForeground                  cssColor
	timestampForeground                    cssColor
	nicknameForeground                     cssColor
	subjectForeground                      cssColor
	infoMessageForeground                  cssColor
	messageForeground                      cssColor
	errorForeground                        cssColor
	configurationForeground                cssColor
	roomMessagesBackground                 cssColor
	roomMessagesBoxShadow                  cssColor
	roomNameDisabledForeground             cssColor
	roomSubjectForeground                  cssColor
	roomOverlaySolidBackground             cssColor
	roomOverlayContentSolidBackground      cssColor
	roomOverlayContentBackground           cssColor
	roomOverlayBackground                  cssColor
	roomOverlayContentForeground           cssColor
	roomOverlayContentBoxShadow            cssColor
	roomWarningsDialogBackground           cssColor
	roomWarningsDialogDecorationBackground cssColor
	roomWarningsDialogDecorationShadow     cssColor
	roomWarningsDialogHeaderBackground     cssColor
	roomWarningsDialogContentBackground    cssColor
	roomWarningsCurrentInfoForeground      cssColor
	roomNotificationsBackground            cssColor
	rosterGroupBackground                  cssColor
	rosterGroupForeground                  cssColor
	rosterOccupantRoleForeground           cssColor
	occupantStatusAvailableForeground      cssColor
	occupantStatusAvailableBackground      cssColor
	occupantStatusAvailableBorder          cssColor
	occupantStatusNotAvailableForeground   cssColor
	occupantStatusNotAvailableBackground   cssColor
	occupantStatusNotAvailableBorder       cssColor
	occupantStatusAwayForeground           cssColor
	occupantStatusAwayBackground           cssColor
	occupantStatusAwayBorder               cssColor
	occupantStatusBusyForeground           cssColor
	occupantStatusBusyBackground           cssColor
	occupantStatusBusyBorder               cssColor
	occupantStatusFreeForChatForeground    cssColor
	occupantStatusFreeForChatBackground    cssColor
	occupantStatusFreeForChatBorder        cssColor
	occupantStatusExtendedAwayForeground   cssColor
	occupantStatusExtendedAwayBackground   cssColor
	occupantStatusExtendedAwayBorder       cssColor
	infoBarDefaultBorderColor              cssColor
	infoBarTypeInfoBackgroundStart         cssColor
	infoBarTypeInfoBackgroundStop          cssColor
	infoBarTypeInfoTitle                   cssColor
	infoBarTypeInfoTime                    cssColor
	infoBarTypeWarningBackgroundStart      cssColor
	infoBarTypeWarningBackgroundStop       cssColor
	infoBarTypeWarningTitle                cssColor
	infoBarTypeWarningTime                 cssColor
	infoBarTypeQuestionBackgroundStart     cssColor
	infoBarTypeQuestionBackgroundStop      cssColor
	infoBarTypeQuestionTitle               cssColor
	infoBarTypeQuestionTime                cssColor
	infoBarTypeErrorBackgroundStart        cssColor
	infoBarTypeErrorBackgroundStop         cssColor
	infoBarTypeErrorTitle                  cssColor
	infoBarTypeErrorTime                   cssColor
	infoBarTypeOtherBackgroundStart        cssColor
	infoBarTypeOtherBackgroundStop         cssColor
	infoBarTypeOtherTitle                  cssColor
	infoBarTypeOtherTime                   cssColor
	infoBarButtonBackground                cssColor
	infoBarButtonForeground                cssColor
	infoBarButtonHoverBackground           cssColor
	infoBarButtonHoverForeground           cssColor
	infoBarButtonActiveBackground          cssColor
	infoBarButtonActiveForeground          cssColor
	entryErrorBackground                   cssColor
	entryErrorBorderShadow                 cssColor
	entryErrorBorder                       cssColor
	entryErrorLabel                        cssColor
	occupantLostConnection                 cssColor
	occupantRestablishConnection           cssColor
}

func (u *gtkUI) currentMUCColorSet() mucColorSet {
	if u.isDarkThemeVariant() {
		return u.defaultMUCDarkColorSet()
	}
	return u.defaultMUCLightColorSet()
}

func (u *gtkUI) defaultMUCLightColorSet() mucColorSet {
	return mucColorSet{
		warningForeground:                      rgbFrom(194, 23, 29),
		warningBackground:                      rgbFrom(254, 202, 202),
		errorForeground:                        rgbFrom(163, 7, 7),
		someoneJoinedForeground:                rgbFrom(41, 115, 22),
		someoneLeftForeground:                  rgbFrom(115, 22, 41),
		timestampForeground:                    colorThemeInsensitiveForeground,
		nicknameForeground:                     rgbFrom(57, 91, 163),
		subjectForeground:                      rgbFrom(0, 0, 128),
		infoMessageForeground:                  rgbFrom(57, 91, 163),
		messageForeground:                      rgbFrom(0, 0, 0),
		configurationForeground:                rgbFrom(154, 4, 191),
		roomMessagesBackground:                 colorThemeBase,
		roomMessagesBoxShadow:                  rgbaFrom(0, 0, 0, 0.35),
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBackground,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBackground,
		roomOverlayBackground:                  rgbaFrom(0, 0, 0, 0.5),
		roomOverlayContentForeground:           colorThemeForeground,
		roomOverlayContentBoxShadow:            rgbaFrom(0, 0, 0, 0.5),
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogDecorationShadow:     rgbaFrom(0, 0, 0, 0.15),
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		roomNotificationsBackground:            colorThemeBackground,
		rosterGroupBackground:                  rgbFrom(245, 245, 244),
		rosterGroupForeground:                  rgbFrom(28, 25, 23),
		rosterOccupantRoleForeground:           rgbFrom(168, 162, 158),
		occupantStatusAvailableForeground:      rgbFrom(22, 101, 52),
		occupantStatusAvailableBackground:      rgbFrom(240, 253, 244),
		occupantStatusAvailableBorder:          rgbFrom(22, 163, 74),
		occupantStatusNotAvailableForeground:   rgbFrom(30, 41, 59),
		occupantStatusNotAvailableBackground:   rgbFrom(248, 250, 252),
		occupantStatusNotAvailableBorder:       rgbFrom(71, 85, 105),
		occupantStatusAwayForeground:           rgbFrom(154, 52, 18),
		occupantStatusAwayBackground:           rgbFrom(255, 247, 237),
		occupantStatusAwayBorder:               rgbFrom(234, 88, 12),
		occupantStatusBusyForeground:           rgbFrom(159, 18, 57),
		occupantStatusBusyBackground:           rgbFrom(255, 241, 242),
		occupantStatusBusyBorder:               rgbFrom(190, 18, 60),
		occupantStatusFreeForChatForeground:    rgbFrom(30, 64, 175),
		occupantStatusFreeForChatBackground:    rgbFrom(239, 246, 255),
		occupantStatusFreeForChatBorder:        rgbFrom(29, 78, 216),
		occupantStatusExtendedAwayForeground:   rgbFrom(146, 64, 14),
		occupantStatusExtendedAwayBackground:   rgbFrom(255, 251, 235),
		occupantStatusExtendedAwayBorder:       rgbFrom(217, 119, 6),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         rgbFrom(63, 98, 18),
		infoBarTypeInfoBackgroundStop:          rgbFrom(77, 124, 15),
		infoBarTypeInfoTitle:                   rgbFrom(236, 254, 255),
		infoBarTypeInfoTime:                    rgbaFrom(236, 254, 255, 0.5),
		infoBarTypeWarningBackgroundStart:      rgbFrom(195, 149, 7),
		infoBarTypeWarningBackgroundStop:       rgbFrom(222, 173, 20),
		infoBarTypeWarningTitle:                rgbFrom(255, 247, 237),
		infoBarTypeWarningTime:                 rgbaFrom(255, 247, 237, 0.5),
		infoBarTypeQuestionBackgroundStart:     rgbFrom(234, 88, 12),
		infoBarTypeQuestionBackgroundStop:      rgbFrom(249, 115, 22),
		infoBarTypeQuestionTitle:               rgbFrom(254, 252, 232),
		infoBarTypeQuestionTime:                rgbaFrom(254, 252, 232, 0.5),
		infoBarTypeErrorBackgroundStart:        rgbFrom(185, 28, 28),
		infoBarTypeErrorBackgroundStop:         rgbFrom(203, 35, 35),
		infoBarTypeErrorTitle:                  rgbFrom(255, 241, 242),
		infoBarTypeErrorTime:                   rgbaFrom(255, 241, 242, 0.5),
		infoBarTypeOtherBackgroundStart:        rgbFrom(7, 89, 133),
		infoBarTypeOtherBackgroundStop:         rgbFrom(3, 105, 161),
		infoBarTypeOtherTitle:                  rgbFrom(240, 253, 250),
		infoBarTypeOtherTime:                   rgbaFrom(240, 253, 250, 0.5),
		infoBarButtonBackground:                rgbaFrom(0, 0, 0, 0.25),
		infoBarButtonForeground:                rgbFrom(255, 255, 255),
		infoBarButtonHoverBackground:           rgbaFrom(0, 0, 0, 0.35),
		infoBarButtonHoverForeground:           rgbFrom(255, 255, 255),
		infoBarButtonActiveBackground:          rgbaFrom(0, 0, 0, 0.45),
		infoBarButtonActiveForeground:          rgbFrom(255, 255, 255),
		entryErrorBackground:                   rgbFrom(255, 245, 246),
		entryErrorBorderShadow:                 rgbFrom(255, 127, 80),
		entryErrorBorder:                       rgbFrom(228, 70, 53),
		entryErrorLabel:                        rgbFrom(228, 70, 53),
		occupantLostConnection:                 rgbFrom(115, 22, 41),
		occupantRestablishConnection:           rgbFrom(41, 115, 22),
	}
}

func (u *gtkUI) defaultMUCDarkColorSet() mucColorSet {
	return mucColorSet{
		warningForeground:                      rgbFrom(194, 23, 29),
		warningBackground:                      rgbFrom(254, 202, 202),
		errorForeground:                        rgbFrom(209, 104, 96),
		someoneJoinedForeground:                rgbFrom(41, 115, 22),
		someoneLeftForeground:                  rgbFrom(115, 22, 41),
		timestampForeground:                    colorThemeInsensitiveForeground,
		nicknameForeground:                     rgbFrom(57, 91, 163),
		subjectForeground:                      rgbFrom(0, 0, 128),
		infoMessageForeground:                  rgbFrom(227, 66, 103),
		messageForeground:                      rgbFrom(0, 0, 0),
		configurationForeground:                rgbFrom(154, 4, 191),
		roomMessagesBackground:                 colorThemeBase,
		roomMessagesBoxShadow:                  rgbaFrom(0, 0, 0, 0.35),
		roomNameDisabledForeground:             colorThemeInsensitiveForeground,
		roomSubjectForeground:                  colorThemeInsensitiveForeground,
		roomOverlaySolidBackground:             colorThemeBase,
		roomOverlayContentSolidBackground:      colorTransparent,
		roomOverlayContentBackground:           colorThemeBase,
		roomOverlayBackground:                  rgbaFrom(0, 0, 0, 0.5),
		roomOverlayContentForeground:           rgbFrom(51, 51, 51),
		roomOverlayContentBoxShadow:            rgbaFrom(0, 0, 0, 0.5),
		roomWarningsDialogBackground:           colorNone,
		roomWarningsDialogDecorationBackground: colorThemeBackground,
		roomWarningsDialogDecorationShadow:     rgbaFrom(0, 0, 0, 0.15),
		roomWarningsDialogHeaderBackground:     colorNone,
		roomWarningsDialogContentBackground:    colorNone,
		roomWarningsCurrentInfoForeground:      colorThemeInsensitiveForeground,
		roomNotificationsBackground:            colorThemeBackground,
		rosterGroupBackground:                  rgbFrom(28, 25, 23),
		rosterGroupForeground:                  rgbFrom(250, 250, 249),
		rosterOccupantRoleForeground:           rgbFrom(231, 229, 228),
		occupantStatusAvailableForeground:      rgbFrom(22, 101, 52),
		occupantStatusAvailableBackground:      rgbFrom(240, 253, 244),
		occupantStatusAvailableBorder:          rgbFrom(22, 163, 74),
		occupantStatusNotAvailableForeground:   rgbFrom(30, 41, 59),
		occupantStatusNotAvailableBackground:   rgbFrom(248, 250, 252),
		occupantStatusNotAvailableBorder:       rgbFrom(71, 85, 105),
		occupantStatusAwayForeground:           rgbFrom(154, 52, 18),
		occupantStatusAwayBackground:           rgbFrom(255, 247, 237),
		occupantStatusAwayBorder:               rgbFrom(234, 88, 12),
		occupantStatusBusyForeground:           rgbFrom(159, 18, 57),
		occupantStatusBusyBackground:           rgbFrom(255, 241, 242),
		occupantStatusBusyBorder:               rgbFrom(190, 18, 60),
		occupantStatusFreeForChatForeground:    rgbFrom(30, 64, 175),
		occupantStatusFreeForChatBackground:    rgbFrom(239, 246, 255),
		occupantStatusFreeForChatBorder:        rgbFrom(29, 78, 216),
		occupantStatusExtendedAwayForeground:   rgbFrom(146, 64, 14),
		occupantStatusExtendedAwayBackground:   rgbFrom(255, 251, 235),
		occupantStatusExtendedAwayBorder:       rgbFrom(217, 119, 6),
		infoBarDefaultBorderColor:              colorThemeBackground,
		infoBarTypeInfoBackgroundStart:         rgbFrom(63, 98, 18),
		infoBarTypeInfoBackgroundStop:          rgbFrom(77, 124, 15),
		infoBarTypeInfoTitle:                   rgbFrom(236, 254, 255),
		infoBarTypeInfoTime:                    rgbaFrom(236, 254, 255, 0.5),
		infoBarTypeWarningBackgroundStart:      rgbFrom(195, 149, 7),
		infoBarTypeWarningBackgroundStop:       rgbFrom(222, 173, 20),
		infoBarTypeWarningTitle:                rgbFrom(255, 247, 237),
		infoBarTypeWarningTime:                 rgbaFrom(255, 247, 237, 0.5),
		infoBarTypeQuestionBackgroundStart:     rgbFrom(234, 88, 12),
		infoBarTypeQuestionBackgroundStop:      rgbFrom(249, 115, 22),
		infoBarTypeQuestionTitle:               rgbFrom(254, 252, 232),
		infoBarTypeQuestionTime:                rgbaFrom(254, 252, 232, 0.5),
		infoBarTypeErrorBackgroundStart:        rgbFrom(185, 28, 28),
		infoBarTypeErrorBackgroundStop:         rgbFrom(203, 35, 35),
		infoBarTypeErrorTitle:                  rgbFrom(255, 241, 242),
		infoBarTypeErrorTime:                   rgbaFrom(255, 241, 242, 0.5),
		infoBarTypeOtherBackgroundStart:        rgbFrom(7, 89, 133),
		infoBarTypeOtherBackgroundStop:         rgbFrom(3, 105, 161),
		infoBarTypeOtherTitle:                  rgbFrom(240, 253, 250),
		infoBarTypeOtherTime:                   rgbaFrom(240, 253, 250, 0.5),
		infoBarButtonBackground:                rgbaFrom(0, 0, 0, 0.25),
		infoBarButtonForeground:                rgbFrom(255, 255, 255),
		infoBarButtonHoverBackground:           rgbaFrom(0, 0, 0, 0.35),
		infoBarButtonHoverForeground:           rgbFrom(255, 255, 255),
		infoBarButtonActiveBackground:          rgbaFrom(0, 0, 0, 0.45),
		infoBarButtonActiveForeground:          rgbFrom(255, 255, 255),
		entryErrorBackground:                   rgbFrom(255, 245, 246),
		entryErrorBorderShadow:                 rgbFrom(255, 127, 80),
		entryErrorBorder:                       rgbFrom(228, 70, 53),
		entryErrorLabel:                        rgbFrom(228, 70, 53),
		occupantLostConnection:                 rgbFrom(115, 22, 41),
		occupantRestablishConnection:           rgbFrom(41, 115, 22),
	}
}

type colorValue float64

type rgb struct {
	red   colorValue
	green colorValue
	blue  colorValue
}

type rgba struct {
	*rgb
	alpha colorValue
}

type cssColorReference struct {
	ref string
}

func cssColorReferenceFrom(ref string) *cssColorReference {
	return &cssColorReference{
		ref: ref,
	}
}

func createColorValueFrom(v uint8) colorValue {
	return colorValue(float64(v) / 255)
}

func (v colorValue) toScaledValue() uint8 {
	return uint8(v * 255)
}

func colorValueFromHex(s string) colorValue {
	value, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return colorValue(0)
	}
	return createColorValueFrom(uint8(value))
}

// rgbFromHex will return an rgb object from either #xxxxxx or #xxx representation
// it returns nil if parsing fails
func rgbFromHex(spec string) *rgb {
	s := strings.TrimPrefix(spec, "#")
	switch len(s) {
	case 3:
		return &rgb{
			red:   colorValueFromHex(s[0:1]),
			green: colorValueFromHex(s[1:2]),
			blue:  colorValueFromHex(s[2:3]),
		}
	case 6:
		return &rgb{
			red:   colorValueFromHex(s[0:2]),
			green: colorValueFromHex(s[2:4]),
			blue:  colorValueFromHex(s[4:6]),
		}
	}
	return nil
}

func rgbFrom(r, g, b uint8) *rgb {
	return &rgb{
		red:   createColorValueFrom(r),
		green: createColorValueFrom(g),
		blue:  createColorValueFrom(b),
	}
}

func rgbaFrom(r, g, b uint8, a float64) *rgba {
	return &rgba{
		rgb:   rgbFrom(r, g, b),
		alpha: colorValue(a),
	}
}

func rgbFromPercent(r, g, b float64) *rgb {
	return &rgb{
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

func rgbFromGetters(v rgbaGetters) *rgb {
	return rgbFromPercent(v.GetRed(), v.GetGreen(), v.GetBlue())
}

func (r *rgb) toScaledColorValues() (uint8, uint8, uint8) {
	return r.red.toScaledValue(), r.green.toScaledValue(), r.blue.toScaledValue()
}

type cssColor interface {
	toCSS() string
}

type hexColor interface {
	toHex() string
}

type color interface {
	cssColor
	hexColor
}

func (r *rgba) String() string {
	return r.String()
}

func (r *rgba) toCSS() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %f)",
		r.red.toScaledValue(),
		r.green.toScaledValue(),
		r.blue.toScaledValue(),
		float64(r.alpha),
	)
}

func (r *rgb) toHex() string {
	return fmt.Sprintf("#%02x%02x%02x", r.red.toScaledValue(), r.green.toScaledValue(), r.blue.toScaledValue())
}

func (r *rgb) String() string {
	return r.toCSS()
}

func (r *rgb) toCSS() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", r.red.toScaledValue(), r.green.toScaledValue(), r.blue.toScaledValue())
}

func (r *cssColorReference) String() string {
	return r.toCSS()
}

func (r *cssColorReference) toCSS() string {
	return r.ref
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
