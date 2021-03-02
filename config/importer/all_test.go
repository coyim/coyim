package importer

import (
	"io/ioutil"
	"os"

	. "gopkg.in/check.v1"
)

type AllSuite struct{}

var _ = Suite(&AllSuite{})

func (s *AllSuite) Test_TryImportAll(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", origHome)
	}()
	os.Setenv("HOME", dir)

	res := TryImportAll()
	c.Assert(res["Adium"], HasLen, 0)
	c.Assert(res["Gajim"], HasLen, 0)
	c.Assert(res["Pidgin"], HasLen, 0)
	c.Assert(res["xmpp-client"], HasLen, 0)
}
