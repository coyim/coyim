package otr3

import "testing"

func Test_MessageEvent_hasValidStringImplementation(t *testing.T) {
	assertEquals(t, MessageEventEncryptionRequired.String(), "MessageEventEncryptionRequired")
	assertEquals(t, MessageEventEncryptionError.String(), "MessageEventEncryptionError")
	assertEquals(t, MessageEventConnectionEnded.String(), "MessageEventConnectionEnded")
	assertEquals(t, MessageEventSetupError.String(), "MessageEventSetupError")
	assertEquals(t, MessageEventMessageReflected.String(), "MessageEventMessageReflected")
	assertEquals(t, MessageEventMessageResent.String(), "MessageEventMessageResent")
	assertEquals(t, MessageEventReceivedMessageNotInPrivate.String(), "MessageEventReceivedMessageNotInPrivate")
	assertEquals(t, MessageEventReceivedMessageUnreadable.String(), "MessageEventReceivedMessageUnreadable")
	assertEquals(t, MessageEventReceivedMessageMalformed.String(), "MessageEventReceivedMessageMalformed")
	assertEquals(t, MessageEventLogHeartbeatReceived.String(), "MessageEventLogHeartbeatReceived")
	assertEquals(t, MessageEventLogHeartbeatSent.String(), "MessageEventLogHeartbeatSent")
	assertEquals(t, MessageEventReceivedMessageGeneralError.String(), "MessageEventReceivedMessageGeneralError")
	assertEquals(t, MessageEventReceivedMessageUnencrypted.String(), "MessageEventReceivedMessageUnencrypted")
	assertEquals(t, MessageEventReceivedMessageUnrecognized.String(), "MessageEventReceivedMessageUnrecognized")
	assertEquals(t, MessageEventReceivedMessageForOtherInstance.String(), "MessageEventReceivedMessageForOtherInstance")
	assertEquals(t, MessageEvent(20000).String(), "MESSAGE EVENT: (THIS SHOULD NEVER HAPPEN)")
}

func Test_combinedMessageEventHandler_callsAllErrorMessageHandlersGiven(t *testing.T) {
	var called1, called2, called3 bool
	f1 := dynamicMessageEventHandler{func(event MessageEvent, message []byte, err error, trace ...interface{}) {
		called1 = true
	}}
	f2 := dynamicMessageEventHandler{func(event MessageEvent, message []byte, err error, trace ...interface{}) {
		called2 = true
	}}
	f3 := dynamicMessageEventHandler{func(event MessageEvent, message []byte, err error, trace ...interface{}) {
		called3 = true
	}}
	d := CombineMessageEventHandlers(f1, f2, nil, f3)
	d.HandleMessageEvent(MessageEventSetupError, []byte("something"), nil)

	assertEquals(t, called1, true)
	assertEquals(t, called2, true)
	assertEquals(t, called3, true)
}

func Test_debugMessageEventHandler_writesTheEventToStderr(t *testing.T) {
	ss := captureStderr(func() {
		DebugMessageEventHandler{}.HandleMessageEvent(MessageEventLogHeartbeatSent, []byte("A message"), newOtrError("hello world"))
	})
	assertEquals(t, ss, "[DEBUG] HandleMessageEvent(MessageEventLogHeartbeatSent, message: \"A message\", error: otr: hello world, trace: [])\n")
}
