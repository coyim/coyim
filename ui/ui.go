package ui

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var tagsToAvoid = make(map[string]bool)

func init() {
	tagsToAvoid["blockquote"] = true
	tagsToAvoid["br"] = true
	tagsToAvoid["cite"] = true
	tagsToAvoid["em"] = true
	tagsToAvoid["font"] = true
	tagsToAvoid["p"] = true
	tagsToAvoid["span"] = true
	tagsToAvoid["strong"] = true
	tagsToAvoid["a"] = true
	tagsToAvoid["i"] = true
	tagsToAvoid["b"] = true
	tagsToAvoid["u"] = true
	tagsToAvoid["img"] = true
}

// UnescapeNewlineTags will remove all newline tags and replace them with actual newlines
func UnescapeNewlineTags(msg []byte) (out []byte) {
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
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			raw := z.Raw()
			name, _ := z.TagName()

			if string(name) == "br" {
				out = append(out, byte('\n'))
			} else {
				out = append(out, raw...)
			}
		case html.CommentToken:
			out = append(out, z.Raw()...)
		case html.DoctypeToken:
			out = append(out, z.Raw()...)
		}
	}

	return
}

// StripSomeHTML removes the most common html presentation tags from the text
func StripSomeHTML(msg []byte) (out []byte) {
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
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			raw := z.Raw()
			name, _ := z.TagName()

			if !tagsToAvoid[string(name)] {
				out = append(out, raw...)
			}
		case html.CommentToken:
			out = append(out, z.Raw()...)
		case html.DoctypeToken:
			out = append(out, z.Raw()...)
		}
	}

	return
}

// StripHTML removes all html in the text
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

// EscapeAllHTMLTags will escape all html tags in the text
func EscapeAllHTMLTags(in string) string {
	in = strings.Replace(in, "<", "&lt;", -1)
	in = strings.Replace(in, ">", "&gt;", -1)
	return in
}

var (
	hexTable = "0123456789abcdef"
	// NewLine contains a new line
	NewLine = []byte{'\n'}
)

// EscapeNonASCII replaces tabs and other non-printable characters with a
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

// UnescapeNonASCII undoes the transformation of escapeNonASCII.
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
