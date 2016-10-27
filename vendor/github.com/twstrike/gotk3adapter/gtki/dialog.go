package gtki

type Dialog interface {
	Window

	Run() int
	SetDefaultResponse(ResponseType)
}

func AssertDialog(_ Dialog) {}
