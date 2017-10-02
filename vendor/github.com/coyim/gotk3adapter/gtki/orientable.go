package gtki

// Orientable is an interface for objects that can change their orientation
type Orientable interface {
	GetOrientation() Orientation
	SetOrientation(Orientation)
}

// Orientation is a layout of an Orientable
type Orientation int

const (
	HorizontalOrientation Orientation = iota
	VerticalOrientation
)
