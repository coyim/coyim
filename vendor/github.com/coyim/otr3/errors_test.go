package otr3

import "testing"

func Test_OtrError_Error_returnsAValidErrorString(t *testing.T) {
	e := newOtrError("hello world")
	assertEquals(t, e.Error(), "otr: hello world")
}
