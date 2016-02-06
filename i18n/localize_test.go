package i18n

import (
	"io/ioutil"
	"log"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
}

type I18NSuite struct{}

var _ = Suite(&I18NSuite{})

func (s *I18NSuite) Test_Local_willReturnTheString(c *C) {
	c.Assert(Local("hello"), Equals, "hello")
	c.Assert(Local("helllo"), Equals, "helllo")
}

func (s *I18NSuite) Test_Localf_willReturnTheString(c *C) {
	c.Assert(Localf("hello"), Equals, "hello")
	c.Assert(Localf("helllo %d", 42), Equals, "helllo 42")
}
