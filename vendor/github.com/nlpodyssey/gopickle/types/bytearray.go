// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// ByteArray represents a Python "bytearray" (builtin type).
type ByteArray []byte

// NewByteArray makes and returns a new empty ByteArray.
func NewByteArray() *ByteArray {
	b := make(ByteArray, 0)
	return &b
}

// NewByteArrayFromSlice makes and returns a new ByteArray initialized with
// the elements of the given slice.
//
// The new ByteArray is a simple type cast of the input slice; the slice is
// _not_ copied.
func NewByteArrayFromSlice(slice []byte) *ByteArray {
	b := ByteArray(slice)
	return &b
}

// Get returns the element of the ByteArray at the given index.
//
// It panics if the index is out of range.
func (b *ByteArray) Get(i int) byte {
	return (*b)[i]
}

// Len returns the length of the ByteArray.
func (b *ByteArray) Len() int {
	return len(*b)
}
