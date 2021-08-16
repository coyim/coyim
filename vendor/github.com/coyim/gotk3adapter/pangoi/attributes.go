package pangoi

type AttrList interface {
	Insert(Attribute)
	GetAttributes() []Attribute
}

func AssertAttrList(_ AttrList) {}
