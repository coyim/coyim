package gui

import (
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/net"
	"github.com/twstrike/gotk3adapter/gtki"
)

var proxyTypes = [][]string{
	[]string{"tor-auto", "Automatic Tor"},
	[]string{"socks5", "SOCKS5"},
}

// findProxyTypeFor returns the index of the proxy type given
func findProxyTypeFor(s string) int {
	for ix, px := range proxyTypes {
		if px[0] == s {
			return ix
		}
	}

	return -1
}

// getProxyTypeNames will yield all i18n proxy names to the function
func getProxyTypeNames(f func(string)) {
	for _, px := range proxyTypes {
		f(i18n.Local(px[1]))
	}
}

// getProxyTypeFor will return the proxy type for the given i18n proxy name
func getProxyTypeFor(act string) string {
	for _, px := range proxyTypes {
		if act == i18n.Local(px[1]) {
			return px[0]
		}
	}
	return ""
}

func getScheme(s gtki.ComboBoxText) string {
	return getProxyTypeFor(s.GetActiveText())
}

func orNil(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}

func (u *gtkUI) editProxy(proxy string, w gtki.Dialog, onSave func(net.Proxy), onCancel func()) {
	prox := net.ParseProxy(proxy)

	b := builderForDefinition("EditProxy")
	dialog := getObjIgnoringErrors(b, "EditProxy").(gtki.Dialog)
	scheme := getObjIgnoringErrors(b, "protocol-type").(gtki.ComboBoxText)
	user := getObjIgnoringErrors(b, "user").(gtki.Entry)
	pass := getObjIgnoringErrors(b, "password").(gtki.Entry)
	server := getObjIgnoringErrors(b, "server").(gtki.Entry)
	port := getObjIgnoringErrors(b, "port").(gtki.Entry)

	getProxyTypeNames(func(name string) {
		scheme.AppendText(name)
	})
	scheme.SetActive(findProxyTypeFor(prox.Scheme))

	if prox.User != nil {
		user.SetText(*prox.User)
	}

	if prox.Pass != nil {
		pass.SetText(*prox.Pass)
	}

	server.SetText(prox.Host)

	if prox.Port != nil {
		port.SetText(*prox.Port)
	}

	b.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			userTxt, _ := user.GetText()
			passTxt, _ := pass.GetText()
			servTxt, _ := server.GetText()
			portTxt, _ := port.GetText()

			prox.Scheme = getProxyTypeFor(scheme.GetActiveText())

			prox.User = orNil(userTxt)
			prox.Pass = orNil(passTxt)
			prox.Host = servTxt
			prox.Port = orNil(portTxt)

			go onSave(prox)
			dialog.Destroy()
		},
		"on_cancel_signal": func() {
			go onCancel()
			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(w)
	dialog.ShowAll()
}
