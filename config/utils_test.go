package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) Test_randomString_returnsErrorIfReadingFails(c *C) {
	origReadFullFunc := ioReadFull
	defer func() {
		ioReadFull = origReadFullFunc
	}()

	ioReadFull = func(io.Reader, []byte) (int, error) {
		return 0, errors.New("something wrong")
	}

	e := randomString(nil)
	c.Assert(e, ErrorMatches, "something wrong")
}

func (s *UtilsSuite) Test_FindFile_returnsNothingIfNoPlacesGiven(c *C) {
	res, ok := FindFile(nil)
	c.Assert(res, Equals, "")
	c.Assert(ok, Equals, false)
}

func (s *UtilsSuite) Test_FindFile_returnsTheFirstFileFound(c *C) {
	tmpfile1, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile1.Name())

	tmpfile2, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile2.Name())

	res, ok := FindFile([]string{tmpfile1.Name(), tmpfile2.Name()})
	c.Assert(res, Equals, tmpfile1.Name())
	c.Assert(ok, Equals, true)
}

func (s *UtilsSuite) Test_FindFile_returnsTheSecondFileIfTheFirstDoesntExist(c *C) {
	tmpfile1, _ := ioutil.TempFile("", "")
	os.Remove(tmpfile1.Name())

	tmpfile2, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile2.Name())

	res, ok := FindFile([]string{tmpfile1.Name(), tmpfile2.Name()})
	c.Assert(res, Equals, tmpfile2.Name())
	c.Assert(ok, Equals, true)
}

func (s *UtilsSuite) Test_XdgCacheDir(c *C) {
	home := os.Getenv("HOME")
	oldEnv := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		os.Setenv("XDG_CACHE_HOME", oldEnv)
	}()
	os.Setenv("XDG_CACHE_HOME", "")

	dir := XdgCacheDir()
	c.Assert(dir, Equals, filepath.Join(home, ".cache"))
}

func (s *UtilsSuite) Test_XdgDataHome(c *C) {
	home := os.Getenv("HOME")
	oldEnv := os.Getenv("XDG_DATA_HOME")
	defer func() {
		os.Setenv("XDG_DATA_HOME", oldEnv)
	}()
	os.Setenv("XDG_DATA_HOME", "")

	dir := XdgDataHome()
	c.Assert(dir, Equals, filepath.Join(home, ".local/share"))
}

func (s *UtilsSuite) Test_XdgDataDirs(c *C) {
	oldEnv := os.Getenv("XDG_DATA_DIRS")
	defer func() {
		os.Setenv("XDG_DATA_DIRS", oldEnv)
	}()
	os.Setenv("XDG_DATA_DIRS", "one:two: three")

	dirs := XdgDataDirs()
	c.Assert(dirs, DeepEquals, []string{"one", "two", " three"})
}
