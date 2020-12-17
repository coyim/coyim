package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigListControllerData struct {
	addOccupantButton        gtki.Button
	removeOccupantButton     gtki.Button
	occupantsTreeView        gtki.TreeView
	occupantsTreeViewColumns []glibi.Type

	addOccupantDialogTitle string
	addOccupantDescription string
	addOccupantForm        func(onAnyValueChanges func()) mucRoomConfigListForm
}

type mucRoomConfigListController struct {
	u *gtkUI

	addOccupantButton        gtki.Button
	removeOccupantButton     gtki.Button
	ocuppantsTreeView        gtki.TreeView
	occupantsTreeViewColumns []glibi.Type

	listComponent    *mucRoomConfigListComponent
	addComponentForm mucRoomConfigListForm
	addComponent     *mucRoomConfigListAddComponent
}

func (u *gtkUI) newMUCRoomConfigListController(d *mucRoomConfigListControllerData) *mucRoomConfigListController {
	c := &mucRoomConfigListController{
		u:                        u,
		addOccupantButton:        d.addOccupantButton,
		removeOccupantButton:     d.removeOccupantButton,
		ocuppantsTreeView:        d.occupantsTreeView,
		occupantsTreeViewColumns: d.occupantsTreeViewColumns,
	}

	c.initListComponent(d)
	c.initListAddFormComponent(d)
	c.initListAddComponent(d)

	return c
}

func (c *mucRoomConfigListController) initListComponent(d *mucRoomConfigListControllerData) {
	c.listComponent = c.u.newMUCRoomConfigListComponent(
		c.ocuppantsTreeView,
		c.occupantsTreeViewColumns,
		c.addOccupantButton,
		c.removeOccupantButton,
		c.onAddOccupantsToList,
	)
}

func (c *mucRoomConfigListController) initListAddFormComponent(d *mucRoomConfigListControllerData) {
	c.addComponentForm = d.addOccupantForm(nil)
}

func (c *mucRoomConfigListController) initListAddComponent(d *mucRoomConfigListControllerData) {
	c.addComponent = c.u.newMUCRoomConfigListAddComponent(
		d.addOccupantDialogTitle,
		d.addOccupantDescription,
		c.addComponentForm,
		c.listComponent.addListItem,
	)
}

func (c *mucRoomConfigListController) onAddOccupantsToList() {
	c.addComponent.show()
}

func (c *mucRoomConfigListController) listItems() [][]string {
	return c.listComponent.items
}
