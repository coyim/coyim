package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type textView struct {
	*container
	internal *gtk.TextView
}

func WrapTextViewSimple(v *gtk.TextView) gtki.TextView {
	if v == nil {
		return nil
	}
	return &textView{WrapContainerSimple(&v.Container).(*container), v}
}

func WrapTextView(v *gtk.TextView, e error) (gtki.TextView, error) {
	return WrapTextViewSimple(v), e
}

func UnwrapTextView(v gtki.TextView) *gtk.TextView {
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
	v.internal.SetBuffer(UnwrapTextBuffer(v1))
}

func (v *textView) GetBuffer() (gtki.TextBuffer, error) {
	return WrapTextBuffer(v.internal.GetBuffer())
}

func (v *textView) ForwardDisplayLine(iter gtki.TextIter) bool {
	return v.internal.ForwardDisplayLine(UnwrapTextIter(iter))
}

func (v *textView) BackwardDisplayLine(iter gtki.TextIter) bool {
	return v.internal.BackwardDisplayLine(UnwrapTextIter(iter))
}

func (v *textView) ForwardDisplayLineEnd(iter gtki.TextIter) bool {
	return v.internal.ForwardDisplayLineEnd(UnwrapTextIter(iter))
}

func (v *textView) BackwardDisplayLineStart(iter gtki.TextIter) bool {
	return v.internal.BackwardDisplayLineStart(UnwrapTextIter(iter))
}

func (v *textView) StartsDisplayLine(iter gtki.TextIter) bool {
	return v.internal.StartsDisplayLine(UnwrapTextIter(iter))
}

func (v *textView) SetJustification(justify gtki.Justification) {
	v.internal.SetJustification(gtk.Justification(justify))
}

func (v *textView) GetJustification() gtki.Justification {
	return gtki.Justification(v.internal.GetJustification())
}

func (v *textView) MoveVisually(iter gtki.TextIter, count int) bool {
	return v.internal.MoveVisually(UnwrapTextIter(iter), count)
}
