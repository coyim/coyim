package pangoi

type Pango interface {
	AsFontDescription(interface{}) FontDescription
	AttrListNew() AttrList
}

func AssertPango(_ Pango) {}
