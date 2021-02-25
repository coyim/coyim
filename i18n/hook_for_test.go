package i18n

import (
	"io/ioutil"
	"testing"

	"github.com/coyim/gotk3adapter/glib_mock"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	InitLocalization(&localGlibMock{&glib_mock.Mock{}})
}
