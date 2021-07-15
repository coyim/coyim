package gtk_mock

import "github.com/coyim/gotk3adapter/pangoi"

type MockLabel struct {
	MockWidget
}

func (*MockLabel) GetLabel() string {
	return ""
}

func (*MockLabel) SetLabel(v1 string) {
}

func (*MockLabel) SetText(v1 string) {
}

func (*MockLabel) SetMarkup(v1 string) {
}

func (*MockLabel) SetSelectable(v1 bool) {
}

func (*MockLabel) GetMnemonicKeyval() uint {
	return 0
}

func (*MockLabel) GetPangoAttributes() pangoi.PangoAttrList {
	return nil
}

func (*MockLabel) SetPangoAttributes(pangoi.PangoAttrList) {
}
