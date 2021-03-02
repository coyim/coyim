// +build darwin

package main

import (
	"github.com/coyim/coyim/gui"
	. "gopkg.in/check.v1"
)

type HooksSuite struct{}

var _ = Suite(&HooksSuite{})

func (s *HooksSuite) Test_hooks(c *C) {
	h := hooks()
	c.Assert(h, Not(IsNil))
	c.Assert(hooks(), Not(FitsTypeOf), &gui.NoHooks{})
}
