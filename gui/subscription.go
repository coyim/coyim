package gui

import (
	"github.com/gotk3/gotk3/gdk"
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

func presenceSubscriptionDialog(accounts []*account) *gtk.Dialog {
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
		model.Set(iter, []int{0, 1}, []interface{}{acc.session.CurrentAccount.Account, acc})
	}

	accountInput, _ := gtk.ComboBoxNewWithModel(&model.TreeModel)
	if len(accounts) > 0 {
		accountInput.SetActive(0)
	}

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

	onAdd := func() {
		//TODO: validate contact
		contact, _ := contactInput.GetText()

		//TODO error
		iter, _ := accountInput.GetActiveIter()

		val, _ := model.GetValue(iter, 1)
		account := (*account)(val.GetPointer())

		if !account.connected() {
			//TODO error
		}

		//TODO: validate
		// - validate if the account is connected
		err := account.session.Conn.SendPresence(contact, "subscribe", "" /* generate id */)
		if err != nil {
			//TODO: error
		}

		dialog.Destroy()
	}

	button.Connect("clicked", onAdd)
	contactInput.Connect("key-press-event", func(_ *gtk.Entry, ev *gdk.Event) bool {
		evKey := gdk.EventKey{ev}
		if (evKey.State()&gdk.GDK_MODIFIER_MASK) == 0 && evKey.KeyVal() == 0xff0d {
			onAdd()
			return true
		}
		return false
	})

	contactInput.GrabFocus()

	return dialog
}
