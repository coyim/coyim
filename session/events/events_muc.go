package events

// MUCInfo description
type MUCInfo struct {
	From string
}

// MUCOccupant description
type MUCOccupant struct {
	MUCInfo
	Nickname string
}

// MUCOccupantJoined description
type MUCOccupantJoined struct {
	MUCOccupantUpdated
	Jid    string
	Status string
	Joined bool
}

// MUCOccupantUpdated description
type MUCOccupantUpdated struct {
	MUCOccupant
	Affiliation string
	Role        string
}

// MUCEventType represents the type of MUC event
type MUCEventType EventType

// MUCEventErrorType represents the type of MUC error event
type MUCEventErrorType MUCEventType

// MUC event types
const (
	MUCOccupantUpdate MUCEventType = iota
	MUCOccupantJoin

	MUCNotAuthorized MUCEventErrorType = iota
	MUCForbidden
	MUCItemNotFound
	MUCNotAllowed
	MUCNotAceptable
	MUCRegistrationRequired
	MUCConflict
	MUCServiceUnavailable
)

// MUC contains information related to MUC session event
type MUC struct {
	EventInfo interface{}
	EventType MUCEventType
}

// MUCError contains information related to MUC-error session event
type MUCError struct {
	EventInfo MUCInfo
	EventType MUCEventErrorType
}
