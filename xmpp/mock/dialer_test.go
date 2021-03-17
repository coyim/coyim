package mock

import (
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	. "gopkg.in/check.v1"
)

type DialerSuite struct{}

var _ = Suite(&DialerSuite{})

func (s *DialerSuite) Test_DialerMock(c *C) {
	var d interfaces.Dialer = &Dialer{}

	c.Assert(d.Config(), DeepEquals, data.Config{})
	_, _ = d.Dial()
	c.Assert(d.GetServer(), Equals, "")
	_, _ = d.RegisterAccount(func(string, string, []interface{}) error { return nil })
	c.Assert(d.ServerAddress(), Equals, "")
	d.SetConfig(data.Config{})
	d.SetJID("")
	d.SetPassword("")
	d.SetProxy(nil)
	d.SetResource("")
	d.SetServerAddress("")
	d.SetShouldConnectTLS(false)
	d.SetShouldSendALPN(false)
	d.SetLogger(nil)
	d.SetKnown(nil)
}
