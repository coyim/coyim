package gui

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

type infoBarColorInfo struct {
	background             string
	titleColor             string
	buttonBackground       string
	buttonColor            string
	buttonHoverBackground  string
	buttonHoverColor       string
	buttonActiveBackground string
	buttonActiveColor      string
}

type infoBarColorStyles map[gtki.MessageType]infoBarColorInfo

func newInfoBarColorStyles(c mucColorSet) infoBarColorStyles {
	return infoBarColorStyles{
		gtki.MESSAGE_INFO:     infoBarTypeColorsFromSet(gtki.MESSAGE_INFO, c),
		gtki.MESSAGE_WARNING:  infoBarTypeColorsFromSet(gtki.MESSAGE_WARNING, c),
		gtki.MESSAGE_QUESTION: infoBarTypeColorsFromSet(gtki.MESSAGE_QUESTION, c),
		gtki.MESSAGE_ERROR:    infoBarTypeColorsFromSet(gtki.MESSAGE_ERROR, c),
		gtki.MESSAGE_OTHER:    infoBarTypeColorsFromSet(gtki.MESSAGE_OTHER, c),
	}
}

func infoBarTypeColorsFromSet(t gtki.MessageType, c mucColorSet) infoBarColorInfo {
	bgStart := c.infoBarTypeOtherBackgroundStart
	bgStop := c.infoBarTypeOtherBackgroundStop
	tc := c.infoBarTypeOtherTitle

	switch t {
	case gtki.MESSAGE_INFO:
		bgStart = c.infoBarTypeInfoBackgroundStart
		bgStop = c.infoBarTypeInfoBackgroundStop
		tc = c.infoBarTypeInfoTitle

	case gtki.MESSAGE_WARNING:
		bgStart = c.infoBarTypeWarningBackgroundStart
		bgStop = c.infoBarTypeWarningBackgroundStop
		tc = c.infoBarTypeWarningTitle

	case gtki.MESSAGE_QUESTION:
		bgStart = c.infoBarTypeQuestionBackgroundStart
		bgStop = c.infoBarTypeQuestionBackgroundStop
		tc = c.infoBarTypeQuestionTitle

	case gtki.MESSAGE_ERROR:
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

func (s *mucStylesProvider) setInfoBarStyle(ib gtki.InfoBar) {
	if st, ok := s.infoBarColorStyles[ib.GetMessageType()]; ok {
		s.setWidgetStyles(ib, styles{
			".infobar": style{
				"background":  st.background,
				"text-shadow": "none",
				"font-weight": "500",
				"padding":     "8px 10px",
				"border":      fmt.Sprintf("2px solid %s", s.colors.infoBarDefaultBorderColor),
			},
			".infobar .content": style{
				"text-shadow": "none",
			},
			".infobar .title": style{
				"color":       st.titleColor,
				"text-shadow": "none",
			},
			".infobar .actions button": {
				"background":    st.buttonBackground,
				"color":         st.buttonColor,
				"box-shadow":    "none",
				"padding":       "4px 12px",
				"border-radius": "200px",
				"border":        "none",
				"text-shadow":   "none",
				"font-size":     "small",
			},
			".infobar .actions button:hover": {
				"background": st.buttonHoverBackground,
				"color":      st.buttonHoverColor,
			},
			".infobar .actions button:active": {
				"background": st.buttonActiveBackground,
				"color":      st.buttonActiveColor,
			},
			".infobar button.close": {
				"padding":     "0",
				"background":  "none",
				"border":      "none",
				"box-shadow":  "none",
				"text-shadow": "none",
				"outline":     "none",
			},
			".infobar button.close:hover": {
				"background": "none",
				"border":     "none",
				"box-shadow": "none",
			},
			".infobar button.close:active": {
				"background": "none",
				"border":     "none",
				"box-shadow": "none",
				"outline":    "none",
			},
		})
	}
}
