package sexp

import (
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %v to equal %v", actual, expected)
	}
}

func assertDeepEquals(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v to equal %v", actual, expected)
	}
}

func checkForPanic(t *testing.T, s string) {
	if r := recover(); r != nil {
		assertEquals(t, r, s)
	} else {
		t.Errorf("Expected panic with message %v to be invoked", s)
	}
}
