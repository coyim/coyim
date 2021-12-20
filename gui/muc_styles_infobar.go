package gui

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

const (
	infoBarComponentClassName = "infobar-component"
	infoBarClassName          = "infobar"
	infoBarInfoClassName      = "info"
	infoBarWarningClassName   = "warning"
	infoBarQuestionClassName  = "question"
	infoBarErrorClassName     = "error"
	infoBarOtherClassName     = "other"
)

const (
	infoBarComponentClassNameSelector = "." + infoBarComponentClassName
	infoBarClassNameSelector          = "." + infoBarClassName
	infoBarInfoClassNameSelector      = "." + infoBarInfoClassName
	infoBarWarningClassNameSelector   = "." + infoBarWarningClassName
	infoBarQuestionClassNameSelector  = "." + infoBarQuestionClassName
	infoBarErrorClassNameSelector     = "." + infoBarErrorClassName
	infoBarOtherClassNameSelector     = "." + infoBarOtherClassName
)

var infoBarClassNames = map[infoBarType]string{
	infoBarTypeInfo:     infoBarInfoClassNameSelector,
	infoBarTypeWarning:  infoBarWarningClassNameSelector,
	infoBarTypeQuestion: infoBarQuestionClassNameSelector,
	infoBarTypeError:    infoBarErrorClassNameSelector,
	infoBarTypeOther:    infoBarOtherClassNameSelector,
}

type infoBarColorInfo struct {
	background             cssColor
	titleColor             cssColor
	timeColor              cssColor
	borderColor            cssColor
	buttonBackground       cssColor
	buttonColor            cssColor
	buttonHoverBackground  cssColor
	buttonHoverColor       cssColor
	buttonActiveBackground cssColor
	buttonActiveColor      cssColor
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
	infoBarTimeSelector                = ".time"
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
		"padding":       "0",
		"border-radius": "0",
	}

	infoBarRevealerBoxCommonStyle = style{
		"padding":    "0px 6px",
		"background": "none",
		"border":     "none",
		"box-shadow": "none",
	}

	infoBarContentCommonStyle = style{
		"text-shadow": "none",
		"padding":     "6px 0px 6px 0px",
		"background":  "none",
	}

	infoBarTitleCommonStyle = style{
		"text-shadow": "none",
	}

	infoBarTimeCommonStyle = style{
		"font-style":  "italic",
		"font-size":   "x-small",
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

	// The following classname selector allows us to identify the infobars used in MUC.
	// In this way we avoid visual changes in the notifications used in other places.
	selector := infoBarClassName + infoBarComponentClassNameSelector
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

	infoBarTimeStyle := mergeStyles(infoBarTimeCommonStyle, style{
		"color": ibst.colors.timeColor,
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
		infoBarTimeSelector:                infoBarTimeStyle,
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
		"border":     fmt.Sprintf("1px solid %s", ibst.colors.borderColor.toCSS()),
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
	addClassStyle(infoBarComponentClassName, ib)
	if ib.GetMessageType() == gtki.MESSAGE_OTHER {
		addClassStyle(infoBarOtherClassName, ib)
	}

	st := s.infoBarStyles.stylesFor(ib)
	s.setWidgetStyles(ib, st)
}

func infoBarTypeColorsFromSet(tp infoBarType, c mucColorSet) infoBarColorInfo {
	bgStart := c.infoBarTypeOtherBackgroundStart
	bgStop := c.infoBarTypeOtherBackgroundStop
	titleColor := c.infoBarTypeOtherTitle
	timeColor := c.infoBarTypeOtherTime

	switch tp {
	case infoBarTypeInfo:
		bgStart = c.infoBarTypeInfoBackgroundStart
		bgStop = c.infoBarTypeInfoBackgroundStop
		titleColor = c.infoBarTypeInfoTitle
		timeColor = c.infoBarTypeInfoTime

	case infoBarTypeWarning:
		bgStart = c.infoBarTypeWarningBackgroundStart
		bgStop = c.infoBarTypeWarningBackgroundStop
		titleColor = c.infoBarTypeWarningTitle
		timeColor = c.infoBarTypeWarningTime

	case infoBarTypeQuestion:
		bgStart = c.infoBarTypeQuestionBackgroundStart
		bgStop = c.infoBarTypeQuestionBackgroundStop
		titleColor = c.infoBarTypeQuestionTitle
		timeColor = c.infoBarTypeQuestionTime

	case infoBarTypeError:
		bgStart = c.infoBarTypeErrorBackgroundStart
		bgStop = c.infoBarTypeErrorBackgroundStop
		titleColor = c.infoBarTypeErrorTitle
		timeColor = c.infoBarTypeErrorTime
	}

	return infoBarColorInfo{
		background:             cssColorReferenceFrom(fmt.Sprintf("linear-gradient(0deg, %s 0%%, %s 100%%)", bgStart.toCSS(), bgStop.toCSS())),
		titleColor:             titleColor,
		timeColor:              timeColor,
		borderColor:            colorTransparent,
		buttonBackground:       c.infoBarButtonBackground,
		buttonColor:            c.infoBarButtonForeground,
		buttonHoverBackground:  c.infoBarButtonHoverBackground,
		buttonHoverColor:       c.infoBarButtonHoverForeground,
		buttonActiveBackground: c.infoBarButtonActiveBackground,
		buttonActiveColor:      c.infoBarButtonActiveForeground,
	}
}
