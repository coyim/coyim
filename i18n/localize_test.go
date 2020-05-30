package i18n

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type localGlibMock struct {
	*glib_mock.Mock
}

func (*localGlibMock) Local(vx string) string {
	return "[local]" + vx
}

func init() {
	log.SetOutput(ioutil.Discard)
	InitLocalization(&localGlibMock{&glib_mock.Mock{}})
}

type I18NSuite struct{}

var _ = Suite(&I18NSuite{})

func (s *I18NSuite) Test_Local_willReturnTheString(c *C) {
	c.Assert(Local("hello"), Equals, "[local]hello")
	c.Assert(Local("helllo"), Equals, "[local]helllo")
}

func (s *I18NSuite) Test_Localf_willReturnTheString(c *C) {
	c.Assert(Localf("hello"), Equals, "[local]hello")
	c.Assert(Localf("helllo %d", 42), Equals, "[local]helllo 42")
}
