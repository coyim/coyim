package pangoi

type Attribute interface {
	SetStartIndex(int)
	SetEndIndex(int)
}

func AssertAttribute(_ Attribute) {}
