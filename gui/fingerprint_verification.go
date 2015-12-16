package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"../client"
	"../config"
	"../i18n"
)

func buildVerifyFingerprintDialog(accountName string, ourFp []byte, uid string, theirFp []byte) *gtk.Dialog {
	var message string
	var builderName string

	if theirFp == nil {
		builderName = "VerifyFingerprintUnknown"
		message = fmt.Sprintf(i18n.Local(
			"You can't verify the fingerprint for %s yet.\n"+
				"You first have to start an encrypted conversation with them.",
		), uid)

	} else {
		m := i18n.Local(`
Is this the correct fingerprint for %[1]s?

Fingerprint for you (%[3]s):
  %[4]s

Purported fingerprint for %[1]s:
  %[2]s
	`)

		message = fmt.Sprintf(m,
			uid,
			config.FormatFingerprint(theirFp),
			accountName,
			config.FormatFingerprint(ourFp),
		)

		builderName = "VerifyFingerprint"
	}

	builder := builderForDefinition(builderName)

	obj, _ := builder.GetObject("dialog")
	dialog := obj.(*gtk.Dialog)

	obj, _ = builder.GetObject("message")
	l := obj.(*gtk.Label)
	l.SetText(message)

	dialog.SetTitle(fmt.Sprintf(i18n.Local("Verify fingerprint for %s"), uid))
	return dialog
}

func (u *gtkUI) verifyFingerprintDialog(account *account, uid string, parent *gtk.Window) gtk.ResponseType {
	accountConfig := account.session.CurrentAccount
	//TODO: review whether it should create new conversations
	//Anyway, if it has created the conversation this function could return
	//(there is no theirFP in this case)
	conversation, _ := account.session.EnsureConversationWith(uid)
	ourFp := conversation.OurFingerprint()
	theirFp := conversation.TheirFingerprint()

	dialog := buildVerifyFingerprintDialog(accountConfig.Account, ourFp, uid, theirFp)
	u.displaySettings.globalFontSettingOn(&dialog.Window.Bin.Container.Widget)

	defer dialog.Destroy()

	dialog.SetTransientFor(parent)
	dialog.ShowAll()

	responseType := gtk.ResponseType(dialog.Run())
	switch responseType {
	case gtk.RESPONSE_YES:
		account.executeCmd(client.AuthorizeFingerprintCmd{
			Account:     accountConfig,
			Peer:        uid,
			Fingerprint: theirFp,
		})
	}

	return responseType
}
