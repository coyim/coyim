package gtki

import "github.com/coyim/gotk3adapter/glibi"

type Application interface {
	glibi.Application

	GetActiveWindow() Window
	AddWindow(Window)
	RemoveWindow(Window)
	PrefersAppMenu() bool

	GetAppMenu() glibi.MenuModel
	SetAppMenu(glibi.MenuModel)
	GetMenubar() glibi.MenuModel
	SetMenubar(glibi.MenuModel)
}

func AssertApplication(_ Application) {}
