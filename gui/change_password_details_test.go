package gui

import (
	. "gopkg.in/check.v1"
)

type ChangePasswordDetailsSuite struct{}

var _ = Suite(&ChangePasswordDetailsSuite{})

func (s *ChangePasswordDetailsSuite) Test_validatePasswords_SuccesfulValidation(c *C) {
	err := validateNewPassword("pass", "pass")
	c.Assert(err, IsNil)
}

func (s *ChangePasswordDetailsSuite) Test_validatePasswords_FailedMistmatch(c *C) {
	err := validateNewPassword("passone", "passtwo")
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "The passwords do not match")
}

func (s *ChangePasswordDetailsSuite) Test_validatePasswords_FailedZeroLength(c *C) {
	err := validateNewPassword("", "")
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "The field can't be empty")
}
