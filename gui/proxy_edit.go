package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/net"
)

func getScheme(s *gtk.ComboBoxText) string {
	return net.GetProxyTypeFor(s.GetActiveText())
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

	net.GetProxyTypeNames(func(name string) {
		scheme.AppendText(name)
	})
	scheme.SetActive(net.FindProxyTypeFor(prox.Scheme))

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

			prox.Scheme = net.GetProxyTypeFor(scheme.GetActiveText())

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
