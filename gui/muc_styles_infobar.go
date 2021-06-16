package gui

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

const (
	infoBarClassName = ".infobar"
)

type infoBarType int

var (
	infoBarTypeInfo     infoBarType
	infoBarTypeWarning  infoBarType
	infoBarTypeQuestion infoBarType
	infoBarTypeError    infoBarType
	infoBarTypeOther    infoBarType
)

func initMUCInfoBarData() {
	infoBarTypeInfo = infoBarType(gtki.MESSAGE_INFO)
	infoBarTypeWarning = infoBarType(gtki.MESSAGE_WARNING)
	infoBarTypeQuestion = infoBarType(gtki.MESSAGE_QUESTION)
	infoBarTypeError = infoBarType(gtki.MESSAGE_ERROR)
	infoBarTypeOther = infoBarType(gtki.MESSAGE_OTHER)
}

type infoBarColorInfo struct {
	background             string
	titleColor             string
	borderColor            string
	buttonBackground       string
	buttonColor            string
	buttonHoverBackground  string
	buttonHoverColor       string
	buttonActiveBackground string
	buttonActiveColor      string
}

type infoBarColorStyles map[infoBarType]infoBarColorInfo

func newInfoBarColorStyles(c mucColorSet) infoBarColorStyles {
	return infoBarColorStyles{
		infoBarTypeInfo:     infoBarTypeColorsFromSet(infoBarTypeInfo, c),
		infoBarTypeWarning:  infoBarTypeColorsFromSet(infoBarTypeWarning, c),
		infoBarTypeQuestion: infoBarTypeColorsFromSet(infoBarTypeQuestion, c),
		infoBarTypeError:    infoBarTypeColorsFromSet(infoBarTypeError, c),
		infoBarTypeOther:    infoBarTypeColorsFromSet(infoBarTypeOther, c),
	}
}

type infoBarStyles struct {
	colorStyles infoBarColorStyles
}

func newInfoBarStyles(c mucColorSet) *infoBarStyles {
	return &infoBarStyles{
		colorStyles: newInfoBarColorStyles(c),
	}
}

func (st *infoBarStyles) colorInfoBasedOnType(tp infoBarType) infoBarColorInfo {
	if colors, ok := st.colorStyles[tp]; ok {
		return colors
	}
	return st.colorStyles[infoBarTypeOther]
}

func (st *infoBarStyles) stylesFor(ib gtki.InfoBar) styles {
	colors := st.colorInfoBasedOnType(infoBarType(ib.GetMessageType()))

	return styles{
		infoBarClassName: style{
			"background":  colors.background,
			"text-shadow": "none",
			"font-weight": "500",
			"padding":     "8px 10px",
			"border":      fmt.Sprintf("2px solid %s", colors.borderColor),
		},
		nestedCSSRules(infoBarClassName, ".content"): style{
			"text-shadow": "none",
		},
		nestedCSSRules(infoBarClassName, ".title"): style{
			"color":       colors.titleColor,
			"text-shadow": "none",
		},
		nestedCSSRules(infoBarClassName, ".actions button"): {
			"background":    colors.buttonBackground,
			"color":         colors.buttonColor,
			"box-shadow":    "none",
			"padding":       "4px 12px",
			"border-radius": "200px",
			"border":        "none",
			"text-shadow":   "none",
			"font-size":     "small",
		},
		nestedCSSRules(infoBarClassName, ".actions button:hover"): {
			"background": colors.buttonHoverBackground,
			"color":      colors.buttonHoverColor,
		},
		nestedCSSRules(infoBarClassName, ".actions button:active"): {
			"background": colors.buttonActiveBackground,
			"color":      colors.buttonActiveColor,
		},
		nestedCSSRules(infoBarClassName, "button.close"): {
			"padding":     "0",
			"background":  "none",
			"border":      "none",
			"box-shadow":  "none",
			"text-shadow": "none",
			"outline":     "none",
		},
		nestedCSSRules(infoBarClassName, "button.close:hover"): {
			"background": "none",
			"border":     "none",
			"box-shadow": "none",
		},
		nestedCSSRules(infoBarClassName, "button.close:active"): {
			"background": "none",
			"border":     "none",
			"box-shadow": "none",
			"outline":    "none",
		},
	}
}

func (s *mucStylesProvider) setInfoBarStyle(ib gtki.InfoBar) {
	st := s.infoBarStyles.stylesFor(ib)
	s.setWidgetStyles(ib, st)
}

func infoBarTypeColorsFromSet(tp infoBarType, c mucColorSet) infoBarColorInfo {
	bgStart := c.infoBarTypeOtherBackgroundStart
	bgStop := c.infoBarTypeOtherBackgroundStop
	tc := c.infoBarTypeOtherTitle

	switch tp {
	case infoBarTypeInfo:
		bgStart = c.infoBarTypeInfoBackgroundStart
		bgStop = c.infoBarTypeInfoBackgroundStop
		tc = c.infoBarTypeInfoTitle

	case infoBarTypeWarning:
		bgStart = c.infoBarTypeWarningBackgroundStart
		bgStop = c.infoBarTypeWarningBackgroundStop
		tc = c.infoBarTypeWarningTitle

	case infoBarTypeQuestion:
		bgStart = c.infoBarTypeQuestionBackgroundStart
		bgStop = c.infoBarTypeQuestionBackgroundStop
		tc = c.infoBarTypeQuestionTitle

	case infoBarTypeError:
		bgStart = c.infoBarTypeErrorBackgroundStart
		bgStop = c.infoBarTypeErrorBackgroundStop
		tc = c.infoBarTypeErrorTitle
	}

	return infoBarColorInfo{
		background:             fmt.Sprintf("linear-gradient(0deg, %s 0%%, %s 100%%)", bgStart, bgStop),
		titleColor:             tc,
		buttonBackground:       c.infoBarButtonBackground,
		buttonColor:            c.infoBarButtonForeground,
		buttonHoverBackground:  c.infoBarButtonHoverBackground,
		buttonHoverColor:       c.infoBarButtonHoverForeground,
		buttonActiveBackground: c.infoBarButtonActiveBackground,
		buttonActiveColor:      c.infoBarButtonActiveForeground,
	}
}
