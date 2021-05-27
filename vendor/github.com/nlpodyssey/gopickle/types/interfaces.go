// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// Callable is implemented by any value that can be directly called to get a
// new value.
//
// It is usually implemented by Python-like functions (returning a value
// given some arguments), or classes (typically returning an instance given
// some constructor arguments).
type Callable interface {
	// Call mimics a direct invocation on a Python value, such as a function
	// or class (constructor).
	Call(args ...interface{}) (interface{}, error)
}

// PyNewable is implemented by any value that has a Python-like
// "__new__" method.
//
// It is usually implemented by values representing Python classes.
type PyNewable interface {
	// PyNew mimics Python invocation of the "__new__" method, usually
	// provided by classes.
	//
	// See: https://docs.python.org/3/reference/datamodel.html#object.__new__
	PyNew(args ...interface{}) (interface{}, error)
}

// PyStateSettable is implemented by any value that has a Python-like
// "__setstate__" method.
type PyStateSettable interface {
	// PySetState mimics Python invocation of the "__setstate__" method.
	//
	// See: https://docs.python.org/3/library/pickle.html#object.__setstate__
	PySetState(state interface{}) error
}

// PyDictSettable is implemented by any value that can store dictionary-like
// key/value pairs. It reflects Python behavior of setting a key/value pair on
// an object's "__dict__" attribute.
type PyDictSettable interface {
	// PyDictSet mimics the setting of a key/value pair on an object's
	//"__dict__" attribute.
	//
	// See: https://docs.python.org/3/library/stdtypes.html#object.__dict__
	PyDictSet(key, value interface{}) error
}

// PyAttrSettable is implemented by any value on which an existing or new
// Python-like attribute can be set. In Python this is done with "setattr"
// builtin function.
type PyAttrSettable interface {
	// PySetAttr mimics the setting of an arbitrary value to an object's
	// attribute.
	//
	// In Python this is done with "setattr" function, to which object,
	// attribute name, and value are passed. For an easy and clear
	// implementation, here instead we require this method to be implemented
	// on the "object" itself.
	//
	// See: https://docs.python.org/3/library/functions.html#setattr
	PySetAttr(key string, value interface{}) error
}
