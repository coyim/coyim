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
	Server{"cloak.dk", "m2dsl4banuimpm6c.onion"}.register()
	Server{"darkness.su", "darknesswn664fcx.onion"}.register()
	Server{"dukgo.com", "wlcpmruglhxp6quz.onion"}.register()
	Server{"im.koderoot.net", "ihkw7qy3tok45dun.onion"}.register()
	Server{"jabber.calyxinstitute.org", "ijeeynrc6x2uy5ob.onion"}.register()
	Server{"jabber.ccc.de", "okj7xc6j2szr2y75.onion"}.register()
	Server{"jabber.cryptoparty.is", "cryjabkbdljzohnp.onion"}.register()
	Server{"jabber.ipredator.se", "3iffdebkzzkpgipa.onion"}.register()
	Server{"jabber.lqdn.fr", "jabber63t4r2qi57.onion"}.register()
	Server{"jabber.otr.im", "5rgdtlawqkcplz75.onion"}.register()
	Server{"jabber.so36.net", "s4fgy24e2b5weqdb.onion"}.register()
	Server{"jabber.systemli.org", "x5tno6mwkncu5m3h.onion"}.register()
	Server{"kjabber.de", "JABBERthelv5p7qv.onion"}.register()
	Server{"kode.im", "ihkw7qy3tok45dun.onion"}.register()
	Server{"riseup.net", "4cjw6cwpeaeppfqz.onion"}.register()
	Server{"securejabber.me", "giyvshdnojeivkom.onion"}.register()
	Server{"wtfismyip.com", "ofkztxcohimx34la.onion"}.register()
	Server{"xmpp.rows.io", "yz6yiv2hxyagvwy6.onion"}.register()
}

// Get returns the given server information if it is known, and not ok otherwise
func Get(s string) (serv Server, ok bool) {
	serv, ok = known[s]
	return
}
