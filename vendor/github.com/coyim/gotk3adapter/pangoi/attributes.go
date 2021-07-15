package pangoi

type PangoAttrList interface {
	GetAttributes() []PangoAttribute
	InsertPangoAttribute(PangoAttribute)
	Insert(Attribute)
}

func AssertPangoAttrList(_ PangoAttrList) {}
