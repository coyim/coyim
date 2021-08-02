package gtki

type SizeGroup interface {
	SetMode(SizeGroupMode)
}

func AssertSizeGroup(_ SizeGroup) {}
