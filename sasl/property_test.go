package sasl

import . "gopkg.in/check.v1"

type PropertySuite struct{}

var _ = Suite(&PropertySuite{})

func (s *PropertySuite) Test_PropertyMissingError_reportsCorrectError(c *C) {
	result := PropertyMissingError{66}.Error()
	c.Assert(result, DeepEquals, "missing property 'B'")
}
