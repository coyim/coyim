package filetransfer

import (
	"fmt"
	"net"
	"strconv"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/socks5"
	"golang.org/x/net/proxy"
)

func tryStreamhost(s access.Session, sh data.BytestreamStreamhost, dstAddr string, k func(net.Conn)) bool {
	port := sh.Port
	if port == 0 {
		port = 1080
	}

	p, err := s.GetConfig().CreateTorProxy()
	if err != nil {
		s.Warn(fmt.Sprintf("Had error when trying to connect: %v", err))
		return false
	}

	if p == nil {
		p = proxy.Direct
	}

	dialer, e := socks5.XMPP("tcp", net.JoinHostPort(sh.Host, strconv.Itoa(port)), nil, p)
	if e != nil {
		s.Info(fmt.Sprintf("Error setting up socks5 for %v: %v", sh, e))
		return false
	}

	conn, e2 := dialer.Dial("tcp", net.JoinHostPort(dstAddr, "0"))
	if e2 != nil {
		s.Info(fmt.Sprintf("Error connecting socks5 for %v: %v", sh, e2))
		return false
	}

	k(conn)
	return true
}
