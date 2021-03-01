package scram

import (
	"github.com/coyim/coyim/sasl"

	. "gopkg.in/check.v1"
)

func (s *ScramSuite) Test_Register(c *C) {
	Register()
	c.Assert(sasl.ClientSupport("SCRAM-SHA-1"), Equals, true)
	c.Assert(sasl.ClientSupport("SCRAM-SHA-1-PLUS"), Equals, true)
	c.Assert(sasl.ClientSupport("SCRAM-SHA-256"), Equals, true)
	c.Assert(sasl.ClientSupport("SCRAM-SHA-256-PLUS"), Equals, true)
	c.Assert(sasl.ClientSupport("SCRAM-SHA-512"), Equals, true)
	c.Assert(sasl.ClientSupport("SCRAM-SHA-512-PLUS"), Equals, true)
}
