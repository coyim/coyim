package net

import (
	"testing"

	"github.com/miekg/dns"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DnsXmppSuite struct{}

var _ = Suite(&DnsXmppSuite{})

func (s *DnsXmppSuite) Test_createCName_createsAValidCnameForAService(c *C) {
	ret := createCName("foo", "bar", "bax.com")
	c.Assert(ret, Equals, "_foo._bar.bax.com.")
}

func (s *DnsXmppSuite) Test_convertAnswerToSRV_returnsNilForNonSRVEntry(c *C) {
	cn := new(dns.CNAME)
	res := convertAnswerToSRV(cn)
	c.Assert(res, IsNil)
}

func (s *DnsXmppSuite) Test_convertAnswerToSRV_returnsAValidNetSRV(c *C) {
	srv := new(dns.SRV)
	srv.Target = "foo.com"
	srv.Port = 123
	srv.Priority = 5
	srv.Weight = 42
	res := convertAnswerToSRV(srv)
	c.Assert(res, Not(IsNil))
	c.Assert(res.Target, Equals, "foo.com")
	c.Assert(res.Port, Equals, uint16(123))
	c.Assert(res.Priority, Equals, uint16(5))
	c.Assert(res.Weight, Equals, uint16(42))
}

func (s *DnsXmppSuite) Test_convertAnswersToSRV_convertsAnswers(c *C) {
	cn := new(dns.CNAME)
	srv := new(dns.SRV)
	srv.Target = "foo2.com"

	in := make([]dns.RR, 2)
	in[0] = cn
	in[1] = srv
	res := convertAnswersToSRV(in)
	c.Assert(res[0], IsNil)
	c.Assert(res[1].Target, Equals, "foo2.com")
}

func (s *DnsXmppSuite) Test_msgSRV_createsMessage(c *C) {
	res := msgSRV("foo.com")
	c.Assert(res.Question[0].Name, Equals, "foo.com")
	c.Assert(res.Question[0].Qtype, Equals, dns.TypeSRV)
}

func (s *DnsXmppSuite) Test_convertAnswersToSRV_sortsByPriority(c *C) {
	srv1 := &dns.SRV{
		Target:   "foo1.com",
		Priority: 5,
		Weight:   1,
	}
	srv2 := &dns.SRV{
		Target:   "foo2.com",
		Priority: 3,
		Weight:   1,
	}
	srv3 := &dns.SRV{
		Target:   "foo3.com",
		Priority: 6,
		Weight:   1,
	}
	srv4 := &dns.SRV{
		Target:   "foo4.com",
		Priority: 1,
		Weight:   1,
	}

	in := []dns.RR{
		srv1,
		srv2,
		srv3,
		srv4,
	}
	res := convertAnswersToSRV(in)
	c.Assert(res[0].Target, Equals, "foo4.com")
	c.Assert(res[1].Target, Equals, "foo2.com")
	c.Assert(res[2].Target, Equals, "foo1.com")
	c.Assert(res[3].Target, Equals, "foo3.com")
}

func (s *DnsXmppSuite) Test_convertAnswersToSRV_sortsByWeightIfPriotityIsTheSame(c *C) {
	srv1 := &dns.SRV{
		Target:   "foo1.com",
		Priority: 1,
		Weight:   5,
	}
	srv2 := &dns.SRV{
		Target:   "foo2.com",
		Priority: 1,
		Weight:   3,
	}
	srv3 := &dns.SRV{
		Target:   "foo3.com",
		Priority: 1,
		Weight:   6,
	}
	srv4 := &dns.SRV{
		Target:   "foo4.com",
		Priority: 1,
		Weight:   1,
	}

	in := []dns.RR{
		srv1,
		srv2,
		srv3,
		srv4,
	}
	res := convertAnswersToSRV(in)
	c.Assert(res[0].Target, Equals, "foo4.com")
	c.Assert(res[1].Target, Equals, "foo2.com")
	c.Assert(res[2].Target, Equals, "foo1.com")
	c.Assert(res[3].Target, Equals, "foo3.com")
}
