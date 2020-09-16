package gtka

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type RealGtk struct{}

var Real = &RealGtk{}

func (*RealGtk) AboutDialogNew() (gtki.AboutDialog, error) {
	return WrapAboutDialog(gtk.AboutDialogNew())
}

func (*RealGtk) AccelGroupNew() (gtki.AccelGroup, error) {
	return WrapAccelGroup(gtk.AccelGroupNew())
}

func (*RealGtk) AcceleratorParse(acc string) (uint, gdki.ModifierType) {
	res, res2 := gtk.AcceleratorParse(acc)
	return res, gdki.ModifierType(res2)
}

func (*RealGtk) AddProviderForScreen(s gdki.Screen, provider gtki.StyleProvider, prio uint) {
	gtk.AddProviderForScreen(gdka.UnwrapScreen(s), UnwrapStyleProvider(provider), prio)
}

func (*RealGtk) ApplicationNew(appId string, flags glibi.ApplicationFlags) (gtki.Application, error) {
	return WrapApplication(gtk.ApplicationNew(appId, glib.ApplicationFlags(flags)))
}

func (*RealGtk) AssistantNew() (gtki.Assistant, error) {
	return WrapAssistant(gtk.AssistantNew())
}

func (*RealGtk) BuilderNew() (gtki.Builder, error) {
	return WrapBuilder(gtk.BuilderNew())
}

func (*RealGtk) BuilderNewFromResource(s string) (gtki.Builder, error) {
	return WrapBuilder(gtk.BuilderNewFromResource(s))
}

func (*RealGtk) CellRendererTextNew() (gtki.CellRendererText, error) {
	return WrapCellRendererText(gtk.CellRendererTextNew())
}

func (*RealGtk) CheckButtonNew() (gtki.CheckButton, error) {
	return WrapCheckButton(gtk.CheckButtonNew())
}

func (*RealGtk) CheckButtonNewWithMnemonic(label string) (gtki.CheckButton, error) {
	return WrapCheckButton(gtk.CheckButtonNewWithMnemonic(label))
}

func (*RealGtk) CheckMenuItemNewWithMnemonic(label string) (gtki.CheckMenuItem, error) {
	return WrapCheckMenuItem(gtk.CheckMenuItemNewWithMnemonic(label))
}

func (*RealGtk) CheckVersion(major, minor, micro uint) error {
	return gtk.CheckVersion(major, minor, micro)
}

func (*RealGtk) ComboBoxNew() (gtki.ComboBox, error) {
	return WrapComboBox(gtk.ComboBoxNew())
}

func (*RealGtk) ComboBoxTextNew() (gtki.ComboBoxText, error) {
	return WrapComboBoxText(gtk.ComboBoxTextNew())
}

func (*RealGtk) CssProviderNew() (gtki.CssProvider, error) {
	return WrapCssProvider(gtk.CssProviderNew())
}

func (*RealGtk) CssProviderGetNamed(name, variant string) (gtki.CssProvider, error) {
	return WrapCssProvider(gtk.CssProviderGetNamed(name, variant))
}

func (*RealGtk) EntryNew() (gtki.Entry, error) {
	return WrapEntry(gtk.EntryNew())
}

func (*RealGtk) EventBoxNew() (gtki.EventBox, error) {
	return WrapEventBox(gtk.EventBoxNew())
}

func (*RealGtk) FileChooserDialogNewWith2Buttons(title string, parent gtki.Window, action gtki.FileChooserAction, first_button_text string, first_button_id gtki.ResponseType, second_button_text string, second_button_id gtki.ResponseType) (gtki.FileChooserDialog, error) {
	return WrapFileChooserDialog(gtk.FileChooserDialogNewWith2Buttons(title, UnwrapWindow(parent), gtk.FileChooserAction(action), first_button_text, gtk.ResponseType(first_button_id), second_button_text, gtk.ResponseType(second_button_id)))
}

func (*RealGtk) GetMajorVersion() uint {
	return gtk.GetMajorVersion()
}

func (*RealGtk) GetMinorVersion() uint {
	return gtk.GetMinorVersion()
}

func (*RealGtk) GetMicroVersion() uint {
	return gtk.GetMicroVersion()
}

func (*RealGtk) ImageNewFromFile(filename string) (gtki.Image, error) {
	return WrapImage(gtk.ImageNewFromFile(filename))
}

func (*RealGtk) ImageNewFromResource(path string) (gtki.Image, error) {
	return WrapImage(gtk.ImageNewFromResource(path))
}

func (*RealGtk) ImageNewFromPixbuf(v1 gdki.Pixbuf) (gtki.Image, error) {
	return WrapImage(gtk.ImageNewFromPixbuf(gdka.UnwrapPixbuf(v1)))
}

func (*RealGtk) ImageNewFromIconName(name string, v2 gtki.IconSize) (gtki.Image, error) {
	return WrapImage(gtk.ImageNewFromIconName(name, gtk.IconSize(v2)))
}

func (*RealGtk) InfoBarNew() (gtki.InfoBar, error) {
	return WrapInfoBar(gtk.InfoBarNew())
}

func (*RealGtk) Init(args *[]string) {
	gtk.Init(args)
}

func (*RealGtk) LabelNew(str string) (gtki.Label, error) {
	return WrapLabel(gtk.LabelNew(str))
}

func unwrapTypes(ts []glibi.Type) []glib.Type {
	result := make([]glib.Type, len(ts))
	for ix, rr := range ts {
		result[ix] = glib.Type(rr)
	}
	return result
}

func (*RealGtk) ListStoreNew(types ...glibi.Type) (gtki.ListStore, error) {
	return WrapListStore(gtk.ListStoreNew(unwrapTypes(types)...))
}

func (*RealGtk) TreeStoreNew(types ...glibi.Type) (gtki.TreeStore, error) {
	return WrapTreeStore(gtk.TreeStoreNew(unwrapTypes(types)...))
}

func (*RealGtk) MenuBarNew() (gtki.MenuBar, error) {
	return WrapMenuBar(gtk.MenuBarNew())
}

func (*RealGtk) MenuItemNew() (gtki.MenuItem, error) {
	return WrapMenuItem(gtk.MenuItemNew())
}

func (*RealGtk) MenuItemNewWithMnemonic(label string) (gtki.MenuItem, error) {
	return WrapMenuItem(gtk.MenuItemNewWithMnemonic(label))
}

func (*RealGtk) MenuItemNewWithLabel(label string) (gtki.MenuItem, error) {
	return WrapMenuItem(gtk.MenuItemNewWithLabel(label))
}

func (*RealGtk) MenuNew() (gtki.Menu, error) {
	return WrapMenu(gtk.MenuNew())
}

func (*RealGtk) SeparatorMenuItemNew() (gtki.SeparatorMenuItem, error) {
	return WrapSeparatorMenuItem(gtk.SeparatorMenuItemNew())
}

func (*RealGtk) SearchBarNew() (gtki.SearchBar, error) {
	return WrapSearchBar(gtk.SearchBarNew())
}

func (*RealGtk) SearchEntryNew() (gtki.SearchEntry, error) {
	return WrapSearchEntry(gtk.SearchEntryNew())
}

func (*RealGtk) TextBufferNew(table gtki.TextTagTable) (gtki.TextBuffer, error) {
	return WrapTextBuffer(gtk.TextBufferNew(UnwrapTextTagTable(table)))
}

func (*RealGtk) TextTagNew(name string) (gtki.TextTag, error) {
	return WrapTextTag(gtk.TextTagNew(name))
}

func (*RealGtk) TextTagTableNew() (gtki.TextTagTable, error) {
	return WrapTextTagTable(gtk.TextTagTableNew())
}

func (*RealGtk) TextViewNew() (gtki.TextView, error) {
	return WrapTextView(gtk.TextViewNew())
}

func (*RealGtk) TreePathNew() gtki.TreePath {
	var tp gtk.TreePath
	return WrapTreePathSimple(&tp)
}

func (*RealGtk) WindowSetDefaultIcon(icon gdki.Pixbuf) {
	gtk.WindowSetDefaultIcon(gdka.UnwrapPixbuf(icon))
}

func (*RealGtk) SettingsGetDefault() (gtki.Settings, error) {
	return WrapSettings(gtk.SettingsGetDefault())
}

func (*RealGtk) StatusIconNew() (gtki.StatusIcon, error) {
	return WrapStatusIcon(gotk3extra.StatusIconNew())
}

func (*RealGtk) StatusIconNewFromFile(filename string) (gtki.StatusIcon, error) {
	return WrapStatusIcon(gotk3extra.StatusIconNewFromFile(filename))
}

func (*RealGtk) StatusIconNewFromIconName(iconName string) (gtki.StatusIcon, error) {
	return WrapStatusIcon(gotk3extra.StatusIconNewFromIconName(iconName))
}

func (*RealGtk) StatusIconNewFromPixbuf(pixbuf gdki.Pixbuf) (gtki.StatusIcon, error) {
	return WrapStatusIcon(gotk3extra.StatusIconNewFromPixbuf(gdka.UnwrapPixbuf(pixbuf)))
}
