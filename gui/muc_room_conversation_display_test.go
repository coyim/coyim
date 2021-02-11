package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

type MUCRoomConversationDisplaySuite struct{}

var _ = Suite(&MUCRoomConversationDisplaySuite{})

type mucRoomConversationDisplayMockGlib struct {
	glib_mock.Mock
}

func (*mucRoomConversationDisplayMockGlib) Local(vx string) string {
	return vx
}

func (*mucRoomConversationDisplayMockGlib) Localf(vx string, args ...interface{}) string {
	return fmt.Sprintf(vx, args...)
}

func initMUCRoomConversationDisplayI18n() {
	i18n.InitLocalization(&mucRoomConversationDisplayMockGlib{})
}

func (*MUCRoomConversationDisplaySuite) Test_mucRoomConversationDisplay_getDisplayRoomSubjectForNickname(c *C) {
	c.Assert(getDisplayRoomSubjectForNickname("", "room subject"), Equals,
		"Someone has updated the room subject to: \"room subject\"")

	c.Assert(getDisplayRoomSubjectForNickname("nick", "new subject"), Equals,
		"nick updated the room subject to: \"new subject\"")

	c.Assert(getDisplayRoomSubjectForNickname("batman", "the batman cave"), Equals,
		"batman updated the room subject to: \"the batman cave\"")
}

func (*MUCRoomConversationDisplaySuite) Test_mucRoomConversationDisplay_getDisplayRoomSubject(c *C) {
	c.Assert(getDisplayRoomSubject(""), Equals,
		"The room does not have a subject")

	c.Assert(getDisplayRoomSubject("the batman cave"), Equals,
		"The room subject is \"the batman cave\"")
}
