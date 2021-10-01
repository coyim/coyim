package filetransfer

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"

	sdata "github.com/coyim/coyim/session/data"
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
