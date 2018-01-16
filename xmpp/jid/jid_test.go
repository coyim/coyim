package jid

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type JidXMPPSuite struct{}

var _ = Suite(&JidXMPPSuite{})

func (s *JidXMPPSuite) Test_ParseJID(c *C) {
	c.Assert(Parse(""), DeepEquals, Domain(""))
	c.Assert(Parse("foo.bar"), DeepEquals, Domain("foo.bar"))
	c.Assert(Parse("ola@foo.bar2"), DeepEquals, bare("ola@foo.bar2"))
	c.Assert(Parse("ola@foo.bar3/foo"), DeepEquals, full("ola@foo.bar3/foo"))
	c.Assert(Parse("foo/bar"), DeepEquals, domainWithResource("foo/bar"))
	c.Assert(Parse("foo3/bar/flux"), DeepEquals, domainWithResource("foo3/bar/flux"))
	c.Assert(Parse("foo3/bar/flux").NoResource(), DeepEquals, Domain("foo3"))
	c.Assert(Parse("foo3/bar/flux").(WithResource).Resource(), DeepEquals, Resource("bar/flux"))
	c.Assert(Parse("ola@foo.bar3").MaybeWithResource("one"), DeepEquals, full("ola@foo.bar3/one"))
	c.Assert(Parse("ola@foo.bar3/foo").MaybeWithResource("zero"), DeepEquals, full("ola@foo.bar3/zero"))
	c.Assert(Parse("ola@foo.bar3").MaybeWithResource(""), DeepEquals, bare("ola@foo.bar3"))
	c.Assert(Parse("ola@foo.bar3/baz").MaybeWithResource(""), DeepEquals, bare("ola@foo.bar3"))
}

func testInterface_Any(a Any) Any {
	return a
}

func testInterface_WithResource(a WithResource) WithResource {
	return a
}

func testInterface_WithoutResource(a WithoutResource) WithoutResource {
	return a
}

func testInterface_Bare(a Bare) Bare {
	return a
}

func testInterface_Full(a Full) Full {
	return a
}

func (s *JidXMPPSuite) Test_interfaceImplementations(c *C) {
	// There are no assertions here - if it compiles, we are fine.
	testInterface_Any(Domain("bla.com"))
	testInterface_Any(bare("bla@bla.com"))
	testInterface_Any(full("bla@bla.com/blu"))
	testInterface_Any(domainWithResource("bla.com/blu"))

	testInterface_WithResource(full("bla@bla.com/blu"))
	testInterface_WithResource(domainWithResource("bla.com/blu"))

	testInterface_WithoutResource(Domain("bla.com"))
	testInterface_WithoutResource(bare("bla@bla.com"))

	testInterface_Bare(bare("bla@bla.com"))

	testInterface_Full(full("bla@bla.com/blu"))
}
