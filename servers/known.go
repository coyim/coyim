package servers

import "sort"

// Server represent a known server
type Server struct {
	Name        string
	Onion       string
	CanRegister bool
	Recommended bool
	BrokenSCRAM bool
}

var known = make(map[string]Server)

func (s Server) register() {
	if _, already := known[s.Name]; already {
		panic("double registration of " + s.Name)
	}
	known[s.Name] = s
}

func init() {
	// These are the servers who we know the onion address for:

	Server{"5222.de", "jtovcabr2vhflcqg.onion", false, false, false}.register()
	Server{"adamas.ai", "gujl27qb6saznwv4.onion", false, false, false}.register()
	Server{"bommboo.de", "ujvdniabz53upqfx.onion", false, false, false}.register()
	Server{"chatme.im", "hdtp7t4nwifcfkjg.onion", true, false, false}.register()
	Server{"cloak.dk", "m2dsl4banuimpm6c.onion", false, false, false}.register()
	Server{"creep.im", "creep7nissfumwyx.onion", true, false, false}.register()
	Server{"darkness.su", "darknesswn664fcx.onion", false, false, false}.register()
	Server{"deshalbfrei.org", "jfel5icoxf3nmftl.onion", false, false, false}.register()
	Server{"dismail.de", "l2epd2e4g2hx3tdf.onion", true, false, false}.register()
	Server{"draugr.de", "jfel5icoxf3nmftl.onion", true, false, false}.register()
	Server{"evil.im", "evilxro6nvjuvxqo.onion", true, false, false}.register()
	Server{"hot-chilli.net", "c2aaokzwkwkct543.onion", true, false, false}.register()
	Server{"im.koderoot.net", "ihkw7qy3tok45dun.onion", false, false, false}.register()
	Server{"jabber-germany.de", "dbbrphko5tqcpar3.onion", true, false, false}.register()
	Server{"jabber.calyxinstitute.org", "ijeeynrc6x2uy5ob.onion", true, true, false}.register()
	Server{"jabber.cat", "sybzodlxacch7st7.onion", true, false, false}.register()
	Server{"jabber.ccc.de", "okj7xc6j2szr2y75.onion", true, true, false}.register()
	Server{"jabber.cryptoparty.is", "cryjabkbdljzohnp.onion", false, false, false}.register()
	Server{"jabber.frozenstar.info", "potu7aaoitlajnxc.onion", false, false, false}.register()
	Server{"jabber.ipredator.se", "3iffdebkzzkpgipa.onion", false, false, false}.register()
	Server{"jabber.lqdn.fr", "jabber63t4r2qi57.onion", false, false, false}.register()
	Server{"jabber.otr.im", "5rgdtlawqkcplz75.onion", true, true, true}.register()
	Server{"jabber.s7t.de", "jabberip5hpbrafx.onion", false, false, false}.register()
	Server{"jabber.so36.net", "s4fgy24e2b5weqdb.onion", false, false, false}.register()
	Server{"jabber.systemausfall.org", "clciwvt5qnxoqykx.onion", true, false, false}.register()
	Server{"jabber.systemli.org", "x5tno6mwkncu5m3h.onion", false, false, false}.register()
	Server{"jabberforum.de", "jabjabdevfoob7hl.onion", false, false, false}.register()
	Server{"jabberwiki.de", "jabjabdevfoob7hl.onion", false, false, false}.register()
	Server{"jabjab.de", "jabjabdevfoob7hl.onion", true, false, false}.register()
	Server{"kjabber.de", "JABBERthelv5p7qv.onion", false, false, false}.register()
	Server{"kode.im", "ihkw7qy3tok45dun.onion", false, false, false}.register()
	Server{"lethyro.net", "jabjabdevfoob7hl.onion", false, false, false}.register()
	Server{"pad7.de", "jabjabdevfoob7hl.onion", false, false, false}.register()
	Server{"pad7.net", "jabjabdevfoob7hl.onion", false, false, false}.register()
	Server{"patchcord.be", "xsydhi3dnbjuatpz.onion", true, false, false}.register()
	Server{"pimux.de", "maspm2xs6xavmpo6.onion", false, false, false}.register()
	Server{"planetjabber.de", "jabjabdevfoob7hl.onion", false, false, false}.register()
	Server{"riseup.net", "4cjw6cwpeaeppfqz.onion", false, false, false}.register()
	Server{"securejabber.me", "giyvshdnojeivkom.onion", false, false, false}.register()
	Server{"securetalks.biz", "wsjabberhzuots2e.onion", false, false, false}.register()
	Server{"suchat.org", "un3v2tzz4eplttzq.onion", true, false, false}.register()
	Server{"tchncs.de", "duvfmyqmdlyvc3mi.onion", true, false, false}.register()
	Server{"trashserver.net", "m4c722bvc2r7brnn.onion", true, false, false}.register()
	Server{"ubuntu-jabber.de", "jfel5icoxf3nmftl.onion", false, false, false}.register()
	Server{"ubuntu-jabber.net", "jfel5icoxf3nmftl.onion", false, false, false}.register()
	Server{"verdammung.org", "jfel5icoxf3nmftl.onion", false, false, false}.register()
	Server{"wallstreetjabber.biz", "wsjabberhzuots2e.onion", false, false, false}.register()
	Server{"wallstreetjabber.com", "wsjabberhzuots2e.onion", false, false, false}.register()
	Server{"wiuwiu.de", "jrbiogs6dv5txt5s.onion", false, false, false}.register()
	Server{"wtfismyip.com", "ofkztxcohimx34la.onion", false, false, false}.register()
	Server{"xabber.de", "jfel5icoxf3nmftl.onion", false, false, false}.register()
	Server{"xmpp.is", "y2qmqomqpszzryei.onion", true, true, false}.register()
	Server{"xmpp.rows.io", "yz6yiv2hxyagvwy6.onion", false, false, false}.register()
	Server{"xor.li", "nt3rf3kjsrle4vtf.onion", false, false, false}.register()
	Server{"ybgood.de", "jabjabdevfoob7hl.onion", false, false, false}.register()

	// These are the servers with public registration with A, A ranking from
	// https://xmpp.net/directory.php

	Server{"blah.im", "", true, false, false}.register()
	Server{"ch3kr.net", "", true, false, false}.register()
	Server{"chinwag.im", "", true, false, false}.register()
	Server{"core.mx", "", true, false, false}.register()
	Server{"datenknoten.me", "", true, false, false}.register()
	Server{"im.apinc.org", "", true, false, false}.register()
	Server{"is-a-furry.org", "", true, false, false}.register()
	Server{"jabber-hosting.de", "", true, false, false}.register()
	Server{"jabber.at", "", true, false, false}.register()
	Server{"jabber.chaos-darmstadt.de", "", true, false, false}.register()
	Server{"jabber.de", "", true, false, false}.register()
	Server{"jabber.hot-chilli.net", "", true, false, false}.register()
	Server{"jabber.meta.net.nz", "", true, false, false}.register()
	Server{"jabber.no", "", true, false, false}.register()
	Server{"jabber.no-sense.net", "", true, false, false}.register()
	Server{"jabber.schnied.net", "", true, false, false}.register()
	Server{"jabber.zone", "", true, false, false}.register()
	Server{"jabberzac.org", "", true, false, false}.register()
	Server{"jabbim.cz", "", true, false, false}.register()
	Server{"jabbim.hu", "", true, false, false}.register()
	Server{"jabbim.pl", "", true, false, false}.register()
	Server{"jabbim.sk", "", true, false, false}.register()
	Server{"jabster.pl", "", true, false, false}.register()
	Server{"jappix.com", "", true, false, false}.register()
	Server{"krautspace.de", "", true, false, false}.register()
	Server{"kwoh.de", "", true, false, false}.register()
	Server{"lightwitch.org", "", true, false, false}.register()
	Server{"miqote.com", "", true, false, false}.register()
	Server{"neko.im", "", true, false, false}.register()
	Server{"njs.netlab.cz", "", true, false, false}.register()
	Server{"psjb.me", "", true, false, false}.register()
	Server{"richim.org", "", true, false, false}.register()
	Server{"tigase.im", "", true, false, false}.register()
	Server{"twattle.net", "", true, false, false}.register()
	Server{"univers-libre.net", "", true, false, false}.register()
	Server{"wusz.org", "", true, false, false}.register()
	Server{"xmpp-hosting.de", "", true, false, false}.register()
	Server{"xmpp.jp", "", true, false, false}.register()
	Server{"xmpp.zone", "", true, false, false}.register()
	Server{"yax.im", "", true, false, false}.register()
}

// Get returns the given server information if it is known, and not ok otherwise
func Get(s string) (serv Server, ok bool) {
	serv, ok = known[s]
	return
}

// GetOnion returns the onion server address for the given host, if we know of it
func GetOnion(s string) (string, bool) {
	serv, ok := known[s]
	if !ok {
		return "", ok
	}
	ok = serv.Onion != ""
	onion := serv.Onion
	return onion, ok
}

// We sort the servers first based on if they are recommended, second, if they have
// onion, third, by name

type sortedServers []Server

func (s sortedServers) Len() int { return len(s) }
func (s sortedServers) Less(i, j int) bool {
	if s[i].Recommended != s[j].Recommended {
		return s[i].Recommended
	}

	if s[i].Onion != s[j].Onion {
		return s[j].Onion == ""
	}

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
