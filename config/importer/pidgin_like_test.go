package importer

import (
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
	tmpFile := s.tempFile([]byte("hopefully something that isn't valid"), c)
	_ = tmpFile.Close()

	res, ok := ImportKeysFromPidginStyle(tmpFile.Name(), nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importAccountsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importAccountsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importAccountsPidginStyle_failsWithBadContent(c *C) {
	tmpFile := s.tempFile([]byte("hopefully something that isn't valid"), c)
	_ = tmpFile.Close()

	res, ok := importAccountsPidginStyle(tmpFile.Name())
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importPeerPrefsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importPeerPrefsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importPeerPrefsPidginStyle_failsWithBadContent(c *C) {
	tmpFile := s.tempFile([]byte("hopefully something that isn't valid"), c)
	_ = tmpFile.Close()

	res, ok := importPeerPrefsPidginStyle(tmpFile.Name())
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importGlobalPrefsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importGlobalPrefsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importGlobalPrefsPidginStyle_failsWithBadContent(c *C) {
	tmpFile := s.tempFile([]byte("hopefully something that isn't valid"), c)
	_ = tmpFile.Close()

	res, ok := importGlobalPrefsPidginStyle(tmpFile.Name())
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_failsWithBadFile(c *C) {
	res, ok := ImportFingerprintsFromPidginStyle("file-that-hopefully-doesn't-exist", nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_failsWithBadContent(c *C) {
	tmpFile := s.tempFile([]byte("hopefully something that isn't valid"), c)
	_ = tmpFile.Close()

	res, ok := ImportFingerprintsFromPidginStyle(tmpFile.Name(), nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_ignoresBadFingerprintHEx(c *C) {
	tmpFile := s.tempFile([]byte("line\tsomething\tbla\tqqqq"), c)
	_ = tmpFile.Close()

	res, ok := ImportFingerprintsFromPidginStyle(tmpFile.Name(), func(string) bool { return true })
	c.Assert(res, HasLen, 0)
	c.Assert(ok, Equals, true)
}

func (s *PidginLikeSuite) tempFile(content []byte, c *C) *os.File {
	tmpfile := tempFile(c)

	_, _ = tmpfile.Write(content)

	return tmpfile
}
