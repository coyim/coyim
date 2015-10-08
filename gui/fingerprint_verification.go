package gui

import "github.com/twstrike/go-gtk/gtk"

func verifyFingerprintDialog() *gtk.Dialog {
	dialog := gtk.NewDialog()
	dialog.SetTitle("Fingerprint verification")
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox := dialog.GetVBox()

	button := gtk.NewButtonWithLabel("Verify")
	vbox.Add(button)

	button.Connect("clicked", func() {
		dialog.Destroy()
	})

	return dialog
}
