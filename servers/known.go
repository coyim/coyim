package servers

import "sort"

// Server represent a known server
type Server struct {
	Name        string
	Onion       string
	CanRegister bool
}

var known = make(map[string]Server)

func (s Server) register() {
	known[s.Name] = s
}

func init() {
	// https://adamas.ai/cgi-bin/wiki.pl/Services
	Server{"adamas.ai", "gujl27qb6saznwv4.onion", false}.register()

	// http://cloak.dk/jabber.html
	Server{"cloak.dk", "m2dsl4banuimpm6c.onion", false}.register()

	// https://darkness.su/en/xmpp.html
	Server{"darkness.su", "darknesswn664fcx.onion", false}.register()

	// https://www.cryptoparty.in/connect/contact/jabber
	Server{"jabber.calyxinstitute.org", "ijeeynrc6x2uy5ob.onion", true}.register()
	Server{"jabber.ccc.de", "okj7xc6j2szr2y75.onion", true}.register()
	Server{"jabber.cryptoparty.is", "cryjabkbdljzohnp.onion", false}.register()
	Server{"jabber.ipredator.se", "3iffdebkzzkpgipa.onion", false}.register()
	Server{"jabber.otr.im", "5rgdtlawqkcplz75.onion", true}.register()
	Server{"jabber.so36.net", "s4fgy24e2b5weqdb.onion", false}.register()
	Server{"jabber.systemli.org", "x5tno6mwkncu5m3h.onion", false}.register()
	Server{"riseup.net", "4cjw6cwpeaeppfqz.onion", false}.register()
	Server{"securejabber.me", "giyvshdnojeivkom.onion", false}.register()
	Server{"wtfismyip.com", "ofkztxcohimx34la.onion", false}.register()
	Server{"xmpp.rows.io", "yz6yiv2hxyagvwy6.onion", false}.register()

	// http://kode.im
	Server{"kode.im", "ihkw7qy3tok45dun.onion", false}.register()

	// https://space.koderoot.net/
	Server{"im.koderoot.net", "ihkw7qy3tok45dun.onion", false}.register()

	// https://jabber.lqdn.fr/
	Server{"jabber.lqdn.fr", "jabber63t4r2qi57.onion", false}.register()

	// https://www.kjabber.de/onion.htm
	Server{"kjabber.de", "JABBERthelv5p7qv.onion", false}.register()

	// http://xor.li
	Server{"xor.li", "nt3rf3kjsrle4vtf.onion", false}.register()
}

// Get returns the given server information if it is known, and not ok otherwise
func Get(s string) (serv Server, ok bool) {
	serv, ok = known[s]
	return
}

type sortedServers []Server

func (s sortedServers) Len() int { return len(s) }
func (s sortedServers) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
func (s sortedServers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetServersForRegistration returns all known servers that we can register at
func GetServersForRegistration() []Server {
	res := []Server{}
	for _, s := range known {
		if s.CanRegister {
			res = append(res, s)
		}
	}
	sort.Sort(sortedServers(res))
	return res
}
