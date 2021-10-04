package filetransfer

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"

	"github.com/coyim/coyim/coylog"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/mock"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

type ReceiverSuite struct{}

var _ = Suite(&ReceiverSuite{})

func (s *ReceiverSuite) Test_receiver_simpleReceiptWorks(c *C) {
	destDir := c.MkDir()
	ctx := &recvContext{
		size:        5,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
	}
	recv := ctx.createReceiver()

	go func() {
		_, _ = recv.Write([]byte{0x01, 0x02, 0x03})
		_, _ = recv.Write([]byte{0x04, 0x05})
	}()

	toSend, fileName, ok, err := recv.wait()

	c.Assert(ok, Equals, true)
	c.Assert(err, IsNil)
	c.Assert(toSend, IsNil)
	c.Assert(strings.HasPrefix(fileName, filepath.Join(destDir, "simple_receipt_test_file")), Equals, true)

	content, _ := ioutil.ReadFile(fileName)
	c.Assert(content, DeepEquals, []byte{0x01, 0x02, 0x03, 0x04, 0x05})
}

type sessionMockWithCustomLog struct {
	mock.SessionMock

	log coylog.Logger
}

func (s *sessionMockWithCustomLog) Log() coylog.Logger {
	return s.log
}

func (s *ReceiverSuite) Test_receiver_anErrorHappensEarlyInTheProcess(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}
	ctx := &recvContext{
		s:           sess,
		size:        5,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "a directory that doesnt exist", "simple_receipt_test_file"),
	}

	go ctx.control.WaitForError(func(error) {})

	recv := ctx.createReceiver()

	_, _ = recv.Write([]byte{0x01, 0x02, 0x03})
	toSend, fileName, ok, err := recv.wait()

	c.Assert(ok, Equals, false)
	c.Assert(err, ErrorMatches, ".*(no such file or directory|cannot find the (file|path) specified).*")
	c.Assert(toSend, IsNil)
	c.Assert(fileName, Equals, "")

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Failed to open temporary file")
}
