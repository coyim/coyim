package gui

import (
	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

func buildVerifyFingerprintDialog(accountName string, ourFp []byte, peer jid.WithoutResource, theirFp []byte) (gtki.Dialog, func() string) {
	var message string
	var builderName string
	getTag := func() string { return "" }

	if theirFp == nil {
		builderName = "VerifyFingerprintUnknown"
		message = i18n.Localf(
			"You can't verify the fingerprint for %s yet.\n\n"+
				"You first have to start an encrypted conversation with them.", peer)

	} else {
		m := `
Is this the correct fingerprint for %[1]s?

Fingerprint for you (%[3]s):
  %[4]s

Purported fingerprint for %[1]s:
  %[2]s
	`

		message = i18n.Localf(m,
			peer,
			config.FormatFingerprint(theirFp),
			accountName,
			config.FormatFingerprint(ourFp),
		)

		builderName = "VerifyFingerprint"
	}

	builder := newBuilder(builderName)

	if theirFp != nil {
		entry := builder.getObj("tag-entry").(gtki.Entry)
		getTag = func() string {
			tag, _ := entry.GetText()
			return tag
		}
	}

	obj := builder.getObj("dialog")
	dialog := obj.(gtki.Dialog)

	obj = builder.getObj("message")
	l := obj.(gtki.Label)
	l.SetText(message)
	l.SetSelectable(true)

	dialog.SetTitle(i18n.Localf("Verify fingerprint for %s", peer))
	return dialog, getTag
}

func verifyFingerprintDialog(account *account, peer jid.Any, parent gtki.Window) gtki.ResponseType {
	accountConfig := account.session.GetConfig()
	conversation, _ := account.session.ConversationManager().EnsureConversationWith(peer, nil)
	ourFp := conversation.OurFingerprint()
	theirFp := conversation.TheirFingerprint()

	dialog, gettag := buildVerifyFingerprintDialog(accountConfig.Account, ourFp, peer.NoResource(), theirFp)
	defer dialog.Destroy()

	dialog.SetTransientFor(parent)
	dialog.ShowAll()

	responseType := gtki.ResponseType(dialog.Run())
	switch responseType {
	case gtki.RESPONSE_YES:
		account.executeCmd(otrclient.AuthorizeFingerprintCmd{
			Account:     accountConfig,
			Session:     account.session,
			Peer:        peer.NoResource(),
			Fingerprint: theirFp,
			Tag:         gettag(),
		})
	}

	return responseType
}
