package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockSpinner struct {
	MockWidget
}

func (*MockSpinner) Start() {
}

func (*MockSpinner) Stop() {
}

func (*Mock) SpinnerNew() (gtki.Spinner, error) {
	return nil, nil
}
