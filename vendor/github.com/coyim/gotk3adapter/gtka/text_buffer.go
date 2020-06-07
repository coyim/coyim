package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type textBuffer struct {
	*gliba.Object
	internal *gtk.TextBuffer
}

func WrapTextBufferSimple(v *gtk.TextBuffer) gtki.TextBuffer {
	if v == nil {
		return nil
	}
	return &textBuffer{gliba.WrapObjectSimple(v.Object), v}
}

func WrapTextBuffer(v *gtk.TextBuffer, e error) (gtki.TextBuffer, error) {
	return WrapTextBufferSimple(v), e
}

func UnwrapTextBuffer(v gtki.TextBuffer) *gtk.TextBuffer {
	if v == nil {
		return nil
	}
	return v.(*textBuffer).internal
}

func (v *textBuffer) ApplyTagByName(v1 string, v2, v3 gtki.TextIter) {
	v.internal.ApplyTagByName(v1, UnwrapTextIter(v2), UnwrapTextIter(v3))
}

func (v *textBuffer) GetCharCount() int {
	return v.internal.GetCharCount()
}

func (v *textBuffer) GetLineCount() int {
	return v.internal.GetLineCount()
}

func (v *textBuffer) GetEndIter() gtki.TextIter {
	return WrapTextIterSimple(v.internal.GetEndIter())
}

func (v *textBuffer) GetIterAtOffset(v1 int) gtki.TextIter {
	return WrapTextIterSimple(v.internal.GetIterAtOffset(v1))
}

func (v *textBuffer) GetStartIter() gtki.TextIter {
	return WrapTextIterSimple(v.internal.GetStartIter())
}

func (v *textBuffer) Insert(v1 gtki.TextIter, v2 string) {
	v.internal.Insert(UnwrapTextIter(v1), v2)
}

func (v *textBuffer) InsertAtCursor(v1 string) {
	v.internal.InsertAtCursor(v1)
}

func (v *textBuffer) GetText(v1, v2 gtki.TextIter, v3 bool) string {
	vx1, _ := v.internal.GetText(UnwrapTextIter(v1), UnwrapTextIter(v2), v3)
	return vx1
}

func (v *textBuffer) SetText(v1 string) {
	v.internal.SetText(v1)
}

func (v *textBuffer) Delete(v1, v2 gtki.TextIter) {
	v.internal.Delete(UnwrapTextIter(v1), UnwrapTextIter(v2))
}

func (v *textBuffer) CreateMark(v1 string, v2 gtki.TextIter, v3 bool) gtki.TextMark {
	return WrapTextMarkSimple(v.internal.CreateMark(v1, UnwrapTextIter(v2), v3))
}

func (v *textBuffer) GetIterAtMark(v1 gtki.TextMark) gtki.TextIter {
	return WrapTextIterSimple(v.internal.GetIterAtMark(UnwrapTextMark(v1)))
}
