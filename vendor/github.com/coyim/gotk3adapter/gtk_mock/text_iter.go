package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockTextIter struct {
}

func (*MockTextIter) BackwardChar() bool {
	return false
}

func (*MockTextIter) BackwardChars(int) bool {
	return false
}

func (*MockTextIter) BackwardCursorPosition() bool {
	return false
}

func (*MockTextIter) BackwardCursorPositions(int) bool {
	return false
}

func (*MockTextIter) BackwardLine() bool {
	return false
}

func (*MockTextIter) BackwardLines(int) bool {
	return false
}

func (*MockTextIter) BackwardToTagToggle(gtki.TextTag) bool {
	return false
}

func (*MockTextIter) BackwardVisibleCursorPosition() bool {
	return false
}

func (*MockTextIter) BackwardVisibleCursorPositions(int) bool {
	return false
}

func (*MockTextIter) BackwardVisibleLine() bool {
	return false
}

func (*MockTextIter) BackwardVisibleLines(int) bool {
	return false
}

func (*MockTextIter) CanInsert(bool) bool {
	return false
}

func (*MockTextIter) Compare(gtki.TextIter) int {
	return 0
}

func (*MockTextIter) Editable(bool) bool {
	return false
}

func (*MockTextIter) EndsLine() bool {
	return false
}

func (*MockTextIter) EndsSentence() bool {
	return false
}

func (*MockTextIter) EndsTag(gtki.TextTag) bool {
	return false
}

func (*MockTextIter) EndsWord() bool {
	return false
}

func (*MockTextIter) Equal(gtki.TextIter) bool {
	return false
}

func (*MockTextIter) ForwardChar() bool {
	return false
}

func (*MockTextIter) ForwardChars(int) bool {
	return false
}

func (*MockTextIter) ForwardCursorPosition() bool {
	return false
}

func (*MockTextIter) ForwardCursorPositions(int) bool {
	return false
}

func (*MockTextIter) ForwardLine() bool {
	return false
}

func (*MockTextIter) ForwardLines(int) bool {
	return false
}

func (*MockTextIter) ForwardSentenceEnd() bool {
	return false
}

func (*MockTextIter) ForwardSentenceEnds(int) bool {
	return false
}

func (*MockTextIter) ForwardToEnd() {

}

func (*MockTextIter) ForwardToLineEnd() bool {
	return false
}

func (*MockTextIter) ForwardToTagToggle(gtki.TextTag) bool {
	return false
}

func (*MockTextIter) ForwardVisibleCursorPosition() bool {
	return false
}

func (*MockTextIter) ForwardVisibleCursorPositions(int) bool {
	return false
}

func (*MockTextIter) ForwardVisibleLine() bool {
	return false
}

func (*MockTextIter) ForwardVisibleLines(int) bool {
	return false
}

func (*MockTextIter) ForwardVisibleWordEnd() bool {
	return false
}

func (*MockTextIter) ForwardVisibleWordEnds(int) bool {
	return false
}

func (*MockTextIter) ForwardWordEnd() bool {
	return false
}

func (*MockTextIter) ForwardWordEnds(int) bool {
	return false
}

func (*MockTextIter) GetBuffer() gtki.TextBuffer {
	return nil
}

func (*MockTextIter) GetBytesInLine() int {
	return 0
}

func (*MockTextIter) GetChar() rune {
	return 0
}

func (*MockTextIter) GetCharsInLine() int {
	return 0
}

func (*MockTextIter) GetLine() int {
	return 0
}

func (*MockTextIter) GetLineIndex() int {
	return 0
}

func (*MockTextIter) GetLineOffset() int {
	return 0
}

func (*MockTextIter) GetOffset() int {
	return 0
}

func (*MockTextIter) GetSlice(gtki.TextIter) string {
	return ""
}

func (*MockTextIter) GetText(gtki.TextIter) string {
	return ""
}

func (*MockTextIter) GetVisibleLineIndex() int {
	return 0
}

func (*MockTextIter) GetVisibleLineOffset() int {
	return 0
}

func (*MockTextIter) GetVisibleSlice(gtki.TextIter) string {
	return ""
}

func (*MockTextIter) GetVisibleText(gtki.TextIter) string {
	return ""
}

func (*MockTextIter) HasTag(gtki.TextTag) bool {
	return false
}

func (*MockTextIter) InRange(gtki.TextIter, gtki.TextIter) bool {
	return false
}

func (*MockTextIter) InsideSentence() bool {
	return false
}

func (*MockTextIter) InsideWord() bool {
	return false
}

func (*MockTextIter) IsCursorPosition() bool {
	return false
}

func (*MockTextIter) IsEnd() bool {
	return false
}

func (*MockTextIter) IsStart() bool {
	return false
}

func (*MockTextIter) SetLine(int) {

}

func (*MockTextIter) SetLineIndex(int) {

}

func (*MockTextIter) SetLineOffset(int) {

}

func (*MockTextIter) SetOffset(int) {

}

func (*MockTextIter) SetVisibleLineIndex(int) {

}

func (*MockTextIter) SetVisibleLineOffset(int) {

}

func (*MockTextIter) StartsLine() bool {
	return false
}

func (*MockTextIter) StartsSentence() bool {
	return false
}

func (*MockTextIter) StartsWord() bool {
	return false
}

func (*MockTextIter) TogglesTag(gtki.TextTag) bool {
	return false
}
