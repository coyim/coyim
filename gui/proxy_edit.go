package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/net"
)

// this is an array because it must be sorted
var proxyTypes = []string{
	"tor-auto",
	"socks5",
}

var proxyTypesNames = map[string]i18n.T{
	"tor-auto": i18n.T("Automatic Tor"),
	"socks5":   i18n.T("SOCKS5"),
}

// findProxyTypeFor returns the index of the proxy type given
func findProxyTypeFor(s string) int {
	for ix, px := range proxyTypes {
		if px == s {
			return ix
		}
	}

	return -1
}

// getProxyTypeNames will yield all i18n proxy names to the function
func getProxyTypeNames(f func(string)) {
	for _, px := range proxyTypes {
		l := i18n.Local(string(proxyTypesNames[px]))
		f(l)
	}
}

// getProxyTypeFor will return the proxy type for the given i18n proxy name
// we are using a GtkComboBoxText "that hides the model-view complexity for simple text-only use cases"
// this function implements our (proxy type, proxy label) model.
func getProxyTypeFor(act string) string {
	for _, px := range proxyTypes {
		l := i18n.Local(string(proxyTypesNames[px]))
		if act == l {
			return px
		}
	}
	return ""
}

func getScheme(s *gtk.ComboBoxText) string {
	return getProxyTypeFor(s.GetActiveText())
}

func orNil(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}

func (u *gtkUI) editProxy(proxy string, w *gtk.Dialog, onSave func(net.Proxy), onCancel func()) {
	prox := net.ParseProxy(proxy)

	b := builderForDefinition("EditProxy")
	dialog := getObjIgnoringErrors(b, "EditProxy").(*gtk.Dialog)
	scheme := getObjIgnoringErrors(b, "protocol-type").(*gtk.ComboBoxText)
	user := getObjIgnoringErrors(b, "user").(*gtk.Entry)
	pass := getObjIgnoringErrors(b, "password").(*gtk.Entry)
	server := getObjIgnoringErrors(b, "server").(*gtk.Entry)
	port := getObjIgnoringErrors(b, "port").(*gtk.Entry)

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
