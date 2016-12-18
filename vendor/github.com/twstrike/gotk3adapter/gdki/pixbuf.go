package gdki

type Pixbuf interface {
	SavePNG(string, int) error
}

func AssertPixbuf(_ Pixbuf) {}
