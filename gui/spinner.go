package gui

import "github.com/coyim/gotk3adapter/gtki"

type spinner struct {
	widget gtki.Spinner
}

func newSpinner() *spinner {
	s, _ := g.gtk.SpinnerNew()

	return &spinner{
		widget: s,
	}
}

func (s *spinner) getWidget() gtki.Spinner {
	return s.widget
}

func (s *spinner) show() {
	s.widget.Start()
	s.widget.Show()
}

func (s *spinner) hide() {
	s.widget.Stop()
	s.widget.Hide()
}
