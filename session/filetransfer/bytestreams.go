package filetransfer

import (
	"net"
	"strconv"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/socks5"
	"golang.org/x/net/proxy"
)

func tryStreamhost(s hasConfigAndLog, sh data.BytestreamStreamhost, dstAddr string, k func(net.Conn)) bool {
	port := sh.Port
	if port == 0 {
		port = 1080
	}

	p, err := s.GetConfig().CreateTorProxy()
	if err != nil {
		s.Log().WithError(err).Warn("Had error when trying to connect")
		return false
	}

	if p == nil {
		p = proxy.Direct
	}

	dialer, e := socks5.XMPP("tcp", net.JoinHostPort(sh.Host, strconv.Itoa(port)), nil, p)
	if e != nil {
		s.Log().WithError(e).WithField("streamhost", sh).Info("Error setting up socks5")
		return false
	}

	conn, e2 := dialer.Dial("tcp", net.JoinHostPort(dstAddr, "0"))
	if e2 != nil {
		s.Log().WithError(e2).WithField("streamhost", sh).Info("Error connecting socks5")
		return false
	}

	k(conn)
	return true
}
