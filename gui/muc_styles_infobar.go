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

var infoBarClassNames = map[infoBarType]string{
	infoBarTypeInfo:     infoBarInfoClassName,
	infoBarTypeWarning:  infoBarWarningClassName,
	infoBarTypeQuestion: infoBarQuestionClassName,
	infoBarTypeError:    infoBarErrorClassName,
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
		"padding":       "4px 4px 4px 0",
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
		"border-radius": "200px",
		"border":        "none",
		"box-shadow":    "none",
		"text-shadow":   "none",
		"outline":       "none",
	}

	infoBarButtonCloseHoverCommonStyle = style{
		"border":     "none",
		"box-shadow": "none",
	}

	infoBarButtonCloseActiveCommonStyle = style{
		"border":     "none",
		"box-shadow": "none",
		"outline":    "none",
	}
)

type infoBarStyle struct {
	tp       infoBarType
	selector string
	colors   infoBarColorInfo
}

func (st *infoBarStyles) newInfoBarStyle(ib gtki.InfoBar) *infoBarStyle {
	tp := infoBarTypeForMessageType(ib.GetMessageType())

	selector := infoBarClassName
	if definedTypeClass, ok := infoBarClassNames[tp]; ok {
		selector = selector + definedTypeClass
	}

	return &infoBarStyle{
		tp:       tp,
		selector: selector,
		colors:   st.colorInfoBasedOnType(tp),
	}
}

func (ibst *infoBarStyle) childStyles() styles {
	infoBarTitleStyle := mergeStyles(infoBarTitleCommonStyle, style{
		"color": ibst.colors.titleColor,
	})

	infoBarActionsButtonStyle := mergeStyles(infoBarActionsButtonCommonStyle, style{
		"background": ibst.colors.buttonBackground,
		"color":      ibst.colors.buttonColor,
	})

	infoBarActionsButtonHoverStyle := mergeStyles(infoBarActionsButtonCommonStyle, style{
		"background": ibst.colors.buttonHoverBackground,
		"color":      ibst.colors.buttonHoverColor,
	})

	infoBarActionsButtonActiveStyle := mergeStyles(infoBarActionsButtonCommonStyle, style{
		"background": ibst.colors.buttonActiveBackground,
		"color":      ibst.colors.buttonActiveColor,
	})

	infoBarButtonCloseStyle := mergeStyles(infoBarButtonCloseCommonStyle, style{
		"background": ibst.colors.buttonBackground,
		"color":      ibst.colors.buttonColor,
	})

	infoBarButtonCloseHoverStyle := mergeStyles(infoBarButtonCloseHoverCommonStyle, style{
		"background": ibst.colors.buttonHoverBackground,
		"color":      ibst.colors.buttonHoverColor,
	})

	infoBarButtonCloseActiveStyle := mergeStyles(infoBarButtonCloseActiveCommonStyle, style{
		"background": ibst.colors.buttonActiveBackground,
		"color":      ibst.colors.buttonActiveColor,
	})

	return styles{
		infoBarRevealerBoxSelector:         infoBarRevealerBoxCommonStyle,
		infoBarContentSelector:             infoBarContentCommonStyle,
		infoBarTitleSelector:               infoBarTitleStyle,
		infoBarActionsButtonSelector:       infoBarActionsButtonStyle,
		infoBarActionsButtonLabelSelector:  infoBarActionsButtonLabelCommonStyle,
		infoBarActionsButtonHoverSelector:  infoBarActionsButtonHoverStyle,
		infoBarActionsButtonActiveSelector: infoBarActionsButtonActiveStyle,
		infoBarButtonCloseSelector:         infoBarButtonCloseStyle,
		infoBarButtonCloseHoverSelector:    infoBarButtonCloseHoverStyle,
		infoBarButtonCloseActiveSelector:   infoBarButtonCloseActiveStyle,
	}
}

func (ibst *infoBarStyle) styles() styles {
	infoBarStyle := mergeStyles(infoBarCommonStyle, style{
		"background": ibst.colors.background,
		"border":     fmt.Sprintf("1px solid %s", ibst.colors.borderColor),
	})

	nested := &nestedStyles{
		rootSelector: ibst.selector,
		rootStyle:    infoBarStyle,
		nestedStyles: ibst.childStyles(),
	}

	return nested.toStyles()
}

func (st *infoBarStyles) stylesFor(ib gtki.InfoBar) styles {
	ibStyle := st.newInfoBarStyle(ib)
	return ibStyle.styles()
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
		borderColor:            colorTransparent,
		buttonBackground:       c.infoBarButtonBackground,
		buttonColor:            c.infoBarButtonForeground,
		buttonHoverBackground:  c.infoBarButtonHoverBackground,
		buttonHoverColor:       c.infoBarButtonHoverForeground,
		buttonActiveBackground: c.infoBarButtonActiveBackground,
		buttonActiveColor:      c.infoBarButtonActiveForeground,
	}
}
