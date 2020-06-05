package sasl

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

func addRuneToRangeTable(rt *unicode.RangeTable, r rune) {
	rt.R16 = append(rt.R16, unicode.Range16{Lo: uint16(r), Hi: uint16(r), Stride: 1})
}

func addRuneRangeToRangeTable(rt *unicode.RangeTable, rl, rh rune) {
	rt.R16 = append(rt.R16, unicode.Range16{Lo: uint16(rl), Hi: uint16(rh), Stride: 1})
}

func createRangeTableToUnmap() *unicode.RangeTable {
	rt := &unicode.RangeTable{}

	addRuneToRangeTable(rt, '\u00AD')
	addRuneToRangeTable(rt, '\u034F')
	addRuneToRangeTable(rt, '\u1806')
	addRuneRangeToRangeTable(rt, '\u180B', '\u180D')
	addRuneToRangeTable(rt, '\u200C')
	addRuneToRangeTable(rt, '\u200D')
	addRuneToRangeTable(rt, '\u2060')
	addRuneRangeToRangeTable(rt, '\uFE00', '\uFE0F')
	addRuneToRangeTable(rt, '\uFEFF')

	return rt
}

var unmappedRangeTable = createRangeTableToUnmap()
var mappedToNothing = runes.Remove(runes.In(unmappedRangeTable))

// Stringprep implements Stringprep Profile for User Names and Passwords (RFC 4013)
// as a transform.Transformer
var Stringprep = transform.Chain(nonASCIISpaceTransformer, mappedToNothing, norm.NFKC)
