package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gliba"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
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
