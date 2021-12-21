package gui

import "github.com/coyim/gotk3adapter/gtki"

type backgroundColorDetectionInvisibleListBox struct {
	lb gtki.ListBox `gtk-widget:"bg-color-detection-invisible-listbox"`
}

func newBackgroundColorDetectionInvisibleListBox() *backgroundColorDetectionInvisibleListBox {
	bgcd := &backgroundColorDetectionInvisibleListBox{}

	bgcd.initBuilder()

	return bgcd
}

func (bgcd *backgroundColorDetectionInvisibleListBox) initBuilder() {
	b := newBuilder("BackgroundColorDetectionInvisibleListBox")
	panicOnDevError(b.bindObjects(bgcd))
}
