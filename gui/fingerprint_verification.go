package gui

import (
	"fmt"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3adapter/gtki"
)

func buildVerifyFingerprintDialog(accountName string, ourFp []byte, uid string, theirFp []byte) gtki.Dialog {
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

	builder := newBuilder(builderName)

	obj := builder.getObj("dialog")
	dialog := obj.(gtki.Dialog)

	obj = builder.getObj("message")
	l := obj.(gtki.Label)
	l.SetText(message)
	l.SetSelectable(true)

	dialog.SetTitle(fmt.Sprintf(i18n.Local("Verify fingerprint for %s"), uid))
	return dialog
}

func verifyFingerprintDialog(account *account, uid, resource string, parent gtki.Window) gtki.ResponseType {
	accountConfig := account.session.GetConfig()
	conversation, _ := account.session.ConversationManager().EnsureConversationWith(uid, resource)
	ourFp := conversation.OurFingerprint()
	theirFp := conversation.TheirFingerprint()

	dialog := buildVerifyFingerprintDialog(accountConfig.Account, ourFp, uid, theirFp)
	defer dialog.Destroy()

	dialog.SetTransientFor(parent)
	dialog.ShowAll()

	responseType := gtki.ResponseType(dialog.Run())
	switch responseType {
	case gtki.RESPONSE_YES:
		account.executeCmd(client.AuthorizeFingerprintCmd{
			Account:     accountConfig,
			Session:     account.session,
			Peer:        uid,
			Fingerprint: theirFp,
		})
	}

	return responseType
}
