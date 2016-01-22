package gui

import "github.com/gotk3/gotk3/glib"

var weAreInUIThread = false

func assertInUIThread() {
	if !weAreInUIThread {
		panic("This function have to be called from the UI thread")
	}
}

func doInUIThread(f func()) {
	glib.IdleAdd(func() {
		weAreInUIThread = true
		defer func() {
			weAreInUIThread = false
		}()
		f()
	})
}
