package gui

import (
	. "gopkg.in/check.v1"
)

type MUCRoomConversationDisplaySuite struct{}

var _ = Suite(&MUCRoomConversationDisplaySuite{})

func (s *MUCRoomConversationDisplaySuite) SetUpSuite(c *C) {
	initMUCI18n()
}

func (*MUCRoomConversationDisplaySuite) Test_mucRoomConversationDisplay_getDisplayRoomSubjectForNickname(c *C) {
	c.Assert(getDisplayRoomSubjectForNickname("", "room subject"), Equals,
		"[localized] Someone has updated the room subject to: \"room subject\"")

	c.Assert(getDisplayRoomSubjectForNickname("nick", "new subject"), Equals,
		"[localized] nick updated the room subject to: \"new subject\"")

	c.Assert(getDisplayRoomSubjectForNickname("batman", "the batman cave"), Equals,
		"[localized] batman updated the room subject to: \"the batman cave\"")
}

func (*MUCRoomConversationDisplaySuite) Test_mucRoomConversationDisplay_getDisplayRoomSubject(c *C) {
	c.Assert(getDisplayRoomSubject(""), Equals, "[localized] The room does not have a subject")

	c.Assert(getDisplayRoomSubject("the batman cave"), Equals, "[localized] The room subject is \"the batman cave\"")
}
