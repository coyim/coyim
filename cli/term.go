package cli

import (
	"bytes"
	"fmt"
	"time"

	"github.com/twstrike/coyim/cli/terminal"
)

func updateTerminalSize(term terminal.Terminal, tc terminal.Control) {
	width, height, err := tc.GetSize(0)
	if err != nil {
		return
	}
	term.SetSize(width, height)
}

func info(term terminal.Terminal, tc terminal.Control, msg string) {
	terminalMessage(term, tc, tc.Escape(term).Blue, msg, false)
}

func warn(term terminal.Terminal, tc terminal.Control, msg string) {
	terminalMessage(term, tc, tc.Escape(term).Magenta, msg, false)
}

func alert(term terminal.Terminal, tc terminal.Control, msg string) {
	terminalMessage(term, tc, tc.Escape(term).Red, msg, false)
}

func critical(term terminal.Terminal, tc terminal.Control, msg string) {
	terminalMessage(term, tc, tc.Escape(term).Red, msg, true)
}

func terminalMessage(term terminal.Terminal, tc terminal.Control, color []byte, msg string, critical bool) {
	line := make([]byte, 0, len(msg)+16)

	line = append(line, ' ')
	line = append(line, color...)
	line = append(line, '*')
	line = append(line, tc.Escape(term).Reset...)
	line = append(line, []byte(fmt.Sprintf(" (%s) ", time.Now().Format(time.Kitchen)))...)
	if critical {
		line = append(line, tc.Escape(term).Red...)
	}
	line = appendTerminalEscaped(line, []byte(msg))
	if critical {
		line = append(line, tc.Escape(term).Reset...)
	}
	line = append(line, '\n')
	term.Write(line)
}

type lineLogger struct {
	term terminal.Terminal
	tc   terminal.Control
	buf  []byte
}

func (l *lineLogger) logLines(in []byte) []byte {
	for len(in) > 0 {
		if newLine := bytes.IndexByte(in, '\n'); newLine >= 0 {
			info(l.term, l.tc, string(in[:newLine]))
			in = in[newLine+1:]
		} else {
			break
		}
	}
	return in
}

func (l *lineLogger) Write(data []byte) (int, error) {
	origLen := len(data)

	if len(l.buf) == 0 {
		data = l.logLines(data)
	}

	if len(data) > 0 {
		l.buf = append(l.buf, data...)
	}

	l.buf = l.logLines(l.buf)
	return origLen, nil
}

// appendTerminalEscaped acts like append(), but breaks terminal escape
// sequences that may be in msg.

func appendTerminalEscaped(out, msg []byte) []byte {
	for _, c := range msg {
		if c == 127 || (c < 32 && c != '\t') {
			out = append(out, '?')
		} else {
			out = append(out, c)
		}
	}
	return out
}
