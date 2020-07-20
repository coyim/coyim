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
	dialog                 gtki.Window    `gtk-widget:"AddContact"`
	model                  gtki.ListStore `gtk-widget:"accounts-model"`
	accountInput           gtki.ComboBox  `gtk-widget:"accounts"`
	contactInput           gtki.Entry     `gtk-widget:"address"`
	notificationArea       gtki.Box       `gtk-widget:"notification-area"`
	notification           gtki.InfoBar
	subscriptionAskMessage gtki.TextBuffer  `gtk-widget:"subscriptionAskMessage"`
	nickname               gtki.Entry       `gtk-widget:"nickname"`
	autoAuth               gtki.CheckButton `gtk-widget:"auto_authorize_checkbutton"`
	accounts               map[string]*account
}

func (acd *addContactDialog) getVerifiedContact(errorNotif *errorNotification) (string, bool) {
	contact, _ := acd.contactInput.GetText()
	isJid, err := verifyXMPPAddress(contact)

	if !isJid {
		if acd.notification != nil {
			acd.notificationArea.Remove(acd.notification)
		}

		errorNotif.ShowMessage(err)
		log.WithField("error", err).Warn("Bad XMPP address entered")

		return "", false
	}

	errorNotif.Hide() // no errors
	return contact, true
}

func (acd *addContactDialog) getCurrentAccount() (string, error) {
	iter, err := acd.accountInput.GetActiveIter()
	if err != nil {
		return "", err
	}
	val, err := acd.model.GetValue(iter, 1)
	if err != nil {
		return "", err
	}
	return val.GetString()
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

func (acd *addContactDialog) initAccounts(accounts []*account) {
	acd.accounts = make(map[string]*account)
	for _, acc := range accounts {
		iter := acd.model.Append()
		_ = acd.model.SetValue(iter, 0, acc.session.GetConfig().Account)
		_ = acd.model.SetValue(iter, 1, acc.session.GetConfig().ID())
		acd.accounts[acc.session.GetConfig().ID()] = acc
	}

	if len(accounts) > 0 {
		acd.accountInput.SetActive(0)
	}
}

func (acd *addContactDialog) init() {
	acd.builder = newBuilder("AddContact")
	panicOnDevError(acd.builder.bindObjects(acd))
}

func presenceSubscriptionDialog(accounts []*account, sendSubscription func(accountID string, peer jid.WithoutResource, msg, nick string, autoauth bool) error) gtki.Window {
	//TODO: this can be opened before a account is connected.
	//In this case the window is useless: cant add a contact and cant see an error
	acd := &addContactDialog{}
	acd.init()
	acd.initAccounts(accounts)

	errorNotif := newErrorNotification(acd.notificationArea)

	acd.builder.ConnectSignals(map[string]interface{}{
		"on_cancel_signal": acd.dialog.Destroy,
		"on_save_signal": func() {
			contact, ok := acd.getVerifiedContact(errorNotif)
			if !ok {
				return
			}

			accountID, err := acd.getCurrentAccount()
			if err != nil {
				//TODO: report error, and close?
				log.WithError(err).Warn("Error encountered when getting account")
				return
			}

			acc := acd.accounts[accountID]

			err = sendSubscription(accountID, jid.NR(contact), acd.getCurrentMessage(), acd.getCurrentNickname(), acd.getAutoAuthorize())
			if err != nil {
				//TODO: report error
				acc.log.WithError(err).Warn("Error encountered when sending subscription")
				return
			}

			acd.dialog.Destroy()
		},
	})

	return acd.dialog
}
