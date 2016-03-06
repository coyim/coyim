package gui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/net"
	"github.com/twstrike/coyim/session/access"
)

type accountDetailsData struct {
	builder             gtki.Builder
	dialog              gtki.Dialog
	notebook            gtki.Notebook
	otherSettings       gtki.CheckButton
	acc                 gtki.Entry
	pass                gtki.Entry
	server              gtki.Entry
	port                gtki.Entry
	proxies             gtki.ListStore
	notificationArea    gtki.Box
	proxiesView         gtki.TreeView
	fingerprintsMessage gtki.Label
}

func getObjIgnoringErrors(b gtki.Builder, name string) glibi.Object {
	obj, _ := b.GetObject(name)
	return obj
}

func (d *accountDetailsData) getObjIgnoringErrors(name string) glibi.Object {
	return getObjIgnoringErrors(d.builder, name)
}

func getBuilderAndAccountDialogDetails() *accountDetailsData {
	data := &accountDetailsData{}

	dialogID := "AccountDetails"
	data.builder = builderForDefinition(dialogID)

	data.dialog = data.getObjIgnoringErrors(dialogID).(gtki.Dialog)
	data.notebook = data.getObjIgnoringErrors("notebook1").(gtki.Notebook)
	data.otherSettings = data.getObjIgnoringErrors("otherSettings").(gtki.CheckButton)
	data.acc = data.getObjIgnoringErrors("account").(gtki.Entry)
	data.pass = data.getObjIgnoringErrors("password").(gtki.Entry)
	data.server = data.getObjIgnoringErrors("server").(gtki.Entry)
	data.port = data.getObjIgnoringErrors("port").(gtki.Entry)
	data.proxies = data.getObjIgnoringErrors("proxies-model").(gtki.ListStore)
	data.notificationArea = data.getObjIgnoringErrors("notification-area").(gtki.Box)
	data.proxiesView = data.getObjIgnoringErrors("proxies-view").(gtki.TreeView)
	data.fingerprintsMessage = data.getObjIgnoringErrors("fingerprintsMessage").(gtki.Label)

	return data
}

func formattedFingerprintsFor(s access.Session) string {
	result := ""

	if s != nil {
		for _, sk := range s.PrivateKeys() {
			pk := sk.PublicKey()
			if pk != nil {
				result = fmt.Sprintf("%s%s%s\n", result, "    ", config.FormatFingerprint(pk.Fingerprint()))
			}
		}
	}

	return result
}

func (u *gtkUI) accountDialog(s access.Session, account *config.Account, saveFunction func()) {
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

	data.fingerprintsMessage.SetSelectable(true)
	m := i18n.Local("Your fingerprints for %s:\n%s")
	message := fmt.Sprintf(m, account.Account, formattedFingerprintsFor(s))
	data.fingerprintsMessage.SetText(message)

	p2, _ := data.notebook.GetNthPage(1)
	p3, _ := data.notebook.GetNthPage(2)
	p4, _ := data.notebook.GetNthPage(3)

	failures := 0

	editProxy := func(iter gtki.TreeIter, onCancel func()) {
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
				p4.Show()
			} else {
				p2.Hide()
				p3.Hide()
				p4.Hide()
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
			if _, iter, ok := ts.GetSelected(); ok {
				editProxy(iter, func() {})
			}
		},
		"on_remove_proxy_signal": func() {
			ts, _ := data.proxiesView.GetSelection()
			if _, iter, ok := ts.GetSelected(); ok {
				data.proxies.Remove(iter)
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
		"on_edit_activate_proxy_signal": func(_ gtki.TreeView, path gtki.TreePath) {
			iter, err := data.proxies.GetIter(path)
			if err == nil {
				editProxy(iter, func() {})
			}
		},
		"on_cancel_signal": func() {
			u.buildAccountsMenu()
			data.dialog.Destroy()
		},
		"on_import_key_signal": func() {
			u.importKeysForDialog(account, data.dialog)
		},
		"on_import_fpr_signal": func() {
			u.importFingerprintsForDialog(account, data.dialog)
		},
		"on_export_key_signal": func() {
			u.exportKeysForDialog(account, data.dialog)
		},
		"on_export_fpr_signal": func() {
			u.exportFingerprintsForDialog(account, data.dialog)
		},
	})

	data.dialog.SetTransientFor(u.window)
	data.dialog.ShowAll()
	data.notebook.SetCurrentPage(0)

	data.notebook.SetShowTabs(u.config.AdvancedOptions)
	if !u.config.AdvancedOptions {
		p2.Hide()
		p3.Hide()
		p4.Hide()
	}
}

func buildBadUsernameNotification(msg string) gtki.InfoBar {
	b := builderForDefinition("BadUsernameNotification")

	infoBar := getObjIgnoringErrors(b, "infobar").(gtki.InfoBar)
	message := getObjIgnoringErrors(b, "message").(gtki.Label)

	message.SetSelectable(true)
	message.SetText(fmt.Sprintf(i18n.Local(msg)))

	return infoBar
}
