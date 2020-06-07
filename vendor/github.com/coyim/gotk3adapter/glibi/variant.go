package glibi

type Variant interface {
	TypeString() string
	IsContainer() bool
	GetBoolean() bool
	GetString() string
	GetStrv() []string
	GetInt() (int64, error)
	Type() VariantType
	IsType(t VariantType) bool
	String() string
	AnnotatedString() string
}

func AssertVariant(_ Variant) {}
