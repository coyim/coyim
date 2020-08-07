package events

// MUCType description
type MUCType struct {
	From string
}

// MUCOccupantType description
type MUCOccupantType struct {
	*MUCType
	Nickname string
}

// MUCOccupantJoinedType description
type MUCOccupantJoinedType struct {
	*MUCOccupantType
	Joined bool
}

// MUCOccupantUpdatedType description
type MUCOccupantUpdatedType struct {
	*MUCOccupantType
	Affiliation string
	Role        string
}

// MUC event types
const (
	MUCOccupantUpdate EventType = iota

	MUCNotAuthorized
	MUCForbidden
	MUCItemNotFound
	MUCNotAllowed
	MUCNotAceptable
	MUCRegistrationRequired
	MUCConflict
	MUCServiceUnavailable
)
