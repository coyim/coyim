package importer

import (
	"io/ioutil"
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
	tmpfile, _ := ioutil.TempFile(c.MkDir(), "")
	return tmpfile
}
