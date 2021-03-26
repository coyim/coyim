package gui

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

// mucStylesProvider is a representation of the styles that can be applied to specific muc-related interfaces.
// Please note that all methods of this struct MUST be called from the UI thread.
type mucStylesProvider struct {
	colors mucColorSet
}

var mucStyles *mucStylesProvider

func initMUCStyles(c mucColorSet) {
	mucStyles = &mucStylesProvider{
		colors: c,
	}
}

func (s *mucStylesProvider) setMessageScrolledWindowStyle(msw gtki.ScrolledWindow) {
	updateWithStyle(msw, providerWithStyle("scrolledwindow", style{
		"border": "none",
	}))
}

func (s *mucStylesProvider) setRoomLoadingInfoBarLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-size":   "16px",
		"font-weight": "bold",
	})
}

func (s *mucStylesProvider) setRoomRosterInfoStyle(b gtki.Box) {
	s.setWidgetStyles(b, styles{
		".occupant-nickname": style{
			"font-weight": "bold",
			"font-size":   "large",
		},
		".occupant-jid": style{
			"font-weight": "normal",
		},
		".occupant-status": {
			"font-size":     "small",
			"font-weight":   "bold",
			"padding":       "2px 6px 2px 6px",
			"border-width":  "1px",
			"border-style":  "solid",
			"border-radius": "200px",
		},
		".occupant-status-available": {
			"color":            s.colors.occupantStatusAvailableForeground,
			"background-color": s.colors.occupantStatusAvailableBackground,
			"border-color":     s.colors.occupantStatusAvailableBorder,
		},
		".occupant-status-not-available": {
			"color":            s.colors.occupantStatusNotAvailableForeground,
			"background-color": s.colors.occupantStatusNotAvailableBackground,
			"border-color":     s.colors.occupantStatusNotAvailableBorder,
		},
		".occupant-status-away": {
			"color":            s.colors.occupantStatusAwayForeground,
			"background-color": s.colors.occupantStatusAwayBackground,
			"border-color":     s.colors.occupantStatusAwayBorder,
		},
		".occupant-status-busy": {
			"color":            s.colors.occupantStatusBusyForeground,
			"background-color": s.colors.occupantStatusBusyBackground,
			"border-color":     s.colors.occupantStatusBusyBorder,
		},
		".occupant-status-free-for-chat": {
			"color":            s.colors.occupantStatusFreeForChatForeground,
			"background-color": s.colors.occupantStatusFreeForChatBackground,
			"border-color":     s.colors.occupantStatusFreeForChatBorder,
		},
		".occupant-status-extended-away": {
			"color":            s.colors.occupantStatusExtendedAwayForeground,
			"background-color": s.colors.occupantStatusExtendedAwayBackground,
			"border-color":     s.colors.occupantStatusExtendedAwayBorder,
		},
		".occupant-role-disabled-help": {
			"opacity":    0.5,
			"font-style": "italic",
		},
	})
}

func (s *mucStylesProvider) setRoomToolbarNameLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-size":   "22px",
		"font-weight": "bold",
	})
}

func (s *mucStylesProvider) setRoomToolbarSubjectLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-size":  "14px",
		"font-style": "italic",
		"color":      s.colors.roomSubjectForeground,
	})
}

func (s *mucStylesProvider) setRoomToolbarNameLabelDisabledStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"color": s.colors.roomNameDisabledForeground,
	})
}

func (s *mucStylesProvider) setRoomWarningsBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"padding": "12px",
	})
}

func (s *mucStylesProvider) setRoomWarningsMessageBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"color":            s.colors.roomWarningForeground,
		"background-color": s.colors.roomWarningBackground,
		"border":           s.border(1, "solid", s.colors.roomWarningBorder),
		"border-radius":    "4px",
		"padding":          "10px",
	})
}

func (s *mucStylesProvider) setRoomMessagesBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.colors.roomMessagesBackground,
		"box-shadow":       s.boxShadow("0 10px 20px", s.rgba(0, 0, 0, 0.35)),
	})
}

func (s *mucStylesProvider) setLabelBoldStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-weight": "bold",
	})
}

func (s *mucStylesProvider) setRoomOverlayMessagesBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.rgba(0, 0, 0, 0.5),
	})
}

func (s *mucStylesProvider) setRoomLoadingViewOverlayTransparentStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.hexToRGBA(s.colors.roomOverlayBackground, 0.5),
	})
}

func (s *mucStylesProvider) setRoomLoadingViewOverlayContentTransparentStyle(b gtki.Box) {
	s.setRoomLoadingViewOverlayContentBoxStyle(b)
}

func (s *mucStylesProvider) setRoomLoadingViewOverlayContentBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.colors.roomOverlayContentBackground,
		"color":            s.colors.roomOverlayContentForeground,
		"border-radius":    "12px",
		"padding":          "18px 24px",
		"box-shadow":       s.boxShadow("0 10px 20px", s.rgba(0, 0, 0, 0.5)),
	})
}

func (s *mucStylesProvider) setRoomLoadingViewOverlaySolidStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.colors.roomOverlaySolidBackground,
	})
}

func (s *mucStylesProvider) setRoomLoadingViewOverlayContentSolidStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.colors.roomOverlayContentSolidBackground,
		"color":            s.colors.roomOverlayContentForeground,
		"border-radius":    "0",
		"box-shadow":       "none",
	})
}

func (s *mucStylesProvider) setRoomConfigFormHelpLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-style": "italic",
	})
}

func (s *mucStylesProvider) setRoomConfigSummarySectionLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-weight": "bold",
	})
}

func (s *mucStylesProvider) setRoomConfigSummarySectionLinkButtonStyle(b gtki.LinkButton) {
	s.setWidgetStyles(b, styles{
		"button.link": {
			"padding":   "0px",
			"font-size": "medium",
		},
	})
}

func (s *mucStylesProvider) setRoomConfigSummaryRoomDescriptionLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-style": "italic",
	})
}

func (s *mucStylesProvider) setRoomConfigPageStyle(p gtki.Box) {
	s.setWidgetStyles(p, styles{
		".config-field-help": style{
			"font-style": "italic",
			"opacity":    "0.7",
		},
	})
}

func (s *mucStylesProvider) setHelpTextStyle(p gtki.Box) {
	s.setWidgetStyles(p, styles{
		".help-text": style{
			"font-style": "italic",
			"opacity":    "0.7",
		},
	})
}

func (s *mucStylesProvider) setDisableRoomStyle(p gtki.Box) {
	s.setBoxStyle(p, style{
		"opacity": "0.5",
	})
}

func (s *mucStylesProvider) setFormSectionLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-weight": "bold",
	})
}

func (s *mucStylesProvider) setRoomDialogErrorComponentHeaderStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-size":   "large",
		"font-weight": "bold",
	})
}

func (s *mucStylesProvider) setWidgetStyles(w gtki.Widget, st styles) {
	updateWithStyles(w, providerWithStyles(st))
}

func (s *mucStylesProvider) setWidgetStyle(w gtki.Widget, se string, st style) {
	updateWithStyle(w, providerWithStyle(se, st))
}

func (s *mucStylesProvider) setLabelStyle(l gtki.Label, st style) {
	s.setWidgetStyle(l, "label", st)
}

func (s *mucStylesProvider) setBoxStyle(b gtki.Box, st style) {
	s.setWidgetStyle(b, "box", st)
}

func (s *mucStylesProvider) setOverlayStyle(o gtki.Overlay, st style) {
	s.setWidgetStyle(o, "overlay", st)
}

func (s *mucStylesProvider) setEntryErrorStyle(e gtki.Entry) {
	s.setWidgetStyles(e, styles{
		".entry-error": style{
			"background-color": s.colors.entryErrorBackground,
			"border-color":     s.colors.entryErrorBorder,
			"box-shadow":       s.boxShadow("0 0 0 1px", s.colors.entryErrorBorderShadow),
		},
	})
}

func (s *mucStylesProvider) setNotificationTimeLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-style": "italic",
		"font-size":  "12px",
		"opacity":    "0.7",
	})
}

func (s *mucStylesProvider) setErrorLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"color": s.colors.entryErrorBorder,
	})
}

func (s *mucStylesProvider) border(size int, style, color string) string {
	return fmt.Sprintf("%dpx %s %s", size, style, color)
}

func (s *mucStylesProvider) rgba(r, g, b uint8, a float64) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %f)", r, g, b, a)
}

func (s *mucStylesProvider) hexToRGBA(hex string, a float64) string {
	rgb, err := s.colors.hexToRGB(hex)
	if err != nil {
		return s.rgba(0, 0, 0, 0.5)
	}

	return s.rgba(rgb.red, rgb.green, rgb.blue, a)
}

func (s *mucStylesProvider) boxShadow(shadowStyle, color string) string {
	return fmt.Sprintf("%s %s", shadowStyle, color)
}
