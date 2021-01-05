package gui

import "github.com/coyim/gotk3adapter/gtki"

type spinner struct {
	s gtki.Spinner
}

func (u *gtkUI) newSpinnerComponent() *spinner {
	s, _ := g.gtk.SpinnerNew()
	return &spinner{s}
}

func (sp *spinner) widget() gtki.Spinner {
	return sp.s
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
