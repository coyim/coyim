package gui

import (
	. "gopkg.in/check.v1"
)

type MUCRoomConversationDisplaySuite struct{}

var _ = Suite(&MUCRoomConversationDisplaySuite{})

func (s *MUCRoomConversationDisplaySuite) SetUpSuite(c *C) {
	initMUCi18n()
}

func (*MUCRoomConversationDisplaySuite) Test_mucRoomConversationDisplay_messageFromRoomSubjectUpdate(c *C) {
	c.Assert(messageFromRoomSubjectUpdate("", "room subject"), Equals,
		"[localized] Someone has updated the room subject to: \"room subject\".")

	c.Assert(messageFromRoomSubjectUpdate("nick", "new subject"), Equals,
		"[localized] $nickname{nick} updated the room subject to: \"new subject\".")

	c.Assert(messageFromRoomSubjectUpdate("batman", "the batman cave"), Equals,
		"[localized] $nickname{batman} updated the room subject to: \"the batman cave\".")
}

func (*MUCRoomConversationDisplaySuite) Test_mucRoomConversationDisplay_messageForRoomSubject(c *C) {
	c.Assert(messageForRoomSubject(""), Equals, "[localized] The room does not have a subject.")

	c.Assert(messageForRoomSubject("the batman cave"), Equals, "[localized] The room subject is \"the batman cave\".")
}
