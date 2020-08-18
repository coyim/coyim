package jid

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"

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
	c.Assert(Parse(""), DeepEquals, Domain{""})
	c.Assert(Parse("foo.bar"), DeepEquals, Domain{"foo.bar"})
	c.Assert(Parse("ola@foo.bar2"), DeepEquals, ParseBare("ola@foo.bar2"))
	c.Assert(Parse("ola@foo.bar3/foo"), DeepEquals, ParseFull("ola@foo.bar3/foo"))
	c.Assert(Parse("foo/bar"), DeepEquals, domainWithResource{NewDomain("foo"), NewResource("bar")})
	c.Assert(Parse("foo3/bar/flux"), DeepEquals, domainWithResource{NewDomain("foo3"), NewResource("bar/flux")})
	c.Assert(Parse("foo3/bar/flux").NoResource(), DeepEquals, Domain{"foo3"})
	c.Assert(Parse("foo3/bar/flux").(WithResource).Resource(), DeepEquals, Resource{"bar/flux"})
	c.Assert(Parse("ola@foo.bar3").MaybeWithResource(Resource{"one"}), DeepEquals, ParseFull("ola@foo.bar3/one"))
	c.Assert(Parse("ola@foo.bar3/foo").MaybeWithResource(Resource{"zero"}), DeepEquals, ParseFull("ola@foo.bar3/zero"))
	c.Assert(Parse("ola@foo.bar3").MaybeWithResource(Resource{""}), DeepEquals, ParseBare("ola@foo.bar3"))
	c.Assert(Parse("ola@foo.bar3/baz").MaybeWithResource(Resource{""}), DeepEquals, ParseBare("ola@foo.bar3"))
}

func testInterfaceAny(a Any) Any {
	return a
}

func testInterfaceWithResource(a WithResource) WithResource {
	return a
}

func testInterfaceWithoutResource(a WithoutResource) WithoutResource {
	return a
}

func testInterfaceBare(a Bare) Bare {
	return a
}

func testInterfaceFull(a Full) Full {
	return a
}

func (s *JidXMPPSuite) Test_interfaceImplementations(c *C) {
	// There are no assertions here - if it compiles, we are fine.
	testInterfaceAny(Domain{"bla.com"})
	testInterfaceAny(ParseBare("bla@bla.com"))
	testInterfaceAny(ParseFull("bla@bla.com/blu"))
	testInterfaceAny(domainWithResource{NewDomain("bla.com"), NewResource("blu")})

	testInterfaceWithResource(ParseFull("bla@bla.com/blu"))
	testInterfaceWithResource(domainWithResource{NewDomain("bla.com"), NewResource("blu")})

	testInterfaceWithoutResource(Domain{"bla.com"})
	testInterfaceWithoutResource(ParseBare("bla@bla.com"))

	testInterfaceBare(ParseBare("bla@bla.com"))

	testInterfaceFull(ParseFull("bla@bla.com/blu"))
}
