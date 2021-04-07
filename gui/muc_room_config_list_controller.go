package gui

import "github.com/coyim/gotk3adapter/gtki"

type mucRoomConfigListControllerData struct {
	addOccupantButton      gtki.Button
	removeOccupantButton   gtki.Button
	removeOccupantLabel    gtki.Label
	occupantsTreeView      gtki.TreeView
	parentWindow           gtki.Window
	addOccupantDialogTitle string
	addOccupantDescription string
	onListUpdated          func() // onListUpdated will be called from the UI thread
}

type mucRoomConfigListController struct {
	u                  *gtkUI
	listComponent      *mucRoomConfigListComponent
	doAfterListUpdated func() // doAfterListUpdated will be called from the UI thread
}

func (u *gtkUI) newMUCRoomConfigListController(d *mucRoomConfigListControllerData) *mucRoomConfigListController {
	c := &mucRoomConfigListController{
		u:                  u,
		doAfterListUpdated: d.onListUpdated,
	}

	c.initListComponent(d)

	return c
}

func (c *mucRoomConfigListController) initListComponent(d *mucRoomConfigListControllerData) {
	c.listComponent = c.u.newMUCRoomConfigListComponent(
		d.occupantsTreeView,
		d.addOccupantButton,
		d.removeOccupantButton,
		d.removeOccupantLabel,
		c.onAddOccupantsToList(d),
		c.doAfterListUpdated,
	)
}

// onAddOccupantsToList MUST be called from the UI thread
func (c *mucRoomConfigListController) onAddOccupantsToList(d *mucRoomConfigListControllerData) func() {
	return func() {
		addToList := c.u.newMUCRoomConfigListAddComponent(
			d.addOccupantDialogTitle,
			d.addOccupantDescription,
			c.onAddListItems,
			d.parentWindow,
		)

		addToList.show()
	}
}

// onAddListItems MUST be called from the UI thread
func (c *mucRoomConfigListController) onAddListItems(jidList []string) {
	c.listComponent.addListItems(jidList)
	c.onListUpdated()
}

// onUpdated MUST be called from the UI thread
func (c *mucRoomConfigListController) onListUpdated() {
	if c.doAfterListUpdated != nil {
		c.doAfterListUpdated()
	}
}

// updateCellForString MUST be called from the UI thread
func (c *mucRoomConfigListController) updateCellForString(column int, path string, newValue string) bool {
	iter, err := c.listComponent.listModel.GetIterFromString(path)
	if err != nil {
		return false
	}

	for _, v := range c.listItems() {
		if v == newValue {
			return false
		}
	}

	err = c.listComponent.listModel.SetValue(iter, column, newValue)
	if err != nil {
		return false
	}

	return true
}

func (c *mucRoomConfigListController) listItems() []string {
	return toArray(c.listComponent.listModel)
}

func (c *mucRoomConfigListController) hasItems() bool {
	_, ok := c.listComponent.listModel.GetIterFirst()
	return ok
}
