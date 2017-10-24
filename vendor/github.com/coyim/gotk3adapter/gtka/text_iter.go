package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type textIter struct {
	internal *gtk.TextIter
}

func wrapTextIterSimple(v *gtk.TextIter) *textIter {
	if v == nil {
		return nil
	}
	return &textIter{v}
}

func wrapTextIter(v *gtk.TextIter, e error) (*textIter, error) {
	return wrapTextIterSimple(v), e
}

func unwrapTextIter(v gtki.TextIter) *gtk.TextIter {
	if v == nil {
		return nil
	}
	return v.(*textIter).internal
}

func (v *textIter) BackwardChar() bool {
	return v.internal.BackwardChar()
}

func (v *textIter) BackwardChars(v1 int) bool {
	return v.internal.BackwardChars(v1)
}

func (v *textIter) BackwardCursorPosition() bool {
	return v.internal.BackwardCursorPosition()
}

func (v *textIter) BackwardCursorPositions(v1 int) bool {
	return v.internal.BackwardCursorPositions(v1)
}

func (v *textIter) BackwardLine() bool {
	return v.internal.BackwardLine()
}

func (v *textIter) BackwardLines(v1 int) bool {
	return v.internal.BackwardLines(v1)
}

func (v *textIter) BackwardToTagToggle(v1 gtki.TextTag) bool {
	return v.internal.BackwardToTagToggle(unwrapTextTag(v1))
}

func (v *textIter) BackwardVisibleCursorPosition() bool {
	return v.internal.BackwardVisibleCursorPosition()
}

func (v *textIter) BackwardVisibleCursorPositions(v1 int) bool {
	return v.internal.BackwardVisibleCursorPositions(v1)
}

func (v *textIter) BackwardVisibleLine() bool {
	return v.internal.BackwardVisibleLine()
}

func (v *textIter) BackwardVisibleLines(v1 int) bool {
	return v.internal.BackwardVisibleLines(v1)
}

func (v *textIter) CanInsert(v1 bool) bool {
	return v.internal.CanInsert(v1)
}

func (v *textIter) Compare(v1 gtki.TextIter) int {
	return v.internal.Compare(unwrapTextIter(v1))
}

func (v *textIter) Editable(v1 bool) bool {
	return v.internal.Editable(v1)
}

func (v *textIter) EndsLine() bool {
	return v.internal.EndsLine()
}

func (v *textIter) EndsSentence() bool {
	return v.internal.EndsSentence()
}

func (v *textIter) EndsTag(v1 gtki.TextTag) bool {
	return v.internal.EndsTag(unwrapTextTag(v1))
}

func (v *textIter) EndsWord() bool {
	return v.internal.EndsWord()
}

func (v *textIter) Equal(v1 gtki.TextIter) bool {
	return v.internal.Equal(unwrapTextIter(v1))
}

func (v *textIter) ForwardChar() bool {
	return v.internal.ForwardChar()
}

func (v *textIter) ForwardChars(v1 int) bool {
	return v.internal.ForwardChars(v1)
}

func (v *textIter) ForwardCursorPosition() bool {
	return v.internal.ForwardCursorPosition()
}

func (v *textIter) ForwardCursorPositions(v1 int) bool {
	return v.internal.ForwardCursorPositions(v1)
}

func (v *textIter) ForwardLine() bool {
	return v.internal.ForwardLine()
}

func (v *textIter) ForwardLines(v1 int) bool {
	return v.internal.ForwardLines(v1)
}

func (v *textIter) ForwardSentenceEnd() bool {
	return v.internal.ForwardSentenceEnd()
}

func (v *textIter) ForwardSentenceEnds(v1 int) bool {
	return v.internal.ForwardSentenceEnds(v1)
}

func (v *textIter) ForwardToEnd() {
	v.internal.ForwardToEnd()
}

func (v *textIter) ForwardToLineEnd() bool {
	return v.internal.ForwardToLineEnd()
}

func (v *textIter) ForwardToTagToggle(v1 gtki.TextTag) bool {
	return v.internal.ForwardToTagToggle(unwrapTextTag(v1))
}

func (v *textIter) ForwardVisibleCursorPosition() bool {
	return v.internal.ForwardVisibleCursorPosition()
}

func (v *textIter) ForwardVisibleCursorPositions(v1 int) bool {
	return v.internal.ForwardVisibleCursorPositions(v1)
}

func (v *textIter) ForwardVisibleLine() bool {
	return v.internal.ForwardVisibleLine()
}

func (v *textIter) ForwardVisibleLines(v1 int) bool {
	return v.internal.ForwardVisibleLines(v1)
}

func (v *textIter) ForwardVisibleWordEnd() bool {
	return v.internal.ForwardVisibleWordEnd()
}

func (v *textIter) ForwardVisibleWordEnds(v1 int) bool {
	return v.internal.ForwardVisibleWordEnds(v1)
}

func (v *textIter) ForwardWordEnd() bool {
	return v.internal.ForwardWordEnd()
}

func (v *textIter) ForwardWordEnds(v1 int) bool {
	return v.internal.ForwardWordEnds(v1)
}

func (v *textIter) GetBuffer() gtki.TextBuffer {
	return wrapTextBufferSimple(v.internal.GetBuffer())
}

func (v *textIter) GetBytesInLine() int {
	return v.internal.GetBytesInLine()
}

func (v *textIter) GetChar() rune {
	return v.internal.GetChar()
}

func (v *textIter) GetCharsInLine() int {
	return v.internal.GetCharsInLine()
}

func (v *textIter) GetLine() int {
	return v.internal.GetLine()
}

func (v *textIter) GetLineIndex() int {
	return v.internal.GetLineIndex()
}

func (v *textIter) GetLineOffset() int {
	return v.internal.GetLineOffset()
}

func (v *textIter) GetOffset() int {
	return v.internal.GetOffset()
}

func (v *textIter) GetSlice(v1 gtki.TextIter) string {
	return v.internal.GetSlice(unwrapTextIter(v1))
}

func (v *textIter) GetText(v1 gtki.TextIter) string {
	return v.internal.GetText(unwrapTextIter(v1))
}

func (v *textIter) GetVisibleLineIndex() int {
	return v.internal.GetVisibleLineIndex()
}

func (v *textIter) GetVisibleLineOffset() int {
	return v.internal.GetVisibleLineOffset()
}

func (v *textIter) GetVisibleSlice(v1 gtki.TextIter) string {
	return v.internal.GetVisibleSlice(unwrapTextIter(v1))
}

func (v *textIter) GetVisibleText(v1 gtki.TextIter) string {
	return v.internal.GetVisibleText(unwrapTextIter(v1))
}

func (v *textIter) HasTag(v1 gtki.TextTag) bool {
	return v.internal.HasTag(unwrapTextTag(v1))
}

func (v *textIter) InRange(v1 gtki.TextIter, v2 gtki.TextIter) bool {
	return v.internal.InRange(unwrapTextIter(v1), unwrapTextIter(v2))
}

func (v *textIter) InsideSentence() bool {
	return v.internal.InsideSentence()
}

func (v *textIter) InsideWord() bool {
	return v.internal.InsideWord()
}

func (v *textIter) IsCursorPosition() bool {
	return v.internal.IsCursorPosition()
}

func (v *textIter) IsEnd() bool {
	return v.internal.IsEnd()
}

func (v *textIter) IsStart() bool {
	return v.internal.IsStart()
}

func (v *textIter) SetLine(v1 int) {
	v.internal.SetLine(v1)
}

func (v *textIter) SetLineIndex(v1 int) {
	v.internal.SetLineIndex(v1)
}

func (v *textIter) SetLineOffset(v1 int) {
	v.internal.SetLineOffset(v1)
}

func (v *textIter) SetOffset(v1 int) {
	v.internal.SetOffset(v1)
}

func (v *textIter) SetVisibleLineIndex(v1 int) {
	v.internal.SetVisibleLineIndex(v1)
}

func (v *textIter) SetVisibleLineOffset(v1 int) {
	v.internal.SetVisibleLineOffset(v1)
}

func (v *textIter) StartsLine() bool {
	return v.internal.StartsLine()
}

func (v *textIter) StartsSentence() bool {
	return v.internal.StartsSentence()
}

func (v *textIter) StartsWord() bool {
	return v.internal.StartsWord()
}

func (v *textIter) TogglesTag(v1 gtki.TextTag) bool {
	return v.internal.TogglesTag(unwrapTextTag(v1))
}
