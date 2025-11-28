package importer

import (
	"os"
	"path"
	"runtime"

	. "gopkg.in/check.v1"
)

func testResourceFilename(s string) string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), s)
}

func tempFile(c *C) *os.File {
	tmpfile, e := os.CreateTemp(c.MkDir(), "coyim-config-18-")
	c.Assert(e, IsNil)
	return tmpfile
}
