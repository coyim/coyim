package xmpp

import (
	"errors"
	"log"
	"net"
	"time"

	ourNet "github.com/twstrike/coyim/net"
	"golang.org/x/net/proxy"
)

var (
	// ErrConnectionFailed indicates a failure to connect to the server provided.
	ErrConnectionFailed = errors.New("could not connect to XMPP server")
)

const defaultDialTimeout = 30 * time.Second

func (d *Dialer) newTCPConn() (net.Conn, error) {
	if d.Proxy == nil {
		d.Proxy = proxy.Direct
	}

	addr := d.GetServer()

	//libpurple and xmpp-client are strict to this section and skip SRV lookup
	//whenever the user has configured a custom server address. We decided to not
	//do this, and as a consequence we can only set a custom server to an address
	//containing SRV records, otherwise it will fail.
	//On libpurple, a custom server with value "xmpp.google.com" or "talk.google.com"
	//will work, but for us it will fail unless config.SkipSRVLookup is set.
	//Our users must use "google.com" or "gmail.com" when asked for a custom server address
	//RFC 6120, Section 3.2.3
	//See: https://xmpp.org/rfcs/rfc6120.html#tcp-resolution-srvnot
	if d.Config.SkipSRVLookup {
		log.Println("Skipping SRV lookup")
		return connectWithProxy(addr, d.Proxy)
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	log.Println("Make SRV lookup to:", host)
	xmppAddrs, err := ResolveSRVWithProxy(d.Proxy, host)

	//Every other error means
	//"the initiating entity [did] not receive a response to its SRV query" and
	//we should use the fallback method
	//See RFC 6120, Section 3.2.1, item 9
	if err == ErrServiceNotAvailable {
		return nil, err
	}

	//RFC 6120, Section 3.2.1, item 9
	//If the SRV has no response, we fallback to use
	//the domain at default port
	if len(xmppAddrs) == 0 {
		//TODO: in this case, a failure to connect might be recovered using HTTP binding
		//See: RFC 6120, Section 3.2.2
		xmppAddrs = []string{
			net.JoinHostPort(d.getJIDDomainpart(), "5222"),
		}
	}

	conn, _, err := connectToFirstAvailable(xmppAddrs, d.Proxy)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func connectToFirstAvailable(xmppAddrs []string, dialer proxy.Dialer) (net.Conn, string, error) {
	for _, addr := range xmppAddrs {
		conn, err := connectWithProxy(addr, dialer)
		if err == nil {
			return conn, addr, nil
		}
	}

	return nil, "", ErrConnectionFailed
}

func dialTimeout(network, addr string, dialer proxy.Dialer, t time.Duration) (c net.Conn, err error) {
	result := make(chan bool, 1)

	go func() {
		c, err = dialer.Dial(network, addr)
		result <- true
	}()

	select {
	case <-time.After(t):
		log.Println("tcp: dial timed out")
		return nil, ourNet.ErrTimeout
	case <-result:
		return
	}
}

func connectWithProxy(addr string, dialer proxy.Dialer) (conn net.Conn, err error) {
	log.Printf("Connecting to %s\n", addr)

	//TODO: It is not clear to me if this follows
	//RFC 6120, Section 3.2.1, item 6
	//See: https://xmpp.org/rfcs/rfc6120.html#tcp-resolution
	conn, err = dialTimeout("tcp", addr, dialer, defaultDialTimeout)
	if err != nil {
		log.Printf("Failed to connect to %s: %s\n", addr, err)
		return
	}

	return
}
