package net

import (
	"errors"
	"net"
	"sort"
	"strconv"

	"github.com/miekg/dns"
	"golang.org/x/net/proxy"
)

const defaultDNSServer = "208.67.222.222:53"

// LookupSRV mirrors net.LookupSRV but uses the provided proxy dialer in order to do the lookup instead.
// By default it uses the OpenDNS server
func LookupSRV(dialer proxy.Dialer, service, proto, name string) (cname string, addrs []*net.SRV, err error) {
	return LookupSRVWith(dialer, defaultDNSServer, service, proto, name)
}

// LookupSRVWith looks up the provided service and protocol on the given name using the proxy dialer given and the dns server provided
func LookupSRVWith(dialer proxy.Dialer, dnsServer, service, proto, name string) (cname string, addrs []*net.SRV, err error) {
	cname = createCName(service, proto, name)

	conn, err := dialer.Dial("tcp", dnsServer)
	if err != nil {
		return
	}

	dnsConn := &dns.Conn{Conn: conn}
	defer dnsConn.Close()

	r, err := exchange(dnsConn, msgSRV(cname))
	if err != nil {
		return
	}

	addrs = convertAnswersToSRV(r.Answer)
	return
}

func createCName(service, proto, name string) string {
	return "_" + service + "._" + proto + "." + name + "."
}

func msgSRV(cname string) *dns.Msg {
	m := &dns.Msg{}
	m.SetQuestion(cname, dns.TypeSRV)
	m.RecursionDesired = true
	return m
}

func exchange(conn *dns.Conn, m *dns.Msg) (r *dns.Msg, err error) {
	if err = conn.WriteMsg(m); err != nil {
		return
	}

	if r, err = conn.ReadMsg(); err != nil {
		return
	}

	if r.Rcode != dns.RcodeSuccess {
		err = errors.New("got return: " + strconv.Itoa(r.Rcode))
	}

	return
}

type byPriorityWeight []*net.SRV

func (s byPriorityWeight) Len() int { return len(s) }
func (s byPriorityWeight) Less(i, j int) bool {
	if s[i] == nil {
		return true
	}

	if s[j] == nil {
		return false
	}

	if s[i].Priority == s[j].Priority {
		return s[i].Weight < s[j].Weight
	}

	return s[i].Priority < s[j].Priority
}
func (s byPriorityWeight) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func convertAnswersToSRV(in []dns.RR) []*net.SRV {
	result := make([]*net.SRV, len(in))

	for ix, a := range in {
		result[ix] = convertAnswerToSRV(a)
	}

	sort.Sort(byPriorityWeight(result))

	return result
}

func convertAnswerToSRV(in dns.RR) *net.SRV {
	srv, ok := in.(*dns.SRV)
	if ok {
		return &net.SRV{
			Target:   srv.Target,
			Port:     srv.Port,
			Priority: srv.Priority,
			Weight:   srv.Weight,
		}
	}
	return nil
}
