package i18n

import (
	"github.com/coyim/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

type localGlibMock struct {
	*glib_mock.Mock
}

func (*localGlibMock) Local(vx string) string {
	return "[local]" + vx
}

type I18NSuite struct{}

var _ = Suite(&I18NSuite{})

func (s *I18NSuite) Test_Local_willReturnTheString(c *C) {
	c.Check(Local("hello"), Equals, "[local]hello")
	c.Check(Local("helllo"), Equals, "[local]helllo")
}

func (s *I18NSuite) Test_Localf_willReturnTheString(c *C) {
	c.Check(Localf("hello"), Equals, "[local]hello")
	c.Check(Localf("helllo %d", 42), Equals, "[local]helllo 42")
}

func (s *I18NSuite) Test_nullLocalizer_Local_returnsStringWithMarker(c *C) {
	c.Check((&nullLocalizer{}).Local("foo"), Equals, "[NULL LOCALIZER] - foo")
	c.Check((&nullLocalizer{}).Local("bar"), Equals, "[NULL LOCALIZER] - bar")
}
