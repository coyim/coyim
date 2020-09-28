package muc

// TODO: We need to find a better way to do this

// MUC is a marker interface that is used to differentiate MUC "things"
type MUC interface {
	MarkAsMUCInterface()
}
