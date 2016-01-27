package sasl

import (
	"io/ioutil"
	"log"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
}

type SASLSuite struct{}

var _ = Suite(&SASLSuite{})
