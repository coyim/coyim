package otr3

import "testing"

func Test_ErrorCode_hasValidStringImplementation(t *testing.T) {
	assertEquals(t, ErrorCodeEncryptionError.String(), "ErrorCodeEncryptionError")
	assertEquals(t, ErrorCodeMessageUnreadable.String(), "ErrorCodeMessageUnreadable")
	assertEquals(t, ErrorCodeMessageMalformed.String(), "ErrorCodeMessageMalformed")
	assertEquals(t, ErrorCodeMessageNotInPrivate.String(), "ErrorCodeMessageNotInPrivate")
	assertEquals(t, ErrorCode(20000).String(), "ERROR CODE: (THIS SHOULD NEVER HAPPEN)")
}

func Test_combinedErrorMessageHandler_callsAllErrorMessageHandlersGiven(t *testing.T) {
	var called1, called2, called3 bool
	f1 := dynamicErrorMessageHandler{func(error ErrorCode) []byte {
		called1 = true
		return nil
	}}
	f2 := dynamicErrorMessageHandler{func(error ErrorCode) []byte {
		called2 = true
		return nil
	}}
	f3 := dynamicErrorMessageHandler{func(error ErrorCode) []byte {
		called3 = true
		return nil
	}}
	d := CombineErrorMessageHandlers(nil, f1, f2, f3)
	d.HandleErrorMessage(ErrorCodeMessageMalformed)

	assertEquals(t, called1, true)
	assertEquals(t, called2, true)
	assertEquals(t, called3, true)
}

func Test_combinedErrorMessageHandler_returnsTheLastResult(t *testing.T) {
	f1 := dynamicErrorMessageHandler{func(error ErrorCode) []byte {
		return []byte("result1")
	}}
	f2 := dynamicErrorMessageHandler{func(error ErrorCode) []byte {
		return []byte("result2")
	}}
	d := CombineErrorMessageHandlers(f1, f2)
	res := d.HandleErrorMessage(ErrorCodeMessageMalformed)

	assertEquals(t, string(res), "result2")
}

func Test_debugErrorMessageHandler_writesTheErrorCodeToStderr(t *testing.T) {
	ss := captureStderr(func() {
		DebugErrorMessageHandler{}.HandleErrorMessage(ErrorCodeMessageMalformed)
	})
	assertEquals(t, ss, "[DEBUG] HandleErrorMessage(ErrorCodeMessageMalformed)\n")
}
