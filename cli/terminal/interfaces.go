package terminal

import "io"

// Terminal represents terminal control
type Terminal interface {
	ReadLine() (string, error)
	ReadPassword(string) (string, error)
	SetBracketedPasteMode(bool)
	SetPrompt(string)
	SetSize(int, int) error
	Write([]byte) (int, error)
}

// EscapeCodes represents the escape codes used for a specific terminal
type EscapeCodes struct {
	Black, Red, Green, Yellow, Blue, Magenta, Cyan, White []byte
	Reset                                                 []byte
}

// Control represents terminal control
type Control interface {
	ErrPasteIndicator() error
	Escape(Terminal) EscapeCodes
	GetSize(int) (int, int, error)
	MakeRaw(int) (interface{}, error)
	NewTerminal(io.ReadWriter, string) Terminal
	Restore(int, interface{}) error
	SetAutoCompleteCallback(Terminal, func(string, int, rune) (string, int, bool))
}

// ControlFactory represents a function that can create a Control
type ControlFactory func() Control
