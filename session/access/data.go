package access

// OfflineError is returned when we try to do an action that can only be done when online
type OfflineError struct {
	Msg string
}

func (v *OfflineError) Error() string {
	return v.Msg
}
