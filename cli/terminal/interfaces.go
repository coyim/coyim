package terminal

import "io"

type Terminal interface {
	ReadLine() (string, error)
	SetPrompt(string)
	Write([]byte) (int, error)
	SetSize(int, int) error
	SetBracketedPasteMode(bool)
	ReadPassword(string) (string, error)
}

type EscapeCodes struct {
	Black, Red, Green, Yellow, Blue, Magenta, Cyan, White []byte
	Reset                                                 []byte
}

type Control interface {
	NewTerminal(io.ReadWriter, string) Terminal
	ErrPasteIndicator() error
	GetSize(int) (int, int, error)
	MakeRaw(int) (interface{}, error)
	Restore(int, interface{}) error
	SetAutoCompleteCallback(Terminal, func(string, int, rune) (string, int, bool))
	Escape(Terminal) EscapeCodes
}

type ControlFactory func() Control
