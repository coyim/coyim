package sasl

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SASLSuite struct{}

var _ = Suite(&SASLSuite{})
