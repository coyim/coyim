package pangoi

type Pango interface {
	AsFontDescription(interface{}) FontDescription
	PangoAttrListNew() PangoAttrList
}

func AssertPango(_ Pango) {}
