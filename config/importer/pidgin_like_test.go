package importer

import (
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
	tmpFileName := writeToTempFileAndReturnName([]byte("hopefully something that isn't valid"), c)

	res, ok := ImportKeysFromPidginStyle(tmpFileName, nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importAccountsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importAccountsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importAccountsPidginStyle_failsWithBadContent(c *C) {
	tmpFileName := writeToTempFileAndReturnName([]byte("hopefully something that isn't valid"), c)

	res, ok := importAccountsPidginStyle(tmpFileName)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importPeerPrefsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importPeerPrefsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importPeerPrefsPidginStyle_failsWithBadContent(c *C) {
	tmpFileName := writeToTempFileAndReturnName([]byte("hopefully something that isn't valid"), c)

	res, ok := importPeerPrefsPidginStyle(tmpFileName)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importGlobalPrefsPidginStyle_failsWithBadFile(c *C) {
	res, ok := importGlobalPrefsPidginStyle("file-that-hopefully-doesn't-exist")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_importGlobalPrefsPidginStyle_failsWithBadContent(c *C) {
	tmpFileName := writeToTempFileAndReturnName([]byte("hopefully something that isn't valid"), c)

	res, ok := importGlobalPrefsPidginStyle(tmpFileName)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_failsWithBadFile(c *C) {
	res, ok := ImportFingerprintsFromPidginStyle("file-that-hopefully-doesn't-exist", nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_failsWithBadContent(c *C) {
	tmpFileName := writeToTempFileAndReturnName([]byte("hopefully something that isn't valid"), c)

	res, ok := ImportFingerprintsFromPidginStyle(tmpFileName, nil)
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginLikeSuite) Test_ImportFingerprintsFromPidginStyle_ignoresBadFingerprintHEx(c *C) {
	tmpFileName := writeToTempFileAndReturnName([]byte("line\tsomething\tbla\tqqqq"), c)

	res, ok := ImportFingerprintsFromPidginStyle(tmpFileName, func(string) bool { return true })
	c.Assert(res, HasLen, 0)
	c.Assert(ok, Equals, true)
}

func writeToTempFileAndReturnName(content []byte, c *C) string {
	tmpfile := tempFile(c)

	_, _ = tmpfile.Write(content)
	_ = tmpfile.Close()

	return tmpfile.Name()
}
