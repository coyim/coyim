package gui

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

const (
	infoBarClassName         = ".infobar"
	infoBarInfoClassName     = ".info"
	infoBarWarningClassName  = ".warning"
	infoBarQuestionClassName = ".question"
	infoBarErrorClassName    = ".error"
)

type infoBarType int

var (
	infoBarTypeInfo     infoBarType
	infoBarTypeWarning  infoBarType
	infoBarTypeQuestion infoBarType
	infoBarTypeError    infoBarType
	infoBarTypeOther    infoBarType
)

var infoBarClassNames = map[infoBarType]string{
	infoBarTypeInfo:     infoBarInfoClassName,
	infoBarTypeWarning:  infoBarWarningClassName,
	infoBarTypeQuestion: infoBarQuestionClassName,
	infoBarTypeError:    infoBarErrorClassName,
}

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
	tp := infoBarType(ib.GetMessageType())
	colors := st.colorInfoBasedOnType(tp)

	nested := &nestedStyles{
		root: style{
			"background":    colors.background,
			"text-shadow":   "none",
			"font-weight":   "500",
			"padding":       "0 4px 0 0",
			"border-radius": "0",
			"border":        fmt.Sprintf("1px solid %s", colors.borderColor),
		},
		nested: styles{
			"revealer > box": style{
				"padding":    "0",
				"background": "none",
				"border":     "none",
				"box-shadow": "none",
			},
			".content": style{
				"text-shadow": "none",
				"padding":     "0",
				"background":  "none",
			},
			".title": style{
				"color":       colors.titleColor,
				"text-shadow": "none",
			},
			".actions button": {
				"background":    colors.buttonBackground,
				"color":         colors.buttonColor,
				"box-shadow":    "none",
				"padding":       "3px 12px",
				"border-radius": "200px",
				"border":        "none",
				"text-shadow":   "none",
				"font-size":     "small",
			},
			".actions button label": {
				"color": "inherit",
			},
			".actions button:hover": {
				"background": colors.buttonHoverBackground,
				"color":      colors.buttonHoverColor,
			},
			".actions button:active": {
				"background": colors.buttonActiveBackground,
				"color":      colors.buttonActiveColor,
			},
			"button.close": {
				"padding":     "0",
				"background":  "none",
				"border":      "none",
				"box-shadow":  "none",
				"text-shadow": "none",
				"outline":     "none",
			},
			"button.close:hover": {
				"background": "none",
				"border":     "none",
				"box-shadow": "none",
			},
			"button.close:active": {
				"background": "none",
				"border":     "none",
				"box-shadow": "none",
				"outline":    "none",
			},
		},
	}

	infoBarSelector := infoBarClassName
	if definedTypeClass, ok := infoBarClassNames[tp]; ok {
		infoBarSelector = infoBarSelector + definedTypeClass
	}

	return nested.toStyles(infoBarSelector)
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
