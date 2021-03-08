package session

import (
	. "gopkg.in/check.v1"
)

type PresenceSuite struct{}

var _ = Suite(&PresenceSuite{})

func (s *PresenceSuite) Test_session_AutoApprove(c *C) {
	appr := map[string]bool{}
	sess := &session{autoApproves: appr}

	sess.AutoApprove("someone@where.org")

	c.Assert(appr["someone@where.org"], Equals, true)
}
