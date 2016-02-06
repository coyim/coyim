package sasl

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

type SASLSuite struct{}

var _ = Suite(&SASLSuite{})

type testMechanism struct {
	s Session
}

func (tm *testMechanism) NewClient() Session {
	return tm.s
}

func (s *SASLSuite) Test_RegisterMechanism_registeringTheSameMechanismTwiceReturnsAnError(c *C) {
	val := RegisterMechanism("foo", &testMechanism{nil})
	val2 := RegisterMechanism("foo", &testMechanism{nil})
	c.Assert(val, IsNil)
	c.Assert(val2, Equals, ErrMechanismAlreadyRegistered)
}

func (s *SASLSuite) Test_ClientSupport_returnsFalseForUnsupportedMechansim(c *C) {
	val := ClientSupport("foo2")
	c.Assert(val, Equals, false)
}

func (s *SASLSuite) Test_ClientSupport_returnsFalseForRegisteredMechanism(c *C) {
	RegisterMechanism("foo3", &testMechanism{nil})
	val := ClientSupport("foo3")
	c.Assert(val, Equals, true)
}

func (s *SASLSuite) Test_NewClient_returnsErrorForUnsupportedMechanism(c *C) {
	ss, err := NewClient("foo4")
	c.Assert(ss, IsNil)
	c.Assert(err, Equals, ErrUnsupportedMechanism)
}

type testSession struct{}

func (*testSession) SetProperty(Property, string) error {
	return nil
}

func (*testSession) Step(Token) (Token, error) {
	return nil, nil
}

func (*testSession) NeedsMore() bool {
	return false
}

func (s *SASLSuite) Test_NewClient_returnsTheResultOfNewClientForSupportedMechanism(c *C) {
	sess := &testSession{}
	RegisterMechanism("foo5", &testMechanism{sess})
	ss, err := NewClient("foo5")
	c.Assert(err, IsNil)
	c.Assert(ss, Equals, sess)
}
