package gdki

type EventKey interface {
	Event

	KeyVal() uint
	State() uint
}

func AssertEventKey(_ EventKey) {}
