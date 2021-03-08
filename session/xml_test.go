package session

import (
	. "gopkg.in/check.v1"
)

type XMLSuite struct{}

var _ = Suite(&XMLSuite{})

func (s *XMLSuite) Test_tryDecodeXML_works(c *C) {
	nspace, local, ok := tryDecodeXML([]byte(`<foo xmlns="foo:bar:1">`))
	c.Assert(ok, Equals, true)
	c.Assert(nspace, Equals, "foo:bar:1")
	c.Assert(local, Equals, "foo")
}

func (s *XMLSuite) Test_tryDecodeXML_failsIfNoData(c *C) {
	nspace, local, ok := tryDecodeXML([]byte(``))
	c.Assert(ok, Equals, false)
	c.Assert(nspace, Equals, "")
	c.Assert(local, Equals, "")
}

func (s *XMLSuite) Test_tryDecodeXML_failsIfWrongToken(c *C) {
	nspace, local, ok := tryDecodeXML([]byte(`<!-- comment -->`))
	c.Assert(ok, Equals, false)
	c.Assert(nspace, Equals, "")
	c.Assert(local, Equals, "")
}
