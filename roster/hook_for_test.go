package roster

import (
	"io"
	"testing"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"
	log "github.com/sirupsen/logrus"
	gocheck "gopkg.in/check.v1"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

func init() {
	log.SetOutput(io.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}
