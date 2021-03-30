package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigListControllerData struct {
	addOccupantButton    gtki.Button
	removeOccupantButton gtki.Button
	occupantsTreeView    gtki.TreeView
	parentWindow         gtki.Window

	addOccupantDialogTitle string
	addOccupantDescription string
	addOccupantForm        func(onFormFieldValueChanges, onFormFieldValueActivates func()) *roomConfigListForm
	onUpdate               func()
}

type mucRoomConfigListController struct {
	u *gtkUI

	addOccupantButton        gtki.Button
	removeOccupantButton     gtki.Button
	ocuppantsTreeView        gtki.TreeView
	occupantsTreeViewColumns []glibi.Type

	listComponent        *mucRoomConfigListComponent
	onAddOccupantsToList func()
	onUpdate             func()
}

func (u *gtkUI) newMUCRoomConfigListController(d *mucRoomConfigListControllerData) *mucRoomConfigListController {
	c := &mucRoomConfigListController{
		u:                    u,
		addOccupantButton:    d.addOccupantButton,
		removeOccupantButton: d.removeOccupantButton,
		ocuppantsTreeView:    d.occupantsTreeView,
		onUpdate:             d.onUpdate,
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
			c.onAddListItems,
			d.parentWindow,
		)

		addToList.show()
	}
}

func (c *mucRoomConfigListController) initListComponent() {
	c.listComponent = c.u.newMUCRoomConfigListComponent(
		c.ocuppantsTreeView,
		c.addOccupantButton,
		c.removeOccupantButton,
		c.onAddOccupantsToList,
		c.onUpdated,
	)
}

func (c *mucRoomConfigListController) listItems() []string {
	return c.listComponent.jidList
}

func (c *mucRoomConfigListController) hasItems() bool {
	return len(c.listComponent.jidList) > 0
}

func (c *mucRoomConfigListController) onAddListItems(jidList []string) {
	c.listComponent.addListItems(jidList)
	c.onUpdated()
}

func (c *mucRoomConfigListController) onUpdated() {
	if c.onUpdate != nil {
		c.onUpdate()
	}
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
