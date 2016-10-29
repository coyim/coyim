package gdki

type PixbufFormat interface {
	GetName() (string, error)
	GetDescription() (string, error)
	GetLicense() (string, error)
}
