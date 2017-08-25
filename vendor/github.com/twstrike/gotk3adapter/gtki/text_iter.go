package gtki

type TextIter interface {
	BackwardChar() bool
	BackwardChars(int) bool
	BackwardCursorPosition() bool
	BackwardCursorPositions(int) bool
	BackwardLine() bool
	BackwardLines(int) bool
	BackwardToTagToggle(TextTag) bool
	BackwardVisibleCursorPosition() bool
	BackwardVisibleCursorPositions(int) bool
	BackwardVisibleLine() bool
	BackwardVisibleLines(int) bool
	CanInsert(bool) bool
	Compare(TextIter) int
	Editable(bool) bool
	EndsLine() bool
	EndsSentence() bool
	EndsTag(TextTag) bool
	EndsWord() bool
	Equal(TextIter) bool
	ForwardChar() bool
	ForwardChars(int) bool
	ForwardCursorPosition() bool
	ForwardCursorPositions(int) bool
	ForwardLine() bool
	ForwardLines(int) bool
	ForwardSentenceEnd() bool
	ForwardSentenceEnds(int) bool
	ForwardToEnd()
	ForwardToLineEnd() bool
	ForwardToTagToggle(TextTag) bool
	ForwardVisibleCursorPosition() bool
	ForwardVisibleCursorPositions(int) bool
	ForwardVisibleLine() bool
	ForwardVisibleLines(int) bool
	ForwardVisibleWordEnd() bool
	ForwardVisibleWordEnds(v1 int) bool
	ForwardWordEnd() bool
	ForwardWordEnds(int) bool
	GetBuffer() TextBuffer
	GetBytesInLine() int
	GetChar() rune
	GetCharsInLine() int
	GetLine() int
	GetLineIndex() int
	GetLineOffset() int
	GetOffset() int
	GetSlice(TextIter) string
	GetText(TextIter) string
	GetVisibleLineIndex() int
	GetVisibleLineOffset() int
	GetVisibleSlice(TextIter) string
	GetVisibleText(TextIter) string
	HasTag(TextTag) bool
	InRange(TextIter, TextIter) bool
	InsideSentence() bool
	InsideWord() bool
	IsCursorPosition() bool
	IsEnd() bool
	IsStart() bool
	SetLine(int)
	SetLineIndex(int)
	SetLineOffset(int)
	SetOffset(int)
	SetVisibleLineIndex(int)
	SetVisibleLineOffset(int)
	StartsLine() bool
	StartsSentence() bool
	StartsWord() bool
	TogglesTag(TextTag) bool
}

func AssertTextIter(_ TextIter) {}
