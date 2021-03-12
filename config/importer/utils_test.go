package importer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"

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

func windowsIcaclsExec(dir, action, permissions string, done chan bool) {
	c := exec.Command("icacls", dir, action, permissions)
	c.Run()
	done <- true
}

func grantWindowsUserDirPermissions(dir string, done chan bool) {
	windowsIcaclsExec(dir, "/grant", fmt.Sprintf("%s:(RX,W)", os.Getenv("username")), done)
}

func denyWindowsUserDirPermissions(dir string, done chan bool) {
	windowsIcaclsExec(dir, "/deny", fmt.Sprintf("%s:(RX,W)", os.Getenv("username")), done)
}

func (s *UtilsSuite) Test_ifExistsDir_returnsTheValueButNothingElseIfReadingDirFails(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer func() {
		if runtime.GOOS == "windows" {
			done := make(chan bool)
			go grantWindowsUserDirPermissions(dir, done)
			<-done
		}
		os.RemoveAll(dir)
	}()

	os.Mkdir(filepath.Join(dir, "foo"), 0755)
	os.Create(filepath.Join(dir, "hello.conf"))
	os.Create(filepath.Join(dir, "goodbye.conf"))

	switch runtime.GOOS {
	case "windows":
		done := make(chan bool)
		go denyWindowsUserDirPermissions(dir, done)
		<-done
	default:
		os.Chmod(dir, 0000)
	}

	defaultDirs := []string{"foo", "bar"}
	res := ifExistsDir(defaultDirs, dir)

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
