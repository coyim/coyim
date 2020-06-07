package glibi

type MenuModel interface {
	Object

	IsMutable() bool
	GetNItems() int
	GetItemLink(index int, link string) MenuModel
	ItemsChanged(position, removed, added int)
}

func AssertMenuModel(_ MenuModel) {}

type Menu interface {
	MenuModel

	Freeze()
	Insert(position int, label, detailed_action string)
	Prepend(label, detailed_action string)
	Append(label, detailed_action string)
	InsertItem(position int, item MenuItem)
	AppendItem(item MenuItem)
	PrependItem(item MenuItem)
	InsertSection(position int, label string, section MenuModel)
	PrependSection(label string, section MenuModel)
	AppendSection(label string, section MenuModel)
	InsertSectionWithoutLabel(position int, section MenuModel)
	PrependSectionWithoutLabel(section MenuModel)
	AppendSectionWithoutLabel(section MenuModel)
	InsertSubmenu(position int, label string, submenu MenuModel)
	PrependSubmenu(label string, submenu MenuModel)
	AppendSubmenu(label string, submenu MenuModel)
	Remove(position int)
	RemoveAll()
}

func AssertMenu(_ Menu) {}

type MenuItem interface {
	Object

	SetLabel(label string)
	SetDetailedAction(act string)
	SetSection(section MenuModel)
	SetSubmenu(submenu MenuModel)
	GetLink(link string) MenuModel
	SetLink(link string, model MenuModel)
}

func AssertMenuItem(_ MenuItem) {}
