package otrclient

import (
	"bytes"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

type ConversationManagerSuite struct{}

var _ = Suite(&ConversationManagerSuite{})

type testSender struct {
	peer jid.Any
	msg  string
	err  error
}

type testConvBuilder struct {
	fake *otr3.Conversation
}

func (cb *testConvBuilder) NewConversation(peer jid.Any) *otr3.Conversation {
	return cb.fake
}

func (ts *testSender) Send(peer jid.Any, msg string, otr bool) error {
	ts.peer = peer
	ts.msg = msg
	return ts.err
}

func (s *ConversationManagerSuite) Test_TerminateAll_willTerminate(c *C) {
	cb := &testConvBuilder{&otr3.Conversation{}}
	ts := &testSender{err: nil}
	mgr := NewConversationManager(cb.NewConversation, ts, "blarg", func(jid.Any, *EventHandler, chan string, chan int) {}, log.New().WithFields(log.Fields{}))
	conv, created := mgr.EnsureConversationWith(jid.NR("someone@whitehouse.gov"), nil)

	c.Assert(created, Equals, true)
	c.Assert(conv, Not(IsNil))

	mgr.TerminateAll()

	c.Assert(ts.msg, Equals, "")
}

func (s *ConversationManagerSuite) Test_conversationManager_getConversationWithUnlocked_getsExistingConversation(c *C) {
	cm := &conversationManager{
		conversations: map[string]*conversation{},
	}
	conv1 := &conversation{}
	cm.conversations["foo.bar@hello.com/somewhere"] = conv1

	conv, ok := cm.getConversationWithUnlocked(jid.Parse("foo.bar@hello.com/somewhere"))
	c.Assert(ok, Equals, true)
	c.Assert(conv, Equals, conv1)
}

func (s *ConversationManagerSuite) Test_conversationManager_getConversationWithUnlocked_doesntFindConversation(c *C) {
	cm := &conversationManager{
		conversations: map[string]*conversation{},
	}

	_, ok := cm.getConversationWithUnlocked(jid.Parse("foo.bar@hello.com/somewhere"))
	c.Assert(ok, Equals, false)
}

func (s *ConversationManagerSuite) Test_conversationManager_getConversationWithUnlocked_managesConversationForUnresourcedJID(c *C) {
	cm := &conversationManager{
		conversations: map[string]*conversation{},
	}
	conv1 := &conversation{
		eh: &EventHandler{},
	}
	cm.conversations["foo.bar@hello.com"] = conv1

	conv, ok := cm.getConversationWithUnlocked(jid.Parse("foo.bar@hello.com/somewhere"))
	c.Assert(ok, Equals, true)
	c.Assert(conv, Equals, conv1)
	c.Assert(conv1.locked, Equals, true)
	c.Assert(conv1.peer, DeepEquals, jid.Parse("foo.bar@hello.com/somewhere"))
	c.Assert(conv1.eh.peer, DeepEquals, jid.Parse("foo.bar@hello.com/somewhere"))
	_, inpwor := cm.conversations["foo.bar@hello.com"]
	c.Assert(inpwor, Equals, false)
	_, inpwr := cm.conversations["foo.bar@hello.com/somewhere"]
	c.Assert(inpwr, Equals, true)
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func (s *ConversationManagerSuite) Test_conversationManager_getConversationWithUnlocked_warnsOnWeirdLocking(c *C) {
	cm := &conversationManager{
		conversations: map[string]*conversation{},
	}
	conv1 := &conversation{
		eh:     &EventHandler{},
		locked: true,
	}
	cm.conversations["foo.bar@hello.com"] = conv1

	var ok bool

	stdout := captureStdout(func() {
		_, ok = cm.getConversationWithUnlocked(jid.Parse("foo.bar@hello.com/somewhere"))
	})
	c.Assert(ok, Equals, true)
	c.Assert(stdout, Equals, "UNEXPECTED SITUATION OCCURRED - conversation with foo.bar@hello.com already locked to foo.bar@hello.com/somewhere without saved correct - this shouldn't be possible\n")
}

func (s *ConversationManagerSuite) Test_conversationManager_GetConversationWith_works(c *C) {
	cm := &conversationManager{
		conversations: map[string]*conversation{},
	}
	conv1 := &conversation{}
	cm.conversations["foo.bar@hello.com/somewhere"] = conv1

	conv, ok := cm.GetConversationWith(jid.Parse("foo.bar@hello.com/somewhere"))
	c.Assert(ok, Equals, true)
	c.Assert(conv, Equals, conv1)
}

func (s *ConversationManagerSuite) Test_conversationManager_EnsureConversationWith_simplyReturnsConversationIfItExists(c *C) {
	cm := &conversationManager{
		conversations: map[string]*conversation{},
	}
	conv1 := &conversation{}
	cm.conversations["foo.bar@hello.com/somewhere"] = conv1

	conv, created := cm.EnsureConversationWith(jid.Parse("foo.bar@hello.com/somewhere"), nil)
	c.Assert(conv, Equals, conv1)
	c.Assert(created, Equals, false)
}

func (s *ConversationManagerSuite) Test_conversationManager_EnsureConversationWith_triesToExtractInstanceTagsFromMessage(c *C) {
	cb := &testConvBuilder{&otr3.Conversation{}}
	cm := &conversationManager{
		builder:       cb.NewConversation,
		conversations: map[string]*conversation{},
		onCreateEH:    func(jid.Any, *EventHandler, chan string, chan int) {},
	}

	_, created := cm.EnsureConversationWith(jid.Parse("foo.bar@hello.com/somewhere"), []byte("hello"))
	c.Assert(created, Equals, true)
}

func (s *ConversationManagerSuite) Test_conversationManager_EnsureConversationWith_triesToLookForOldConversations_butFailing(c *C) {
	cb := &testConvBuilder{&otr3.Conversation{}}
	cm := &conversationManager{
		builder:       cb.NewConversation,
		conversations: map[string]*conversation{},
		onCreateEH:    func(jid.Any, *EventHandler, chan string, chan int) {},
	}

	cm.conversations["foo@bar.com/something"] = &conversation{}

	_, created := cm.EnsureConversationWith(jid.Parse("foo.bar@hello.com/somewhere"), []byte("?OTR:AAICAAAAxPWaCOvRNycg72w2shQjcSEiYjcTh+w7rq+48UM9mpZIkpN08jtTAPcc8/9fcx9mmlVy/We+n6/G65RvobYWPoY+KD9Si41TFKku34gU4HaBbwwa7XpB/4u1gPCxY6EGe0IjthTUGK2e3qLf9YCkwJ1lm+X9kPOS/Jqu06V0qKysmbUmuynXG8T5Q8rAIRPtA/RYMqSGIvfNcZfrlJRIw6M784YtWlF3i2B6dmtjMrjH/8x5myN++Q2bxh69g6z/WX1rAFoAAAAg7Vwgf3JoiH5MdRznnS3aL66tjxQzN5qiwLtImE+KFnM=."))
	c.Assert(created, Equals, true)
}

func (s *ConversationManagerSuite) Test_conversationManager_EnsureConversationWith_triesToLookForOldConversations_butFailingWithWrongInstanceTag(c *C) {
	origGetTheirInstanceTag := getTheirInstanceTag
	defer func() {
		getTheirInstanceTag = origGetTheirInstanceTag
	}()

	getTheirInstanceTag = func(c *otr3.Conversation) uint32 {
		return uint32(42)
	}

	cb := &testConvBuilder{&otr3.Conversation{}}
	cm := &conversationManager{
		builder:       cb.NewConversation,
		conversations: map[string]*conversation{},
		onCreateEH:    func(jid.Any, *EventHandler, chan string, chan int) {},
	}

	cm.conversations["foo.bar@hello.com/something"] = &conversation{Conversation: &otr3.Conversation{}}

	_, created := cm.EnsureConversationWith(jid.Parse("foo.bar@hello.com/somewhere"), []byte("?OTR:AAICAAAAxPWaCOvRNycg72w2shQjcSEiYjcTh+w7rq+48UM9mpZIkpN08jtTAPcc8/9fcx9mmlVy/We+n6/G65RvobYWPoY+KD9Si41TFKku34gU4HaBbwwa7XpB/4u1gPCxY6EGe0IjthTUGK2e3qLf9YCkwJ1lm+X9kPOS/Jqu06V0qKysmbUmuynXG8T5Q8rAIRPtA/RYMqSGIvfNcZfrlJRIw6M784YtWlF3i2B6dmtjMrjH/8x5myN++Q2bxh69g6z/WX1rAFoAAAAg7Vwgf3JoiH5MdRznnS3aL66tjxQzN5qiwLtImE+KFnM=."))
	c.Assert(created, Equals, true)
}

func (s *ConversationManagerSuite) Test_getTheirInstanceTag_works(c *C) {
	v := &otr3.Conversation{}
	c.Assert(getTheirInstanceTag(v), Equals, uint32(0))
}

func (s *ConversationManagerSuite) Test_conversationManager_EnsureConversationWith_triesToLookForOldConversationsAndChangesWhenFound(c *C) {
	origGetTheirInstanceTag := getTheirInstanceTag
	defer func() {
		getTheirInstanceTag = origGetTheirInstanceTag
	}()

	getTheirInstanceTag = func(c *otr3.Conversation) uint32 {
		return uint32(196)
	}

	cb := &testConvBuilder{&otr3.Conversation{}}
	cm := &conversationManager{
		builder:       cb.NewConversation,
		conversations: map[string]*conversation{},
		onCreateEH:    func(jid.Any, *EventHandler, chan string, chan int) {},
	}

	conv := &conversation{Conversation: &otr3.Conversation{}}
	cm.conversations["foo.bar@hello.com/something"] = conv

	_, created := cm.EnsureConversationWith(jid.Parse("foo.bar@hello.com/somewhere"), []byte("?OTR:AAICAAAAxPWaCOvRNycg72w2shQjcSEiYjcTh+w7rq+48UM9mpZIkpN08jtTAPcc8/9fcx9mmlVy/We+n6/G65RvobYWPoY+KD9Si41TFKku34gU4HaBbwwa7XpB/4u1gPCxY6EGe0IjthTUGK2e3qLf9YCkwJ1lm+X9kPOS/Jqu06V0qKysmbUmuynXG8T5Q8rAIRPtA/RYMqSGIvfNcZfrlJRIw6M784YtWlF3i2B6dmtjMrjH/8x5myN++Q2bxh69g6z/WX1rAFoAAAAg7Vwgf3JoiH5MdRznnS3aL66tjxQzN5qiwLtImE+KFnM=."))
	c.Assert(created, Equals, false)
	c.Assert(conv.peer, DeepEquals, jid.Parse("foo.bar@hello.com/somewhere"))
	_, ex := cm.conversations["foo.bar@hello.com/something"]
	c.Assert(ex, Equals, false)
}
