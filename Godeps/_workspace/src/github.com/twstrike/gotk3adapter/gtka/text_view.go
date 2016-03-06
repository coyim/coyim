package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type textView struct {
	*container
	internal *gtk.TextView
}

func wrapTextViewSimple(v *gtk.TextView) *textView {
	if v == nil {
		return nil
	}
	return &textView{wrapContainerSimple(&v.Container), v}
}

func wrapTextView(v *gtk.TextView, e error) (*textView, error) {
	return wrapTextViewSimple(v), e
}

func unwrapTextView(v gtki.TextView) *gtk.TextView {
	if v == nil {
		return nil
	}
	return v.(*textView).internal
}

func (v *textView) SetEditable(v1 bool) {
	v.internal.SetEditable(v1)
}

func (v *textView) SetCursorVisible(v1 bool) {
	v.internal.SetCursorVisible(v1)
}

func (v *textView) SetBuffer(v1 gtki.TextBuffer) {
	v.internal.SetBuffer(unwrapTextBuffer(v1))
}

func (v *textView) GetBuffer() (gtki.TextBuffer, error) {
	return wrapTextBuffer(v.internal.GetBuffer())
}
