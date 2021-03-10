package importer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/coyim/coyim/config"
	. "gopkg.in/check.v1"
)

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) Test_ifExists_returnsTheValueAndTheFileIfItExists(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	res := ifExists([]string{"foo", "bar"}, tmpfile.Name())

	c.Assert(res, DeepEquals, []string{"foo", "bar", tmpfile.Name()})
}

func (s *UtilsSuite) Test_ifExists_returnsTheValueButNothingElseIfItsADir(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	res := ifExists([]string{"foo", "bar"}, dir)

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExists_returnsTheValueButNothingElseIfDoesntExist(c *C) {
	res := ifExists([]string{"foo", "bar"}, "filename_that_doesnt_exist_hopefully.foo")

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfFile(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	res := ifExistsDir([]string{"foo", "bar"}, tmpfile.Name())

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfDoesntExists(c *C) {
	res := ifExistsDir([]string{"foo", "bar"}, "bla-dir-that-hopefully-doesnt-exist")

	c.Assert(res, DeepEquals, []string{"foo", "bar"})
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfReadingDirFails(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	os.Mkdir(filepath.Join(dir, "foo"), 0755)
	os.Create(filepath.Join(dir, "hello.conf"))
	os.Create(filepath.Join(dir, "goodbye.conf"))

	os.Chmod(dir, 0000)

	defaultDirs := []string{"foo", "bar"}
	res := ifExistsDir(defaultDirs, dir)

	// Windows manage directory permissions in other way
	// We should remove this when we find a way to make a directory in windows unaccessible
	// The next is an article in which the file perms in Windows with Go
	// are described: https://medium.com/@MichalPristas/go-and-file-perms-on-windows-3c944d55dd44
	if config.IsWindows() {
		defaultDirs = append(defaultDirs, filepath.Join(dir, "goodbye.conf"), filepath.Join(dir, "hello.conf"))
	}

	c.Assert(res, DeepEquals, defaultDirs)
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueAndFilesInside(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	os.Mkdir(filepath.Join(dir, "foo"), 0755)
	os.Create(filepath.Join(dir, "hello.conf"))
	os.Create(filepath.Join(dir, "goodbye.conf"))

	res := ifExistsDir([]string{"foo", "bar"}, dir)

	sort.Sort(sort.StringSlice(res))

	c.Assert(res, DeepEquals, []string{filepath.Join(dir, "goodbye.conf"), filepath.Join(dir, "hello.conf"), "bar", "foo"})
}
