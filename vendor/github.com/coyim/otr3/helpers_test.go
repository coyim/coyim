package otr3

import (
	"encoding/hex"
	"io"
	"math/big"
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected:\n%#v \nto equal:\n%#v\n", actual, expected)
	}
}

func assertNotEquals(t *testing.T, actual, expected interface{}) {
	if actual == expected {
		t.Errorf("Expected:\n%#v \nto not equal:\n%#v\n", actual, expected)
	}
}

func assertFuncEquals(t *testing.T, actual, expected interface{}) {
	f1 := reflect.ValueOf(actual)
	f2 := reflect.ValueOf(expected)
	if f1.Pointer() != f2.Pointer() {
		t.Errorf("Expected:\n%#v \nto equal:\n%#v\n", actual, expected)
	}
}

func isNil(actual interface{}) bool {
	val := reflect.ValueOf(actual)
	switch val.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return actual == nil
	}
}

func assertNil(t *testing.T, actual interface{}) {
	if !isNil(actual) {
		t.Errorf("Expected:\n%#v \nto be nil\n", actual)
	}
}

func assertTrue(t *testing.T, actual bool) {
	if !actual {
		t.Errorf("Expected: %#v to be true\n", actual)
	}
}

func assertFalse(t *testing.T, actual bool) {
	if actual {
		t.Errorf("Expected: %#v to be false\n", actual)
	}
}

func assertNotNil(t *testing.T, actual interface{}) {
	if isNil(actual) {
		t.Errorf("Expected:\n%#v \nto not be nil\n", actual)
	}
}

func assertDeepEquals(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected:\n%#v \nto equal:\n%#v\n", actual, expected)
	}
}

func dhMsgType(msg []byte) byte {
	return msg[2]
}

func dhMsgVersion(msg []byte) uint16 {
	_, protocolVersion, _ := extractShort(msg)
	return protocolVersion
}

func bytesFromHex(s string) []byte {
	val, _ := hex.DecodeString(s)
	return val
}

// bnFromHex is a test utility that doesn't take into account possible errors. Thus, make sure to only call it with valid hexadecimal strings (of even length)
func bnFromHex(s string) *big.Int {
	res, _ := new(big.Int).SetString(s, 16)
	return res
}

// parseIntoPrivateKey is a test utility that doesn't take into account possible errors. Thus, make sure to only call it with valid values
// it only parses DSA keys right now
func parseIntoPrivateKey(hexString string) PrivateKey {
	b, _ := hex.DecodeString(hexString)
	pk := new(DSAPrivateKey)
	pk.Parse(b)
	return pk
}

func newConversation(v otrVersion, rand io.Reader) *Conversation {
	var p policy
	switch v {
	case otrV3{}:
		p = allowV3
	case otrV2{}:
		p = allowV2
	}
	akeNotStarted := new(ake)
	akeNotStarted.state = authStateNone{}

	return &Conversation{
		version: v,
		Rand:    rand,
		smp: smp{
			state: smpStateExpect1{},
		},
		ake:              akeNotStarted,
		Policies:         policies(p),
		fragmentSize:     65535, //we are not testing fragmentation by default
		ourInstanceTag:   0x101, //every conversation should be able to talk to each other
		theirInstanceTag: 0x101,
	}
}

func (c *Conversation) expectMessageEvent(t *testing.T, f func(), expectedEvent MessageEvent, expectedMessage []byte, expectedError error) {
	called := false

	c.messageEventHandler = dynamicMessageEventHandler{func(event MessageEvent, message []byte, err error, trace ...interface{}) {
		assertDeepEquals(t, event, expectedEvent)
		assertDeepEquals(t, message, expectedMessage)
		assertDeepEquals(t, err, expectedError)
		called = true
	}}

	f()

	assertEquals(t, called, true)
}

func (c *Conversation) doesntExpectMessageEvent(t *testing.T, f func()) {
	c.messageEventHandler = dynamicMessageEventHandler{func(event MessageEvent, message []byte, err error, trace ...interface{}) {
		t.Errorf("Didn't expect a message event, but got: %v with msg %v and error %#v", event, message, err)
	}}

	f()
}

func (c *Conversation) expectSMPEvent(t *testing.T, f func(), expectedEvent SMPEvent, expectedProgress int, expectedQuestion string) {
	called := false

	c.smpEventHandler = dynamicSMPEventHandler{func(event SMPEvent, progressPercent int, question string) {
		assertEquals(t, event, expectedEvent)
		assertEquals(t, progressPercent, expectedProgress)
		assertEquals(t, question, expectedQuestion)
		called = true
	}}

	f()

	assertEquals(t, called, true)
}

func (c *Conversation) expectSecurityEvent(t *testing.T, f func(), expectedEvent SecurityEvent) {
	called := false

	c.securityEventHandler = dynamicSecurityEventHandler{func(event SecurityEvent) {
		assertEquals(t, event, expectedEvent)
		called = true
	}}

	f()

	assertEquals(t, called, true)
}

func (c *Conversation) doesntExpectSecurityEvent(t *testing.T, f func()) {
	c.securityEventHandler = dynamicSecurityEventHandler{func(event SecurityEvent) {
		t.Errorf("Didn't expect a security event, but got: %v", event)
	}}

	f()
}
