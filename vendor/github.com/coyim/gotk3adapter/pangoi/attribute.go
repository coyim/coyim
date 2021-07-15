package pangoi

type PangoAttribute interface {
	SetStartIndex(int)
	SetEndIndex(int)
}

func AssertPangoAttribute(_ PangoAttribute) {}
