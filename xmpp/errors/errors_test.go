package errors

import (
	"errors"

	. "gopkg.in/check.v1"
)

type ErrorsSuite struct{}

var _ = Suite(&ErrorsSuite{})

func (s *ErrorsSuite) Test_ErrFailedToConnect(c *C) {
	e := CreateErrFailedToConnect("somewhere:84", errors.New("woot"))
	c.Assert(e.Error(), Equals, "Failed to connect to somewhere:84: woot")
}
