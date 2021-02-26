package events

import (
	. "gopkg.in/check.v1"
)

type EventsSuite struct{}

var _ = Suite(&EventsSuite{})

func testMUCInterface(m MUC) {
	m.markAsMUCEventTypeInterface()
}

func (s *EventsSuite) Test_MarkerFunctions(c *C) {
	testMUCInterface(MUCError{})
	testMUCInterface(MUCRoom{})
	testMUCInterface(MUCRoomCreated{})
	testMUCInterface(MUCRoomDestroyed{})
	testMUCInterface(MUCRoomRenamed{})
	testMUCInterface(MUCOccupant{})
	testMUCInterface(MUCOccupantUpdated{})
	testMUCInterface(MUCOccupantJoined{})
	testMUCInterface(MUCSelfOccupantJoined{})
	testMUCInterface(MUCOccupantLeft{})
	testMUCInterface(MUCLiveMessageReceived{})
	testMUCInterface(MUCDelayedMessageReceived{})
	testMUCInterface(MUCSubjectUpdated{})
	testMUCInterface(MUCSubjectReceived{})
	testMUCInterface(MUCLoggingEnabled{})
	testMUCInterface(MUCLoggingDisabled{})
	testMUCInterface(MUCRoomAnonymityChanged{})
	testMUCInterface(MUCDiscussionHistoryReceived{})
	testMUCInterface(MUCRoomDiscoInfoReceived{})
	testMUCInterface(MUCRoomConfigTimeout{})
	testMUCInterface(MUCRoomConfigChanged{})
	testMUCInterface(MUCOccupantRemoved{})
	testMUCInterface(MUCSelfOccupantRemoved{})
	testMUCInterface(MUCOccupantAffiliationUpdated{})
	testMUCInterface(MUCSelfOccupantAffiliationUpdated{})
	testMUCInterface(MUCOccupantRoleUpdated{})
	testMUCInterface(MUCSelfOccupantRoleUpdated{})
	testMUCInterface(MUCOccupantAffiliationRoleUpdated{})
	testMUCInterface(MUCSelfOccupantAffiliationRoleUpdated{})
	testMUCInterface(MUCOccupantKicked{})
	testMUCInterface(MUCSelfOccupantKicked{})
}
