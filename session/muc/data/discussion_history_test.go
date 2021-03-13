package data

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *MucSuite) Test_newDelayedMessage(c *C) {
	res := newDelayedMessage("hello", "someone", time.Now().UTC())
	c.Assert(res.Nickname, Equals, "hello")
	c.Assert(res.Message, Equals, "someone")
	c.Assert(res.Timestamp.Location(), Equals, time.Local)
}

func (s *MucSuite) Test_newDelayedMessages(c *C) {
	res := newDelayedMessages(time.Now().UTC())
	c.Assert(res.GetDate().Location(), Equals, time.Local)
}

func (s *MucSuite) Test_DelayedMessages_GetMessages(c *C) {
	dms := newDelayedMessages(time.Now().UTC())
	dms.add("foo", "bar", time.Now().UTC())

	res := dms.GetMessages()
	c.Assert(res, HasLen, 1)
	c.Assert(res[0].Nickname, Equals, "foo")
}

func (s *MucSuite) Test_DiscussionHistory_GetHistory(c *C) {
	dh := NewDiscussionHistory()
	msg1t := time.Date(2020, 10, 1, 0, 0, 0, 0, time.Local)
	msg2t := time.Date(2020, 10, 1, 0, 0, 0, 0, time.Local)
	msg3t := time.Date(2020, 10, 2, 0, 0, 0, 0, time.Local)
	msg4t := time.Date(2020, 11, 1, 0, 0, 0, 0, time.Local)
	msg5t := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)

	dh.AddMessage("from1", "msg1", msg1t)
	dh.AddMessage("from2", "msg2", msg2t)
	dh.AddMessage("from3", "msg3", msg3t)
	dh.AddMessage("from4", "msg4", msg4t)
	dh.AddMessage("from5", "msg5", msg5t)

	res := dh.GetHistory()
	c.Assert(res, HasLen, 4)
}
