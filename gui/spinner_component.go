package gui

import "github.com/coyim/gotk3adapter/gtki"

type spinnerSize int

const (
	spinnerSizeNormal spinnerSize = 16
	spinnerSizeMedium spinnerSize = 24

	spinnerSizeDefault = spinnerSizeNormal
)

type spinner struct {
	s gtki.Spinner
}

func (u *gtkUI) newSpinnerComponent() *spinner {
	s, _ := g.gtk.SpinnerNew()
	sp := &spinner{s}
	sp.setSize(spinnerSizeDefault)
	return sp
}

func (sp *spinner) spinner() gtki.Spinner {
	return sp.s
}

// setSize MUST be called from the UI thread
func (sp *spinner) setSize(s spinnerSize) {
	sp.s.SetProperty("width_request", int(s))
	sp.s.SetProperty("height_request", int(s))
}

// show MUST be called from the ui thread
func (sp *spinner) show() {
	sp.s.Start()
	sp.s.Show()
}

// hide MUST be called from the ui thread
func (sp *spinner) hide() {
	sp.s.Stop()
	sp.s.Hide()
}
