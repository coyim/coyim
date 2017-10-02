package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
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

func (v *textView) ForwardDisplayLine(iter gtki.TextIter) bool {
	return v.internal.ForwardDisplayLine(unwrapTextIter(iter))
}

func (v *textView) BackwardDisplayLine(iter gtki.TextIter) bool {
	return v.internal.BackwardDisplayLine(unwrapTextIter(iter))
}

func (v *textView) ForwardDisplayLineEnd(iter gtki.TextIter) bool {
	return v.internal.ForwardDisplayLineEnd(unwrapTextIter(iter))
}

func (v *textView) BackwardDisplayLineStart(iter gtki.TextIter) bool {
	return v.internal.BackwardDisplayLineStart(unwrapTextIter(iter))
}

func (v *textView) StartsDisplayLine(iter gtki.TextIter) bool {
	return v.internal.StartsDisplayLine(unwrapTextIter(iter))
}

func (v *textView) MoveVisually(iter gtki.TextIter, count int) bool {
	return v.internal.MoveVisually(unwrapTextIter(iter), count)
}
