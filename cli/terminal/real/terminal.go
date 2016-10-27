package real

import (
	"io"

	"github.com/twstrike/coyim/cli/terminal"
	ssh_terminal "golang.org/x/crypto/ssh/terminal"
)

type realTerminalControl struct{}

func (*realTerminalControl) NewTerminal(c io.ReadWriter, prompt string) terminal.Terminal {
	return ssh_terminal.NewTerminal(c, prompt)
}

func (*realTerminalControl) ErrPasteIndicator() error {
	return ssh_terminal.ErrPasteIndicator
}

func (*realTerminalControl) GetSize(fd int) (width, height int, err error) {
	return ssh_terminal.GetSize(fd)
}

func (*realTerminalControl) MakeRaw(fd int) (interface{}, error) {
	return ssh_terminal.MakeRaw(fd)
}

func (*realTerminalControl) Restore(fd int, state interface{}) error {
	realState := state.(*ssh_terminal.State)
	return ssh_terminal.Restore(fd, realState)
}

func (*realTerminalControl) SetAutoCompleteCallback(t terminal.Terminal, f func(string, int, rune) (string, int, bool)) {
	realT := t.(*ssh_terminal.Terminal)
	realT.AutoCompleteCallback = f
}

func (*realTerminalControl) Escape(t terminal.Terminal) terminal.EscapeCodes {
	realT := t.(*ssh_terminal.Terminal)
	e := realT.Escape

	return terminal.EscapeCodes{
		Black:   e.Black,
		Red:     e.Red,
		Green:   e.Green,
		Yellow:  e.Yellow,
		Blue:    e.Blue,
		Magenta: e.Magenta,
		Cyan:    e.Cyan,
		White:   e.White,
		Reset:   e.Reset,
	}
}

// Factory creates a new terminal.Control that is connected to a real terminal
func Factory() terminal.Control {
	return &realTerminalControl{}
}
