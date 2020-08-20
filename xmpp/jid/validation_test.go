package jid

import (
	. "gopkg.in/check.v1"
)

func (s *JidXMPPSuite) Test_ValidLocal(c *C) {
	c.Assert(ValidLocal(""), Equals, false)
	c.Assert(ValidLocal("a"), Equals, true)
	// this is 1023 characters long
	c.Assert(ValidLocal("abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefg"), Equals, true)
	// this is 1024 characters long
	c.Assert(ValidLocal("abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"), Equals, false)

	c.Assert(ValidLocal("a b"), Equals, false)
	c.Assert(ValidLocal("a2b"), Equals, true)
	c.Assert(ValidLocal("a@b"), Equals, false)
	c.Assert(ValidLocal("a\"b"), Equals, false)
	c.Assert(ValidLocal("a&b"), Equals, false)
	c.Assert(ValidLocal("a'b"), Equals, false)
	c.Assert(ValidLocal("a/b"), Equals, false)
	c.Assert(ValidLocal("a:b"), Equals, false)
	c.Assert(ValidLocal("a<b"), Equals, false)
	c.Assert(ValidLocal("a>b"), Equals, false)
	c.Assert(ValidLocal("a.com"), Equals, true)
}

func (s *JidXMPPSuite) Test_ValidDomain(c *C) {
	c.Assert(ValidDomain(""), Equals, false)
	c.Assert(ValidDomain("a"), Equals, true)
	// this is 1024 characters long
	c.Assert(ValidDomain("abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"), Equals, false)
	c.Assert(ValidDomain("10.0.1.3"), Equals, true)
	c.Assert(ValidDomain("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), Equals, true)
	c.Assert(ValidDomain("a b.com"), Equals, false)
	c.Assert(ValidDomain("a2b.com"), Equals, true)
	c.Assert(ValidDomain("a@b.com"), Equals, false)
	c.Assert(ValidDomain("a\"b.com"), Equals, false)
	c.Assert(ValidDomain("a&b.com"), Equals, false)
	c.Assert(ValidDomain("a'b.com"), Equals, false)
	c.Assert(ValidDomain("a/b.com"), Equals, false)
	c.Assert(ValidDomain("a:b.com"), Equals, false)
	c.Assert(ValidDomain("a<b.com"), Equals, false)
	c.Assert(ValidDomain("a>b.com"), Equals, false)
}

func (s *JidXMPPSuite) Test_ValidResource(c *C) {
	c.Assert(ValidResource(""), Equals, false)
	c.Assert(ValidResource("a"), Equals, true)
	// this is 1023 characters long
	c.Assert(ValidResource("abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefg"), Equals, true)
	// this is 1024 characters long
	c.Assert(ValidResource("abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"+
		"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"), Equals, false)

	c.Assert(ValidResource("a b"), Equals, true)
	c.Assert(ValidResource("a\u06DDb"), Equals, false)
	c.Assert(ValidResource("a2b"), Equals, true)
	c.Assert(ValidResource("a@b"), Equals, true)
	c.Assert(ValidResource("a\"b"), Equals, true)
	c.Assert(ValidResource("a&b"), Equals, true)
	c.Assert(ValidResource("a'b"), Equals, true)
	c.Assert(ValidResource("a/b"), Equals, true)
	c.Assert(ValidResource("a:b"), Equals, true)
	c.Assert(ValidResource("a<b"), Equals, true)
	c.Assert(ValidResource("a>b"), Equals, true)
	c.Assert(ValidResource("a.com"), Equals, true)
}

func (s *JidXMPPSuite) Test_ValidBareJID(c *C) {
	c.Assert(ValidBareJID("abc"), Equals, false)
	c.Assert(ValidBareJID("abc.com/resource"), Equals, false)
	c.Assert(ValidBareJID("foo@abc.com"), Equals, true)
	c.Assert(ValidBareJID("foo@abc.com/resource"), Equals, true)
	c.Assert(ValidBareJID("fo:o@abc.com"), Equals, false)
	c.Assert(ValidBareJID("foo@ab c.com"), Equals, false)
}

func (s *JidXMPPSuite) Test_ValidFullJID(c *C) {
	c.Assert(ValidFullJID("abc"), Equals, false)
	c.Assert(ValidFullJID("abc.com/resource"), Equals, false)
	c.Assert(ValidFullJID("foo@abc.com"), Equals, false)
	c.Assert(ValidFullJID("foo@abc.com/resource"), Equals, true)
	c.Assert(ValidFullJID("fo:o@abc.com"), Equals, false)
	c.Assert(ValidFullJID("foo@ab c.com"), Equals, false)
	c.Assert(ValidFullJID("fo:o@abc.com/resource"), Equals, false)
	c.Assert(ValidFullJID("foo@ab c.com/resource"), Equals, false)
}

func (s *JidXMPPSuite) Test_ValidDomainWithResource(c *C) {
	c.Assert(ValidDomainWithResource("abc"), Equals, false)
	c.Assert(ValidDomainWithResource("abc.com/resource"), Equals, true)
	c.Assert(ValidDomainWithResource("foo@abc.com"), Equals, false)
	c.Assert(ValidDomainWithResource("foo@abc.com/resource"), Equals, true)
	c.Assert(ValidDomainWithResource("fo:o@abc.com"), Equals, false)
	c.Assert(ValidDomainWithResource("foo@ab c.com"), Equals, false)
	c.Assert(ValidDomainWithResource("fo:o@abc.com/resource"), Equals, false)
	c.Assert(ValidDomainWithResource("foo@ab c.com/resource"), Equals, false)
}

func (s *JidXMPPSuite) Test_ValidJID(c *C) {
	c.Assert(ValidJID("abc"), Equals, true)
	c.Assert(ValidJID("abc.com/resource"), Equals, true)
	c.Assert(ValidJID("foo@abc.com"), Equals, true)
	c.Assert(ValidJID("foo@abc.com/resource"), Equals, true)
	c.Assert(ValidJID("fo:o@abc.com"), Equals, false)
	c.Assert(ValidJID("foo@ab c.com"), Equals, false)
	c.Assert(ValidJID("fo:o@abc.com/resource"), Equals, false)
	c.Assert(ValidJID("foo@ab c.com/resource"), Equals, false)
}
