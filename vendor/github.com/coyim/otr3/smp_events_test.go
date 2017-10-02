package otr3

import "testing"

func Test_SMPEvent_hasValidStringImplementation(t *testing.T) {
	assertEquals(t, SMPEventError.String(), "SMPEventError")
	assertEquals(t, SMPEventAbort.String(), "SMPEventAbort")
	assertEquals(t, SMPEventCheated.String(), "SMPEventCheated")
	assertEquals(t, SMPEventAskForAnswer.String(), "SMPEventAskForAnswer")
	assertEquals(t, SMPEventAskForSecret.String(), "SMPEventAskForSecret")
	assertEquals(t, SMPEventInProgress.String(), "SMPEventInProgress")
	assertEquals(t, SMPEventSuccess.String(), "SMPEventSuccess")
	assertEquals(t, SMPEventFailure.String(), "SMPEventFailure")
	assertEquals(t, SMPEvent(20000).String(), "SMP EVENT: (THIS SHOULD NEVER HAPPEN)")
}

func Test_combinedSMPEventHandler_callsAllErrorMessageHandlersGiven(t *testing.T) {
	var called1, called2, called3 bool
	f1 := dynamicSMPEventHandler{func(event SMPEvent, progressPercent int, question string) {
		called1 = true
	}}
	f2 := dynamicSMPEventHandler{func(event SMPEvent, progressPercent int, question string) {
		called2 = true
	}}
	f3 := dynamicSMPEventHandler{func(event SMPEvent, progressPercent int, question string) {
		called3 = true
	}}
	d := CombineSMPEventHandlers(f1, nil, f2, f3)
	d.HandleSMPEvent(SMPEventError, 61, "")

	assertEquals(t, called1, true)
	assertEquals(t, called2, true)
	assertEquals(t, called3, true)
}

func Test_debugSMPEventHandler_writesTheEventToStderr(t *testing.T) {
	ss := captureStderr(func() {
		DebugSMPEventHandler{}.HandleSMPEvent(SMPEventInProgress, 43, "Maybe now?")
	})
	assertEquals(t, ss, "[DEBUG] HandleSMPEvent(SMPEventInProgress, 43, \"Maybe now?\")\n")
}
