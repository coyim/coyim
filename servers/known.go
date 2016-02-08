package servers

// Server represent a known server
type Server struct {
	Name  string
	Onion string
}

var known = make(map[string]Server)

func (s Server) register() {
	known[s.Name] = s
}

func init() {
	// https://adamas.ai/cgi-bin/wiki.pl/Services
	Server{"adamas.ai", "gujl27qb6saznwv4.onion"}.register()

	// http://cloak.dk/jabber.html
	Server{"cloak.dk", "m2dsl4banuimpm6c.onion"}.register()

	// https://darkness.su/en/xmpp.html
	Server{"darkness.su", "darknesswn664fcx.onion"}.register()

	// https://www.cryptoparty.in/connect/contact/jabber
	Server{"dukgo.com", "wlcpmruglhxp6quz.onion"}.register()
	Server{"jabber.calyxinstitute.org", "ijeeynrc6x2uy5ob.onion"}.register()
	Server{"jabber.ccc.de", "okj7xc6j2szr2y75.onion"}.register()
	Server{"jabber.cryptoparty.is", "cryjabkbdljzohnp.onion"}.register()
	Server{"jabber.ipredator.se", "3iffdebkzzkpgipa.onion"}.register()
	Server{"jabber.otr.im", "5rgdtlawqkcplz75.onion"}.register()
	Server{"jabber.so36.net", "s4fgy24e2b5weqdb.onion"}.register()
	Server{"jabber.systemli.org", "x5tno6mwkncu5m3h.onion"}.register()
	Server{"riseup.net", "4cjw6cwpeaeppfqz.onion"}.register()
	Server{"securejabber.me", "giyvshdnojeivkom.onion"}.register()
	Server{"wtfismyip.com", "ofkztxcohimx34la.onion"}.register()
	Server{"xmpp.rows.io", "yz6yiv2hxyagvwy6.onion"}.register()

	// http://kode.im
	Server{"kode.im", "ihkw7qy3tok45dun.onion"}.register()

	// https://space.koderoot.net/
	Server{"im.koderoot.net", "ihkw7qy3tok45dun.onion"}.register()

	// https://jabber.lqdn.fr/
	Server{"jabber.lqdn.fr", "jabber63t4r2qi57.onion"}.register()

	// https://www.kjabber.de/onion.htm
	Server{"kjabber.de", "JABBERthelv5p7qv.onion"}.register()

	// http://xor.li
	Server{"xor.li", "nt3rf3kjsrle4vtf.onion"}.register()
}

// Get returns the given server information if it is known, and not ok otherwise
func Get(s string) (serv Server, ok bool) {
	serv, ok = known[s]
	return
}
