package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/gotk3adapter/gliba"
	"github.com/twstrike/gotk3adapter/gtki"
)

type textBuffer struct {
	*gliba.Object
	internal *gtk.TextBuffer
}

func wrapTextBufferSimple(v *gtk.TextBuffer) *textBuffer {
	if v == nil {
		return nil
	}
	return &textBuffer{gliba.WrapObjectSimple(v.Object), v}
}

func wrapTextBuffer(v *gtk.TextBuffer, e error) (*textBuffer, error) {
	return wrapTextBufferSimple(v), e
}

func unwrapTextBuffer(v gtki.TextBuffer) *gtk.TextBuffer {
	if v == nil {
		return nil
	}
	return v.(*textBuffer).internal
}

func (v *textBuffer) ApplyTagByName(v1 string, v2, v3 gtki.TextIter) {
	v.internal.ApplyTagByName(v1, unwrapTextIter(v2), unwrapTextIter(v3))
}

func (v *textBuffer) GetCharCount() int {
	return v.internal.GetCharCount()
}

func (v *textBuffer) GetLineCount() int {
	return v.internal.GetLineCount()
}

func (v *textBuffer) GetEndIter() gtki.TextIter {
	return wrapTextIterSimple(v.internal.GetEndIter())
}

func (v *textBuffer) GetIterAtOffset(v1 int) gtki.TextIter {
	return wrapTextIterSimple(v.internal.GetIterAtOffset(v1))
}

func (v *textBuffer) GetStartIter() gtki.TextIter {
	return wrapTextIterSimple(v.internal.GetStartIter())
}

func (v *textBuffer) Insert(v1 gtki.TextIter, v2 string) {
	v.internal.Insert(unwrapTextIter(v1), v2)
}

func (v *textBuffer) InsertAtCursor(v1 string) {
	v.internal.InsertAtCursor(v1)
}

func (v *textBuffer) GetText(v1, v2 gtki.TextIter, v3 bool) string {
	vx1, _ := v.internal.GetText(unwrapTextIter(v1), unwrapTextIter(v2), v3)
	return vx1
}

func (v *textBuffer) Delete(v1, v2 gtki.TextIter) {
	v.internal.Delete(unwrapTextIter(v1), unwrapTextIter(v2))
}

func (v *textBuffer) CreateMark(v1 string, v2 gtki.TextIter, v3 bool) gtki.TextMark {
	return wrapTextMarkSimple(v.internal.CreateMark(v1, unwrapTextIter(v2), v3))
}

func (v *textBuffer) GetIterAtMark(v1 gtki.TextMark) gtki.TextIter {
	return wrapTextIterSimple(v.internal.GetIterAtMark(unwrapTextMark(v1)))
}
