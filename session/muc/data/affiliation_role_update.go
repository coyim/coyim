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

// Visit is part of the implementation for the Update interface
func (u AffiliationRoleUpdate) Visit(vis UpdateVisitor) {
	vis.OnAffiliationRoleUpdate(u)
}

// SelfAffiliationRoleUpdate contains information related to a new and previous affiliation and role of the self occupant
type SelfAffiliationRoleUpdate struct {
	AffiliationRoleUpdate
}

// Visit is part of the implementation for the Update interface
func (u SelfAffiliationRoleUpdate) Visit(vis UpdateVisitor) {
	vis.OnSelfAffiliationRoleUpdate(u)
}
