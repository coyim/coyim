package gliba

import "github.com/gotk3/gotk3/glib"
import "github.com/coyim/gotk3adapter/glibi"

type menuModel struct {
	*Object
	*glib.MenuModel
}

func WrapMenuModelSimple(v *glib.MenuModel) glibi.MenuModel {
	if v == nil {
		return nil
	}
	return &menuModel{WrapObjectSimple(v.Object), v}
}

func WrapMenuModel(v *glib.MenuModel, e error) (glibi.MenuModel, error) {
	return WrapMenuModelSimple(v), e
}

func UnwrapMenuModel(v glibi.MenuModel) *glib.MenuModel {
	if v == nil {
		return nil
	}
	return v.(*menuModel).MenuModel
}

func (m *menuModel) IsMutable() bool {
	return m.MenuModel.IsMutable()
}

func (m *menuModel) GetNItems() int {
	return m.MenuModel.GetNItems()
}

func (m *menuModel) GetItemLink(index int, link string) glibi.MenuModel {
	return WrapMenuModelSimple(m.MenuModel.GetItemLink(index, link))
}

func (m *menuModel) ItemsChanged(position, removed, added int) {
	m.MenuModel.ItemsChanged(position, removed, added)
}

type menu struct {
	*menuModel
	*glib.Menu
}

func WrapMenuSimple(v *glib.Menu) glibi.Menu {
	if v == nil {
		return nil
	}
	return &menu{WrapMenuModelSimple(&v.MenuModel).(*menuModel), v}
}

func WrapMenu(v *glib.Menu, e error) (glibi.Menu, error) {
	return WrapMenuSimple(v), e
}

func UnwrapMenu(v glibi.Menu) *glib.Menu {
	if v == nil {
		return nil
	}
	return v.(*menu).Menu
}

func (m *menu) Freeze() {
	m.Menu.Freeze()
}

func (m *menu) Insert(position int, label, detailed_action string) {
	m.Menu.Insert(position, label, detailed_action)
}

func (m *menu) Prepend(label, detailed_action string) {
	m.Menu.Prepend(label, detailed_action)
}

func (m *menu) Append(label, detailed_action string) {
	m.Menu.Append(label, detailed_action)
}

func (m *menu) InsertItem(position int, item glibi.MenuItem) {
	m.Menu.InsertItem(position, UnwrapMenuItem(item))
}

func (m *menu) AppendItem(item glibi.MenuItem) {
	m.Menu.AppendItem(UnwrapMenuItem(item))
}

func (m *menu) PrependItem(item glibi.MenuItem) {
	m.Menu.PrependItem(UnwrapMenuItem(item))
}

func (m *menu) InsertSection(position int, label string, section glibi.MenuModel) {
	m.Menu.InsertSection(position, label, UnwrapMenuModel(section))
}

func (m *menu) PrependSection(label string, section glibi.MenuModel) {
	m.Menu.PrependSection(label, UnwrapMenuModel(section))
}

func (m *menu) AppendSection(label string, section glibi.MenuModel) {
	m.Menu.AppendSection(label, UnwrapMenuModel(section))
}

func (m *menu) InsertSectionWithoutLabel(position int, section glibi.MenuModel) {
	m.Menu.InsertSectionWithoutLabel(position, UnwrapMenuModel(section))
}

func (m *menu) PrependSectionWithoutLabel(section glibi.MenuModel) {
	m.Menu.PrependSectionWithoutLabel(UnwrapMenuModel(section))
}

func (m *menu) AppendSectionWithoutLabel(section glibi.MenuModel) {
	m.Menu.AppendSectionWithoutLabel(UnwrapMenuModel(section))
}

func (m *menu) InsertSubmenu(position int, label string, submenu glibi.MenuModel) {
	m.Menu.InsertSubmenu(position, label, UnwrapMenuModel(submenu))
}

func (m *menu) PrependSubmenu(label string, submenu glibi.MenuModel) {
	m.Menu.PrependSubmenu(label, UnwrapMenuModel(submenu))
}

func (m *menu) AppendSubmenu(label string, submenu glibi.MenuModel) {
	m.Menu.AppendSubmenu(label, UnwrapMenuModel(submenu))
}

func (m *menu) Remove(position int) {
	m.Menu.Remove(position)
}

func (m *menu) RemoveAll() {
	m.Menu.RemoveAll()
}

type menuItem struct {
	*Object
	*glib.MenuItem
}

func WrapMenuItemSimple(v *glib.MenuItem) glibi.MenuItem {
	if v == nil {
		return nil
	}
	return &menuItem{WrapObjectSimple(v.Object), v}
}

func WrapMenuItem(v *glib.MenuItem, e error) (glibi.MenuItem, error) {
	return WrapMenuItemSimple(v), e
}

func UnwrapMenuItem(v glibi.MenuItem) *glib.MenuItem {
	if v == nil {
		return nil
	}
	return v.(*menuItem).MenuItem
}

func (m *menuItem) SetLabel(label string) {
	m.MenuItem.SetLabel(label)
}

func (m *menuItem) SetDetailedAction(act string) {
	m.MenuItem.SetDetailedAction(act)
}

func (m *menuItem) SetSection(section glibi.MenuModel) {
	m.MenuItem.SetSection(UnwrapMenuModel(section))
}

func (m *menuItem) SetSubmenu(submenu glibi.MenuModel) {
	m.MenuItem.SetSubmenu(UnwrapMenuModel(submenu))
}

func (m *menuItem) GetLink(link string) glibi.MenuModel {
	return WrapMenuModelSimple(m.MenuItem.GetLink(link))
}

func (m *menuItem) SetLink(link string, model glibi.MenuModel) {
	m.MenuItem.SetLink(link, UnwrapMenuModel(model))
}
