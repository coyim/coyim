package filetransfer

import (
	"bytes"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/coyim/xmpp/mock"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"

	mck "github.com/stretchr/testify/mock"
)

type SendDirSuite struct{}

var _ = Suite(&SendDirSuite{})

func (s *SendDirSuite) Test_InitSend_works(c *C) {
	orgSupportedSendingMechanisms := supportedSendingMechanisms
	orgIsSendingMechanismCurrentlyValid := isSendingMechanismCurrentlyValid

	defer func() {
		supportedSendingMechanisms = orgSupportedSendingMechanisms
		isSendingMechanismCurrentlyValid = orgIsSendingMechanismCurrentlyValid
	}()

	supportedSendingMechanisms = map[string]func(*sendContext){}
	isSendingMechanismCurrentlyValid = map[string]func(string, hasConnectionAndConfig) bool{}

	done := make(chan bool)

	supportedSendingMechanisms["foo"] = func(sc *sendContext) {
		done <- true
	}

	isSendingMechanismCurrentlyValid["foo"] = func(string, hasConnectionAndConfig) bool {
		return true
	}

	dd := c.MkDir()
	createTemporaryDirectoryStructure(dd)

	conn := new(mock.MockedConn)
	conn.On("DiscoveryFeatures", "some2@one.org").Return([]string{
		"foo",
		"http://jabber.org/protocol/si/profile/hubba",
		"http://jabber.org/protocol/si",
		"http://jabber.org/protocol/si/profile/bubba",
		"http://jabber.org/protocol/si/profile/directory-transfer",
	}, true)
	conn.On("Rand").Return(bytes.NewBufferString("abcdefesdfgsdfgsdfg"))
	siChan := make(chan data.Stanza, 1)
	var ret1 <-chan data.Stanza = siChan
	var args *mck.Arguments = nil
	conn.On("SendIQ", "some2@one.org", "set", mck.AnythingOfType("data.SI")).Return(ret1, data.Cookie(42), nil).Run(func(a mck.Arguments) {
		args = &a
	})

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mc := &allSessionMock{
		mockHasConnectionAndConfigAndLog: &mockHasConnectionAndConfigAndLog{
			c: conn,
			l: l,
		},

		createSymm: func(v jid.Any) []byte { return nil },
	}

	ctrl := InitSendDir(mc, jid.Parse("some2@one.org"), dd, nil, nil)
	c.Assert(ctrl, Not(IsNil))

	siChan <- data.Stanza{
		Value: &data.ClientIQ{
			Type: "result",
			Query: []byte(`<si xmlns="http://jabber.org/protocol/si">
  <feature xmlns="http://jabber.org/protocol/feature-neg">
    <x xmlns="jabber:x:data" type="submit">
      <field var="stream-method">
        <value>foo</value>
      </field>
    </x>
  </feature>
</si>`),
		},
	}

	args = args

	go ctrl.WaitForError(func(e error) {
		done <- false
	})

	c.Assert(<-done, Equals, true)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Data, HasLen, 3)
	c.Assert(hook.Entries[0].Data["method"], Equals, "foo")
	c.Assert(hook.Entries[0].Data["peer"], Equals, "some2@one.org")
	c.Assert(hook.Entries[0].Message, Equals, "Started sending file to peer using method")
}
