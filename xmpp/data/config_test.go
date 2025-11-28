package data

import (
	"io"
	"os"

	. "gopkg.in/check.v1"
)

type ConfigSuite struct{}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) Test_Config_GetLog_returnsDiscardIfLogIsNil(c *C) {
	conf := &Config{}
	c.Assert(conf.GetLog(), Equals, io.Discard)
}

func (s *ConfigSuite) Test_Config_GetLog_returnsTheLogSet(c *C) {
	conf := &Config{
		Log: os.Stdout,
	}
	c.Assert(conf.GetLog(), Equals, os.Stdout)
}
