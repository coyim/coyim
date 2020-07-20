package gui

func panicOnDevError(e error) {
	if e != nil {
		panic(e)
	}
}
