package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
)

var proxyTypes = [][]string{
	[]string{"tor-auto", "Automatic Tor"},
	[]string{"socks4", "SOCKS4"},
	[]string{"socks5", "SOCKS5"},
}

func findProxyTypeFor(s string) int {
	for ix, px := range proxyTypes {
		if px[0] == s {
			return ix
		}
	}

	return -1
}

func getScheme(s *gtk.ComboBoxText) string {
	act := s.GetActiveText()
	for _, px := range proxyTypes {
		if act == i18n.Local(px[1]) {
			return px[0]
		}
	}
	return ""
}

func (u *gtkUI) editProxy(proxy string, onSave func(proxy)) {
	prox := parseProxy(proxy)

	b := builderForDefinition("EditProxy")
	dialog := getObjIgnoringErrors(b, "EditProxy").(*gtk.Dialog)
	scheme := getObjIgnoringErrors(b, "protocol-type").(*gtk.ComboBoxText)
	user := getObjIgnoringErrors(b, "user").(*gtk.Entry)
	pass := getObjIgnoringErrors(b, "password").(*gtk.Entry)
	server := getObjIgnoringErrors(b, "server").(*gtk.Entry)
	port := getObjIgnoringErrors(b, "port").(*gtk.Entry)

	for _, px := range proxyTypes {
		scheme.AppendText(i18n.Local(px[1]))
	}
	scheme.SetActive(findProxyTypeFor(prox.scheme))

	if prox.userSet {
		user.SetText(prox.user)
	}

	if prox.passSet {
		pass.SetText(prox.pass)
	}

	server.SetText(prox.host)

	if prox.portSet {
		port.SetText(prox.port)
	}

	b.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			userTxt, _ := user.GetText()
			passTxt, _ := pass.GetText()
			servTxt, _ := server.GetText()
			portTxt, _ := port.GetText()

			prox.scheme = getScheme(scheme)

			prox.userSet = userTxt != ""
			prox.user = userTxt

			prox.passSet = passTxt != ""
			prox.pass = passTxt

			prox.host = servTxt

			prox.portSet = portTxt != ""
			prox.port = portTxt

			go onSave(prox)
			dialog.Destroy()
		},
		"on_cancel_signal": func() {
			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}
