package servers

import (
	. "gopkg.in/check.v1"
)

type KnownSuite struct{}

var _ = Suite(&KnownSuite{})

func (s *KnownSuite) Test_Get_returnsTheServerIfItExists(c *C) {
	serv, ok := Get("riseup.net")
	c.Assert(ok, Equals, true)
	c.Assert(serv.Onion, Equals, "4cjw6cwpeaeppfqz.onion")

	_, ok2 := Get("blarg.net")
	c.Assert(ok2, Equals, false)
}

func (s *KnownSuite) Test_register_willAddANewServer(c *C) {
	origKnown := known
	defer func() {
		known = origKnown
	}()

	known = map[string]Server{}

	Server{"something.de", "123123123.onion", false, false, false}.register()
	serv, _ := Get("something.de")
	c.Assert(serv.Onion, Equals, "123123123.onion")
}

func (s *KnownSuite) Test_register_willPanicIfDoubleRegistering(c *C) {
	origKnown := known
	defer func() {
		known = origKnown
	}()

	known = map[string]Server{}

	srv := Server{"something.de", "123123123.onion", false, false, false}

	srv.register()

	c.Assert(srv.register, PanicMatches, "double registration of something.de")
}

func (s *KnownSuite) Test_GetOnion_failsWhenServerNotKnown(c *C) {
	res, ok := GetOnion("somewhere.com")
	c.Assert(res, Equals, "")
	c.Assert(ok, Equals, false)
}

func (s *KnownSuite) Test_GetOnion_failsWhenServerDoesntHaveOnion(c *C) {
	res, ok := GetOnion("jabber.at")
	c.Assert(res, Equals, "")
	c.Assert(ok, Equals, false)
}

func (s *KnownSuite) Test_GetOnion_returnsOnion(c *C) {
	res, ok := GetOnion("jabber.ccc.de")
	c.Assert(res, Equals, "okj7xc6j2szr2y75.onion")
	c.Assert(ok, Equals, true)
}

func (s *KnownSuite) Test_GetServersForRegistration(c *C) {
	res := GetServersForRegistration()
	c.Assert(res, HasLen, 58)
}
