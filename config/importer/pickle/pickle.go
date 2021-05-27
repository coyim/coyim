package pickle

import (
	"fmt"
	"path/filepath"

	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/nlpodyssey/gopickle/types"
)

// WrongTypeError represents a wrong pickle type error
type WrongTypeError struct {
	Result  interface{}
	Request string
}

func (wte WrongTypeError) Error() string {
	return fmt.Sprintf("Unpickling returned type %T which cannot be converted to %s", wte.Result, wte.Request)
}

func newWrongTypeError(result interface{}, request interface{}) error {
	return WrongTypeError{Result: result, Request: fmt.Sprintf("%T", request)}
}

// Load a value from a file reader. This function takes a file path and
// attempts to read a complete pickle program from it.
func Load(file string) (interface{}, error) {
	return pickle.Load(filepath.Clean(file))
}

// Dict attempts to convert the return value of Load into a types.Dict.
func Dict(v interface{}) (types.Dict, error) {
	d, ok := v.(*types.Dict)
	if !ok {
		return nil, newWrongTypeError(v, d)
	}

	var d2 types.Dict = *d
	return d2, nil
}

// DictString attempts to convert the return value of Load into map[string]interface{}.
func DictString(v interface{}) (map[string]interface{}, error) {
	d, err := Dict(v)
	if err != nil {
		return nil, err
	}

	return tryDictToDictString(d)
}

func tryDictToDictString(dd types.Dict) (map[string]interface{}, error) {
	r := map[string]interface{}{}

	for _, e := range dd {
		kstr, ok := e.Key.(string)
		if !ok {
			return nil, newWrongTypeError(dd, kstr)
		}
		r[kstr] = e.Value
	}

	return r, nil
}
