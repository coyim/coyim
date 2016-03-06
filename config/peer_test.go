package config

import (
	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
)

type AccountPeerSuite struct{}

var _ = Suite(&AccountPeerSuite{})

func (s *AccountPeerSuite) Test_Peer_ChangesNicknameAndGroups(c *C) {
	pid := "peer@coy.im"
	a := &Account{
		Account: "account@coy.im",
		Peers: []*Peer{
			&Peer{
				UserID: pid,
			},
		},
	}

	a.SavePeerDetails(pid, "bff", []string{"Foes"})

	p, _ := a.GetPeer(pid)
	c.Check(p.Nickname, Equals, "bff")
	c.Check(p.Groups, DeepEquals, []string{"Foes"})
}
