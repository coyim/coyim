package ui

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"golang.org/x/net/html"
)

var tagsToAvoid = make(map[string]bool)
var tagsHTMLToEscape = make(map[string]bool)

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

	tagsHTMLToEscape["a"] = true
	tagsHTMLToEscape["abbr"] = true
	tagsHTMLToEscape["address"] = true
	tagsHTMLToEscape["area"] = true
	tagsHTMLToEscape["article"] = true
	tagsHTMLToEscape["aside"] = true
	tagsHTMLToEscape["audio"] = true
	tagsHTMLToEscape["b"] = true
	tagsHTMLToEscape["base"] = true
	tagsHTMLToEscape["bdi"] = true
	tagsHTMLToEscape["bdo"] = true
	tagsHTMLToEscape["blockquote"] = true
	tagsHTMLToEscape["body"] = true
	tagsHTMLToEscape["br"] = true
	tagsHTMLToEscape["button"] = true
	tagsHTMLToEscape["canvas"] = true
	tagsHTMLToEscape["caption"] = true
	tagsHTMLToEscape["cite"] = true
	tagsHTMLToEscape["code"] = true
	tagsHTMLToEscape["col"] = true
	tagsHTMLToEscape["colgroup"] = true
	tagsHTMLToEscape["data"] = true
	tagsHTMLToEscape["datalist"] = true
	tagsHTMLToEscape["dd"] = true
	tagsHTMLToEscape["del"] = true
	tagsHTMLToEscape["details"] = true
	tagsHTMLToEscape["dfn"] = true
	tagsHTMLToEscape["dialog"] = true
	tagsHTMLToEscape["div"] = true
	tagsHTMLToEscape["dl"] = true
	tagsHTMLToEscape["dt"] = true
	tagsHTMLToEscape["em"] = true
	tagsHTMLToEscape["embed"] = true
	tagsHTMLToEscape["fieldset"] = true
	tagsHTMLToEscape["figcaption"] = true
	tagsHTMLToEscape["figure"] = true
	tagsHTMLToEscape["font"] = true
	tagsHTMLToEscape["footer"] = true
	tagsHTMLToEscape["form"] = true
	tagsHTMLToEscape["h1"] = true
	tagsHTMLToEscape["h2"] = true
	tagsHTMLToEscape["h3"] = true
	tagsHTMLToEscape["h4"] = true
	tagsHTMLToEscape["h5"] = true
	tagsHTMLToEscape["h6"] = true
	tagsHTMLToEscape["head"] = true
	tagsHTMLToEscape["header"] = true
	tagsHTMLToEscape["hgroup"] = true
	tagsHTMLToEscape["hr"] = true
	tagsHTMLToEscape["html"] = true
	tagsHTMLToEscape["i"] = true
	tagsHTMLToEscape["iframe"] = true
	tagsHTMLToEscape["img"] = true
	tagsHTMLToEscape["input"] = true
	tagsHTMLToEscape["ins"] = true
	tagsHTMLToEscape["kbd"] = true
	tagsHTMLToEscape["keygen"] = true
	tagsHTMLToEscape["label"] = true
	tagsHTMLToEscape["legend"] = true
	tagsHTMLToEscape["li"] = true
	tagsHTMLToEscape["link"] = true
	tagsHTMLToEscape["main"] = true
	tagsHTMLToEscape["map"] = true
	tagsHTMLToEscape["mark"] = true
	tagsHTMLToEscape["math"] = true
	tagsHTMLToEscape["menu"] = true
	tagsHTMLToEscape["menuitem"] = true
	tagsHTMLToEscape["meta"] = true
	tagsHTMLToEscape["meter"] = true
	tagsHTMLToEscape["nav"] = true
	tagsHTMLToEscape["noscript"] = true
	tagsHTMLToEscape["object"] = true
	tagsHTMLToEscape["ol"] = true
	tagsHTMLToEscape["optgroup"] = true
	tagsHTMLToEscape["option"] = true
	tagsHTMLToEscape["output"] = true
	tagsHTMLToEscape["p"] = true
	tagsHTMLToEscape["param"] = true
	tagsHTMLToEscape["picture"] = true
	tagsHTMLToEscape["pre"] = true
	tagsHTMLToEscape["progress"] = true
	tagsHTMLToEscape["q"] = true
	tagsHTMLToEscape["rb"] = true
	tagsHTMLToEscape["rp"] = true
	tagsHTMLToEscape["rt"] = true
	tagsHTMLToEscape["rtc"] = true
	tagsHTMLToEscape["ruby"] = true
	tagsHTMLToEscape["s"] = true
	tagsHTMLToEscape["samp"] = true
	tagsHTMLToEscape["script"] = true
	tagsHTMLToEscape["section"] = true
	tagsHTMLToEscape["select"] = true
	tagsHTMLToEscape["slot"] = true
	tagsHTMLToEscape["small"] = true
	tagsHTMLToEscape["source"] = true
	tagsHTMLToEscape["span"] = true
	tagsHTMLToEscape["strong"] = true
	tagsHTMLToEscape["style"] = true
	tagsHTMLToEscape["sub"] = true
	tagsHTMLToEscape["summary"] = true
	tagsHTMLToEscape["sup"] = true
	tagsHTMLToEscape["svg"] = true
	tagsHTMLToEscape["table"] = true
	tagsHTMLToEscape["tbody"] = true
	tagsHTMLToEscape["td"] = true
	tagsHTMLToEscape["template"] = true
	tagsHTMLToEscape["textarea"] = true
	tagsHTMLToEscape["tfoot"] = true
	tagsHTMLToEscape["th"] = true
	tagsHTMLToEscape["thead"] = true
	tagsHTMLToEscape["time"] = true
	tagsHTMLToEscape["title"] = true
	tagsHTMLToEscape["tr"] = true
	tagsHTMLToEscape["track"] = true
	tagsHTMLToEscape["u"] = true
	tagsHTMLToEscape["ul"] = true
	tagsHTMLToEscape["var"] = true
	tagsHTMLToEscape["video"] = true
	tagsHTMLToEscape["wbr"] = true
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
	var out []byte
	msg := []byte(in)
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
				return string(out)
			}
			break loop
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			raw := z.Raw()
			name, _ := z.TagName()
			if tagsHTMLToEscape[string(name)] {
				out = append(out, []byte((html.EscapeString((string(raw)))))...)
			} else {
				out = append(out, raw...)
			}
		case html.CommentToken, html.DoctypeToken:
			raw := z.Raw()
			out = append(out, []byte((html.EscapeString((string(raw)))))...)
		}
	}

	return string(out)
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
