package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"

	"github.com/prashantv/gostub"
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
	tmpfile1, e := ioutil.TempFile("", "")
	c.Assert(e, IsNil)
	defer func() {
		logPotentialError(c, os.Remove(tmpfile1.Name()))
	}()

	tmpfile2, e2 := ioutil.TempFile("", "")
	c.Assert(e2, IsNil)
	defer func() {
		logPotentialError(c, os.Remove(tmpfile2.Name()))
	}()

	res, ok := FindFile([]string{tmpfile1.Name(), tmpfile2.Name()})
	c.Assert(res, Equals, tmpfile1.Name())
	c.Assert(ok, Equals, true)
}

func (s *UtilsSuite) Test_FindFile_returnsTheSecondFileIfTheFirstDoesntExist(c *C) {
	tmpfile1, e := ioutil.TempFile("", "")
	c.Assert(e, IsNil)
	logPotentialError(c, tmpfile1.Close())
	logPotentialError(c, os.Remove(tmpfile1.Name()))

	tmpfile2, e2 := ioutil.TempFile("", "")
	c.Assert(e2, IsNil)
	logPotentialError(c, tmpfile2.Close())
	defer func() {
		logPotentialError(c, os.Remove(tmpfile2.Name()))
	}()

	res, ok := FindFile([]string{tmpfile1.Name(), tmpfile2.Name()})
	c.Assert(res, Equals, tmpfile2.Name())
	c.Assert(ok, Equals, true)
}

func (s *UtilsSuite) Test_XdgCacheDir(c *C) {
	home := os.Getenv("HOME")
	defer gostub.New().SetEnv("XDG_CACHE_HOME", "").Reset()

	dir := XdgCacheDir()
	c.Assert(dir, Equals, filepath.Join(home, ".cache"))
}

func (s *UtilsSuite) Test_XdgDataHome(c *C) {
	home := os.Getenv("HOME")
	defer gostub.New().SetEnv("XDG_DATA_HOME", "").Reset()

	dir := XdgDataHome()
	c.Assert(dir, Equals, filepath.Join(home, ".local/share"))
}

func (s *UtilsSuite) Test_XdgDataDirs(c *C) {
	defer gostub.New().SetEnv("XDG_DATA_DIRS", "one:two: three").Reset()

	dirs := XdgDataDirs()
	c.Assert(dirs, DeepEquals, []string{"one", "two", " three"})
}
