package config

import . "gopkg.in/check.v1"

type AccountXmppSuite struct{}

var _ = Suite(&AccountXmppSuite{})

func (s *AccountXmppSuite) Test_Account_Is_recognizesJids(c *C) {
	a := &Account{Account: "hello@bar.com"}
	c.Check(a.Is("foo"), Equals, false)
	c.Check(a.Is("hello@bar.com"), Equals, true)
	c.Check(a.Is("hello@bar.com/foo"), Equals, true)
}

func (s *AccountXmppSuite) Test_Account_ShouldEncryptTo(c *C) {
	a := &Account{Account: "hello@bar.com", AlwaysEncrypt: false, AlwaysEncryptWith: []string{"one@foo.com", "two@foo.com"}}
	a2 := &Account{Account: "hello@bar.com", AlwaysEncrypt: true, AlwaysEncryptWith: []string{"one@foo.com", "two@foo.com"}}
	c.Check(a.ShouldEncryptTo("foo"), Equals, false)
	c.Check(a.ShouldEncryptTo("hello@bar.com"), Equals, false)
	c.Check(a.ShouldEncryptTo("one@foo.com"), Equals, true)
	c.Check(a.ShouldEncryptTo("two@foo.com"), Equals, true)
	c.Check(a.ShouldEncryptTo("two@foo.com/blarg"), Equals, true)
	c.Check(a2.ShouldEncryptTo("foo"), Equals, true)
	c.Check(a2.ShouldEncryptTo("hello@bar.com"), Equals, true)
}
