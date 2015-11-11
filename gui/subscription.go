package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
)

func authorizePresenceSubscriptionDialog(parent *gtk.Window, from string) *gtk.MessageDialog {
	confirmDialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_QUESTION,
		gtk.BUTTONS_YES_NO,
		i18n.Local("%s wants to talk to you. Is that ok?"), from,
	)
	confirmDialog.SetTitle(i18n.Local("Subscription request"))

	return confirmDialog
}

func presenceSubscriptionDialog(accounts []*account, sendSubscription func(accountID, peer string) error) *gtk.Dialog {
	vars := make(map[string]string)
	vars["$title"] = i18n.Local("Add contact")
	vars["$accountsLabel"] = i18n.Local("Account")
	vars["$addressLabel"] = i18n.Local("Account to add (for example: arnoldsPub@jabber.ccc.de)")
	vars["$addLabel"] = i18n.Local("Add")

	builder, loadErr := loadBuilderWith("AddContactDefinition", vars)
	if loadErr != nil {
		panic(loadErr.Error())
	}

	model, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // account name
		glib.TYPE_STRING, // account_id
	)

	for _, acc := range accounts {
		iter := model.Append()
		//TODO stop passing pointers
		model.Set(iter, []int{0, 1}, []interface{}{acc.session.CurrentAccount.Account, acc.session.CurrentAccount.ID()})
	}

	accountsObj, _ := builder.GetObject("accounts")
	accountInput, _ := accountsObj.(*gtk.ComboBox)
	accountInput.SetModel(&model.TreeModel)

	accountObj, _ := builder.GetObject("address")
	contactInput := accountObj.(*gtk.Entry)

	if len(accounts) > 0 {
		accountInput.SetActive(0)
	}

	renderer, _ := gtk.CellRendererTextNew()
	accountInput.PackStart(renderer, true)
	accountInput.AddAttribute(renderer, "text", 0)

	dialogObj, _ := builder.GetObject("AddContact")
	dialog := dialogObj.(*gtk.Dialog)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			defer dialog.Destroy()

			//TODO: validate contact
			contact, _ := contactInput.GetText()

			//TODO error
			iter, _ := accountInput.GetActiveIter()

			val, _ := model.GetValue(iter, 1)
			accountID, _ := val.GetString()

			//TODO error
			sendSubscription(accountID, contact)
		},
	})

	return dialog
}
