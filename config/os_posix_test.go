// +build !windows

package config

import (
	. "gopkg.in/check.v1"
)

type OSSuite struct{}

var _ = Suite(&OSSuite{})

func (s *OSSuite) Test_IsWindows(c *C) {
	c.Assert(IsWindows(), Equals, false)
}
