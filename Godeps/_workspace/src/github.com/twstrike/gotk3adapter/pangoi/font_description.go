package pangoi

type FontDescription interface {
	GetSize() int
}

func AssertFontDescription(_ FontDescription) {}
