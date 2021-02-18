package muc

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
}

type MucSuite struct{}

var _ = Suite(&MucSuite{})
