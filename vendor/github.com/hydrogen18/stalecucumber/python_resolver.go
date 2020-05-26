package stalecucumber

import "errors"

// A type to convert to a GLOBAL opcode to something meaningful in golang
type PythonResolver interface {
	Resolve(module string, name string, args []interface{}) (interface{}, error)
}

var ErrUnresolvablePythonGlobal = errors.New("Unresolvable Python global value")

type PythonResolverChain []PythonResolver

func MakePythonResolverChain(args... PythonResolver) PythonResolverChain{
	if len(args) == 0 {
		panic("Cannot make chain of length zero")
	}
	return PythonResolverChain(args)
}

func (this PythonResolverChain) Resolve(module string, name string, args []interface{})(interface{}, error){
	var err error
	for _, resolver := range this {
		var result interface{}
		result, err = resolver.Resolve(module, name, args)
		if err == nil {
			return result, nil
		}
		if err != ErrUnresolvablePythonGlobal {
			return nil, err
		}
	}

	return nil, err
}
