package gtk_mock

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type Mock struct{}

func (*Mock) AboutDialogNew() (gtki.AboutDialog, error) {
	return nil, nil
}

func (*Mock) AccelGroupNew() (gtki.AccelGroup, error) {
	return nil, nil
}

func (*Mock) AcceleratorParse(acc string) (uint, gdki.ModifierType) {
	return 0, gdki.ModifierType(0)
}

func (*Mock) AddProviderForScreen(s gdki.Screen, provider gtki.StyleProvider, prio uint) {
}

func (*Mock) ApplicationNew(appId string, flags glibi.ApplicationFlags) (gtki.Application, error) {
	return nil, nil
}

func (*Mock) AssistantNew() (gtki.Assistant, error) {
	return nil, nil
}

func (*Mock) BuilderNew() (gtki.Builder, error) {
	return nil, nil
}

func (*Mock) BuilderNewFromResource(string) (gtki.Builder, error) {
	return nil, nil
}

func (*Mock) CellRendererTextNew() (gtki.CellRendererText, error) {
	return nil, nil
}

func (*Mock) CheckButtonNewWithMnemonic(label string) (gtki.CheckButton, error) {
	return nil, nil
}

func (*Mock) CheckMenuItemNewWithMnemonic(label string) (gtki.CheckMenuItem, error) {
	return nil, nil
}

func (*Mock) CssProviderNew() (gtki.CssProvider, error) {
	return nil, nil
}

func (*Mock) CssProviderGetDefault() (gtki.CssProvider, error) {
	return nil, nil
}

func (*Mock) CssProviderGetNamed(string, string) (gtki.CssProvider, error) {
	return nil, nil
}

func (*Mock) EntryNew() (gtki.Entry, error) {
	return nil, nil
}

func (*Mock) FileChooserDialogNewWith2Buttons(title string, parent gtki.Window, action gtki.FileChooserAction, first_button_text string, first_button_id gtki.ResponseType, second_button_text string, second_button_id gtki.ResponseType) (gtki.FileChooserDialog, error) {
	return nil, nil
}

func (*Mock) Init(args *[]string) {
}

func (*Mock) LabelNew(str string) (gtki.Label, error) {
	return nil, nil
}

func (*Mock) ListStoreNew(types ...glibi.Type) (gtki.ListStore, error) {
	return nil, nil
}

func (*Mock) MenuItemNew() (gtki.MenuItem, error) {
	return nil, nil
}

func (*Mock) MenuItemNewWithMnemonic(label string) (gtki.MenuItem, error) {
	return nil, nil
}

func (*Mock) MenuItemNewWithLabel(label string) (gtki.MenuItem, error) {
	return nil, nil
}

func (*Mock) MenuNew() (gtki.Menu, error) {
	return nil, nil
}

func (*Mock) SeparatorMenuItemNew() (gtki.SeparatorMenuItem, error) {
	return nil, nil
}

func (*Mock) TextBufferNew(table gtki.TextTagTable) (gtki.TextBuffer, error) {
	return nil, nil
}

func (*Mock) TextTagNew(name string) (gtki.TextTag, error) {
	return nil, nil
}

func (*Mock) TextTagTableNew() (gtki.TextTagTable, error) {
	return nil, nil
}

func (*Mock) TextViewNew() (gtki.TextView, error) {
	return nil, nil
}

func (*Mock) TreePathNew() gtki.TreePath {
	return nil
}

func (*Mock) WindowSetDefaultIcon(icon gdki.Pixbuf) {
}

func (*Mock) SettingsGetDefault() (gtki.Settings, error) {
	return nil, nil
}
