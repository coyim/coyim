package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/twstrike/coyim/xmpp"
	"golang.org/x/net/proxy"
)

var (
	errHomeDirNotSet             = errors.New("$HOME not set. Please either export $HOME or use the -config-file option.\n")
	errInvalidConfigFile         = errors.New("Failed to parse config file")
	ConfigFile           *string = flag.String("config-file", "", "Location of the config file")
	CreateAccount        *bool   = flag.Bool("create", false, "If true, attempt to create account")
)

type Config struct {
	Filename string `json:"-"`

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

//TODO return a config with secure defaults
func NewConfig() *Config {
	var torProxy []string = nil
	torAddress := detectTor()

	if len(torAddress) != 0 {
		torProxy = []string{
			newTorProxy(torAddress),
		}
	}

	return &Config{
		//TODO: Should those 2 setting be set on startup (or maybe every connection)?
		Proxies: torProxy,
		UseTor:  torProxy != nil,

		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
	}
}

func NewXMPPConn(config *Config, password string, createCallback xmpp.FormCallback, logger io.Writer) (*xmpp.Conn, error) {
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
