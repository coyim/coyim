package importer

import (
	"os"
	"path/filepath"
	"sort"

	. "gopkg.in/check.v1"
)

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) Test_ifExists_returnsTheValueAndTheFileIfItExists(c *C) {
	tmpfile := tempFile(c)
	_ = tmpfile.Close()

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
	tmpfile := tempFile(c)
	_ = tmpfile.Close()

	res := ifExistsDir([]string{"foo", "bar"}, tmpfile.Name())

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfDoesntExists(c *C) {
	res := ifExistsDir([]string{"foo", "bar"}, "bla-dir-that-hopefully-doesnt-exist")

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfReadingDirFails(c *C) {
	dir := c.MkDir()
	defer makeDirectoryAccessible(dir)

	_ = os.Mkdir(filepath.Join(dir, "foo"), 0755)

	f1, _ := os.Create(filepath.Join(dir, "hello.conf"))
	f2, _ := os.Create(filepath.Join(dir, "goodbye.conf"))

	_ = f1.Close()
	_ = f2.Close()

	makeDirectoryInaccessible(dir)

	res := ifExistsDir([]string{"foo", "bar"}, dir)

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueAndFilesInside(c *C) {
	dir := c.MkDir()

	_ = os.Mkdir(filepath.Join(dir, "foo"), 0755)

	f1, _ := os.Create(filepath.Join(dir, "hello.conf"))
	f2, _ := os.Create(filepath.Join(dir, "goodbye.conf"))

	_ = f1.Close()
	_ = f2.Close()

	res := ifExistsDir([]string{"foo", "bar"}, dir)

	sort.Strings(res)

	c.Assert(res, DeepEquals, []string{filepath.Join(dir, "goodbye.conf"), filepath.Join(dir, "hello.conf"), "bar", "foo"})
}
