package gtki

type SearchEntry interface {
	Entry
}

func AssertSearchEntry(_ SearchEntry) {}
