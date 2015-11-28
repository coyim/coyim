package net

// ErrTimeout is the singleton timeout error instance for CoyIM
var ErrTimeout error = &TimeoutError{}

// TimeoutError represents a timeout error
type TimeoutError struct{}

func (e *TimeoutError) Error() string { return "i/o timeout" }
