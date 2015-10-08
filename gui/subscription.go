package gui

import (
	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
)

func presenceSubscriptionDialog(accounts []*Account) *gtk.Dialog {
	dialog := gtk.NewDialog()
	dialog.SetTitle("Add contact")
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox := dialog.GetVBox()

	accountLabel := gtk.NewLabel("Account")
	vbox.Add(accountLabel)

	model := gtk.NewListStore(
		gtk.TYPE_STRING,  // account name
		gtk.TYPE_POINTER, // *Account
	)

	iter := &gtk.TreeIter{}
	for _, acc := range accounts {
		model.Append(iter)
		model.Set(iter,
			0, acc.Account,
			1, acc,
		)
	}

	accountInput := gtk.NewComboBoxWithModel(&model.TreeModel)
	vbox.Add(accountInput)

	//TODO: ComboBox should have a CellLayout embedded
	cellLayout := accountInput.GetCellLayout()
	renderer := gtk.NewCellRendererText()
	cellLayout.PackStart(renderer, true)
	cellLayout.AddAttribute(renderer, "text", 0)

	vbox.Add(gtk.NewLabel("ID"))
	contactInput := gtk.NewEntry()
	contactInput.SetEditable(true)
	vbox.Add(contactInput)

	//TODO: disable the add button until the form has all the data
	//- an account selected
	//- an ID
	button := gtk.NewButtonWithLabel("Add")
	vbox.Add(button)

	button.Connect("clicked", func() {
		//TODO: validate contact
		contact := contactInput.GetText()

		iter := &gtk.TreeIter{}
		if !accountInput.GetActiveIter(iter) {
			//TODO error
		}

		val := &glib.GValue{}
		model.GetValue(iter, 1, val)
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
