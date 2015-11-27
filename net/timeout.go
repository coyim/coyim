package net

var ErrTimeout error = &TimeoutError{}

type TimeoutError struct{}

func (e *TimeoutError) Error() string { return "i/o timeout" }
