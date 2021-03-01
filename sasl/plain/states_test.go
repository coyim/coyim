package plain

import (
	. "gopkg.in/check.v1"
)

func (s *SASLPlain) Test_finished_challenge(c *C) {
	f := finished{}
	ret, tok, err := f.challenge("one", "two")
	c.Assert(ret, DeepEquals, ret)
	c.Assert(tok, IsNil)
	c.Assert(err, IsNil)
}
