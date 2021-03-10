package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	. "gopkg.in/check.v1"
)

type FileSuite struct{}

var _ = Suite(&FileSuite{})

func isNotWindows() bool {
	return runtime.GOOS != "windows"
}

func (s *FileSuite) Test_ensureDir_createsUnknownDirectories(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	fullName := filepath.Join(dir, "foo", "bar")

	ensureDir(fullName, 0755)

	st, e := os.Stat(filepath.Join(dir, "foo"))
	c.Assert(e, IsNil)
	c.Assert(st.IsDir(), Equals, true)
	if isNotWindows() {
		c.Assert(int(st.Mode().Perm()), Equals, int(0755))
	}

	st, e = os.Stat(filepath.Join(dir, "foo", "bar"))
	c.Assert(e, IsNil)
	c.Assert(st.IsDir(), Equals, true)
	if isNotWindows() {
		c.Assert(int(st.Mode().Perm()), Equals, int(0755))
	}
}

func (s *FileSuite) Test_findConfigFile_returnsBasicPathIfNoFileFound(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origSystemConfigDir := SystemConfigDir
	defer func() {
		SystemConfigDir = origSystemConfigDir
	}()

	SystemConfigDir = func() string { return dir }

	res := findConfigFile("")
	c.Assert(res, Equals, filepath.Join(dir, "coyim", "accounts.json"))
}

func (s *FileSuite) Test_findConfigFile_returnsTheFilanameIfGiven(c *C) {
	res := findConfigFile("somewhere.json")
	c.Assert(res, Equals, "somewhere.json")
}

func (s *FileSuite) Test_findConfigFile_returnsEncryptedFileIfExists(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origSystemConfigDir := SystemConfigDir
	defer func() {
		SystemConfigDir = origSystemConfigDir
	}()

	SystemConfigDir = func() string { return dir }

	ensureDir(filepath.Join(dir, "coyim"), 0700)
	ioutil.WriteFile(filepath.Join(dir, "coyim", "accounts.json.enc"), []byte("hello"), 0666)

	res := findConfigFile("")
	c.Assert(res, Equals, filepath.Join(dir, "coyim", "accounts.json.enc"))
}

func (s *FileSuite) Test_findConfigFile_returnsEncryptedBackupFileIfExists(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origSystemConfigDir := SystemConfigDir
	defer func() {
		SystemConfigDir = origSystemConfigDir
	}()

	SystemConfigDir = func() string { return dir }

	ensureDir(filepath.Join(dir, "coyim"), 0700)
	ioutil.WriteFile(filepath.Join(dir, "coyim", "accounts.json.enc.000~"), []byte("hello"), 0666)

	res := findConfigFile("")
	c.Assert(res, Equals, filepath.Join(dir, "coyim", "accounts.json.enc"))
}

func (s *FileSuite) Test_readFileOrTemporaryBackup_readsBackupFileIfOriginalFileIsEmpty(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	ioutil.WriteFile(filepath.Join(dir, "accounts.json"), []byte(""), 0666)
	ioutil.WriteFile(filepath.Join(dir, "accounts.json.000~"), []byte("hello"), 0666)

	data, e := readFileOrTemporaryBackup(filepath.Join(dir, "accounts.json"))
	c.Assert(e, IsNil)
	c.Assert(string(data), Equals, "hello")
}

func (s *FileSuite) Test_readFileOrTemporaryBackup_readsOriginalFileEvenIfBackupExists(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	ioutil.WriteFile(filepath.Join(dir, "accounts.json"), []byte("who is there?"), 0666)
	ioutil.WriteFile(filepath.Join(dir, "accounts.json.000~"), []byte("hello"), 0666)

	data, e := readFileOrTemporaryBackup(filepath.Join(dir, "accounts.json"))
	c.Assert(e, IsNil)
	c.Assert(string(data), Equals, "who is there?")
}

func (s *FileSuite) Test_safeWrite_doesntAllowWritingOfTooLittleData(c *C) {
	e := safeWrite("somewhere.conf", []byte("123456"), 0700)
	c.Assert(e, ErrorMatches, "data amount too small.*")
}

func (s *FileSuite) Test_safeWrite_removesBackupFileIfItExists(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	ioutil.WriteFile(filepath.Join(dir, "accounts.json.backup.000~"), []byte("backup content"), 0666)
	e := safeWrite(filepath.Join(dir, "accounts.json"), []byte("12345678910111213"), 0700)
	c.Assert(e, IsNil)

	c.Assert(filepath.Join(dir, "accounts.json.backup.000~"), Not(FileExists))
}

type fileExistsChecker struct{}

var FileExists Checker = &fileExistsChecker{}

func (*fileExistsChecker) Info() *CheckerInfo {
	return &CheckerInfo{Name: "FileExists", Params: []string{"filename"}}
}

func (*fileExistsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	filename := fmt.Sprintf("%s", params[0])
	if fileExists(filename) {
		return true, ""
	}

	return false, "file does not exists"
}

func (s *FileSuite) Test_safeWrite_savesABackupFileIfPreviousFileExists(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	ioutil.WriteFile(filepath.Join(dir, "accounts.json"), []byte("previous content"), 0666)
	e := safeWrite(filepath.Join(dir, "accounts.json"), []byte("12345678910111213"), 0700)
	c.Assert(e, IsNil)

	c.Assert(filepath.Join(dir, "accounts.json.backup.000~"), FileExists)
}

func (s *FileSuite) Test_safeWrite_filesIfWeCantSaveTheBackupFile(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origOsRename := osRename
	defer func() {
		osRename = origOsRename
	}()

	osRename = func(string, string) error {
		return errors.New("so wroooooong")
	}

	ioutil.WriteFile(filepath.Join(dir, "accounts.json"), []byte("previous content"), 0666)
	e := safeWrite(filepath.Join(dir, "accounts.json"), []byte("12345678910111213"), 0700)
	c.Assert(e, ErrorMatches, "so wro+ng")
}

func (s *FileSuite) Test_safeWrite_failsIfImpossibleToWriteFile(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	ioutil.WriteFile(filepath.Join(dir, "accounts.json.000~"), []byte("previous content"), 0444)
	e := safeWrite(filepath.Join(dir, "accounts.json"), []byte("12345678910111213"), 0700)
	c.Assert(e, ErrorMatches, ".*denied.*")
}
