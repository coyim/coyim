package session

import (
	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type MUCStatusesSuite struct{}

var _ = Suite(&MUCStatusesSuite{})

func (s *MUCStatusesSuite) Test_mucUserStatuses_contains(c *C) {
	v := mucUserStatuses{
		data.MUCUserStatus{Code: 42},
		data.MUCUserStatus{Code: 55},
		data.MUCUserStatus{Code: 16},
	}

	c.Assert(v.contains(), Equals, true)
	c.Assert(v.contains(42), Equals, true)
	c.Assert(v.contains(42, 55), Equals, true)
	c.Assert(v.contains(42, 17), Equals, false)
	c.Assert(v.contains(1), Equals, false)
	c.Assert(v.contains(1, 2, 3), Equals, false)
}

func (s *MUCStatusesSuite) Test_mucUserStatuses_containsAny(c *C) {
	v := mucUserStatuses{
		data.MUCUserStatus{Code: 42},
		data.MUCUserStatus{Code: 55},
		data.MUCUserStatus{Code: 16},
	}

	c.Assert(v.containsAny(), Equals, false)
	c.Assert(v.containsAny(42), Equals, true)
	c.Assert(v.containsAny(42, 55), Equals, true)
	c.Assert(v.containsAny(42, 17), Equals, true)
	c.Assert(v.containsAny(1, 16), Equals, true)
	c.Assert(v.containsAny(1, 2, 3), Equals, false)
}
