package gui

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

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

func (u *gtkUI) accountDialog(account *config.Account, saveFunction func()) {
	dialogID := "AccountDetails"
	builder := builderForDefinition(dialogID)

	obj, _ := builder.GetObject(dialogID)
	dialog := obj.(*gtk.Dialog)

	obj, _ = builder.GetObject("account")
	accEntry := obj.(*gtk.Entry)
	accEntry.SetText(account.Account)

	obj, _ = builder.GetObject("password")
	passEntry := obj.(*gtk.Entry)

	obj, _ = builder.GetObject("server")
	serverEntry := obj.(*gtk.Entry)
	serverEntry.SetText(account.Server)

	obj, _ = builder.GetObject("port")
	portEntry := obj.(*gtk.Entry)
	if account.Port == 0 {
		account.Port = 5222
	}
	portEntry.SetText(strconv.Itoa(account.Port))

	obj, _ = builder.GetObject("notification-area")
	notificationArea := obj.(*gtk.Box)

	failures := 0

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			accTxt, _ := accEntry.GetText()
			passTxt, _ := passEntry.GetText()
			servTxt, _ := serverEntry.GetText()
			portTxt, _ := portEntry.GetText()

			isEmail, err := isEmail(accTxt)
			if !isEmail && failures > 0 {
				failures++
				log.Printf("authentication has failed %d times", failures)
				return
			}

			if "" == accTxt || !isEmail {
				notification := buildBadUsernameNotification()
				notificationArea.Add(notification)
				notification.ShowAll()
				failures++
				log.Printf(err)
				return
			}

			account.Account = accTxt
			account.Server = servTxt

			if passTxt != "" {
				account.Password = passTxt
			}

			convertedPort, e := strconv.Atoi(portTxt)
			if len(strings.TrimSpace(portTxt)) == 0 || e != nil {
				convertedPort = 5222
			}

			account.Port = convertedPort

			go saveFunction()
			dialog.Destroy()
		},
		"on_cancel_signal": func() {
			u.buildAccountsMenu()
			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func isEmail(address string) (bool, string) {
	matcher := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	matches := matcher.MatchString(address)
	var err string
	if !matches {
		err = fmt.Sprintf("<%s> is not a valid username", address)
	}
	return matches, err
}

func buildBadUsernameNotification() *gtk.InfoBar {
	builder := builderForDefinition("BadUsernameNotification")

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	message := obj.(*gtk.Label)

	text := fmt.Sprintf(i18n.Local("Username is required and should look like an email address"))
	message.SetText(text)

	return infoBar
}

func toggleConnectAndDisconnectMenuItems(s *session.Session, connect, disconnect *gtk.MenuItem) {
	glib.IdleAdd(func() {
		connect.SetSensitive(s.ConnStatus == session.DISCONNECTED)
		disconnect.SetSensitive(s.ConnStatus == session.CONNECTED)
	})
}

var accountsLock sync.Mutex

func (u *gtkUI) buildAccountsMenu() {
	accountsLock.Lock()
	defer accountsLock.Unlock()

	submenu, _ := gtk.MenuNew()

	for _, account := range u.accounts {
		account.appendMenuTo(submenu)
	}

	if len(u.accounts) > 0 {
		sep, _ := gtk.SeparatorMenuItemNew()
		submenu.Append(sep)
	}

	connectAutomaticallyItem, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Connect On _Startup"))
	u.config.WhenLoaded(func(a *config.ApplicationConfig) {
		connectAutomaticallyItem.SetActive(a.ConnectAutomatically)
	})
	connectAutomaticallyItem.Connect("activate", func() {
		u.toggleConnectAllAutomatically()
	})
	submenu.Append(connectAutomaticallyItem)

	connectAllMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect All"))
	connectAllMenu.Connect("activate", func() { u.connectAllAutomatics(true) })
	submenu.Append(connectAllMenu)

	sep2, _ := gtk.SeparatorMenuItemNew()
	submenu.Append(sep2)

	addAccMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Add..."))
	addAccMenu.Connect("activate", func() { u.showAddAccountWindow() })
	submenu.Append(addAccMenu)

	importMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Import..."))
	importMenu.Connect("activate", func() { u.runImporter() })
	submenu.Append(importMenu)

	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
}
