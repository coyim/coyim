package muc

func panicOnDevError(e error) {
	if e != nil {
		panic(e)
	}
}
