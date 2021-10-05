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

func (s *ReceiverSuite) Test_receiver_cancelingWorks(c *C) {
	destDir := c.MkDir()
	l, _ := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}
	ctx := &recvContext{
		s:           sess,
		size:        5,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
	}

	go ctx.control.WaitForError(func(error) {})

	recv := ctx.createReceiver()

	_, _ = recv.Write([]byte{0x01, 0x02, 0x03})
	recv.cancel()
	_, e2 := recv.Write([]byte{0x01, 0x02, 0x03})
	c.Assert(e2, ErrorMatches, "local cancel")

	recv.cancel()
	toSend, fileName, ok, err := recv.wait()

	c.Assert(ok, Equals, false)
	c.Assert(err, Not(IsNil))
	c.Assert(toSend, IsNil)
	c.Assert(fileName, Equals, "")
}

var (
	testDataEncryptedContent = []byte{0x42, 0xf9, 0x6b, 0x9e, 0x70, 0x2d, 0xf8}

	testDataMac = []byte{
		0x2, 0x12, 0xac, 0x1b, 0xc3, 0xf6, 0x66, 0xe1,
		0x54, 0xb9, 0x95, 0xf9, 0xbd, 0x70, 0xf, 0x6a,
		0xad, 0x4a, 0xf3, 0x3c, 0x8d, 0x95, 0x6b, 0x26,
		0xe4, 0x78, 0x26, 0x77, 0x41, 0x81, 0x49, 0xfc,
	}

	testDataIV = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
	}

	testDataEncryptionKey = []byte{
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
	}

	testDataMacKey = []byte{
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F,
	}

	testDataContent = []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01, 0x02}
)

func (s *ReceiverSuite) Test_receiver_receiptOfEncryptedDataWorks(c *C) {
	destDir := c.MkDir()
	ctx := &recvContext{
		size:        7,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
		enc: &encryptionParameters{
			macKey:        testDataMacKey,
			encryptionKey: testDataEncryptionKey,
		},
	}
	recv := ctx.createReceiver()

	go func() {
		_, _ = recv.Write(testDataIV)
		_, _ = recv.Write(testDataEncryptedContent)
		_, _ = recv.Write(testDataMac)
	}()

	toSend, fileName, ok, err := recv.wait()

	c.Assert(ok, Equals, true)
	c.Assert(err, IsNil)
	c.Assert(toSend, DeepEquals, testDataMacKey)
	c.Assert(strings.HasPrefix(fileName, filepath.Join(destDir, "simple_receipt_test_file")), Equals, true)

	content, _ := ioutil.ReadFile(fileName)
	c.Assert(content, DeepEquals, testDataContent)
}

func (s *ReceiverSuite) Test_receiver_receiptOfEncryptedDataWithIncorrectMac(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}
	ctx := &recvContext{
		s:           sess,
		size:        7,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file"),
		enc: &encryptionParameters{
			macKey:        testDataMacKey,
			encryptionKey: testDataEncryptionKey,
		},
	}

	hadError := make(chan error)
	go ctx.control.WaitForError(func(e error) {
		hadError <- e
	})

	recv := ctx.createReceiver()

	go func() {
		_, _ = recv.Write(testDataIV)
		_, _ = recv.Write(testDataEncryptedContent)
		modifiedMac := append([]byte{}, testDataMac...)
		modifiedMac[0] = 0xFF
		_, _ = recv.Write(modifiedMac)
	}()

	toSend, fileName, ok, err := recv.wait()

	c.Assert(ok, Equals, false)
	c.Assert(err, ErrorMatches, "bad MAC - transfer integrity broken")
	c.Assert(<-hadError, ErrorMatches, "Couldn't verify integrity of sent file")
	c.Assert(toSend, IsNil)
	c.Assert(fileName, Equals, "")

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Couldn't verify integrity of sent file")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "bad MAC - transfer integrity broken")
}
