// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// SetAdder is implemented by any value that exhibits a set-like behaviour,
// allowing arbitrary values to be added.
type SetAdder interface {
	Add(v interface{})
}

// Set represents a Python "set" (builtin type).
//
// It is implemented in Go as a map with empty struct values; the actual set
// of generic "interface{}" items is thus represented by all the keys.
type Set map[interface{}]setEmptyStruct

var _ SetAdder = &Set{}

type setEmptyStruct struct{}

// NewSet makes and returns a new empty Set.
func NewSet() *Set {
	s := make(Set)
	return &s
}

// NewSetFromSlice makes and returns a new Set initialized with the elements
// of the given slice.
func NewSetFromSlice(slice []interface{}) *Set {
	s := make(Set, len(slice))
	for _, item := range slice {
		s[item] = setEmptyStruct{}
	}
	return &s
}

// Len returns the length of the Set.
func (s *Set) Len() int {
	return len(*s)
}

// Add adds one element to the Set.
func (s *Set) Add(v interface{}) {
	(*s)[v] = setEmptyStruct{}
}

// Has returns whether the given value is present in the Set (true)
// or not (false).
func (s *Set) Has(v interface{}) bool {
	_, ok := (*s)[v]
	return ok
}
