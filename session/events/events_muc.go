package events

// MUC description
type MUC struct {
	From string
}

// MUCOccupant description
type MUCOccupant struct {
	*MUC
	Nickname string
}

// MUCOccupantJoined description
type MUCOccupantJoined struct {
	*MUCOccupant
	Joined bool
}

// MUCOccupantUpdated description
type MUCOccupantUpdated struct {
	*MUCOccupant
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
