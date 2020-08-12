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
	*MUCOccupantUpdated
	Joined bool
}

// MUCOccupantUpdated description
type MUCOccupantUpdated struct {
	*MUCOccupant
	Affiliation string
	Jid         string
	Role        string
	Status      string
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

// MUCErrorEvent structure
type MUCErrorEvent struct {
	*MUC
	Event EventType
}
