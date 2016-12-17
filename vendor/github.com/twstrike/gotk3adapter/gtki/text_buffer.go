package gtki

import "github.com/twstrike/gotk3adapter/glibi"

type TextBuffer interface {
	glibi.Object

	ApplyTagByName(string, TextIter, TextIter)
	CreateMark(string, TextIter, bool) TextMark
	Delete(TextIter, TextIter)
	GetCharCount() int
	GetEndIter() TextIter
	GetIterAtMark(TextMark) TextIter
	GetIterAtOffset(int) TextIter
	GetLineCount() int
	GetStartIter() TextIter
	GetText(TextIter, TextIter, bool) string
	Insert(TextIter, string)
	InsertAtCursor(string)
}

func AssertTextBuffer(_ TextBuffer) {}
