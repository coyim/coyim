package main

import . "gopkg.in/check.v1"

type MainSuite struct{}

var _ = Suite(&MainSuite{})

func (s *MainSuite) Test_mainInit_setsVersionFromCommit(c *C) {
	origBuildTag := BuildTag
	origBuildCommit := BuildCommit
	origCoyimVersion := coyimVersion

	defer func() {
		BuildTag = origBuildTag
		BuildCommit = origBuildCommit
		coyimVersion = origCoyimVersion
	}()

	BuildTag = "(no tag)"
	BuildCommit = "hello abc"
	mainInit()
	c.Assert(coyimVersion, Equals, "hello abc")

	BuildTag = ""
	BuildCommit = "hello def"
	mainInit()
	c.Assert(coyimVersion, Equals, "hello def")

	BuildTag = "v42"
	BuildCommit = "hello def"
	mainInit()
	c.Assert(coyimVersion, Equals, "v42")
}
