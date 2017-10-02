package otr3

import "testing"

func Test_SecurityEvent_hasValidStringImplementation(t *testing.T) {
	assertEquals(t, GoneInsecure.String(), "GoneInsecure")
	assertEquals(t, GoneSecure.String(), "GoneSecure")
	assertEquals(t, StillSecure.String(), "StillSecure")
	assertEquals(t, SecurityEvent(20000).String(), "SECURITY EVENT: (THIS SHOULD NEVER HAPPEN)")
}

func Test_combinedSecurityEventHandler_callsAllSecurityEventHandlersGiven(t *testing.T) {
	var called1, called2, called3 bool
	f1 := dynamicSecurityEventHandler{func(event SecurityEvent) {
		called1 = true
	}}
	f2 := dynamicSecurityEventHandler{func(event SecurityEvent) {
		called2 = true
	}}
	f3 := dynamicSecurityEventHandler{func(event SecurityEvent) {
		called3 = true
	}}
	d := CombineSecurityEventHandlers(f1, nil, f2, f3)
	d.HandleSecurityEvent(GoneSecure)

	assertEquals(t, called1, true)
	assertEquals(t, called2, true)
	assertEquals(t, called3, true)
}

func Test_debugSecurityEventHandler_writesTheEventToStderr(t *testing.T) {
	ss := captureStderr(func() {
		DebugSecurityEventHandler{}.HandleSecurityEvent(StillSecure)
	})
	assertEquals(t, ss, "[DEBUG] HandleSecurityEvent(StillSecure)\n")
}
