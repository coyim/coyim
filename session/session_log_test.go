package session

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type SessionLogSuite struct{}

var _ = Suite(&SessionLogSuite{})

func (s *SessionLogSuite) Test_openLogFile_returnsNilIfNoFilenameGiven(c *C) {
	c.Assert(openLogFile(""), IsNil)
}

func (s *SessionLogSuite) Test_openLogFile_failsIfCantOpenFile(c *C) {
	tf, _ := ioutil.TempFile("", "")
	defer os.Remove(tf.Name())

	ll := log.StandardLogger()
	orgLevel := ll.Level
	defer func() {
		ll.SetLevel(orgLevel)
	}()

	ll.SetLevel(log.DebugLevel)
	hook := test.NewGlobal()

	c.Assert(openLogFile(filepath.Join(tf.Name(), "something elsE")), IsNil)

	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Failed to open log file.")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, ".*not a directory")
}

func (s *SessionLogSuite) Test_openLogFile_succeeds(c *C) {
	tf, _ := ioutil.TempFile("", "")
	defer os.Remove(tf.Name())

	ll := log.StandardLogger()
	orgLevel := ll.Level
	defer func() {
		ll.SetLevel(orgLevel)
	}()

	ll.SetLevel(log.DebugLevel)
	hook := test.NewGlobal()

	c.Assert(openLogFile(tf.Name()), Not(IsNil))

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.DebugLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Logging XMPP messages to file")
	c.Assert(hook.LastEntry().Data["file"], Equals, tf.Name())
}
