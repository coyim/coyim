package main

import (
	"crypto/rand"
	"flag"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	coyconf "github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/proxy"
)

var configFile *string = flag.String("config-file", "", "Location of the config file")
var createAccount *bool = flag.Bool("create", false, "If true, attempt to create account")

func enroll(config *coyconf.Config, term *terminal.Terminal) bool {
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
	if debugLog, err := term.ReadLine(); err != nil || !coyconf.ParseYes(debugLog) {
		info(term, "Not enabling debug logging...")
	} else {
		info(term, "Debug logging enabled...")
		config.RawLogFile = "/tmp/xmpp-client-debug.log"
	}

	term.SetPrompt("Use Tor?: ")
	if useTorQuery, err := term.ReadLine(); err != nil || len(useTorQuery) == 0 || !coyconf.ParseYes(useTorQuery) {
		info(term, "Not using Tor...")
		config.UseTor = false
	} else {
		info(term, "Using Tor...")
		config.UseTor = true
	}

	term.SetPrompt("File to import libotr private key from (enter to generate): ")

	var priv otr3.PrivateKey
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
