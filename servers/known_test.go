package servers

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
	Server{"something.de", "123123123.onion"}.register()
	serv, _ := Get("something.de")
	c.Assert(serv.Onion, Equals, "123123123.onion")
}
