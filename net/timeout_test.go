package net

import (
	. "gopkg.in/check.v1"
)

type TimeoutSuite struct{}

var _ = Suite(&TimeoutSuite{})

func (s *TimeoutSuite) Test_TimeoutError_Error(c *C) {
	c.Assert(&TimeoutError{}, ErrorMatches, "i/o timeout")
}
