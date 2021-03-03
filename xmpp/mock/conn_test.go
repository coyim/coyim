package mock

import (
	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type ConnSuite struct{}

var _ = Suite(&ConnSuite{})

func (s *ConnSuite) Test_ConnMock(c *C) {
	m := &Conn{}

	c.Assert(m.Authenticate("", ""), IsNil)
	c.Assert(m.AuthenticationFailure(), IsNil)
	c.Assert(m.BindResource(), IsNil)
	c.Assert(m.Cancel(0), Equals, false)
	c.Assert(m.Cache(), IsNil)
	c.Assert(m.ChangePassword("", "", ""), IsNil)
	c.Assert(m.Close(), IsNil)
	c.Assert(m.Config(), IsNil)
	c.Assert(m.CustomStorage(), IsNil)
	_, _ = m.DiscoveryFeatures("")
	_, _, _ = m.DiscoveryFeaturesAndIdentities("")
	c.Assert(m.Features(), DeepEquals, data.StreamFeatures{})
	_, _ = m.GetRosterDelimiter()
	c.Assert(m.In(), IsNil)
	c.Assert(m.Lock(), IsNil)
	_, _ = m.Next()
	c.Assert(m.OriginDomain(), Equals, "")
	c.Assert(m.Out(), IsNil)
	c.Assert(m.Rand(), IsNil)
	c.Assert(m.RawOut(), IsNil)
	c.Assert(m.ReadStanzas(nil), IsNil)
	_, _ = m.RegisterAccount("", "")
	_, _, _ = m.RequestRoster()
	_, _, _ = m.RequestVCard()
	c.Assert(m.Send("", "", false), IsNil)
	c.Assert(m.SendMessage(nil), IsNil)
	_, _, _ = m.SendIQ("", "", nil)
	c.Assert(m.SendIQReply("", "", "", nil), IsNil)
	c.Assert(m.SendInitialStreamHeader(), IsNil)
	_, _, _ = m.SendPing()
	c.Assert(m.SendPresence("", "", "", ""), IsNil)
	c.Assert(m.SendMUCPresence("", nil), IsNil)
	c.Assert(m.ServerAddress(), Equals, "")
	m.SetInOut(nil, nil)
	m.SetKeepaliveOut(nil)
	m.SetRawOut(nil)
	m.SetServerAddress("")
	c.Assert(m.SignalPresence(""), IsNil)

	m.SetChannelBinding(nil)
	c.Assert(m.GetChannelBinding(), IsNil)

	c.Assert(m.GetJIDResource(), Equals, "")
	m.SetJIDResource("")

	c.Assert(m.HasSupportTo(""), Equals, false)

	_, _ = m.QueryServiceInformation("")
	_, _ = m.QueryServiceItems("")
	_, _ = m.EntityExists("")

	c.Assert(m.ServerHasFeature(""), Equals, false)
}
