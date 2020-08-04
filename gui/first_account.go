package gui

import "github.com/coyim/gotk3adapter/gtki"

func (u *gtkUI) showFirstAccountWindow() {
	result := make(chan func())
	var cleanup func()

	doInUIThread(func() {
		b := newBuilder("FirstAccountDialog")
		dialog := b.getObj("dialog").(gtki.Dialog)
		dialog.SetTransientFor(u.window)
		cleanup = dialog.Destroy

		b.ConnectSignals(map[string]interface{}{
			"on_register": func() {
				result <- u.showServerSelectionWindow
			},
			"on_existing": func() {
				result <- u.showAddAccountWindow
			},
			"on_import": func() {
				result <- u.runImporter
			},
			"on_cancel": func() {
				result <- func() {}
			},
		})

		dialog.ShowAll()
	})

	tp := <-result
	doInUIThread(cleanup)
	doInUIThread(tp)
}
