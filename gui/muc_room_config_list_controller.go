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
	parentWindow             gtki.Window

	addOccupantDialogTitle string
	addOccupantDescription string
	addOccupantForm        func(onFormFieldValueChanges, onFormFieldValueActivates func()) mucRoomConfigListForm
}

type mucRoomConfigListController struct {
	u *gtkUI

	addOccupantButton        gtki.Button
	removeOccupantButton     gtki.Button
	ocuppantsTreeView        gtki.TreeView
	occupantsTreeViewColumns []glibi.Type

	listComponent        *mucRoomConfigListComponent
	onAddOccupantsToList func()
}

func (u *gtkUI) newMUCRoomConfigListController(d *mucRoomConfigListControllerData) *mucRoomConfigListController {
	c := &mucRoomConfigListController{
		u:                        u,
		addOccupantButton:        d.addOccupantButton,
		removeOccupantButton:     d.removeOccupantButton,
		ocuppantsTreeView:        d.occupantsTreeView,
		occupantsTreeViewColumns: d.occupantsTreeViewColumns,
	}

	c.initListAddComponent(d)
	c.initListComponent(d)

	return c
}

func (c *mucRoomConfigListController) initListAddComponent(d *mucRoomConfigListControllerData) {
	c.onAddOccupantsToList = func() {
		addToList := c.u.newMUCRoomConfigListAddComponent(
			d.addOccupantDialogTitle,
			d.addOccupantDescription,
			d.addOccupantForm,
			c.listComponent.addListItem,
			d.parentWindow,
		)

		addToList.show()
	}
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

func (c *mucRoomConfigListController) listItems() [][]string {
	return c.listComponent.items
}
