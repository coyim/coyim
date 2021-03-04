package session

import (
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/coyim/xmpp/mock"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type FileTransferSuite struct{}

var _ = Suite(&FileTransferSuite{})

type discoveryFeaturesXMPPConnMock struct {
	*mock.Conn

	discoveryFeatures func(string) ([]string, bool)
}

func (d *discoveryFeaturesXMPPConnMock) DiscoveryFeatures(v1 string) ([]string, bool) {
	return d.discoveryFeatures(v1)
}

func (s *FileTransferSuite) Test_session_SendFileTo(c *C) {
	l, hook := test.NewNullLogger()
	mc := &discoveryFeaturesXMPPConnMock{
		discoveryFeatures: func(string) ([]string, bool) {
			return nil, false
		},
	}
	sess := &session{log: l, conn: mc}
	sess.convManager = otrclient.NewConversationManager(nil, nil, "", nil, nil)

	res := sess.SendFileTo(jid.Parse("hello@bar.com"), "non-existing-file.bin", func() bool { return false }, func(bool) {})
	c.Assert(res, Not(IsNil))
	c.Assert(hook.Entries, HasLen, 1)
}

func (s *FileTransferSuite) Test_session_SendDirTo(c *C) {
	l, hook := test.NewNullLogger()
	mc := &discoveryFeaturesXMPPConnMock{
		discoveryFeatures: func(string) ([]string, bool) {
			return nil, false
		},
	}
	sess := &session{log: l, conn: mc}
	sess.convManager = otrclient.NewConversationManager(nil, nil, "", nil, nil)

	res := sess.SendDirTo(jid.Parse("hello@bar.com"), "non-existing-file.bin", func() bool { return false }, func(bool) {})
	c.Assert(res, Not(IsNil))
	c.Assert(hook.Entries, HasLen, 1)
}
