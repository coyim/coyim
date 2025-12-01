package config

import (
	"testing"

	"github.com/coyim/coyim/testutil"
	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { testutil.InitTest(t) }

func logPotentialError(c *check.C, e error) {
	c.Assert(e, check.IsNil)
}
