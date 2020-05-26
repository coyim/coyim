package gtki

type LinkButton interface {
	Button

	GetUri() string
	SetUri(string)
}

func AssertLinkButton(_ LinkButton) {}
