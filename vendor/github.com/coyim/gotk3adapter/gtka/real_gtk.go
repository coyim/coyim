package gtka

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type RealGtk struct{}

var Real = &RealGtk{}

func (*RealGtk) AboutDialogNew() (gtki.AboutDialog, error) {
	return wrapAboutDialog(gtk.AboutDialogNew())
}

func (*RealGtk) AccelGroupNew() (gtki.AccelGroup, error) {
	return wrapAccelGroup(gtk.AccelGroupNew())
}

func (*RealGtk) AcceleratorParse(acc string) (uint, gdki.ModifierType) {
	res, res2 := gtk.AcceleratorParse(acc)
	return res, gdki.ModifierType(res2)
}

func (*RealGtk) AddProviderForScreen(s gdki.Screen, provider gtki.StyleProvider, prio uint) {
	gtk.AddProviderForScreen(gdka.UnwrapScreen(s), unwrapStyleProvider(provider), prio)
}

func (*RealGtk) ApplicationNew(appId string, flags glibi.ApplicationFlags) (gtki.Application, error) {
	return wrapApplication(gtk.ApplicationNew(appId, glib.ApplicationFlags(flags)))
}

func (*RealGtk) AssistantNew() (gtki.Assistant, error) {
	return wrapAssistant(gtk.AssistantNew())
}

func (*RealGtk) BuilderNew() (gtki.Builder, error) {
	return wrapBuilder(gtk.BuilderNew())
}

func (*RealGtk) BuilderNewFromResource(s string) (gtki.Builder, error) {
	return wrapBuilder(gtk.BuilderNewFromResource(s))
}

func (*RealGtk) CellRendererTextNew() (gtki.CellRendererText, error) {
	return wrapCellRendererText(gtk.CellRendererTextNew())
}

func (*RealGtk) CheckButtonNew() (gtki.CheckButton, error) {
	return wrapCheckButton(gtk.CheckButtonNew())
}

func (*RealGtk) CheckButtonNewWithMnemonic(label string) (gtki.CheckButton, error) {
	return wrapCheckButton(gtk.CheckButtonNewWithMnemonic(label))
}

func (*RealGtk) CheckMenuItemNewWithMnemonic(label string) (gtki.CheckMenuItem, error) {
	return wrapCheckMenuItem(gtk.CheckMenuItemNewWithMnemonic(label))
}

func (*RealGtk) ComboBoxNew() (gtki.ComboBox, error) {
	return wrapComboBox(gtk.ComboBoxNew())
}

func (*RealGtk) ComboBoxTextNew() (gtki.ComboBoxText, error) {
	return wrapComboBoxText(gtk.ComboBoxTextNew())
}

func (*RealGtk) CssProviderNew() (gtki.CssProvider, error) {
	return wrapCssProvider(gtk.CssProviderNew())
}

func (*RealGtk) CssProviderGetNamed(name, variant string) (gtki.CssProvider, error) {
	return wrapCssProvider(gtk.CssProviderGetNamed(name, variant))
}

func (*RealGtk) EntryNew() (gtki.Entry, error) {
	return wrapEntry(gtk.EntryNew())
}

func (*RealGtk) EventBoxNew() (gtki.EventBox, error) {
	return wrapEventBox(gtk.EventBoxNew())
}

func (*RealGtk) FileChooserDialogNewWith2Buttons(title string, parent gtki.Window, action gtki.FileChooserAction, first_button_text string, first_button_id gtki.ResponseType, second_button_text string, second_button_id gtki.ResponseType) (gtki.FileChooserDialog, error) {
	return wrapFileChooserDialog(gtk.FileChooserDialogNewWith2Buttons(title, unwrapWindow(parent), gtk.FileChooserAction(action), first_button_text, gtk.ResponseType(first_button_id), second_button_text, gtk.ResponseType(second_button_id)))
}

func (*RealGtk) ImageNewFromFile(filename string) (gtki.Image, error) {
	return wrapImage(gtk.ImageNewFromFile(filename))
}

func (*RealGtk) ImageNewFromResource(path string) (gtki.Image, error) {
	return wrapImage(gtk.ImageNewFromResource(path))
}

func (*RealGtk) ImageNewFromPixbuf(v1 gdki.Pixbuf) (gtki.Image, error) {
	return wrapImage(gtk.ImageNewFromPixbuf(gdka.UnwrapPixbuf(v1)))
}

func (*RealGtk) ImageNewFromIconName(name string, v2 gtki.IconSize) (gtki.Image, error) {
	return wrapImage(gtk.ImageNewFromIconName(name, gtk.IconSize(v2)))
}

func (*RealGtk) InfoBarNew() (gtki.InfoBar, error) {
	return wrapInfoBar(gtk.InfoBarNew())
}

func (*RealGtk) Init(args *[]string) {
	gtk.Init(args)
}

func (*RealGtk) LabelNew(str string) (gtki.Label, error) {
	return wrapLabel(gtk.LabelNew(str))
}

func unwrapTypes(ts []glibi.Type) []glib.Type {
	result := make([]glib.Type, len(ts))
	for ix, rr := range ts {
		result[ix] = glib.Type(rr)
	}
	return result
}

func (*RealGtk) ListStoreNew(types ...glibi.Type) (gtki.ListStore, error) {
	return wrapListStore(gtk.ListStoreNew(unwrapTypes(types)...))
}

func (*RealGtk) MenuItemNew() (gtki.MenuItem, error) {
	return wrapMenuItem(gtk.MenuItemNew())
}

func (*RealGtk) MenuItemNewWithMnemonic(label string) (gtki.MenuItem, error) {
	return wrapMenuItem(gtk.MenuItemNewWithMnemonic(label))
}

func (*RealGtk) MenuItemNewWithLabel(label string) (gtki.MenuItem, error) {
	return wrapMenuItem(gtk.MenuItemNewWithLabel(label))
}

func (*RealGtk) MenuNew() (gtki.Menu, error) {
	return wrapMenu(gtk.MenuNew())
}

func (*RealGtk) SeparatorMenuItemNew() (gtki.SeparatorMenuItem, error) {
	return wrapSeparatorMenuItem(gtk.SeparatorMenuItemNew())
}

func (*RealGtk) SearchBarNew() (gtki.SearchBar, error) {
	return wrapSearchBar(gtk.SearchBarNew())
}

func (*RealGtk) SearchEntryNew() (gtki.SearchEntry, error) {
	return wrapSearchEntry(gtk.SearchEntryNew())
}

func (*RealGtk) TextBufferNew(table gtki.TextTagTable) (gtki.TextBuffer, error) {
	return wrapTextBuffer(gtk.TextBufferNew(unwrapTextTagTable(table)))
}

func (*RealGtk) TextTagNew(name string) (gtki.TextTag, error) {
	return wrapTextTag(gtk.TextTagNew(name))
}

func (*RealGtk) TextTagTableNew() (gtki.TextTagTable, error) {
	return wrapTextTagTable(gtk.TextTagTableNew())
}

func (*RealGtk) TextViewNew() (gtki.TextView, error) {
	return wrapTextView(gtk.TextViewNew())
}

func (*RealGtk) TreePathNew() gtki.TreePath {
	var tp gtk.TreePath
	return wrapTreePathSimple(&tp)
}

func (*RealGtk) WindowSetDefaultIcon(icon gdki.Pixbuf) {
	gtk.WindowSetDefaultIcon(gdka.UnwrapPixbuf(icon))
}

func (*RealGtk) SettingsGetDefault() (gtki.Settings, error) {
	return wrapSettings(gtk.SettingsGetDefault())
}
