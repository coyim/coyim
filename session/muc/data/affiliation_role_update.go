package data

// Actor represents the occupant that updates the role or affiliation
type Actor struct {
	Nickname    string
	Affiliation Affiliation
	Role        Role
}

// AffiliationRoleUpdate contains information related to a new and previous affiliation and role
type AffiliationRoleUpdate struct {
	Nickname            string
	Reason              string
	NewAffiliation      Affiliation
	PreviousAffiliation Affiliation
	NewRole             Role
	PreviousRole        Role
	Actor               *Actor
}
