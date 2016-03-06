package pangoi

type Pango interface {
	AsFontDescription(interface{}) FontDescription
}

func AssertPango(_ Pango) {}
