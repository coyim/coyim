package gui

import (
	"fmt"
	"log"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/gotk3adapter/gtki"
)

func authorizePresenceSubscriptionDialog(parent gtki.Window, from string) gtki.MessageDialog {
	builder := builderForDefinition("AuthorizeSubscription")

	obj, _ := builder.GetObject("dialog")
	confirmDialog := obj.(gtki.MessageDialog)

	text := fmt.Sprintf(i18n.Local("%s wants to talk to you. Is that ok?"), from)
	confirmDialog.SetProperty("text", text)

	confirmDialog.SetTransientFor(parent)
	return confirmDialog
}

func presenceSubscriptionDialog(accounts []*account, sendSubscription func(accountID, peer string) error) gtki.Dialog {
	builder := builderForDefinition("AddContact")

	//TODO: move model to XML builder
	model, _ := g.gtk.ListStoreNew(
		glibi.TYPE_STRING, // account name
		glibi.TYPE_STRING, // account_id
	)

	for _, acc := range accounts {
		model.Set2(model.Append(), []int{0, 1}, []interface{}{acc.session.GetConfig().Account, acc.session.GetConfig().ID()})
	}

	accountsObj, _ := builder.GetObject("accounts")
	accountInput := accountsObj.(gtki.ComboBox)
	accountInput.SetModel(model)

	accountObj, _ := builder.GetObject("address")
	contactInput := accountObj.(gtki.Entry)

	if len(accounts) > 0 {
		accountInput.SetActive(0)
	}

	renderer, _ := g.gtk.CellRendererTextNew()
	accountInput.PackStart(renderer, true)
	accountInput.AddAttribute(renderer, "text", 0)

	dialogObj, _ := builder.GetObject("AddContact")
	dialog := dialogObj.(gtki.Dialog)

	obj, _ := builder.GetObject("notification-area")
	notificationArea := obj.(gtki.Box)

	failures := 0
	var notification gtki.InfoBar

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			contact, _ := contactInput.GetText()
			isJid, errmsg := verifyXmppAddress(contact)

			if !isJid && failures > 0 {
				notificationArea.Remove(notification)
				notification = buildBadUsernameNotification(errmsg)
				notificationArea.Add(notification)
				notification.ShowAll()
				failures++
				log.Printf(errmsg)
				return
			}

			if !isJid {
				notification = buildBadUsernameNotification(errmsg)
				notificationArea.Add(notification)
				notification.ShowAll()
				failures++
				log.Printf(errmsg)
				return
			}

			iter, err := accountInput.GetActiveIter()
			if err != nil {
				log.Printf("Error encountered when getting account: %v", err)
				return
			}
			val, err := model.GetValue(iter, 1)
			if err != nil {
				log.Printf("Error encountered when getting account: %v", err)
				return
			}
			accountID, err := val.GetString()
			if err != nil {
				log.Printf("Error encountered when getting account: %v", err)
				return
			}

			err = sendSubscription(accountID, contact)
			if err != nil {
				log.Printf("Error encountered when sending subscription: %v", err)
				return
			}

			dialog.Destroy()
		},
	})

	return dialog
}
