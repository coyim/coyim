package filetransfer

import (
	"bytes"
	"io"
	"os"
	"time"

	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type IBBSendSuite struct{}

var _ = Suite(&IBBSendSuite{})

func (s *IBBSendSuite) Test_ibbSendCurrentlyValid_alwaysReturnsTrue(c *C) {
	c.Assert(ibbSendCurrentlyValid("foo", nil), Equals, true)
	c.Assert(ibbSendCurrentlyValid("bar", nil), Equals, true)
}

func (s *IBBSendSuite) Test_ibbSendDoWithBlockSize_sendsProperly(c *C) {
	orgIbbScheduleSendLimit := ibbScheduleSendLimit
	defer func() {
		ibbScheduleSendLimit = orgIbbScheduleSendLimit
	}()
	ibbScheduleSendLimit = time.Duration(1) * time.Millisecond

	tf, ex := os.CreateTemp("", "coyim-session-ibb_send-")
	c.Assert(ex, IsNil)
	_, ex = tf.Write([]byte(`hello again`))
	c.Assert(ex, IsNil)
	ex = tf.Close()
	c.Assert(ex, IsNil)
	defer func() {
		ex2 := os.Remove(tf.Name())
		c.Assert(ex2, IsNil)
	}()

	conn := &sendIQAndRandMock{
		sendIQMock: &sendIQMock{},
		rand: func() io.Reader {
			return bytes.NewBufferString("abcdefesdfgsdfgsdfg")
		},
	}

	done := make(chan bool)

	seqs := []uint16{}
	datas := []string{}

	conn.sendIQ = func(v1 string, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
		switch tp := v3.(type) {
		case data.IBBOpen:
		case data.IBBData:
			seqs = append(seqs, tp.Sequence)
			datas = append(datas, tp.Base64)
		case data.IBBClose:
			done <- true
		}

		ch := make(chan data.Stanza, 1)
		ch <- data.Stanza{
			Value: &data.ClientIQ{
				Type: "result",
			},
		}
		return ch, 0, nil
	}

	mc := &mockHasConnectionAndConfigAndLog{
		c: conn,
		l: nil,
	}

	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	ctx := &sendContext{
		file:    tf.Name(),
		sid:     "baz",
		s:       mc,
		control: ctrl,
	}

	ibbSendDoWithBlockSize(ctx, 3)

	<-done

	c.Assert(seqs, DeepEquals, []uint16{0, 1, 2, 3})
	c.Assert(datas, DeepEquals, []string{"aGVs", "bG8g", "YWdh", "aW4="})
}
