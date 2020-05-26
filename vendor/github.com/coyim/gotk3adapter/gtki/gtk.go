package gtki

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
)

type Gtk interface {
	AboutDialogNew() (AboutDialog, error)
	AccelGroupNew() (AccelGroup, error)
	AcceleratorParse(string) (uint, gdki.ModifierType)
	AddProviderForScreen(gdki.Screen, StyleProvider, uint)
	ApplicationNew(string, glibi.ApplicationFlags) (Application, error)
	AssistantNew() (Assistant, error)
	BuilderNew() (Builder, error)
	BuilderNewFromResource(string) (Builder, error)
	CellRendererTextNew() (CellRendererText, error)
	CheckButtonNew() (CheckButton, error)
	CheckButtonNewWithMnemonic(string) (CheckButton, error)
	CheckMenuItemNewWithMnemonic(string) (CheckMenuItem, error)
	CheckVersion(major, minor, micro uint) error
	ComboBoxNew() (ComboBox, error)
	ComboBoxTextNew() (ComboBoxText, error)
	CssProviderNew() (CssProvider, error)
	CssProviderGetDefault() (CssProvider, error)
	CssProviderGetNamed(string, string) (CssProvider, error)
	EntryNew() (Entry, error)
	EventBoxNew() (EventBox, error)
	FileChooserDialogNewWith2Buttons(string, Window, FileChooserAction, string, ResponseType, string, ResponseType) (FileChooserDialog, error)
	GetMajorVersion() uint
	GetMinorVersion() uint
	GetMicroVersion() uint
	ImageNewFromFile(string) (Image, error)
	ImageNewFromResource(string) (Image, error)
	ImageNewFromPixbuf(gdki.Pixbuf) (Image, error)
	ImageNewFromIconName(string, IconSize) (Image, error)
	Init(*[]string)
	InfoBarNew() (InfoBar, error)
	LabelNew(string) (Label, error)
	ListStoreNew(...glibi.Type) (ListStore, error)
	MenuItemNew() (MenuItem, error)
	MenuItemNewWithLabel(string) (MenuItem, error)
	MenuItemNewWithMnemonic(string) (MenuItem, error)
	MenuNew() (Menu, error)
	SearchBarNew() (SearchBar, error)
	SearchEntryNew() (SearchEntry, error)
	SeparatorMenuItemNew() (SeparatorMenuItem, error)
	TextBufferNew(TextTagTable) (TextBuffer, error)
	TextTagNew(string) (TextTag, error)
	TextTagTableNew() (TextTagTable, error)
	TextViewNew() (TextView, error)
	TreePathNew() TreePath
	WindowSetDefaultIcon(gdki.Pixbuf)
	SettingsGetDefault() (Settings, error)
}

func AssertGtk(_ Gtk) {}
