package gui

import (
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

func authorizePresenceSubscriptionDialog(parent gtki.Window, peer jid.WithoutResource, f func(gtki.ResponseType)) {
	builder := newBuilder("AuthorizeSubscription")

	confirmDialog := builder.getObj("dialog").(gtki.MessageDialog)
	text := i18n.Localf("%s wants to talk to you. Is that ok?", peer)
	_ = confirmDialog.SetProperty("text", text)
	confirmDialog.SetTransientFor(parent)

	_, _ = confirmDialog.Connect("response", func(_ interface{}, tp int) {
		f(gtki.ResponseType(tp))
		confirmDialog.Destroy()
	})

	doInUIThread(func() {
		confirmDialog.ShowAll()
	})
}

type addContactDialog struct {
	builder                *builder
	dialog                 gtki.Window `gtk-widget:"AddContact"`
	contactInput           gtki.Entry  `gtk-widget:"address"`
	notificationArea       gtki.Box    `gtk-widget:"notification-area"`
	notification           gtki.InfoBar
	subscriptionAskMessage gtki.TextBuffer  `gtk-widget:"subscriptionAskMessage"`
	nickname               gtki.Entry       `gtk-widget:"nickname"`
	autoAuth               gtki.CheckButton `gtk-widget:"auto_authorize_checkbutton"`
	errorNotif             *errorNotification
}

func (acd *addContactDialog) clearErrors() {
	acd.errorNotif.Hide()
}

func (acd *addContactDialog) notifyOnError(err string) {
	if acd.notification != nil {
		acd.notificationArea.Remove(acd.notification)
	}

	acd.errorNotif.ShowMessage(err)
}

// getVerifiedContact should ONLY be called from the UI thread
func (acd *addContactDialog) getVerifiedContact() (string, bool) {
	contact, _ := acd.contactInput.GetText()
	isJid, err := verifyXMPPAddress(contact)

	if !isJid {
		acd.notifyOnError(err)
		log.WithField("error", err).Warn("Bad XMPP address entered")

		return "", false
	}

	acd.clearErrors()
	return contact, true
}

func (acd *addContactDialog) getCurrentMessage() string {
	return acd.subscriptionAskMessage.GetText(
		acd.subscriptionAskMessage.GetStartIter(),
		acd.subscriptionAskMessage.GetEndIter(),
		false,
	)
}

func (acd *addContactDialog) getCurrentNickname() string {
	txt, _ := acd.nickname.GetText()
	return txt
}

func (acd *addContactDialog) getAutoAuthorize() bool {
	return acd.autoAuth.GetActive()
}

func (acd *addContactDialog) init() {
	acd.builder = newBuilder("AddContact")
	panicOnDevError(acd.builder.bindObjects(acd))
	acd.errorNotif = newErrorNotification(acd.notificationArea)
}

func (u *gtkUI) presenceSubscriptionDialog(sendSubscription func(accountID string, peer jid.WithoutResource, msg, nick string, autoauth bool) error) gtki.Window {
	acd := &addContactDialog{}
	acd.init()

	accountsInput := acd.builder.get("accounts").(gtki.ComboBox)
	ac := u.createConnectedAccountsComponent(accountsInput, acd,
		func(acc *account) {
			doInUIThread(acd.clearErrors)
		},
		func() {
			doInUIThread(func() { acd.notifyOnError(i18n.Local("There are no currently connected accounts")) })
		},
	)

	acd.builder.ConnectSignals(map[string]interface{}{
		"on_cancel_signal": func() {
			acd.dialog.Destroy()
			ac.onDestroy()
		},
		"on_save_signal": func() {
			contact, ok := acd.getVerifiedContact()
			if !ok {
				return
			}

			acc := ac.currentAccount()
			if acc == nil {
				acd.notifyOnError(i18n.Local("There is no connected account selected"))
				acc.log.Warn("can't send subscription without a current account")
				return
			}

			err := sendSubscription(acc.ID(), jid.NR(contact), acd.getCurrentMessage(), acd.getCurrentNickname(), acd.getAutoAuthorize())
			if err != nil {
				acd.notifyOnError(i18n.Local("We couldn't send a subscription"))
				acc.log.WithError(err).Warn("Error encountered when sending subscription")
				return
			}

			acd.dialog.Destroy()
			ac.onDestroy()
		},
	})

	return acd.dialog
}
