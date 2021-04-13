package filetransfer

import (
	"io/ioutil"
	"os"

	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type RecvSuite struct{}

var _ = Suite(&RecvSuite{})

func (s *RecvSuite) Test_chooseAppropriateFileTransferOptionFrom(c *C) {
	orgSupportedFileTransferMethods := supportedFileTransferMethods
	defer func() {
		supportedFileTransferMethods = orgSupportedFileTransferMethods
	}()

	supportedFileTransferMethods = map[string]int{}

	res, ok := chooseAppropriateFileTransferOptionFrom([]string{})
	c.Assert(res, Equals, "")
	c.Assert(ok, Equals, false)

	res, ok = chooseAppropriateFileTransferOptionFrom([]string{"foo", "bar"})
	c.Assert(res, Equals, "")
	c.Assert(ok, Equals, false)

	supportedFileTransferMethods["bar"] = 2
	res, ok = chooseAppropriateFileTransferOptionFrom([]string{"foo", "bar"})
	c.Assert(res, Equals, "bar")
	c.Assert(ok, Equals, true)

	supportedFileTransferMethods["bar"] = 2
	res, ok = chooseAppropriateFileTransferOptionFrom([]string{"bar", "foo"})
	c.Assert(res, Equals, "bar")
	c.Assert(ok, Equals, true)
}

func (s *RecvSuite) Test_iqResultChosenStreamMethod(c *C) {
	res := iqResultChosenStreamMethod("foo")
	c.Assert(res, DeepEquals, data.SI{
		File: &data.File{},
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "submit",
				Fields: []data.FormFieldX{
					{Var: "stream-method", Values: []string{"foo"}},
				},
			},
		},
	})
}

func (s *RecvSuite) Test_recvContext_finalizeFileTransfer_forFile(c *C) {
	tf, _ := ioutil.TempFile("", "")
	tf.Write([]byte(`hello again`))
	_ = tf.Close()
	defer os.Remove(tf.Name())

	tf2, _ := ioutil.TempFile("", "")
	os.Remove(tf2.Name())
	defer os.Remove(tf2.Name())

	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	notDeclined := make(chan bool)
	go func() {
		ctrl.WaitForFinish(func(v bool) {
			notDeclined <- v
		})
	}()

	ctx := &recvContext{
		directory:   false,
		destination: tf2.Name(),
		control:     ctrl,
	}

	e := ctx.finalizeFileTransfer(tf.Name())
	c.Assert(e, IsNil)
	notDec := <-notDeclined

	c.Assert(notDec, Equals, true)
}

func (s *RecvSuite) Test_recvContext_finalizeFileTransfer_forFile_failsOnRename(c *C) {
	ctrl := sdata.CreateFileTransferControl(func() bool { return false }, func(bool) {})
	ee := make(chan error)
	go func() {
		ctrl.WaitForError(func(eee error) {
			ee <- eee
		})
	}()

	ctx := &recvContext{
		directory:   false,
		destination: "hmm",
		control:     ctrl,
	}

	e := ctx.finalizeFileTransfer("file that hopefully doesn't exist")
	c.Assert(e, ErrorMatches, ".*(no such file or directory|cannot find the path specified).*")
	e2 := <-ee
	c.Assert(e2, ErrorMatches, "Couldn't save final file")
}
