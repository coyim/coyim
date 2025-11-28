package filetransfer

import (
	"os"
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

func createTemporarySessionWithLog() canSendIQAndHasLogAndConnection {
	l, _ := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}
	return sess
}

func (s *ReceiverSuite) Test_receiver_simpleReceiptWorks(c *C) {
	destDir := c.MkDir()
	ctx := &recvContext{
		s:           createTemporarySessionWithLog(),
		size:        5,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp8_"),
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
	c.Assert(strings.HasPrefix(fileName, filepath.Join(destDir, "simple_receipt_test_file_tmp8_")), Equals, true)

	content, _ := os.ReadFile(fileName)
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
		destination: filepath.Join(destDir, "a directory that doesnt exist", "simple_receipt_test_file_tmp9_"),
	}

	go ctx.control.WaitForError(func(error) {})

	recv := ctx.createReceiver()

	toSend, fileName, ok, err := recv.wait()

	c.Assert(ok, Equals, false)
	c.Assert(err, ErrorMatches, ".*(no such file or directory|cannot find the (file|path) specified).*")
	c.Assert(toSend, IsNil)
	c.Assert(fileName, Equals, "")

	c.Assert(len(hook.Entries), Equals, 2)
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
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp10_"),
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

func (s *ReceiverSuite) Test_receiver_receiptOfEncryptedDataWorks(c *C) {
	destDir := c.MkDir()
	ctx := &recvContext{
		s:           createTemporarySessionWithLog(),
		size:        7,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp11_"),
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
	c.Assert(strings.HasPrefix(fileName, filepath.Join(destDir, "simple_receipt_test_file_tmp11_")), Equals, true)

	content, _ := os.ReadFile(fileName)
	c.Assert(content, DeepEquals, testDataContent)
}

func (s *ReceiverSuite) Test_receiver_receiptWithTooLittleDataForIV(c *C) {
	destDir := c.MkDir()
	l, hook := test.NewNullLogger()
	sess := &sessionMockWithCustomLog{
		log: l,
	}
	ctx := &recvContext{
		s:           sess,
		size:        7,
		control:     sdata.CreateFileTransferControl(nil, nil),
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp15_"),
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
		_, _ = recv.Write(testDataIV[0:5])
		recv.Close()
	}()

	toSend, fileName, ok, err := recv.wait()

	c.Assert(ok, Equals, false)
	c.Assert(err, ErrorMatches, "couldn't read the IV")
	c.Assert(<-hadError, ErrorMatches, "Error while reading encryption parameters")
	c.Assert(toSend, IsNil)
	c.Assert(fileName, Equals, "")

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Couldn't read encryption parameters")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "couldn't read the IV")
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
		destination: filepath.Join(destDir, "simple_receipt_test_file_tmp16_"),
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
