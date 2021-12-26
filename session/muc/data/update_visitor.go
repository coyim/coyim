package data

// UpdateVisitor allows type-safe access to all the different types of updates
// even though they don't share any type. This is the partner interface to
// the Update interface
type UpdateVisitor interface {
	OnAffiliationUpdate(AffiliationUpdate)
	OnRoleUpdate(RoleUpdate)
	OnAffiliationRoleUpdate(AffiliationRoleUpdate)
	OnSelfAffiliationUpdate(SelfAffiliationUpdate)
	OnSelfRoleUpdate(SelfRoleUpdate)
	OnSelfAffiliationRoleUpdate(SelfAffiliationRoleUpdate)
}

// Update represents any kind of update message, such as affiliation, role or both
type Update interface {
	Visit(UpdateVisitor)
}
