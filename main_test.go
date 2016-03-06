package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"

	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type MainSuite struct{}

var _ = Suite(&MainSuite{})

func (s *MainSuite) Test_initLog_DoesntSetLogFlags_IfNotDebugging(c *C) {
	log.SetFlags(0)
	*config.DebugFlag = false
	initLog()
	c.Assert(log.Flags(), Equals, 0)
}

func (s *MainSuite) Test_initLog_SetsLogFlagsIfDebugging(c *C) {
	log.SetFlags(0)
	*config.DebugFlag = true
	initLog()
	c.Assert(log.Flags(), Equals, log.Ldate|log.Ltime|log.Llongfile)
	c.Assert(log.Prefix(), Equals, "[CoyIM] ")
}
