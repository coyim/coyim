package gui

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/digests"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/net"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/gotk3adapter/gtki"
)

type accountDetailsData struct {
	builder             *builder
	dialog              gtki.Dialog       `gtk-widget:"AccountDetails"`
	notebook            gtki.Notebook     `gtk-widget:"notebook1"`
	otherSettings       gtki.CheckButton  `gtk-widget:"otherSettings"`
	acc                 gtki.Entry        `gtk-widget:"account"`
	pass                gtki.Entry        `gtk-widget:"password"`
	displayName         gtki.Entry        `gtk-widget:"displayName"`
	server              gtki.Entry        `gtk-widget:"server"`
	port                gtki.Entry        `gtk-widget:"port"`
	proxies             gtki.ListStore    `gtk-widget:"proxies-model"`
	pins                gtki.ListStore    `gtk-widget:"pins-model"`
	notificationArea    gtki.Box          `gtk-widget:"notification-area"`
	proxiesView         gtki.TreeView     `gtk-widget:"proxies-view"`
	fingerprintsMessage gtki.Label        `gtk-widget:"fingerprintsMessage"`
	pinningPolicy       gtki.ComboBoxText `gtk-widget:"pinningPolicyValue"`
	pinsView            gtki.TreeView     `gtk-widget:"pins-view"`
}

func getBuilderAndAccountDialogDetails() *accountDetailsData {
	data := &accountDetailsData{}
	data.builder = newBuilder("AccountDetails")

	//data.proxies needs to be kept beyond the lifespan of the builder.
	//Because this also holds a reference to the builder, we should be fine.

	panicOnDevError(data.builder.bindObjects(data))

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

func findPinningPolicyFor(t string) int {
	switch t {
	case "none":
		return 0
	case "deny":
		return 1
	case "add":
		return 2
	case "", "add-first-ask-rest":
		return 3
	case "add-first-deny-rest":
		return 4
	case "ask":
		return 5
	}
	return -1
}

func filterCertificates(oldCerts []*config.CertificatePin, newList gtki.ListStore) []*config.CertificatePin {
	allPins := make(map[string]bool)

	iter, ok := newList.GetIterFirst()
	for ok {
		vv, _ := newList.GetValue(iter, 2)
		pp, _ := vv.GetString()
		allPins[pp] = true
		ok = newList.IterNext(iter)
	}

	newCerts := []*config.CertificatePin{}

	for _, cc := range oldCerts {
		if allPins[hex.EncodeToString(cc.Fingerprint)] {
			newCerts = append(newCerts, cc)
		}
	}

	return newCerts
}

func (u *gtkUI) connectionInfoDialog(account *account) {
	assertInUIThread()

	dialogID := "ConnectionInformation"
	builder := newBuilder(dialogID)

	var dialog gtki.Dialog
	var server, tlsAlgo, tlsVersion, tlsFingerprint gtki.Label
	var pinCertButton gtki.Button

	builder.getItems(
		dialogID, &dialog,
		"serverValue", &server,
		"tlsAlgoValue", &tlsAlgo,
		"tlsVersionValue", &tlsVersion,
		"tlsFingerprintValue", &tlsFingerprint,
		"pin-cert", &pinCertButton,
	)

	if account.session == nil || account.session.Conn() == nil || account.session.Conn().RawOut() == nil {
		// If these are nil, we bug out early, since there's no connection to give information about yet
		return
	}

	tlsConn := account.session.Conn().RawOut().(*tls.Conn)

	serverAddress := account.session.Conn().ServerAddress()
	if serverAddress == "" {
		parts := strings.SplitN(account.Account(), "@", 2)
		serverAddress = parts[1]
	}
	server.SetLabel(serverAddress)

	tlsAlgo.SetLabel(xmpp.GetCipherSuiteName(tlsConn.ConnectionState()))
	tlsVersion.SetLabel(xmpp.GetTLSVersion(tlsConn.ConnectionState()))

	certs := tlsConn.ConnectionState().PeerCertificates
	chunks := splitStringEvery(fmt.Sprintf("%X", digests.Sha3_256(certs[0].Raw)), chunkingDefaultGrouping)
	tlsFingerprint.SetLabel(fmt.Sprintf("%s %s %s %s\n%s %s %s %s", chunks[0], chunks[1], chunks[2], chunks[3], chunks[4], chunks[5], chunks[6], chunks[7]))

	if checkPinned(account.session.GetConfig(), certs) {
		pinCertButton.SetSensitive(false)
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_close_signal": dialog.Destroy,
		"on_pin_signal": func() {
			account.session.GetConfig().SaveCert(certs[0].Subject.CommonName, certs[0].Issuer.CommonName, digests.Sha3_256(certs[0].Raw))
			u.SaveConfig()
			pinCertButton.SetSensitive(false)
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

type accountDetails struct {
	accTxt  string
	passTxt string
	dispTxt string
	servTxt string
	portTxt string
}

func getAccountDetails(data *accountDetailsData) *accountDetails {
	accTxt, _ := data.acc.GetText()
	passTxt, _ := data.pass.GetText()
	dispTxt, _ := data.displayName.GetText()
	servTxt, _ := data.server.GetText()
	portTxt, _ := data.port.GetText()

	details := &accountDetails{
		accTxt,
		passTxt,
		dispTxt,
		servTxt,
		portTxt,
	}

	return details
}

func renderSessionDetails(s access.Session, data *accountDetailsData, account string) {
	_ = data.displayName.SetProperty("placeholder-text", s.DisplayName())
	nick := s.GetConfig().Nickname
	if nick != "" {
		data.displayName.SetText(nick)
	}

	if s.PrivateKeys() != nil && len(s.PrivateKeys()) > 0 {
		data.fingerprintsMessage.SetSelectable(true)
		m := i18n.Local("The fingerprints for %s:\n%s")
		message := fmt.Sprintf(m, account, formattedFingerprintsFor(s))
		data.fingerprintsMessage.SetText(message)
	}
}

func renderAccountDetails(account *config.Account, data *accountDetailsData) {
	if account.Password != "" {
		_ = data.pass.SetProperty("placeholder-text", "(saved in configuration file)")
	}

	data.server.SetText(account.Server)
	if account.Port == 0 {
		account.Port = 5222
	}
	data.port.SetText(strconv.Itoa(account.Port))

	for _, px := range account.Proxies {
		iter := data.proxies.Append()
		_ = data.proxies.SetValue(iter, 0, net.ParseProxy(px).ForPresentation())
		_ = data.proxies.SetValue(iter, 1, px)
	}

	for _, px := range account.Certificates {
		iter := data.pins.Append()
		_ = data.pins.SetValue(iter, 0, px.Subject)
		_ = data.pins.SetValue(iter, 1, px.Issuer)
		_ = data.pins.SetValue(iter, 2, hex.EncodeToString(px.Fingerprint))
	}

	data.pinningPolicy.SetActive(findPinningPolicyFor(account.PinningPolicy))
}

func addAccount(account *config.Account, accDtails *accountDetails, data *accountDetailsData) {
	account.Account = accDtails.accTxt
	account.Server = accDtails.servTxt

	if accDtails.passTxt != "" {
		account.Password = accDtails.passTxt
	}

	account.Nickname = accDtails.dispTxt

	convertedPort, e := strconv.Atoi(accDtails.portTxt)
	if len(strings.TrimSpace(accDtails.portTxt)) == 0 || e != nil {
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

	account.Certificates = filterCertificates(account.Certificates, data.pins)
	account.PinningPolicy = data.pinningPolicy.GetActiveID()
}

func (u *gtkUI) accountDialog(s access.Session, account *config.Account, saveFunction func()) {
	assertInUIThread()

	data := getBuilderAndAccountDialogDetails()

	data.otherSettings.SetActive(u.config.AdvancedOptions)
	data.acc.SetText(account.Account)

	if s != nil {
		renderSessionDetails(s, data, account.Account)
	}

	renderAccountDetails(account, data)

	p2, _ := data.notebook.GetNthPage(1)
	p3, _ := data.notebook.GetNthPage(2)
	p4, _ := data.notebook.GetNthPage(3)

	editProxy := func(iter gtki.TreeIter, onCancel func()) {
		val, _ := data.proxies.GetValue(iter, 1)
		realProxyData, _ := val.GetString()
		u.editProxy(realProxyData, data.dialog,
			func(p net.Proxy) {
				_ = data.proxies.SetValue(iter, 0, p.ForPresentation())
				_ = data.proxies.SetValue(iter, 1, p.ForProcessing())
			}, onCancel)
	}

	errorNotif := newErrorNotification(data.notificationArea)

	data.builder.ConnectSignals(map[string]interface{}{
		"on_toggle_other_settings": func() {
			otherSettings := data.otherSettings.GetActive()
			u.setShowAdvancedSettings(otherSettings)
			data.notebook.SetShowTabs(otherSettings)
			if !otherSettings {
				p2.Hide()
				p3.Hide()
				p4.Hide()
			}

			p2.Show()
			p3.Show()
			p4.Show()
		},
		"on_save_signal": func() {
			accDtails := getAccountDetails(data)
			if isJid, err := verifyXMPPAddress(accDtails.accTxt); !isJid {
				errorNotif.ShowMessage(err)
				log.WithFields(log.Fields{
					"account": account.Account,
					"error":   err,
				}).Warn("Something went wrong")
				return
			}

			addAccount(account, accDtails, data)
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
		"on_remove_pin_signal": func() {
			ts, _ := data.pinsView.GetSelection()
			if _, iter, ok := ts.GetSelected(); ok {
				data.pins.Remove(iter)
			}
		},
		"on_add_proxy_signal": func() {
			iter := data.proxies.Append()
			_ = data.proxies.SetValue(iter, 0, "tor-auto://")
			_ = data.proxies.SetValue(iter, 1, "tor-auto://")
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

	data.notebook.SetCurrentPage(0)
	data.notebook.SetShowTabs(u.config.AdvancedOptions)
	if !u.config.AdvancedOptions {
		p2.Hide()
		p3.Hide()
		p4.Hide()
	}

	data.dialog.SetTransientFor(u.window)
	data.dialog.ShowAll()
}
