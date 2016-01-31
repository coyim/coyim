package gui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/net"
)

type accountDetailsData struct {
	builder          *gtk.Builder
	dialog           *gtk.Dialog
	notebook         *gtk.Notebook
	otherSettings    *gtk.CheckButton
	acc              *gtk.Entry
	pass             *gtk.Entry
	server           *gtk.Entry
	port             *gtk.Entry
	proxies          *gtk.ListStore
	notificationArea *gtk.Box
	proxiesView      *gtk.TreeView
}

func getObjIgnoringErrors(b *gtk.Builder, name string) glib.IObject {
	obj, _ := b.GetObject(name)
	return obj
}

func (d *accountDetailsData) getObjIgnoringErrors(name string) glib.IObject {
	return getObjIgnoringErrors(d.builder, name)
}

func getBuilderAndAccountDialogDetails() *accountDetailsData {
	data := &accountDetailsData{}

	dialogID := "AccountDetails"
	data.builder = builderForDefinition(dialogID)

	data.dialog = data.getObjIgnoringErrors(dialogID).(*gtk.Dialog)
	data.notebook = data.getObjIgnoringErrors("notebook1").(*gtk.Notebook)
	data.otherSettings = data.getObjIgnoringErrors("otherSettings").(*gtk.CheckButton)
	data.acc = data.getObjIgnoringErrors("account").(*gtk.Entry)
	data.pass = data.getObjIgnoringErrors("password").(*gtk.Entry)
	data.server = data.getObjIgnoringErrors("server").(*gtk.Entry)
	data.port = data.getObjIgnoringErrors("port").(*gtk.Entry)
	data.proxies = data.getObjIgnoringErrors("proxies-model").(*gtk.ListStore)
	data.notificationArea = data.getObjIgnoringErrors("notification-area").(*gtk.Box)
	data.proxiesView = data.getObjIgnoringErrors("proxies-view").(*gtk.TreeView)

	return data
}

func (u *gtkUI) accountDialog(account *config.Account, saveFunction func()) {
	data := getBuilderAndAccountDialogDetails()

	data.otherSettings.SetActive(u.config.AdvancedOptions)
	data.acc.SetText(account.Account)
	data.server.SetText(account.Server)
	if account.Port == 0 {
		account.Port = 5222
	}
	data.port.SetText(strconv.Itoa(account.Port))

	for _, px := range account.Proxies {
		iter := data.proxies.Append()
		data.proxies.SetValue(iter, 0, net.ParseProxy(px).ForPresentation())
		data.proxies.SetValue(iter, 1, px)
	}

	p2, _ := data.notebook.GetNthPage(1)
	p3, _ := data.notebook.GetNthPage(2)

	failures := 0

	editProxy := func(iter *gtk.TreeIter, onCancel func()) {
		val, _ := data.proxies.GetValue(iter, 1)
		realProxyData, _ := val.GetString()
		u.editProxy(realProxyData, data.dialog,
			func(p net.Proxy) {
				data.proxies.SetValue(iter, 0, p.ForPresentation())
				data.proxies.SetValue(iter, 1, p.ForProcessing())
			}, onCancel)
	}

	data.builder.ConnectSignals(map[string]interface{}{
		"on_toggle_other_settings": func() {
			otherSettings := data.otherSettings.GetActive()
			u.setShowAdvancedSettings(otherSettings)
			data.notebook.SetShowTabs(otherSettings)
			if otherSettings {
				p2.Show()
				p3.Show()
			} else {
				p2.Hide()
				p3.Hide()
			}
		},
		"on_save_signal": func() {
			accTxt, _ := data.acc.GetText()
			passTxt, _ := data.pass.GetText()
			servTxt, _ := data.server.GetText()
			portTxt, _ := data.port.GetText()

			isJid, err := verifyXmppAddress(accTxt)
			if !isJid && failures > 0 {
				failures++
				return
			}

			if !isJid {
				notification := buildBadUsernameNotification(err)
				data.notificationArea.Add(notification)
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

			newProxies := []string{}
			iter, ok := data.proxies.GetIterFirst()
			for ok {
				vv, _ := data.proxies.GetValue(iter, 1)
				newProxy, _ := vv.GetString()
				newProxies = append(newProxies, newProxy)
				ok = data.proxies.IterNext(iter)
			}

			account.Proxies = newProxies

			go saveFunction()
			data.dialog.Destroy()
		},
		"on_edit_proxy_signal": func() {
			ts, _ := data.proxiesView.GetSelection()
			var iter gtk.TreeIter
			if ts.GetSelected(nil, &iter) {
				editProxy(&iter, func() {})
			}
		},
		"on_remove_proxy_signal": func() {
			ts, _ := data.proxiesView.GetSelection()
			var iter gtk.TreeIter
			if ts.GetSelected(nil, &iter) {
				data.proxies.Remove(&iter)
			}
		},
		"on_add_proxy_signal": func() {
			iter := data.proxies.Append()
			data.proxies.SetValue(iter, 0, "tor-auto://")
			data.proxies.SetValue(iter, 1, "tor-auto://")
			ts, _ := data.proxiesView.GetSelection()
			ts.SelectIter(iter)
			editProxy(iter, func() {
				data.proxies.Remove(iter)
			})
		},
		"on_edit_activate_proxy_signal": func(_ *gtk.TreeView, path *gtk.TreePath) {
			iter, err := data.proxies.GetIter(path)
			if err == nil {
				editProxy(iter, func() {})
			}
		},
		"on_cancel_signal": func() {
			u.buildAccountsMenu()
			data.dialog.Destroy()
		},
	})

	data.dialog.SetTransientFor(u.window)
	data.dialog.ShowAll()
	data.notebook.SetCurrentPage(0)

	data.notebook.SetShowTabs(u.config.AdvancedOptions)
	if !u.config.AdvancedOptions {
		p2.Hide()
		p3.Hide()
	}
}

func buildBadUsernameNotification(msg string) *gtk.InfoBar {
	b := builderForDefinition("BadUsernameNotification")

	infoBar := getObjIgnoringErrors(b, "infobar").(*gtk.InfoBar)
	message := getObjIgnoringErrors(b, "message").(*gtk.Label)

	message.SetSelectable(true)
	message.SetText(fmt.Sprintf(i18n.Local(msg)))

	return infoBar
}
