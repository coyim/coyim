package glibi

type Value interface {
	GetString() (string, error)
	GoValue() (interface{}, error)
}

func AssertValue(_ Value) {}
