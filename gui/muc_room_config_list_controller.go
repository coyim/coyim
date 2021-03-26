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
	refreshView            func()
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
	c.initListComponent()

	return c
}

func (c *mucRoomConfigListController) initListAddComponent(d *mucRoomConfigListControllerData) {
	c.onAddOccupantsToList = func() {
		addToList := c.u.newMUCRoomConfigListAddComponent(
			d.addOccupantDialogTitle,
			d.addOccupantDescription,
			d.addOccupantForm,
			func(items [][]string) {
				c.listComponent.addListItems(items)
				d.refreshView()
			},
			d.parentWindow,
		)

		addToList.show()
	}
}

func (c *mucRoomConfigListController) initListComponent() {
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

func (c *mucRoomConfigListController) hasItems() bool {
	return len(c.listComponent.items) > 0
}

func (c *mucRoomConfigListController) updateCellForString(column int, path string, newValue string) error {
	iter, err := c.listComponent.iterForString(path)
	if err != nil {
		return err
	}

	err = c.listComponent.updateValueBasedOnIter(iter, column, newValue)
	if err != nil {
		return err
	}

	return nil
}
