package ui

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"sync"

	coyconf "github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
	"golang.org/x/net/html"
)

// OTRWhitespaceTagStart may be appended to plaintext messages to signal to the
// remote client that we support OTR. It should be followed by one of the
// version specific tags, below. See "Tagged plaintext messages" in
// http://www.cypherpunks.ca/otr/Protocol-v3-4.0.0.html.
var OTRWhitespaceTagStart = []byte("\x20\x09\x20\x20\x09\x09\x09\x09\x20\x09\x20\x09\x20\x09\x20\x20")

var OTRWhiteSpaceTagV1 = []byte("\x20\x09\x20\x09\x20\x20\x09\x20")
var OTRWhiteSpaceTagV2 = []byte("\x20\x20\x09\x09\x20\x20\x09\x20")
var OTRWhiteSpaceTagV3 = []byte("\x20\x20\x09\x09\x20\x20\x09\x09")

var OTRWhitespaceTag = append(OTRWhitespaceTagStart, OTRWhiteSpaceTagV2...)

type UI interface {
	Enroll(*coyconf.Config) bool
	AskForPassword(*coyconf.Config) (string, error)
	Loop()

	// callbacks
	RegisterCallback() xmpp.FormCallback
	NewOTRKeys(from string, conversation *otr3.Conversation)
	OTREnded(string)
	MessageReceived(from, timestamp string, encrypted bool, message []byte)
	IQReceived(string)
	RosterReceived(roster []xmpp.RosterEntry)

	ProcessPresence(*xmpp.ClientPresence)

	Info(string)
	Warn(string)
	Alert(string)
}

func StripHTML(msg []byte) (out []byte) {
	z := html.NewTokenizer(bytes.NewReader(msg))

loop:
	for {
		tt := z.Next()
		switch tt {
		case html.TextToken:
			out = append(out, z.Text()...)
		case html.ErrorToken:
			if err := z.Err(); err != nil && err != io.EOF {
				out = msg
				return
			}
			break loop
		}
	}

	return
}

// RosterEdit contains information about a pending roster edit. Roster edits
// occur by writing the roster to a file and inviting the user to edit the
// file.
type RosterEdit struct {
	// FileName is the name of the file containing the roster information.
	FileName string
	// Roster contains the state of the roster at the time of writing the
	// file. It's what we diff against when reading the file.
	Roster []xmpp.RosterEntry
	// isComplete is true if this is the result of reading an edited
	// roster, rather than a report that the file has been written.
	IsComplete bool
	// contents contains the edited roster, if isComplete is true.
	Contents []byte
}

func IsAwayStatus(status string) bool {
	switch status {
	case "xa", "away":
		return true
	}
	return false
}

var hexTable = "0123456789abcdef"

// escapeNonASCII replaces tabs and other non-printable characters with a
// "\x01" form of hex escaping. It works on a byte-by-byte basis.
func EscapeNonASCII(in string) string {
	escapes := 0
	for i := 0; i < len(in); i++ {
		if in[i] < 32 || in[i] > 126 || in[i] == '\\' {
			escapes++
		}
	}

	if escapes == 0 {
		return in
	}

	out := make([]byte, 0, len(in)+3*escapes)
	for i := 0; i < len(in); i++ {
		if in[i] < 32 || in[i] > 126 || in[i] == '\\' {
			out = append(out, '\\', 'x', hexTable[in[i]>>4], hexTable[in[i]&15])
		} else {
			out = append(out, in[i])
		}
	}

	return string(out)
}

// unescapeNonASCII undoes the transformation of escapeNonASCII.
func UnescapeNonASCII(in string) (string, error) {
	needsUnescaping := false
	for i := 0; i < len(in); i++ {
		if in[i] == '\\' {
			needsUnescaping = true
			break
		}
	}

	if !needsUnescaping {
		return in, nil
	}

	out := make([]byte, 0, len(in))
	for i := 0; i < len(in); i++ {
		if in[i] == '\\' {
			if len(in) <= i+3 {
				return "", errors.New("truncated escape sequence at end: " + in)
			}
			if in[i+1] != 'x' {
				return "", errors.New("escape sequence didn't start with \\x in: " + in)
			}
			v, err := strconv.ParseUint(in[i+2:i+4], 16, 8)
			if err != nil {
				return "", errors.New("failed to parse value in '" + in + "': " + err.Error())
			}
			out = append(out, byte(v))
			i += 3
		} else {
			out = append(out, in[i])
		}
	}

	return string(out), nil
}

type rawLogger struct {
	out    io.Writer
	prefix []byte
	lock   *sync.Mutex
	other  *rawLogger
	buf    []byte
}

func (r *rawLogger) Write(data []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.other.flush(); err != nil {
		return 0, nil
	}

	origLen := len(data)
	for len(data) > 0 {
		if newLine := bytes.IndexByte(data, '\n'); newLine >= 0 {
			r.buf = append(r.buf, data[:newLine]...)
			data = data[newLine+1:]
		} else {
			r.buf = append(r.buf, data...)
			data = nil
		}
	}

	return origLen, nil
}

var NewLine = []byte{'\n'}

func (r *rawLogger) flush() error {
	if len(r.buf) == 0 {
		return nil
	}

	if _, err := r.out.Write(r.prefix); err != nil {
		return err
	}
	if _, err := r.out.Write(r.buf); err != nil {
		return err
	}
	if _, err := r.out.Write(NewLine); err != nil {
		return err
	}
	r.buf = r.buf[:0]
	return nil
}
