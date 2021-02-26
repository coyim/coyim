package access

import (
	. "gopkg.in/check.v1"
)

type DataSuite struct{}

var _ = Suite(&DataSuite{})

func (s *DataSuite) Test_OfflineError_Error_returnsMessage(c *C) {
	x := &OfflineError{"something"}

	c.Assert(x.Error(), Equals, "something")
}
