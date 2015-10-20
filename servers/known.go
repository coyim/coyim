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
	Server{"jabber.ccc.de", "okj7xc6j2szr2y75.onion"}.register()
	Server{"riseup.net", "4cjw6cwpeaeppfqz.onion"}.register()
	Server{"jabber.calyxinstitute.org", "ijeeynrc6x2uy5ob.onion"}.register()
	Server{"jabber.otr.im", "5rgdtlawqkcplz75.onion"}.register()
	Server{"wtfismyip.com", "ofkztxcohimx34la.onion"}.register()
	Server{"dukgo.com", "wlcpmruglhxp6quz.onion"}.register()
}

// Get returns the given server information if it is known, and not ok otherwise
func Get(s string) (serv Server, ok bool) {
	serv, ok = known[s]
	return
}
