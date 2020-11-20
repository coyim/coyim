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
	updateWithStyle(l, providerWithStyle("label", style{
		"font-size":   "16px",
		"font-weight": "bold",
	}))
}

func (s *mucStylesProvider) setRoomRosterInfoNicknameLabelStyle(l gtki.Label) {
	updateWithStyle(l, providerWithStyle("label", style{
		"font-size":   "14px",
		"font-weight": "bold",
	}))
}

func (s *mucStylesProvider) setRoomRosterInfoUserJIDLabelStyle(l gtki.Label) {
	updateWithStyle(l, providerWithStyle("label", style{
		"font-size": "12px",
	}))
}

func (s *mucStylesProvider) setRoomRosterInfoStatusLabelStyle(l gtki.Label) {
	updateWithStyle(l, providerWithStyle("label", style{
		"font-size":   "12px",
		"font-style":  "italic",
		"font-weight": "bold",
		"color":       s.colors.gray500,
	}))
}

func (s *mucStylesProvider) setRoomToolbarNameLabelStyle(l gtki.Label) {
	updateWithStyle(l, providerWithStyle("label", style{
		"font-size":   "22px",
		"font-weight": "bold",
	}))
}

func (s *mucStylesProvider) setRoomToolbarSubjectLabelStyle(l gtki.Label) {
	updateWithStyle(l, providerWithStyle("label", style{
		"font-size":  "14px",
		"font-style": "italic",
		"color":      s.colors.gray500,
	}))
}

func (s *mucStylesProvider) setRoomToolbarNameLabelDisabledStyle(l gtki.Label) {
	updateWithStyle(l, providerWithStyle("label", style{
		"color": s.colors.gray300,
	}))
}

func (s *mucStylesProvider) setRoomWarningsBoxStyle(b gtki.Box) {
	updateWithStyle(b, providerWithStyle("box", style{
		"padding": "12px",
	}))
}

func (s *mucStylesProvider) setRoomWarningsMessageBoxStyle(b gtki.Box) {
	updateWithStyle(b, providerWithStyle("box", style{
		"color":            s.colors.brown500,
		"background-color": s.colors.yellow200,
		"border":           border(1, "solid", s.colors.yellow600),
		"border-radius":    "4px",
		"padding":          "10px",
	}))
}

func (s *mucStylesProvider) setRoomMessagesBoxStyle(b gtki.Box) {
	updateWithStyle(b, providerWithStyle("box", style{
		"background-color": s.colors.white,
		"box-shadow":       "0 10px 20px rgba(0, 0, 0, 0.35)",
	}))
}

func (s *mucStylesProvider) setRoomOverlayMessagesBoxStyle(b gtki.Box) {
	updateWithStyle(b, providerWithStyle("box", style{
		"background-color": "rgba(0, 0, 0, 0.5)",
	}))
}

func border(size int, style, color string) string {
	return fmt.Sprintf("%dpx %s %s", size, style, color)
}
