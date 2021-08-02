package importer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	. "gopkg.in/check.v1"
)

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) Test_ifExists_returnsTheValueAndTheFileIfItExists(c *C) {
	tmpfile, ex := ioutil.TempFile("", "")
	c.Assert(ex, IsNil)
	logPotentialError(c, tmpfile.Close())

	defer func() {
		logPotentialError(c, os.Remove(tmpfile.Name()))
	}()

	res := ifExists([]string{"foo", "bar"}, tmpfile.Name())

	c.Assert(res, DeepEquals, []string{"foo", "bar", tmpfile.Name()})
}

func (s *UtilsSuite) Test_ifExists_returnsTheValueButNothingElseIfItsADir(c *C) {
	dir := c.MkDir()

	res := ifExists([]string{"foo", "bar"}, dir)

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExists_returnsTheValueButNothingElseIfDoesntExist(c *C) {
	res := ifExists([]string{"foo", "bar"}, "filename_that_doesnt_exist_hopefully.foo")

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfFile(c *C) {
	tmpfile, ex := ioutil.TempFile("", "")
	c.Assert(ex, IsNil)
	logPotentialError(c, tmpfile.Close())

	defer func() {
		logPotentialError(c, os.Remove(tmpfile.Name()))
	}()

	res := ifExistsDir([]string{"foo", "bar"}, tmpfile.Name())

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfDoesntExists(c *C) {
	res := ifExistsDir([]string{"foo", "bar"}, "bla-dir-that-hopefully-doesnt-exist")

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfReadingDirFails(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer func() {
		makeDirectoryAccessible(dir)
		logPotentialError(c, os.RemoveAll(dir))
	}()

	logPotentialError(c, os.Mkdir(filepath.Join(dir, "foo"), 0755))
	ff, ex := os.Create(filepath.Join(dir, "hello.conf"))
	c.Assert(ex, IsNil)
	logPotentialError(c, ff.Close())
	ff, ex = os.Create(filepath.Join(dir, "goodbye.conf"))
	c.Assert(ex, IsNil)
	logPotentialError(c, ff.Close())

	makeDirectoryInaccessible(dir)

	res := ifExistsDir([]string{"foo", "bar"}, dir)

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueAndFilesInside(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer func() {
		logPotentialError(c, os.RemoveAll(dir))
	}()

	logPotentialError(c, os.Mkdir(filepath.Join(dir, "foo"), 0755))
	ff, ex := os.Create(filepath.Join(dir, "hello.conf"))
	c.Assert(ex, IsNil)
	logPotentialError(c, ff.Close())
	ff, ex = os.Create(filepath.Join(dir, "goodbye.conf"))
	c.Assert(ex, IsNil)
	logPotentialError(c, ff.Close())

	res := ifExistsDir([]string{"foo", "bar"}, dir)

	sort.Strings(res)

	c.Assert(res, DeepEquals, []string{filepath.Join(dir, "goodbye.conf"), filepath.Join(dir, "hello.conf"), "bar", "foo"})
}
