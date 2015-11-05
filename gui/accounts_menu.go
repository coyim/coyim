package gui

import (
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
)

var (
	// TODO: shouldn't this be specific to the account ID in question?
	accountChangedSignal, _ = glib.SignalNew("coyim-account-changed")
)

func firstProxy(account *account) string {
	if len(account.session.CurrentAccount.Proxies) > 0 {
		return account.session.CurrentAccount.Proxies[0]
	}
	return ""
}

func onAccountDialogClicked(account *config.Account, saveFunction func(), reg *widgetRegistry) func() {
	return func() {
		account.Account = reg.getText("account")
		account.Password = reg.getText("password")

		parts := strings.SplitN(account.Account, "@", 2)
		if len(parts) != 2 {
			log.Println("invalid username (want user@domain): " + account.Account)
			return
		}

		go saveFunction()
		reg.dialogDestroy("dialog")
	}
}

func accountDialog(account *config.Account, saveFunction func()) {
	//	onClicked := onAccountDialogClicked(account, saveFunction, reg)
	//
	//d := dialog{
	//	title:    i18n.Local("Account Details"),
	//	position: gtk.WIN_POS_CENTER,
	//	id:       "dialog",
	//	content: []creatable{
	//		label{text: i18n.Local("Your account (for example: kim42@dukgo.com)")},
	//		entry{
	//			text:       account.Account,
	//			editable:   true,
	//			visibility: true,
	//			id:         "account",
	//			onActivate: onClicked,
	//		},

	//		label{text: i18n.Local("Password\nAlert!! Your password is going to be stored as plaintext")},
	//		entry{
	//			text:       account.Password,
	//			editable:   true,
	//			visibility: false,
	//			id:         "password",
	//			onActivate: onClicked,
	//		},

	//		button{
	//			text:      i18n.Local("Save"),
	//			onClicked: onClicked,
	//		},
	//	},
	//}

	////d.create(reg)
	//reg.dialogShowAll("dialog")
	vars := make(map[string]string)
	vars["$title"] = i18n.Local("Account Details")
	vars["$accountMessage"] = i18n.Local("Your account (for example: kim42@dukgo.com)")
	//TODO: is this message still necessary?
	vars["$pswMessage"] = i18n.Local("Password\nAlert!! Your password is going to be stored as plaintext")
	vars["$saveLabel"] = i18n.Local("Save")
	builder, buildError := loadBuilderWith("AccountDetailsDefinition", vars)
	if buildError != nil {
		panic(buildError.Error())
	}
	// reg := createWidgetRegistry()
	obj, _ := builder.GetObject("AccountDetailsDialog")
	dialog := obj.(*gtk.Dialog)
	accObj, _ := builder.GetObject("account")
	accEntry := accObj.(*gtk.Entry)
	accEntry.SetText(account.Account)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passObj, _ := builder.GetObject("password")
			accTxt, _ := accEntry.GetText()
			passTxt, _ := passObj.(*gtk.Entry).GetText()
			account.Account = accTxt
			account.Password = passTxt

			parts := strings.SplitN(account.Account, "@", 2)
			if len(parts) != 2 {
				log.Println("invalid username (want user@domain): " + account.Account)
				return
			}

			fmt.Printf("\nsaving account %s", accTxt)
			go saveFunction()
			dialog.Destroy()
		},
	})
	dialog.ShowAll()
}

func toggleConnectAndDisconnectMenuItems(s *session.Session, connect, disconnect *gtk.MenuItem) {
	connected := s.ConnStatus == session.CONNECTED
	connect.SetSensitive(!connected)
	disconnect.SetSensitive(connected)
}

func (u *gtkUI) buildAccountsMenu() {
	submenu, _ := gtk.MenuNew()

	for _, account := range u.accounts {
		account.appendMenuTo(submenu)
	}

	if len(u.accounts) > 0 {
		sep, _ := gtk.SeparatorMenuItemNew()
		submenu.Append(sep)
	}

	addAccMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Add..."))
	addAccMenu.Connect("activate", func() { u.showAddAccountWindow() })

	submenu.Append(addAccMenu)
	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
	u.accountsMenu.ShowAll()
}
