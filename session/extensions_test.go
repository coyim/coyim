package session

import (
	"bytes"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type ExtensionsSuite struct{}

var _ = Suite(&ExtensionsSuite{})

func (s *ExtensionsSuite) Test_unknownExtension_logsInformation(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{log: l}

	ext := &data.Extension{
		Body: "hello world",
	}

	unknownExtension(sess, nil, ext)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Unknown extension")
	c.Assert(hook.LastEntry().Data["extension"].(*bytes.Buffer).String(), Equals, "hello world")
}

func (s *ExtensionsSuite) Test_getExtensionHandler_returnsAPreviouslyRegisterExtension(c *C) {
	called := false

	registerKnownExtension("temp1 temp2", func(access.Session, *data.ClientMessage, *data.Extension) {
		called = true
	})

	getExtensionHandler("temp1", "temp2")(nil, nil, nil)

	c.Assert(called, Equals, true)
}

func (s *ExtensionsSuite) Test_getExtensionHandler_returnsUnknownExtensionHandler(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{log: l}

	ext := &data.Extension{
		Body: "goodbye world",
	}

	getExtensionHandler("temp3", "temp4")(sess, nil, ext)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Unknown extension")
	c.Assert(hook.LastEntry().Data["extension"].(*bytes.Buffer).String(), Equals, "goodbye world")
}

func (s *ExtensionsSuite) Test_processExtensions_works(c *C) {
	called := false

	registerKnownExtension("foobar query", func(access.Session, *data.ClientMessage, *data.Extension) {
		called = true
	})

	sess := &session{}

	msg := &data.ClientMessage{}
	msg.Extensions = []*data.Extension{
		&data.Extension{Body: `<query xmlns="foobar"/>`},
		&data.Extension{Body: `<query`},
	}

	sess.processExtensions(msg)

	c.Assert(called, Equals, true)
}
