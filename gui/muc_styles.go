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

func (s *mucStylesProvider) setRoomRosterInfoNicknameLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-size":   "14px",
		"font-weight": "bold",
	})
}

func (s *mucStylesProvider) setRoomRosterInfoUserJIDLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-size": "12px",
	})
}

func (s *mucStylesProvider) setRoomRosterInfoStatusLabelStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"font-size":   "12px",
		"font-style":  "italic",
		"font-weight": "bold",
		"color":       s.colors.gray500,
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
		"color":      s.colors.gray500,
	})
}

func (s *mucStylesProvider) setRoomToolbarNameLabelDisabledStyle(l gtki.Label) {
	s.setLabelStyle(l, style{
		"color": s.colors.gray300,
	})
}

func (s *mucStylesProvider) setRoomWarningsBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"padding": "12px",
	})
}

func (s *mucStylesProvider) setRoomWarningsMessageBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"color":            s.colors.brown500,
		"background-color": s.colors.yellow200,
		"border":           s.border(1, "solid", s.colors.yellow600),
		"border-radius":    "4px",
		"padding":          "10px",
	})
}

func (s *mucStylesProvider) setRoomMessagesBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.colors.white,
		"box-shadow":       s.boxShadow("0 10px 20px", s.rgba(0, 0, 0, 0.35)),
	})
}

func (s *mucStylesProvider) setRoomOverlayMessagesBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.rgba(0, 0, 0, 0.5),
	})
}

func (s *mucStylesProvider) setRoomLoadingViewOverlayBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.hexToRGBA(s.colors.dark, 0.25),
	})
}

func (s *mucStylesProvider) setRoomLoadingViewOverlayContentBoxStyle(b gtki.Box) {
	s.setBoxStyle(b, style{
		"background-color": s.colors.light,
		"color":            s.colors.dark,
		"border-radius":    "6px",
		"box-shadow":       s.boxShadow("0 10px 20px", s.rgba(0, 0, 0, 0.5)),
	})
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
