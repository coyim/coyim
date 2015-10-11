package gui

import (
	"fmt"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/go-gtk/gtk"
)

func verifyFingerprintDialog(account *Account, uid string) *gtk.Dialog {
	dialog := gtk.NewDialog()
	dialog.SetTitle(i18n.Local("Fingerprint verification"))
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox := dialog.GetVBox()

	conversation := account.GetConversationWith(uid)
	fpr := conversation.GetTheirKey().DefaultFingerprint()

	// message copied from libpurple
	message := fmt.Sprintf(i18n.Local(`
	Fingerprint for you (%s): %x

	Purported fingerprint for %s: %x

	Is this the verifiably correct fingerprint for %s?
	`), account.Account, account.Session.PrivateKey.DefaultFingerprint(), uid, fpr, uid)
	vbox.Add(gtk.NewLabel(message))

	button := gtk.NewButtonWithLabel(i18n.Local("Verify"))
	vbox.Add(button)

	button.Connect("clicked", func() {
		defer dialog.Destroy()

		err := account.AuthorizeFingerprint(uid, fpr)
		if err != nil {
			//TODO: Error
			return
		}

		//TODO: error
		account.configManager.Save()
	})

	return dialog
}
