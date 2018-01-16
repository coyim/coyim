package gui

import (
	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/otr_client"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

func buildVerifyFingerprintDialog(accountName string, ourFp []byte, uid string, theirFp []byte) gtki.Dialog {
	var message string
	var builderName string

	if theirFp == nil {
		builderName = "VerifyFingerprintUnknown"
		message = i18n.Localf(
			"You can't verify the fingerprint for %s yet.\n\n"+
				"You first have to start an encrypted conversation with them.", uid)

	} else {
		m := `
Is this the correct fingerprint for %[1]s?

Fingerprint for you (%[3]s):
  %[4]s

Purported fingerprint for %[1]s:
  %[2]s
	`

		message = i18n.Localf(m,
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

	dialog.SetTitle(i18n.Localf("Verify fingerprint for %s", uid))
	return dialog
}

func verifyFingerprintDialog(account *account, uid jid.WithoutResource, resource jid.Resource, parent gtki.Window) gtki.ResponseType {
	accountConfig := account.session.GetConfig()
	conversation, _ := account.session.ConversationManager().EnsureConversationWith(uid, resource)
	ourFp := conversation.OurFingerprint()
	theirFp := conversation.TheirFingerprint()

	dialog := buildVerifyFingerprintDialog(accountConfig.Account, ourFp, uid.String(), theirFp)
	defer dialog.Destroy()

	dialog.SetTransientFor(parent)
	dialog.ShowAll()

	responseType := gtki.ResponseType(dialog.Run())
	switch responseType {
	case gtki.RESPONSE_YES:
		account.executeCmd(otr_client.AuthorizeFingerprintCmd{
			Account:     accountConfig,
			Session:     account.session,
			Peer:        uid,
			Fingerprint: theirFp,
		})
	}

	return responseType
}
