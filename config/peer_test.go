package config

import (
	. "gopkg.in/check.v1"
)

type AccountPeerSuite struct{}

var _ = Suite(&AccountPeerSuite{})

func (s *AccountPeerSuite) Test_Peer_ChangesNickname(c *C) {
	pid := "peer@coy.im"
	a := &Account {
		Account: "account@coy.im",
		Peers: []*Peer{
			&Peer{
				UserID: pid,
			},
		},
	}
	
	a.SetPeersNickname(pid, "bff")
	
	p, _ := a.GetPeer(pid)
	c.Check(p.Nickname, Equals, "bff")
}
