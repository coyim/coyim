package gui

import (
	. "gopkg.in/check.v1"
)

type AccountDetailsSuite struct{}

var _ = Suite(&AccountDetailsSuite{})

func (s *AccountDetailsSuite) Test_validatePasswords_SuccesfulValidation(c *C) {
	err := validatePasswords("pass", "pass")
	c.Assert(err, IsNil)
}

func (s *AccountDetailsSuite) Test_validatePasswords_SuccesfulWithSpaces(c *C) {
	err := validatePasswords(" pass", "pass")
	c.Assert(err, IsNil)
}

func (s *AccountDetailsSuite) Test_validatePasswords_FailedMistmatch(c *C) {
	err := validatePasswords("passone", "passtwo")
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "The passwords do not match")
}

func (s *AccountDetailsSuite) Test_validatePasswords_FailedZeroLength(c *C) {
	err := validatePasswords("", "")
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "The password can't be empty")
}
