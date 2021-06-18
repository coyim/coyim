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

var infoBarClassNames map[infoBarType]string

func initMUCInfoBarData() {
	infoBarTypeInfo = infoBarType(gtki.MESSAGE_INFO)
	infoBarTypeWarning = infoBarType(gtki.MESSAGE_WARNING)
	infoBarTypeQuestion = infoBarType(gtki.MESSAGE_QUESTION)
	infoBarTypeError = infoBarType(gtki.MESSAGE_ERROR)
	infoBarTypeOther = infoBarType(gtki.MESSAGE_OTHER)

	infoBarClassNames = map[infoBarType]string{
		infoBarTypeInfo:     infoBarInfoClassName,
		infoBarTypeWarning:  infoBarWarningClassName,
		infoBarTypeQuestion: infoBarQuestionClassName,
		infoBarTypeError:    infoBarErrorClassName,
	}
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

const (
	infoBarRevealerBoxSelector         = "revealer > box"
	infoBarContentSelector             = ".content"
	infoBarTitleSelector               = ".title"
	infoBarActionsButtonSelector       = ".actions button"
	infoBarActionsButtonLabelSelector  = ".actions button label"
	infoBarActionsButtonHoverSelector  = ".actions button:hover"
	infoBarActionsButtonActiveSelector = ".actions button:active"
	infoBarButtonCloseSelector         = "button.close"
	infoBarButtonCloseHoverSelector    = "button.close:hover"
	infoBarButtonCloseActiveSelector   = "button.close:active"
)

var (
	infoBarCommonStyle = style{
		"text-shadow":   "none",
		"font-weight":   "500",
		"padding":       "0 4px 0 0",
		"border-radius": "0",
	}

	infoBarRevealerBoxCommonStyle = style{
		"padding":    "0",
		"background": "none",
		"border":     "none",
		"box-shadow": "none",
	}

	infoBarContentCommonStyle = style{
		"text-shadow": "none",
		"padding":     "0",
		"background":  "none",
	}

	infoBarTitleCommonStyle = style{
		"text-shadow": "none",
	}

	infoBarActionsButtonCommonStyle = style{
		"box-shadow":    "none",
		"padding":       "3px 12px",
		"border-radius": "200px",
		"border":        "none",
		"text-shadow":   "none",
		"font-size":     "small",
	}

	infoBarActionsButtonLabelCommonStyle = style{
		"color": "inherit",
	}

	infoBarButtonCloseCommonStyle = style{
		"padding":     "0",
		"background":  "none",
		"border":      "none",
		"box-shadow":  "none",
		"text-shadow": "none",
		"outline":     "none",
	}

	infoBarButtonCloseHoverCommonStyle = style{
		"background": "none",
		"border":     "none",
		"box-shadow": "none",
	}

	infoBarButtonCloseActiveCommonStyle = style{
		"background": "none",
		"border":     "none",
		"box-shadow": "none",
		"outline":    "none",
	}
)

func (st *infoBarStyles) stylesFor(ib gtki.InfoBar) styles {
	tp := infoBarType(ib.GetMessageType())
	colors := st.colorInfoBasedOnType(tp)

	infoBarStyle := mergeStyles(infoBarCommonStyle, style{
		"background": colors.background,
		"border":     fmt.Sprintf("1px solid %s", colors.borderColor),
	})

	infoBarTitleStyle := mergeStyles(infoBarTitleCommonStyle, style{
		"color": colors.titleColor,
	})

	infoBarActionsButtonStyle := mergeStyles(infoBarActionsButtonCommonStyle, style{
		"background": colors.buttonBackground,
		"color":      colors.buttonColor,
	})

	infoBarActionsButtonHoverStyle := mergeStyles(infoBarActionsButtonCommonStyle, style{
		"background": colors.buttonHoverBackground,
		"color":      colors.buttonHoverColor,
	})

	infoBarActionsButtonActiveStyle := mergeStyles(infoBarActionsButtonCommonStyle, style{
		"background": colors.buttonActiveBackground,
		"color":      colors.buttonActiveColor,
	})

	infoBarSelector := infoBarClassName
	if definedTypeClass, ok := infoBarClassNames[tp]; ok {
		infoBarSelector = infoBarSelector + definedTypeClass
	}

	nested := &nestedStyles{
		rootSelector: infoBarSelector,
		rootStyle:    infoBarStyle,
		nestedStyles: styles{
			infoBarRevealerBoxSelector:         infoBarRevealerBoxCommonStyle,
			infoBarContentSelector:             infoBarContentCommonStyle,
			infoBarTitleSelector:               infoBarTitleStyle,
			infoBarActionsButtonSelector:       infoBarActionsButtonStyle,
			infoBarActionsButtonLabelSelector:  infoBarActionsButtonLabelCommonStyle,
			infoBarActionsButtonHoverSelector:  infoBarActionsButtonHoverStyle,
			infoBarActionsButtonActiveSelector: infoBarActionsButtonActiveStyle,
			infoBarButtonCloseSelector:         infoBarButtonCloseCommonStyle,
			infoBarButtonCloseHoverSelector:    infoBarButtonCloseHoverCommonStyle,
			infoBarButtonCloseActiveSelector:   infoBarButtonCloseActiveCommonStyle,
		},
	}

	return nested.toStyles()
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
