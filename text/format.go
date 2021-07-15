package text

/*
  Simple idea - taking input as text, returning formatted text.  Formatted text is first returned as a slice of
  interface Fragment. Each fragment can be just the simple form, or it can be formatted. This also allows for the
  possibility of interpolation fragments later.

  For simplicity, this package also provides a `Join()` function to create the final output that can be used to send to
  actual rendering.

  Basic syntax:
  - All text is treated normally, until a $ is encountered
  - A $ followed by a $ will generate a single $
  - A $ followed by [a-zA-Z0-9_]+ followed by { followed by any text, followed by } will generate a formatted region
    where the text between $ and { will be the name of the format. A $ followed by } will generate a literal
    } in the output, and won't close the escape sequence
  - If the parsing doesn't match in some way, the `ParseWithFormat()` function will return false, and the
    original string will be returned as a fragment.
  - Examples:
     - "hello world" will generate a simple fragment, saying "hello world"
     - "hello, the cost will be $$42" will generate "hello, the cost will be $42"
     - "hello and welcome, $nick{Luke} - it's time to start" will generate three fragments:
        A starting text, saying "hello and welcome, ", a formatted fragment containing "Luke" with the format
        "nick", and a final text fragment " - it's time to start"
     - "hello and welcome, $role{foo{$}bar$$}" will generate "hello and welcome, " and the formatted fragment "foo{}bar$"

  Possible extension - variable interpolation:
  - This is not currently supported, but can easily be added
  - The only change in the actual API would be to make `Join()` take `...string`, which would be the values to interpolate.
    (These could also be `interface{}` and then run through possible `Stringer`)
  - The syntax would be extended so that the format name could be followed by either { or [. If it's followed by {, things work
    as before, but if followed by [, the next value would be a zero-indexed value, followed only by a ]. No escapes would be necessary.
    Only zero or positive digits would be possible. An example: "hello $nick[1] - your role is $role[0] or $role{bla}"
  - With this extension, it might be necessary to also return an error or boolean from the `Join()` method.
*/

type fragmentWithFormat struct {
	format   string
	fragment Fragment
}

type textFragment struct {
	text string
}

type compositeFragment struct {
	text []string
}

func genTextFragment(txt ...string) Fragment {
	res := []string{}
	for _, val := range txt {
		if val != "" {
			res = append(res, val)
		}
	}

	if len(res) == 1 {
		return &textFragment{res[0]}
	}
	return &compositeFragment{res}
}

// Fragment is one of any type of text fragments - either formatted or unformatted
type Fragment interface{}

// FormattedText is a slice of fragments
type FormattedText []Fragment

func isFormatNameCharacter(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}

func parseFormatName(txt string) (formatName, rest string, ok bool) {
	ix := 0
	rstxt := []rune(txt)
	for isFormatNameCharacter(rstxt[ix]) {
		ix++
	}
	return string(rstxt[0:ix]), string(rstxt[ix:]), true
}

func parseFormattedText(txt string) (f Fragment, rest string, ok bool) {
	currentStart := 0
	ix := 0
	result := []string{}
	end := false
	for ix < len(txt) && !end {
		switch txt[ix] {
		case '$':
			if ix+1 < len(txt) && txt[ix+1] == '}' {
				result = append(result, txt[currentStart:ix])
				result = append(result, "}")
				currentStart = ix + 2
				ix++
			}
			ix++
		case '}':
			result = append(result, txt[currentStart:ix])
			end = true
		default:
			ix++
		}
	}
	return genTextFragment(result...), txt[ix+1:], true
}

func parseNextFormattedFragment(txt string) (f Fragment, rest string, more bool, ok bool) {
	formatName, rest2, ok2 := parseFormatName(txt)
	ok2 = ok2
	if rest2[0] == '{' {
		f2, rest3, ok3 := parseFormattedText(rest2[1:])
		ok3 = ok3
		return &fragmentWithFormat{formatName, f2}, rest3, true, true
	}
	return nil, "", false, false
}

func parseNextEscapeOrFormattedFragment(txt string) (f Fragment, rest string, more bool, ok bool) {
	if txt == "" {
		return genTextFragment("$"), "", false, false
	}

	if txt[0] == '$' {
		return genTextFragment("$"), txt[1:], len(txt) > 1, true
	}

	return parseNextFormattedFragment(txt)
}

func parseNext(txt string) (f Fragment, rest string, more bool, ok bool) {
	if txt == "" {
		return nil, "", false, true
	}

	if txt[0] == '$' {
		return parseNextEscapeOrFormattedFragment(txt[1:])
	}

	ix := 0
	for ix < len(txt) && txt[ix] != '$' {
		ix++
	}

	return genTextFragment(txt[0:ix]), txt[ix:], true, true
}

// ParseWithFormat parses the given text following the description in the package documentation
// It return false if the format is incorrect. In this case it will return a
// one-slice formatted text containing the original text.
func ParseWithFormat(txt string) (FormattedText, bool) {
	result := FormattedText{}

	f, rest, more, ok := parseNext(txt)
	if f != nil {
		result = append(result, f)
	}

	for more && ok {
		f, rest, more, ok = parseNext(rest)
		if f != nil {
			result = append(result, f)
		}
	}

	if !ok {
		return FormattedText{genTextFragment(txt)}, false
	}

	return result, true
}

// Join will generate a final text, making the text into one segment, and calculating all
// indices and lengths for format strings
func (ft FormattedText) Join() (text string, starts []int, lengths []int, formats []string) {
	return "", nil, nil, nil
}
