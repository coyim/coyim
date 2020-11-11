package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewMenuItem interface {
	getRoomMenuItem() gtki.Widget
}

type roomViewMenuButton struct {
	button gtki.ModelButton `gtk-widget:"room-menu-button"`
}

func newRoomViewMenuButton(l string, onClick func()) *roomViewMenuButton {
	mb := &roomViewMenuButton{}

	b := newBuilder("MUCRoomMenuButton")
	panicOnDevError(b.bindObjects(mb))

	b.ConnectSignals(map[string]interface{}{
		"on_clicked": onClick,
	})

	mb.button.SetLabel(l)

	return mb
}

// implements roomViewMenuItem interface
func (b *roomViewMenuButton) getRoomMenuItem() gtki.Widget {
	return b.button
}

type roomViewMenuDivider struct {
	s gtki.Separator
}

// implements roomViewMenuItem interface
func (d *roomViewMenuDivider) getRoomMenuItem() gtki.Widget {
	return d.s
}

func newRoomViewMenuDivider() *roomViewMenuDivider {
	d, _ := g.gtk.SeparatorNew(gtki.HorizontalOrientation)
	return &roomViewMenuDivider{d}
}

type roomViewMenu struct {
	items map[string]roomViewMenuItem

	popover gtki.Popover `gtk-widget:"room-menu"`
	view    gtki.Box     `gtk-widget:"room-menu-box"`
}

// newRoomViewMenu MUST be called from the UI thread
func newRoomViewMenu() *roomViewMenu {
	m := &roomViewMenu{
		items: make(map[string]roomViewMenuItem),
	}

	m.initBuilder()

	return m
}

func (m *roomViewMenu) initBuilder() {
	b := newBuilder("MUCRoomMenu")
	panicOnDevError(b.bindObjects(m))
}

// setMenuItem MUST always be called from the UI thread
func (m *roomViewMenu) setMenuItem(id string, item roomViewMenuItem) {
	m.items[id] = item
	m.redraw()
}

// setButtonItem MUST always be called from the UI thread
func (m *roomViewMenu) setButtonItem(id, l string, f func()) {
	m.setMenuItem(id, newRoomViewMenuButton(l, f))
}

// addDividerItem MUST always be called from the UI thread
func (m *roomViewMenu) addDividerItem() {
	m.setMenuItem(fmt.Sprintf("divider-%d", len(m.items)+1), newRoomViewMenuDivider())
}

// redraw MUST be called from the UI thread
func (m *roomViewMenu) redraw() {
	m.removeAll()

	for _, i := range m.items {
		m.view.Add(i.getRoomMenuItem())
	}
}

// removeAll MUST be called from the UI thread
func (m *roomViewMenu) removeAll() {
	for _, i := range m.items {
		m.view.Remove(i.getRoomMenuItem())
	}
}

// reset MUST be called from the UI thread
//
// The difference between this method and "removeAll" is
// that this method will remove all items from the list and from the view,
// while "removeAll" only will remove elements from the view
func (m *roomViewMenu) reset() {
	m.removeAll()
	m.items = make(map[string]roomViewMenuItem)
}

// initRoomMenu MUST be called from the UI thread
func (v *roomView) initRoomMenu() {
	v.menu = newRoomViewMenu()
	v.refreshRoomMenu()
}

// refreshRoomMenu MUST be called from the UI thread
func (v *roomView) refreshRoomMenu() {
	v.menu.reset()

	if v.isSelfOccupantAnOwner() {
		v.menu.setButtonItem("destroy-room", i18n.Local("Destroy room"), v.onDestroyRoom)
		v.menu.addDividerItem()
	}

	if v.isSelfOccupantJoined() {
		v.menu.setButtonItem("leave-room", i18n.Local("Leave room"), v.onLeaveRoom)
	}
}

func (v *roomView) getRoomMenuWidget() gtki.Popover {
	return v.menu.popover
}
