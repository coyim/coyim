package stalecucumber

import "fmt"
import "strings"

type UnparseablePythonGlobalError struct {
	Args interface{}
	Message string
}

func(this UnparseablePythonGlobalError) Error() string {
	return fmt.Sprintf("%s; arguments (%T): %v", this.Message, this.Args, this.Args)
}

type PythonBuiltinResolver struct {}

func (this PythonBuiltinResolver) Resolve(module string, name string, args []interface{}) (interface{}, error) {
	// Up to version 2 this is always "__builtin__"
	// In version 3+ it becomes "builtin" but that is not supported yet
	if module != "__builtin__" {
		return nil, ErrUnresolvablePythonGlobal
	}

	if name == "set" {
		return this.handlePythonSet(args)
	}

	if name == "bytearray" {
		return this.handlePythonByteArray(args)
	}


	return nil, ErrUnresolvablePythonGlobal
}

func (this PythonBuiltinResolver) handlePythonSet(args []interface{}) (interface{}, error){
	
	if len(args) != 1 {
		return nil, UnparseablePythonGlobalError{
			Args: args, 
			Message: "Expected args to be of length 1",
		}
	}

	tuple, ok := args[0].([]interface{})
	if !ok {
		return nil, UnparseablePythonGlobalError{
			Args: args, 
			Message: fmt.Sprintf("Expected first argument of args to be of type %T", tuple),
		}
	}

	// A map is the equivalent golang type for a python set
	set := make(map[interface{}]bool, len(tuple))
	for _, item := range tuple {
		set[item] = true
	}

	return set, nil
}

func (this PythonBuiltinResolver) handlePythonByteArray(args []interface{}) (interface{}, error){	
	// Up to version 2 the implementation of bytearray always pickles as a tuple like
	// (theStringValue, 'latin-1', )
	// version 3+ of pickle is different but we're not supporting that presently
	if len(args) != 2{
		return nil, UnparseablePythonGlobalError{
			Args: args,
			Message: "Expected args to be of length 2",
		}
	}

	const magic = `latin-1`
	magicValue, ok := args[1].(string)
	if !ok || magicValue != magic{
		return nil, UnparseablePythonGlobalError{
			Args: args,
			Message: fmt.Sprintf("Expected second arg to be string %q", magic),
		}
	}

	value, ok := args[0].(string)
	if !ok {
		return nil, UnparseablePythonGlobalError{
			Args: args,
			Message: "Expected first arg to be a string",
		}
	}
	return strings.NewReader(value), nil
}
 