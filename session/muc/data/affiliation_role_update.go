package data

// OccupantUpdateActor represents the occupant that updates the role or affiliation
type OccupantUpdateActor struct {
	Nickname    string
	Affiliation Affiliation
	Role        Role
}

// OccupantUpdateAffiliationRole represents the common properties for updating occupant role or affiliation
type OccupantUpdateAffiliationRole struct {
	Nickname string
	Reason   string
	Actor    *OccupantUpdateActor
}
