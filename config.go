package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/twstrike/coyim/xmpp"
	otr "github.com/twstrike/otr3"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/proxy"
)

var (
	errHomeDirNotSet = errors.New("$HOME not set. Please either export $HOME or use the -config-file option.\n")
)

type MultiAccountConfig struct {
	filename string `json:"-"`
	Accounts []Config
}

type Config struct {
	filename                      string `json:"-"`
	Account                       string
	Server                        string   `json:",omitempty"`
	Proxies                       []string `json:",omitempty"`
	Password                      string   `json:",omitempty"`
	Port                          int      `json:",omitempty"`
	PrivateKey                    []byte
	KnownFingerprints             []KnownFingerprint
	RawLogFile                    string   `json:",omitempty"`
	NotifyCommand                 []string `json:",omitempty"`
	IdleSecondsBeforeNotification int      `json:",omitempty"`
	Bell                          bool
	HideStatusUpdates             bool
	UseTor                        bool
	OTRAutoTearDown               bool
	OTRAutoAppendTag              bool
	OTRAutoStartSession           bool
	ServerCertificateSHA256       string   `json:",omitempty"`
	AlwaysEncrypt                 bool     `json:",omitempty"`
	AlwaysEncryptWith             []string `json:",omitempty"`
}

type KnownFingerprint struct {
	UserId         string
	FingerprintHex string
	fingerprint    []byte `json:"-"`
}

func findConfigFile(homeDir string) (*string, error) {
	if len(homeDir) == 0 {
		return nil, errHomeDirNotSet
	}

	persistentDir := filepath.Join(homeDir, "Persistent")
	if stat, err := os.Lstat(persistentDir); err == nil && stat.IsDir() {
		// Looks like Tails.
		homeDir = persistentDir
	}

	configFile := filepath.Join(homeDir, ".xmpp-client")
	return &configFile, nil
}

func ParseMultiConfig(filename string) (*MultiAccountConfig, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c, err := parseMultiConfig(contents)
	if err != nil {
		return nil, err
	}

	c.filename = filename

	return c, nil
}

func parseMultiConfig(conf []byte) (m *MultiAccountConfig, err error) {
	m = &MultiAccountConfig{}
	if err = json.Unmarshal(conf, &m); err != nil {
		return
	}

	if m.Accounts == nil {
		return fallbackToSingleAccountConfig(conf)
	}

	return
}

func fallbackToSingleAccountConfig(conf []byte) (*MultiAccountConfig, error) {
	c, err := parseConfig(conf)
	if err != nil {
		return nil, err
	}

	//TODO: Convert from single to multi account format

	return &MultiAccountConfig{
		Accounts: []Config{*c},
	}, nil
}

func ParseConfig(filename string) (*Config, error) {
	m, err := ParseMultiConfig(filename)
	if err != nil {
		return nil, err
	}

	if len(m.Accounts) == 0 {
		return nil, errors.New("account config is missing")
	}

	c := &m.Accounts[0]
	c.filename = filename
	return c, parseFingerprints(c)
}

func parseConfig(contents []byte) (c *Config, err error) {
	c = new(Config)
	if err = json.Unmarshal(contents, &c); err != nil {
		return
	}

	return
}

func parseFingerprints(c *Config) error {
	var err error
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].fingerprint, err = hex.DecodeString(known.FingerprintHex)
		if err != nil {
			return errors.New("xmpp: failed to parse hex fingerprint for " + known.UserId + ": " + err.Error())
		}
	}

	return nil
}

func (c *Config) Save() error {
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].FingerprintHex = hex.EncodeToString(known.fingerprint)
	}

	contents, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.filename, contents, 0600)
}

func (c *Config) UserIdForFingerprint(fpr []byte) string {
	for _, known := range c.KnownFingerprints {
		if bytes.Equal(fpr, known.fingerprint) {
			return known.UserId
		}
	}

	return ""
}

func (c *Config) HasFingerprint(uid string) bool {
	for _, known := range c.KnownFingerprints {
		if uid == known.UserId {
			return true
		}
	}

	return false
}

func (c *Config) ShouldEncryptTo(uid string) bool {
	if c.AlwaysEncrypt {
		return true
	}

	for _, contact := range c.AlwaysEncryptWith {
		if contact == uid {
			return true
		}
	}
	return false
}

func parseYes(input string) bool {
	switch strings.ToLower(input) {
	case "y", "yes":
		return true
	}

	return false
}

func enroll(config *Config, term *terminal.Terminal) bool {
	var err error
	warn(term, "Enrolling new config file")

	var domain string
	for {
		term.SetPrompt("Account (i.e. user@example.com, enter to quit): ")
		if config.Account, err = term.ReadLine(); err != nil || len(config.Account) == 0 {
			return false
		}

		parts := strings.SplitN(config.Account, "@", 2)
		if len(parts) != 2 {
			alert(term, "invalid username (want user@domain): "+config.Account)
			continue
		}
		domain = parts[1]
		break
	}

	term.SetPrompt("Enable debug logging to /tmp/xmpp-client-debug.log? ")
	if debugLog, err := term.ReadLine(); err != nil || !parseYes(debugLog) {
		info(term, "Not enabling debug logging...")
	} else {
		info(term, "Debug logging enabled...")
		config.RawLogFile = "/tmp/xmpp-client-debug.log"
	}

	term.SetPrompt("Use Tor?: ")
	if useTorQuery, err := term.ReadLine(); err != nil || len(useTorQuery) == 0 || !parseYes(useTorQuery) {
		info(term, "Not using Tor...")
		config.UseTor = false
	} else {
		info(term, "Using Tor...")
		config.UseTor = true
	}

	term.SetPrompt("File to import libotr private key from (enter to generate): ")

	var priv otr.PrivateKey
	for {
		importFile, err := term.ReadLine()
		if err != nil {
			return false
		}
		if len(importFile) > 0 {
			privKeyBytes, err := ioutil.ReadFile(importFile)
			if err != nil {
				alert(term, "Failed to open private key file: "+err.Error())
				continue
			}

			if !priv.Import(privKeyBytes) {
				alert(term, "Failed to parse libotr private key file (the parser is pretty simple I'm afraid)")
				continue
			}
			break
		} else {
			info(term, "Generating private key...")
			priv.Generate(rand.Reader)
			break
		}
	}
	config.PrivateKey = priv.Serialize()

	config.OTRAutoAppendTag = true
	config.OTRAutoStartSession = true
	config.OTRAutoTearDown = false

	// List well known Tor hidden services.
	knownTorDomain := map[string]string{
		"jabber.ccc.de":             "okj7xc6j2szr2y75.onion",
		"riseup.net":                "4cjw6cwpeaeppfqz.onion",
		"jabber.calyxinstitute.org": "ijeeynrc6x2uy5ob.onion",
		"jabber.otr.im":             "5rgdtlawqkcplz75.onion",
		"wtfismyip.com":             "ofkztxcohimx34la.onion",
	}

	// Autoconfigure well known Tor hidden services.
	if hiddenService, ok := knownTorDomain[domain]; ok && config.UseTor {
		const torProxyURL = "socks5://127.0.0.1:9050"
		info(term, "It appears that you are using a well known server and we will use its Tor hidden service to connect.")
		config.Server = hiddenService
		config.Port = 5222
		config.Proxies = []string{torProxyURL}
		term.SetPrompt("> ")
		return true
	}

	var proxyStr string
	proxyDefaultPrompt := ", enter for none"
	if config.UseTor {
		proxyDefaultPrompt = ", which is the default"
	}
	term.SetPrompt("Proxy (i.e socks5://127.0.0.1:9050" + proxyDefaultPrompt + "): ")

	for {
		if proxyStr, err = term.ReadLine(); err != nil {
			return false
		}
		if len(proxyStr) == 0 {
			if !config.UseTor {
				break
			} else {
				proxyStr = "socks5://127.0.0.1:9050"
			}
		}
		u, err := url.Parse(proxyStr)
		if err != nil {
			alert(term, "Failed to parse "+proxyStr+" as a URL: "+err.Error())
			continue
		}
		if _, err = proxy.FromURL(u, proxy.Direct); err != nil {
			alert(term, "Failed to parse "+proxyStr+" as a proxy: "+err.Error())
			continue
		}
		break
	}

	if len(proxyStr) > 0 {
		config.Proxies = []string{proxyStr}

		info(term, "Since you selected a proxy, we need to know the server and port to connect to as a SRV lookup would leak information every time.")
		term.SetPrompt("Server (i.e. xmpp.example.com, enter to lookup using unproxied DNS): ")
		if config.Server, err = term.ReadLine(); err != nil {
			return false
		}
		if len(config.Server) == 0 {
			var port uint16
			info(term, "Performing SRV lookup")
			if config.Server, port, err = xmpp.Resolve(domain); err != nil {
				alert(term, "SRV lookup failed: "+err.Error())
				return false
			}
			config.Port = int(port)
			info(term, "Resolved "+config.Server+":"+strconv.Itoa(config.Port))
		} else {
			for {
				term.SetPrompt("Port (enter for 5222): ")
				portStr, err := term.ReadLine()
				if err != nil {
					return false
				}
				if len(portStr) == 0 {
					portStr = "5222"
				}
				if config.Port, err = strconv.Atoi(portStr); err != nil || config.Port <= 0 || config.Port > 65535 {
					info(term, "Port numbers must be 0 < port <= 65535")
					continue
				}
				break
			}
		}
	}

	term.SetPrompt("> ")

	return true
}

func loadConfig(ui UI) (*Config, string, error) {
	var err error

	if len(*configFile) == 0 {
		if configFile, err = findConfigFile(os.Getenv("HOME")); err != nil {
			ui.Alert(err.Error())
			return nil, "", err
		}
	}

	config, err := ParseConfig(*configFile)
	if err != nil {
		ui.Alert("Failed to parse config file: " + err.Error())
		config = new(Config)
		if !ui.Enroll(config) {
			return config, "", errors.New("Failed to create config")
		}

		config.filename = *configFile
		config.Save()
	}

	password := config.Password
	if len(password) == 0 {
		if password, err = ui.AskForPassword(config); err != nil {
			ui.Alert("Failed to read password: " + err.Error())
			return config, "", err
		}
	}

	return config, password, err
}

func NewXMPPConn(ui UI, config *Config, password string, createCallback xmpp.FormCallback, logger io.Writer) (*xmpp.Conn, error) {
	parts := strings.SplitN(config.Account, "@", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid username (want user@domain): " + config.Account)
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
			return nil, errors.New("Cannot connect via a proxy without Server and Port being set in the config file as an SRV lookup would leak information.")
		}

		host, port, err := xmpp.Resolve(domain)
		if err != nil {
			return nil, errors.New("Failed to resolve XMPP server: " + err.Error())
		}
		addr = fmt.Sprintf("%s:%d", host, port)
	}

	var dialer proxy.Dialer
	for i := len(config.Proxies) - 1; i >= 0; i-- {
		u, err := url.Parse(config.Proxies[i])
		if err != nil {
			return nil, errors.New("Failed to parse " + config.Proxies[i] + " as a URL: " + err.Error())
		}

		if dialer == nil {
			dialer = proxy.Direct
		}

		if dialer, err = proxy.FromURL(u, dialer); err != nil {
			return nil, errors.New("Failed to parse " + config.Proxies[i] + " as a proxy: " + err.Error())
		}
	}

	var certSHA256 []byte
	var err error
	if len(config.ServerCertificateSHA256) > 0 {
		certSHA256, err = hex.DecodeString(config.ServerCertificateSHA256)
		if err != nil {
			return nil, errors.New("Failed to parse ServerCertificateSHA256 (should be hex string): " + err.Error())
		}

		if len(certSHA256) != 32 {
			return nil, errors.New("ServerCertificateSHA256 is not 32 bytes long")
		}
	}

	xmppConfig := &xmpp.Config{
		Log:                     logger,
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
			//TODO: UI should have a Alert() method
			//alert(term, "Temporarily trusting only CACert root for CCC Jabber server")
			roots.AddCert(caCertRoot)
			xmppConfig.TLSConfig.RootCAs = roots
		} else {
			//TODO
			//alert(term, "Tried to add CACert root for jabber.ccc.de but failed: "+err.Error())
		}
	}

	//TODO: It may be locking
	//Also, move this defered functions
	//if len(config.RawLogFile) > 0 {
	//	rawLog, err := os.OpenFile(config.RawLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	//	if err != nil {
	//		return nil, errors.New("Failed to open raw log file: " + err.Error())
	//	}

	//	lock := new(sync.Mutex)
	//	in := rawLogger{
	//		out:    rawLog,
	//		prefix: []byte("<- "),
	//		lock:   lock,
	//	}
	//	out := rawLogger{
	//		out:    rawLog,
	//		prefix: []byte("-> "),
	//		lock:   lock,
	//	}
	//	in.other, out.other = &out, &in

	//	xmppConfig.InLog = &in
	//	xmppConfig.OutLog = &out

	//	defer in.flush()
	//	defer out.flush()
	//}

	if dialer != nil {
		//TODO
		//info(term, "Making connection to "+addr+" via proxy")
		if xmppConfig.Conn, err = dialer.Dial("tcp", addr); err != nil {
			return nil, errors.New("Failed to connect via proxy: " + err.Error())
		}
	}

	conn, err := xmpp.Dial(addr, user, domain, password, xmppConfig)
	if err != nil {
		return nil, errors.New("Failed to connect to XMPP server: " + err.Error())
	}

	return conn, err
}
