package tls

import (
	"crypto/tls"

	. "gopkg.in/check.v1"
)

type FactorySuite struct{}

var _ = Suite(&FactorySuite{})

func (s *FactorySuite) Test_Real(c *C) {
	v := Real(nil, nil)
	c.Assert(v, FitsTypeOf, &tls.Conn{})
}
