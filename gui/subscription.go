package gui

import (
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3/glib"
	"github.com/twstrike/gotk3/gtk"
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

func presenceSubscriptionDialog(accounts []*Account) *gtk.Dialog {
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle(i18n.Local("Add contact"))
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox, _ := dialog.GetContentArea()

	accountLabel, _ := gtk.LabelNew(i18n.Local("Account"))
	vbox.Add(accountLabel)

	model, _ := gtk.ListStoreNew(
		glib.TYPE_STRING,  // account name
		glib.TYPE_POINTER, // *Account
	)

	for _, acc := range accounts {
		iter := model.Append()
		//TODO stop passing pointers
		model.Set(iter, []int{0, 1}, []interface{}{acc.Account, acc})
	}

	accountInput, _ := gtk.ComboBoxNewWithModel(&model.TreeModel)
	vbox.Add(accountInput)

	renderer, _ := gtk.CellRendererTextNew()
	accountInput.PackStart(renderer, true)
	accountInput.AddAttribute(renderer, "text", 0)

	l, _ := gtk.LabelNew(i18n.Local("ID"))
	vbox.Add(l)

	contactInput, _ := gtk.EntryNew()
	contactInput.SetEditable(true)
	vbox.Add(contactInput)

	//TODO: disable the add button until the form has all the data
	//- an account selected
	//- an ID
	button, _ := gtk.ButtonNewWithLabel(i18n.Local("Add"))
	vbox.Add(button)

	button.Connect("clicked", func() {
		//TODO: validate contact
		contact, _ := contactInput.GetText()

		//TODO error
		iter, _ := accountInput.GetActiveIter()

		val, _ := model.GetValue(iter, 1)
		account := (*Account)(val.GetPointer())

		if !account.Connected() {
			//TODO error
		}

		//TODO: validate
		// - validate if the account is connected
		err := account.Conn.SendPresence(contact, "subscribe", "" /* generate id */)
		if err != nil {
			//TODO: error
		}

		dialog.Destroy()
	})

	return dialog
}
