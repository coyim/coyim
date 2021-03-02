package importer

import (
	"io/ioutil"
	"os"

	. "gopkg.in/check.v1"
)

type PidginLikeSuite struct{}

var _ = Suite(&PidginLikeSuite{})

func (s *PidginLikeSuite) Test_ImportKeysFromPidginStyle_failsWithBadFile(c *C) {
	res, ok := ImportKeysFromPidginStyle("file-that-hopefully-doesn't-exist", nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportKeysFromPidginStyle_failsWithBadContent(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte("hopefully something that isn't valid"))

	res, ok := ImportKeysFromPidginStyle(tmpfile.Name(), nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importAccountsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importAccountsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importAccountsPidginStyle_failsWithBadContent(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte("hopefully something that isn't valid"))

	res, ok := importAccountsPidginStyle(tmpfile.Name())
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importPeerPrefsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importPeerPrefsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importPeerPrefsPidginStyle_failsWithBadContent(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte("hopefully something that isn't valid"))

	res, ok := importPeerPrefsPidginStyle(tmpfile.Name())
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importGlobalPrefsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importGlobalPrefsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importGlobalPrefsPidginStyle_failsWithBadContent(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte("hopefully something that isn't valid"))

	res, ok := importGlobalPrefsPidginStyle(tmpfile.Name())
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_failsWithBadFile(c *C) {
	res, ok := ImportFingerprintsFromPidginStyle("file-that-hopefully-doesn't-exist", nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_failsWithBadContent(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte("hopefully something that isn't valid"))

	res, ok := ImportFingerprintsFromPidginStyle(tmpfile.Name(), nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_ignoresBadFingerprintHEx(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte("line\tsomething\tbla\tqqqq"))

	res, ok := ImportFingerprintsFromPidginStyle(tmpfile.Name(), func(string) bool { return true })
	c.Assert(res, HasLen, 0)
	c.Assert(ok, Equals, true)
}
