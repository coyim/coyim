package xmpp

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	ourNet "github.com/coyim/coyim/net"
	"github.com/coyim/coyim/xmpp/errors"
	"golang.org/x/net/proxy"
)

const defaultDialTimeout = 60 * time.Second

func (d *dialer) newTCPConn() (net.Conn, error) {
	if d.proxy == nil {
		d.proxy = proxy.Direct
	}

	//libpurple and xmpp-client are strict to section 3.2.3 and skip SRV lookup
	//whenever the user has configured a custom server address.
	//This is necessary to keep imported accounts from Adium/Pidgin/xmpp-client
	//working as expected.
	if d.hasCustomServer() {
		d.config.SkipSRVLookup = true
	}

	//RFC 6120, Section 3.2.3
	//See: https://xmpp.org/rfcs/rfc6120.html#tcp-resolution-srvnot
	if d.config.SkipSRVLookup {
		d.log.Info("Skipping SRV lookup")
		return d.connectWithProxy(d.GetServer(), false, d.proxy)
	}

	return d.srvLookupAndFallback()
}

func (d *dialer) srvLookupAndFallback() (net.Conn, error) {
	host := d.getJIDDomainpart()

	log.WithFields(log.Fields{
		"host": host,
	}).Info("Making SRV lookup")

	xmppsAddrs, xmppAddrs, err := ResolveSRVWithProxy(d.proxy, host)

	log.WithFields(log.Fields{
		"xmpp":  xmppAddrs,
		"xmpps": xmppsAddrs,
	}).Info("Received SRV records")

	//Every other error means
	//"the initiating entity [did] not receive a response to its SRV query" and
	//we should use the fallback method
	//See RFC 6120, Section 3.2.1, item 9
	if err == ErrServiceNotAvailable {
		return nil, err
	}

	//RFC 6120, Section 3.2.1, item 9
	//If the SRV has no response, we fallback to use the origin domain
	//at default port.
	if len(xmppsAddrs) == 0 && len(xmppAddrs) == 0 {
		err = errors.ErrTCPBindingFailed

		//TODO: in this case, a failure to connect might be recovered using HTTP binding
		//See: RFC 6120, Section 3.2.2
		xmppAddrs = []string{d.getFallbackServer()}
	} else {
		//The SRV lookup succeeded but we failed to connect
		err = errors.ErrConnectionFailed
	}

	conn, _, e := d.connectToFirstAvailable(xmppAddrs, false, d.proxy)
	if e != nil {
		return nil, err
	}

	return conn, nil
}

func (d *dialer) connectToFirstAvailable(xmppAddrs []string, tls bool, dialer proxy.Dialer) (net.Conn, string, error) {
	for _, addr := range xmppAddrs {
		conn, err := d.connectWithProxy(addr, tls, dialer)
		if err == nil {
			return conn, addr, nil
		}
	}

	return nil, "", errors.ErrConnectionFailed
}

func (d *dialer) dialTimeout(network, addr string, dialer proxy.Dialer, t time.Duration) (c net.Conn, err error) {
	result := make(chan bool, 1)

	go func() {
		c, err = dialer.Dial(network, addr)
		result <- true
	}()

	select {
	case <-time.After(t):
		d.log.Warn("tcp: dial timed out")
		return nil, ourNet.ErrTimeout
	case <-result:
		return
	}
}

func (d *dialer) connectWithProxy(addr string, tls bool, dialer proxy.Dialer) (conn net.Conn, err error) {
	d.log.WithField("addr", addr).Info("Connecting")

	//TODO: It is not clear to me if this follows
	//RFC 6120, Section 3.2.1, item 6
	//See: https://xmpp.org/rfcs/rfc6120.html#tcp-resolution
	conn, err = d.dialTimeout("tcp", addr, dialer, defaultDialTimeout)
	if err != nil {
		if err == ourNet.ErrTimeout {
			return nil, err
		}

		return nil, errors.CreateErrFailedToConnect(addr, err)
	}

	return
}
