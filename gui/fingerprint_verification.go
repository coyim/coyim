package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/otr3"
)

func formatFingerprint(fpr []byte) string {
	str := fmt.Sprintf("%X", fpr)
	result := ""

	sep := ""
	for len(str) > 0 {
		result = result + sep + str[0:10]
		sep = " "
		str = str[10:]
	}

	return result
}

func verifyFingerprintDialog(account *account, uid string, parent *gtk.Window) {
	//TODO: errors
	dialog, _ := gtk.DialogNew()
	defer dialog.Destroy()
	dialog.SetTransientFor(parent)
	dialog.SetTitle(fmt.Sprintf(i18n.Local("Verify fingerprint for %s"), uid))
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox, _ := dialog.GetContentArea()
	vbox.SetBorderWidth(10)

	conversation := account.session.GetConversationWith(uid)
	if conversation == nil || conversation.GetTheirKey() == nil {
		pkey := account.session.PrivateKeys[0]
		if conversation == nil {
			conversation = otr3.NewConversationWithVersion(3)
		} else {
			pkey = conversation.GetOurCurrentKey()
		}

		message := fmt.Sprintf(i18n.Local(`
You can't verify the fingerprint for %s yet.
You first have to start an encrypted conversation with them.

Fingerprint for you (%s):
  %s
	`), uid, account.session.CurrentAccount.Account, formatFingerprint(conversation.DefaultFingerprintFor(pkey.PublicKey())))

		l, _ := gtk.LabelNew(message)
		vbox.Add(l)
		dialog.AddButton(i18n.Local("OK"), gtk.RESPONSE_OK)
		dialog.SetDefaultResponse(gtk.RESPONSE_OK)
		dialog.ShowAll()
		dialog.Run()
	} else {
		fpr := conversation.DefaultFingerprintFor(conversation.GetTheirKey())
		message := fmt.Sprintf(i18n.Local(`
Is this the correct fingerprint for %s?

Fingerprint for you (%s):
  %s

Purported fingerprint for %s:
  %s
	`), uid, account.session.CurrentAccount.Account, formatFingerprint(conversation.DefaultFingerprintFor(conversation.GetOurCurrentKey().PublicKey())), uid, formatFingerprint(fpr))

		l, _ := gtk.LabelNew(message)
		vbox.Add(l)
		dialog.AddButton(i18n.Local("Cancel"), gtk.RESPONSE_NO)
		dialog.AddButton(i18n.Local("Verify"), gtk.RESPONSE_YES)
		dialog.SetDefaultResponse(gtk.RESPONSE_NO)

		dialog.ShowAll()

		responseType := gtk.ResponseType(dialog.Run())
		switch responseType {
		case gtk.RESPONSE_YES:
			account.ExecuteCmd(client.AuthorizeFingerprintCmd{
				Account:     account.session.CurrentAccount,
				Peer:        uid,
				Fingerprint: fpr,
			})
		}
	}
}
