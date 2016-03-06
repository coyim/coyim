package sasl

import (
	"unicode"
	"unicode/utf8"

	"github.com/twstrike/coyim/Godeps/_workspace/src/golang.org/x/text/transform"
	"github.com/twstrike/coyim/Godeps/_workspace/src/golang.org/x/text/unicode/norm"
)

type replaceTransformer func(r rune) rune

func (t replaceTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	var runeBytes [utf8.UTFMax]byte
	for r, sz := rune(0), 0; len(src) > 0; src = src[sz:] {

		if r = rune(src[0]); r < utf8.RuneSelf {
			sz = 1
		} else {
			r, sz = utf8.DecodeRune(src)

			if sz == 1 {
				// Invalid rune.
				if !atEOF && !utf8.FullRune(src) {
					err = transform.ErrShortSrc
					break
				}
			}
		}

		dsz := utf8.EncodeRune(runeBytes[:], t(r))
		if nDst+dsz > len(dst) {
			err = transform.ErrShortDst
			break
		}

		nDst += copy(dst[nDst:], runeBytes[:dsz])
		nSrc += sz
	}
	return
}

func (t replaceTransformer) Reset() {}

var nonASCIISpaceTransformer = replaceTransformer(func(r rune) rune {
	space := '\u0020'

	if unicode.Is(unicode.Zs, r) {
		return space
	}

	//TODO: replace by a unicode.RangeTable but I could not understand the Stride thing
	if r == '\u200B' {
		return space
	}

	return r
})

var mappedToNothing = transform.RemoveFunc(func(r rune) bool {
	//TODO: replace by a unicode.RangeTable
	switch r {
	case '\u00AD':
	case '\u034F':
	case '\u1806':
	case '\u180B':
	case '\u180C':
	case '\u180D':
	case '\u200C':
	case '\u200D':
	case '\u2060':
	case '\uFE00':
	case '\uFE01':
	case '\uFE02':
	case '\uFE03':
	case '\uFE04':
	case '\uFE05':
	case '\uFE06':
	case '\uFE07':
	case '\uFE08':
	case '\uFE09':
	case '\uFE0A':
	case '\uFE0B':
	case '\uFE0C':
	case '\uFE0D':
	case '\uFE0E':
	case '\uFE0F':
	case '\uFEFF':
	default:
		return false
	}

	return true
})

// Stringprep implements Stringprep Profile for User Names and Passwords (RFC 4013)
// as a transform.Transformer
var Stringprep = transform.Chain(nonASCIISpaceTransformer, mappedToNothing, norm.NFKC)
