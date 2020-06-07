package glib_mock

import "github.com/coyim/gotk3adapter/glibi"

type MockMenuModel struct {
	MockObject
}

func (*MockMenuModel) IsMutable() bool {
	return false
}

func (*MockMenuModel) GetNItems() int {
	return 0
}

func (*MockMenuModel) GetItemLink(index int, link string) glibi.MenuModel {
	return nil
}

func (*MockMenuModel) ItemsChanged(position, removed, added int) {
}

type MockMenu struct {
	MockMenuModel
}

func (*MockMenu) Freeze() {
}

func (*MockMenu) Insert(position int, label, detailed_action string) {
}

func (*MockMenu) Prepend(label, detailed_action string) {
}

func (*MockMenu) Append(label, detailed_action string) {
}

func (*MockMenu) InsertItem(position int, item glibi.MenuItem) {
}

func (*MockMenu) AppendItem(item glibi.MenuItem) {
}

func (*MockMenu) PrependItem(item glibi.MenuItem) {
}

func (*MockMenu) InsertSection(position int, label string, section glibi.MenuModel) {
}

func (*MockMenu) PrependSection(label string, section glibi.MenuModel) {
}

func (*MockMenu) AppendSection(label string, section glibi.MenuModel) {
}

func (*MockMenu) InsertSectionWithoutLabel(position int, section glibi.MenuModel) {
}

func (*MockMenu) PrependSectionWithoutLabel(section glibi.MenuModel) {
}

func (*MockMenu) AppendSectionWithoutLabel(section glibi.MenuModel) {
}

func (*MockMenu) InsertSubmenu(position int, label string, submenu glibi.MenuModel) {
}

func (*MockMenu) PrependSubmenu(label string, submenu glibi.MenuModel) {
}

func (*MockMenu) AppendSubmenu(label string, submenu glibi.MenuModel) {
}

func (*MockMenu) Remove(position int) {
}

func (*MockMenu) RemoveAll() {
}

type MockMenuItem struct {
	MockObject
}

func (*MockMenuItem) SetLabel(label string) {
}

func (*MockMenuItem) SetDetailedAction(act string) {
}

func (*MockMenuItem) SetSection(section glibi.MenuModel) {
}

func (*MockMenuItem) SetSubmenu(submenu glibi.MenuModel) {
}

func (*MockMenuItem) GetLink(link string) glibi.MenuModel {
	return nil
}

func (*MockMenuItem) SetLink(link string, model glibi.MenuModel) {
}
