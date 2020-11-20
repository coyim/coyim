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

func (s *JidXMPPSuite) Test_NewResource(c *C) {
	c.Assert(NewResource("foo").Valid(), Equals, true)
	c.Assert(NewResource("a\u06DDb").Valid(), Equals, false)
}

func (s *JidXMPPSuite) Test_NewBare(c *C) {
	c.Assert(NewBare(NewLocal("hello"), NewDomain("goodbye.com")).String(), Equals, "hello@goodbye.com")
}

func (s *JidXMPPSuite) Test_NewBareFromStrings(c *C) {
	c.Assert(NewBareFromStrings("", "").String(), Equals, "@")
	c.Assert(NewBareFromStrings("hello", "goodbye.com").String(), Equals, "hello@goodbye.com")
	c.Assert(NewBareFromStrings("hello", "").String(), Equals, "@")
	c.Assert(NewBareFromStrings("@", "").String(), Equals, "@")
	c.Assert(NewBareFromStrings("#", "#").String(), Equals, "@")
}

func (s *JidXMPPSuite) Test_NewFull(c *C) {
	c.Assert(NewFull(NewLocal("hello"), NewDomain("goodbye.com"), NewResource("somewhere")), DeepEquals,
		full{
			l: NewLocal("hello"),
			d: NewDomain("goodbye.com"),
			r: NewResource("somewhere"),
		},
	)
}

func (s *JidXMPPSuite) Test_ParseDomain(c *C) {
	c.Assert(ParseDomain("foo@bar.com/res"), Equals, NewDomain("bar.com"))
}

func (s *JidXMPPSuite) Test_Domain_Host(c *C) {
	c.Assert(NewDomain("bar.com").Host(), Equals, NewDomain("bar.com"))
}

func (s *JidXMPPSuite) Test_Domain_PotentialResource(c *C) {
	c.Assert(NewDomain("bar.com").PotentialResource(), Equals, Resource{""})
}

func (s *JidXMPPSuite) Test_Domain_PotentialSplit(c *C) {
	l, r := NewDomain("bar.com").PotentialSplit()
	c.Assert(l, Equals, NewDomain("bar.com"))
	c.Assert(r, Equals, Resource{""})
}

func (s *JidXMPPSuite) Test_bare_Host(c *C) {
	c.Assert(bare{NewLocal("foo"), NewDomain("bar.com")}.Host(), Equals, NewDomain("bar.com"))
}

func (s *JidXMPPSuite) Test_bare_PotentialResource(c *C) {
	c.Assert(bare{NewLocal("foo"), NewDomain("bar.com")}.PotentialResource(), Equals, Resource{""})
}

func (s *JidXMPPSuite) Test_bare_PotentialSplit(c *C) {
	l, r := bare{NewLocal("foo"), NewDomain("bar.com")}.PotentialSplit()
	c.Assert(l, Equals, bare{NewLocal("foo"), NewDomain("bar.com")})
	c.Assert(r, Equals, Resource{""})
}

func (s *JidXMPPSuite) Test_bare_Local(c *C) {
	c.Assert(bare{NewLocal("foo"), NewDomain("bar.com")}.Local(), Equals, NewLocal("foo"))
}

func (s *JidXMPPSuite) Test_bare_WithResource(c *C) {
	c.Assert(bare{NewLocal("foo"), NewDomain("bar.com")}.WithResource(NewResource("someone")).String(), Equals, "foo@bar.com/someone")
}

func (s *JidXMPPSuite) Test_bare_Bare(c *C) {
	c.Assert(bare{NewLocal("foo"), NewDomain("bar.com")}.Bare(), Equals, bare{NewLocal("foo"), NewDomain("bar.com")})
}

func (s *JidXMPPSuite) Test_full_Host(c *C) {
	c.Assert(full{NewLocal("foo"), NewDomain("bar.com"), NewResource("someone")}.Host(), Equals, NewDomain("bar.com"))
}

func (s *JidXMPPSuite) Test_full_String(c *C) {
	c.Assert(full{NewLocal("foo"), NewDomain("bar.com"), NewResource("someone")}.String(), Equals, "foo@bar.com/someone")
}

func (s *JidXMPPSuite) Test_full_WithResource(c *C) {
	c.Assert(full{NewLocal("foo"), NewDomain("bar.com"), NewResource("someone")}.WithResource(NewResource("elsewhere")).String(), Equals, "foo@bar.com/elsewhere")
}

func (s *JidXMPPSuite) Test_full_PotentialResource(c *C) {
	c.Assert(full{NewLocal("foo"), NewDomain("bar.com"), NewResource("someone")}.PotentialResource(), Equals, Resource{"someone"})
}

func (s *JidXMPPSuite) Test_full_PotentialSplit(c *C) {
	l, r := full{NewLocal("foo"), NewDomain("bar.com"), NewResource("someone")}.PotentialSplit()
	c.Assert(l, Equals, bare{NewLocal("foo"), NewDomain("bar.com")})
	c.Assert(r, Equals, NewResource("someone"))
}

func (s *JidXMPPSuite) Test_full_Local(c *C) {
	c.Assert(full{NewLocal("foo"), NewDomain("bar.com"), NewResource("someone")}.Local(), Equals, NewLocal("foo"))
}

func (s *JidXMPPSuite) Test_full_Bare(c *C) {
	c.Assert(full{NewLocal("foo"), NewDomain("bar.com"), NewResource("someone")}.Bare(), Equals, bare{NewLocal("foo"), NewDomain("bar.com")})
}

func (s *JidXMPPSuite) Test_domainWithResource_Host(c *C) {
	c.Assert(domainWithResource{NewDomain("bar.com"), NewResource("someone")}.Host(), Equals, NewDomain("bar.com"))
}

func (s *JidXMPPSuite) Test_domainWithResource_String(c *C) {
	c.Assert(domainWithResource{NewDomain("bar.com"), NewResource("someone")}.String(), Equals, "bar.com/someone")
}

func (s *JidXMPPSuite) Test_domainWithResource_MaybeWithResource(c *C) {
	c.Assert(domainWithResource{NewDomain("bar.com"), NewResource("someone")}.MaybeWithResource(NewResource("elsewhere")), Equals, domainWithResource{NewDomain("bar.com"), NewResource("elsewhere")})
}

func (s *JidXMPPSuite) Test_domainWithResource_WithResource(c *C) {
	c.Assert(domainWithResource{NewDomain("bar.com"), NewResource("someone")}.WithResource(NewResource("elsewhere")), Equals, domainWithResource{NewDomain("bar.com"), NewResource("elsewhere")})
}

func (s *JidXMPPSuite) Test_domainWithResource_PotentialResource(c *C) {
	c.Assert(domainWithResource{NewDomain("bar.com"), NewResource("someone")}.PotentialResource(), Equals, NewResource("someone"))
}

func (s *JidXMPPSuite) Test_domainWithResource_PotentialSplit(c *C) {
	l, r := domainWithResource{NewDomain("bar.com"), NewResource("someone")}.PotentialSplit()
	c.Assert(l, Equals, NewDomain("bar.com"))
	c.Assert(r, Equals, NewResource("someone"))
}

func (s *JidXMPPSuite) Test_domainWithResource_Split(c *C) {
	l, r := domainWithResource{NewDomain("bar.com"), NewResource("someone")}.Split()
	c.Assert(l, Equals, NewDomain("bar.com"))
	c.Assert(r, Equals, NewResource("someone"))
}

func (s *JidXMPPSuite) Test_Domain_WithResource(c *C) {
	c.Assert(NewDomain("foo.com").WithResource(NewResource("somewhere")), Equals, domainWithResource{NewDomain("foo.com"), NewResource("somewhere")})
}

func (s *JidXMPPSuite) Test_Domain_MaybeWithResource(c *C) {
	c.Assert(NewDomain("foo.com").MaybeWithResource(NewResource("somewhere")), Equals, domainWithResource{NewDomain("foo.com"), NewResource("somewhere")})
}

func (s *JidXMPPSuite) Test_MaybeLocal(c *C) {
	c.Assert(MaybeLocal(NewDomain("foo.com")), Equals, Local{""})
	c.Assert(MaybeLocal(bare{NewLocal("someone"), NewDomain("foo.com")}), Equals, Local{"someone"})
}

func (s *JidXMPPSuite) Test_WithAndWithout(c *C) {
	wr, wnr := WithAndWithout(Domain{"foo.bar"})
	c.Assert(wr, IsNil)
	c.Assert(wnr, Equals, Domain{"foo.bar"})

	wr, wnr = WithAndWithout(full{Local{"someone"}, Domain{"foo.bar"}, Resource{"bla"}})
	c.Assert(wr, Equals, full{Local{"someone"}, Domain{"foo.bar"}, Resource{"bla"}})
	c.Assert(wnr, Equals, bare{Local{"someone"}, Domain{"foo.bar"}})
}
