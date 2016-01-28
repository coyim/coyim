package gui

import (
	"fmt"
	"log"
	"net"
	"net/url"
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
	if len(account.session.GetConfig().Proxies) > 0 {
		return account.session.GetConfig().Proxies[0]
	}

	return ""
}

func checkIsLikelyRandom(protocol, username, password string) bool {
	return protocol == "socks5" && ((strings.HasPrefix(username, "randomTor-") &&
		strings.HasPrefix(password, "randomTor-")) ||
		(len(username) == 10 &&
			len(password) == 10))
}

func parseProxy(s string) (protocol, host, port, username, password string, ok bool) {
	p, err := url.Parse(s)
	if err != nil {
		return "", "", "", "", "", false
	}

	host, port, err = net.SplitHostPort(p.Host)
	if err != nil {
		nerr, ok := err.(*net.AddrError)
		if !ok || nerr.Err != "missing port in address" {
			return "", "", "", "", "", false
		}

		port = ""
	}

	ok = true
	protocol = p.Scheme
	username = p.User.Username()
	password, _ = p.User.Password()

	return
}

func (u *gtkUI) accountDialog(account *config.Account, saveFunction func()) {
	dialogID := "AccountDetails"
	builder := builderForDefinition(dialogID)

	obj, _ := builder.GetObject(dialogID)
	dialog := obj.(*gtk.Dialog)

	obj, _ = builder.GetObject("notebook1")
	notebook := obj.(*gtk.Notebook)

	obj, _ = builder.GetObject("otherSettings")
	otherSettingsToggle := obj.(*gtk.CheckButton)
	otherSettingsToggle.SetActive(u.config.AdvancedOptions)

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

	obj, _ = builder.GetObject("proxies-model")
	proxiesModel := obj.(*gtk.ListStore)

	for _, px := range account.Proxies {
		p, _ := url.Parse(px)
		us := ""
		ps := ""
		compose := ""
		if p.User != nil {
			us = p.User.Username()
			compose = "@"
			_, passSet := p.User.Password()
			if passSet {
				ps = ":*****"
			}
		}

		iter := proxiesModel.Append()
		proxiesModel.SetValue(iter, 0,
			fmt.Sprintf("%s://%s%s%s%s", p.Scheme, us, ps, compose, p.Host))
		proxiesModel.SetValue(iter, 1, px)
	}

	obj, _ = builder.GetObject("notification-area")
	notificationArea := obj.(*gtk.Box)

	p2, _ := notebook.GetNthPage(1)
	p3, _ := notebook.GetNthPage(2)

	failures := 0

	builder.ConnectSignals(map[string]interface{}{
		"on_toggle_other_settings": func() {
			otherSettings := otherSettingsToggle.GetActive()
			u.setShowAdvancedSettings(otherSettings)
			if otherSettings {
				p2.Show()
				p3.Show()
			} else {
				p2.Hide()
				p3.Hide()
			}
		},
		"on_save_signal": func() {
			accTxt, _ := accEntry.GetText()
			passTxt, _ := passEntry.GetText()
			servTxt, _ := serverEntry.GetText()
			portTxt, _ := portEntry.GetText()

			isJid, err := verifyXmppAddress(accTxt)
			if !isJid && failures > 0 {
				failures++
				return
			}

			if !isJid {
				notification := buildBadUsernameNotification(err)
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
	notebook.SetCurrentPage(0)

	if !u.config.AdvancedOptions {
		p2.Hide()
		p3.Hide()
	}
}

func buildBadUsernameNotification(msg string) *gtk.InfoBar {
	builder := builderForDefinition("BadUsernameNotification")

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	message := obj.(*gtk.Label)
	message.SetSelectable(true)

	text := fmt.Sprintf(i18n.Local(msg))
	message.SetText(text)

	return infoBar
}

func toggleConnectAndDisconnectMenuItems(s *session.Session, connect, disconnect *gtk.MenuItem) {
	doInUIThread(func() {
		connect.SetSensitive(s.IsDisconnected())
		disconnect.SetSensitive(s.IsConnected())
	})
}

var accountsLock sync.Mutex

func (u *gtkUI) buildStaticAccountsMenu(submenu *gtk.Menu) {
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
	addAccMenu.Connect("activate", u.showAddAccountWindow)
	submenu.Append(addAccMenu)

	importMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Import..."))
	importMenu.Connect("activate", u.runImporter)
	submenu.Append(importMenu)

	registerAccMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Register..."))
	registerAccMenu.Connect("activate", u.showServerSelectionWindow)
	submenu.Append(registerAccMenu)
}

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

	u.buildStaticAccountsMenu(submenu)

	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
}
