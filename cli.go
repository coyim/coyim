// +build !nocli

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/proxy"
)

func main() {
	flag.Parse()

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err.Error())
	}
	defer terminal.Restore(0, oldState)
	term := terminal.NewTerminal(os.Stdin, "")
	updateTerminalSize(term)
	term.SetBracketedPasteMode(true)
	defer term.SetBracketedPasteMode(false)

	resizeChan := make(chan os.Signal)
	go func() {
		for _ = range resizeChan {
			updateTerminalSize(term)
		}
	}()
	signal.Notify(resizeChan, syscall.SIGWINCH)

	if len(*configFile) == 0 {
		if configFile, err = findConfigFile(os.Getenv("HOME")); err != nil {
			alert(term, err.Error())
			return
		}
	}

	config, err := ParseConfig(*configFile)
	if err != nil {
		alert(term, "Failed to parse config file: "+err.Error())
		config = new(Config)
		if !enroll(config, term) {
			return
		}
		config.filename = *configFile
		config.Save()
	}

	password := config.Password
	if len(password) == 0 {
		if password, err = term.ReadPassword(fmt.Sprintf("Password for %s (will not be saved to disk): ", config.Account)); err != nil {
			alert(term, "Failed to read password: "+err.Error())
			return
		}
	}
	term.SetPrompt("> ")

	parts := strings.SplitN(config.Account, "@", 2)
	if len(parts) != 2 {
		alert(term, "invalid username (want user@domain): "+config.Account)
		return
	}
	user := parts[0]
	domain := parts[1]

	var addr string
	addrTrusted := false

	if len(config.Server) > 0 && config.Port > 0 {
		addr = fmt.Sprintf("%s:%d", config.Server, config.Port)
		addrTrusted = true
	} else {
		if len(config.Proxies) > 0 {
			alert(term, "Cannot connect via a proxy without Server and Port being set in the config file as an SRV lookup would leak information.")
			return
		}
		host, port, err := xmpp.Resolve(domain)
		if err != nil {
			alert(term, "Failed to resolve XMPP server: "+err.Error())
			return
		}
		addr = fmt.Sprintf("%s:%d", host, port)
	}

	var dialer proxy.Dialer
	for i := len(config.Proxies) - 1; i >= 0; i-- {
		u, err := url.Parse(config.Proxies[i])
		if err != nil {
			alert(term, "Failed to parse "+config.Proxies[i]+" as a URL: "+err.Error())
			return
		}
		if dialer == nil {
			dialer = proxy.Direct
		}
		if dialer, err = proxy.FromURL(u, dialer); err != nil {
			alert(term, "Failed to parse "+config.Proxies[i]+" as a proxy: "+err.Error())
			return
		}
	}

	var certSHA256 []byte
	if len(config.ServerCertificateSHA256) > 0 {
		certSHA256, err = hex.DecodeString(config.ServerCertificateSHA256)
		if err != nil {
			alert(term, "Failed to parse ServerCertificateSHA256 (should be hex string): "+err.Error())
			return
		}
		if len(certSHA256) != 32 {
			alert(term, "ServerCertificateSHA256 is not 32 bytes long")
			return
		}
	}

	var createCallback xmpp.FormCallback
	if *createAccount {
		createCallback = func(title, instructions string, fields []interface{}) error {
			return promptForForm(term, user, password, title, instructions, fields)
		}
	}

	xmppConfig := &xmpp.Config{
		Log:                     &lineLogger{term, nil},
		CreateCallback:          createCallback,
		TrustedAddress:          addrTrusted,
		Archive:                 false,
		ServerCertificateSHA256: certSHA256,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS10,
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
		},
	}

	if domain == "jabber.ccc.de" {
		// jabber.ccc.de uses CACert but distros are removing that root
		// certificate.
		roots := x509.NewCertPool()
		caCertRoot, err := x509.ParseCertificate(caCertRootDER)
		if err == nil {
			alert(term, "Temporarily trusting only CACert root for CCC Jabber server")
			roots.AddCert(caCertRoot)
			xmppConfig.TLSConfig.RootCAs = roots
		} else {
			alert(term, "Tried to add CACert root for jabber.ccc.de but failed: "+err.Error())
		}
	}

	if len(config.RawLogFile) > 0 {
		rawLog, err := os.OpenFile(config.RawLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			alert(term, "Failed to open raw log file: "+err.Error())
			return
		}

		lock := new(sync.Mutex)
		in := rawLogger{
			out:    rawLog,
			prefix: []byte("<- "),
			lock:   lock,
		}
		out := rawLogger{
			out:    rawLog,
			prefix: []byte("-> "),
			lock:   lock,
		}
		in.other, out.other = &out, &in

		xmppConfig.InLog = &in
		xmppConfig.OutLog = &out

		defer in.flush()
		defer out.flush()
	}

	if dialer != nil {
		info(term, "Making connection to "+addr+" via proxy")
		if xmppConfig.Conn, err = dialer.Dial("tcp", addr); err != nil {
			alert(term, "Failed to connect via proxy: "+err.Error())
			return
		}
	}

	conn, err := xmpp.Dial(addr, user, domain, password, xmppConfig)
	if err != nil {
		alert(term, "Failed to connect to XMPP server: "+err.Error())
		return
	}

	//TODO support one session per account
	s := Session{
		account:           config.Account,
		conn:              conn,
		term:              term,
		conversations:     make(map[string]*otr3.Conversation),
		eh:                make(map[string]*eventHandler),
		knownStates:       make(map[string]string),
		privateKey:        new(otr3.PrivateKey),
		config:            config,
		pendingRosterChan: make(chan *rosterEdit),
		pendingSubscribes: make(map[string]string),
		lastActionTime:    time.Now(),
	}
	info(term, "Fetching roster")

	//var rosterReply chan xmpp.Stanza
	rosterReply, _, err := s.conn.RequestRoster()
	if err != nil {
		alert(term, "Failed to request roster: "+err.Error())
		return
	}

	conn.SignalPresence("")

	s.input = Input{
		term:        term,
		uidComplete: new(priorityList),
	}
	commandChan := make(chan interface{})
	go s.input.ProcessCommands(commandChan)

	stanzaChan := make(chan xmpp.Stanza)
	go s.readMessages(stanzaChan)

	s.privateKey.Parse(config.PrivateKey)
	s.timeouts = make(map[xmpp.Cookie]time.Time)

	info(term, fmt.Sprintf("Your fingerprint is %x", s.privateKey.DefaultFingerprint()))

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan bool)

	// commandLoop would not be necessary on a GUI client
	go commandLoop(term, &s, config, commandChan, quit)

	go stanzaLoop(term, &s, stanzaChan, quit)
	go rosterLoop(term, &s, rosterReply, quit)
	go timeoutLoop(&s, ticker.C)

	<-quit // wait
	os.Stdout.Write([]byte("\n"))
}
