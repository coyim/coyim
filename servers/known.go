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
	// These are the servers for which we know the v3 onion address:

	Server{"anrc.mooo.com", "6w5iasklrbr2kw53zqrsjktgjapvjebxodoki3gjnmvb4dvcbmz7n3qd.onion", false, false, false}.register()
	Server{"dismail.de", "4colmnerbjz3xtsjmqogehtpbt5upjzef57huilibbq3wfgpsylub7yd.onion", true, false, false}.register()
	Server{"jabber.cat", "7drfpncjeom3svqkyjitif26ezb3xvmtgyhgplcvqa7wwbb4qdbsjead.onion", true, false, false}.register()
	Server{"jabber.de", "uoj2xiqxk25p36wbpufiyuhluvxakhpqum7frembhoiuq7a5735ay3qd.onion", true, false, false}.register()
	Server{"jabber.nr18.space", "szd7r26dbcrrrn4jthercrdypxfdmzzrysusyjohn4mpv2zbwcgmeqqd.onion", false, false, false}.register()
	Server{"jabber.otr.im", "ynnuxkbbiy5gicdydekpihmpbqd4frruax2mqhpc35xqjxp5ayvrjuqd.onion", true, true, true}.register()
	Server{"jabber.so36.net", "yxkc2uu3rlwzzhxf2thtnzd7obsdd76vtv7n34zwald76g5ogbvjbbqd.onion", false, false, false}.register()
	Server{"jabber.systemausfall.org", "jaswtrycaot3jzkr7znje4ebazzvbxtzkyyox67frgvgemwfbzzi6uqd.onion", true, false, false}.register()
	Server{"jabber.systemli.org", "razpihro3mgydaiykvxwa44l57opvktqeqfrsg3vvwtmvr2srbkcihyd.onion", false, false, false}.register()
	Server{"krautspace.de", "jeirlvruhz22jqduzixi6li4xyoweytqglwjons4mbuif76fgslg5uad.onion", true, false, false}.register()
	Server{"talk36.net", "yxkc2uu3rlwzzhxf2thtnzd7obsdd76vtv7n34zwald76g5ogbvjbbqd.onion", false, false, false}.register()
	Server{"trashserver.net", "xiynxwxxpw7olq76uhrbvx2ts3i7jagqnqix7arfbknmleuoiwsmt5yd.onion", false, false, false}.register()
	Server{"wiuwiu.de", "qawb5xl3mxiixobjsw2d45dffngyyacp4yd3wjpmhdrazwvt4ytxvayd.onion", false, false, false}.register()
	Server{"xmpp.is", "6voaf7iamjpufgwoulypzwwecsm2nu7j5jpgadav2rfqixmpl4d65kid.onion", true, true, false}.register()
	Server{"xmpp.riseup.net", "jukrlvyhgguiedqswc5lehrag2fjunfktouuhi4wozxhb6heyzvshuyd.onion", false, false, false}.register()

	// These are all the hosts controlled by jabjab

	Server{"jabjab.de", "jabjabdea2eewo3gzfurscj2sjqgddptwumlxi3wur57rzf5itje2rid.onion", true, false, false}.register()
	Server{"jabjabj.de", "jabjabdea2eewo3gzfurscj2sjqgddptwumlxi3wur57rzf5itje2rid.onion", false, false, false}.register()
	Server{"lethyro.net", "jabjabdea2eewo3gzfurscj2sjqgddptwumlxi3wur57rzf5itje2rid.onion", false, false, false}.register()
	Server{"pad7.de", "jabjabdea2eewo3gzfurscj2sjqgddptwumlxi3wur57rzf5itje2rid.onion", false, false, false}.register()
	Server{"pad7.net", "jabjabdea2eewo3gzfurscj2sjqgddptwumlxi3wur57rzf5itje2rid.onion", false, false, false}.register()
	Server{"planetjabber.de", "jabjabdea2eewo3gzfurscj2sjqgddptwumlxi3wur57rzf5itje2rid.onion", false, false, false}.register()

	// These are all the hosts controlled by hot-chilli

	Server{"jabber.hot-chilli.net", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", true, false, false}.register()
	Server{"jabber.hot-chilli.eu", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"hot-chilli.net", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"hot-chilli.eu", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"im.hot-chilli.net", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"im.hot-chilli.eu", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"jabb3r.de", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"jabb3r.org", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"jabber-hosting.de", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()
	Server{"xmpp-hosting.de", "chillingguw3yu2rmrkqsog4554egiry6fmy264l5wblyadds3c2lnyd.onion", false, false, false}.register()

	// These are servers that used to have onion services, but only v2 versions are available at this time.
	// For this reason, the onion services have been removed, until we can find v3 versions for these servers.

	Server{"5222.de", "", false, false, false}.register()
	Server{"adamas.ai", "", false, false, false}.register()
	Server{"bommboo.de", "", false, false, false}.register()
	Server{"chatme.im", "", true, false, false}.register()
	Server{"cloak.dk", "", false, false, false}.register()
	Server{"creep.im", "", true, false, false}.register()
	Server{"darkness.su", "", false, false, false}.register()
	Server{"deshalbfrei.org", "", false, false, false}.register()
	Server{"draugr.de", "", true, false, false}.register()
	Server{"evil.im", "", true, false, false}.register()
	Server{"im.koderoot.net", "", false, false, false}.register()
	Server{"jabber-germany.de", "", true, false, false}.register()
	Server{"jabber.calyxinstitute.org", "", true, true, false}.register()
	Server{"jabber.ccc.de", "", true, true, false}.register()
	Server{"jabber.cryptoparty.is", "", false, false, false}.register()
	Server{"jabber.frozenstar.info", "", false, false, false}.register()
	Server{"jabber.ipredator.se", "", false, false, false}.register()
	Server{"jabber.lqdn.fr", "", false, false, false}.register()
	Server{"jabber.s7t.de", "", false, false, false}.register()
	Server{"jabberforum.de", "", false, false, false}.register()
	Server{"jabberwiki.de", "", false, false, false}.register()
	Server{"kjabber.de", "", false, false, false}.register()
	Server{"kode.im", "", false, false, false}.register()
	Server{"patchcord.be", "", true, false, false}.register()
	Server{"pimux.de", "", false, false, false}.register()
	Server{"securejabber.me", "", false, false, false}.register()
	Server{"securetalks.biz", "", false, false, false}.register()
	Server{"suchat.org", "", true, false, false}.register()
	Server{"tchncs.de", "", true, false, false}.register()
	Server{"ubuntu-jabber.de", "", false, false, false}.register()
	Server{"ubuntu-jabber.net", "", false, false, false}.register()
	Server{"verdammung.org", "", false, false, false}.register()
	Server{"wallstreetjabber.biz", "", false, false, false}.register()
	Server{"wallstreetjabber.com", "", false, false, false}.register()
	Server{"wtfismyip.com", "", false, false, false}.register()
	Server{"xabber.de", "", false, false, false}.register()
	Server{"xmpp.rows.io", "", false, false, false}.register()
	Server{"xor.li", "", false, false, false}.register()
	Server{"ybgood.de", "", false, false, false}.register()

	// These are the servers with public registration with A, A ranking from
	// https://xmpp.net/directory.php

	Server{"blah.im", "", true, false, false}.register()
	Server{"ch3kr.net", "", true, false, false}.register()
	Server{"chinwag.im", "", true, false, false}.register()
	Server{"core.mx", "", true, false, false}.register()
	Server{"datenknoten.me", "", true, false, false}.register()
	Server{"im.apinc.org", "", true, false, false}.register()
	Server{"is-a-furry.org", "", true, false, false}.register()
	Server{"jabber.at", "", true, false, false}.register()
	Server{"jabber.chaos-darmstadt.de", "", true, false, false}.register()
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
