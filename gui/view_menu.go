package gui

import (
	"github.com/coyim/coyim/config"
	"github.com/coyim/gotk3adapter/gtki"
)

type viewMenu struct {
	merge      gtki.CheckMenuItem `gtk-widget:"CheckItemMerge"`
	offline    gtki.CheckMenuItem `gtk-widget:"CheckItemShowOffline"`
	waiting    gtki.CheckMenuItem `gtk-widget:"CheckItemShowWaiting"`
	sortStatus gtki.CheckMenuItem `gtk-widget:"CheckItemSortStatus"`
}

func (v *viewMenu) setFromConfig(c *config.ApplicationConfig) {
	doInUIThread(func() {
		v.merge.SetActive(c.Display.MergeAccounts)
		v.offline.SetActive(!c.Display.ShowOnlyOnline)
		v.waiting.SetActive(!c.Display.ShowOnlyConfirmed)
		v.sortStatus.SetActive(c.Display.SortByStatus)
	})
}

func (u *gtkUI) toggleMergeAccounts() {
	if u.config != nil {
		u.config.Display.MergeAccounts = u.viewMenu.merge.GetActive()
		u.saveConfigOnly()
	}

	u.roster.redraw()
}

func (u *gtkUI) toggleShowOffline() {
	if u.config != nil {
		u.config.Display.ShowOnlyOnline = !u.viewMenu.offline.GetActive()
		u.saveConfigOnly()
	}

	u.roster.redraw()
}

func (u *gtkUI) toggleShowWaiting() {
	if u.config != nil {
		u.config.Display.ShowOnlyConfirmed = !u.viewMenu.waiting.GetActive()
		u.saveConfigOnly()
	}

	u.roster.redraw()
}

func (u *gtkUI) toggleSortByStatus() {
	if u.config != nil {
		u.config.Display.SortByStatus = u.viewMenu.sortStatus.GetActive()
		u.saveConfigOnly()
	}

	u.roster.redraw()
}
