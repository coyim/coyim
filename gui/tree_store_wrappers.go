package gui

import "github.com/coyim/gotk3adapter/gtki"

type stringStoreField struct {
	store     gtki.TreeStore
	index     int
	doOnError func(error)
}

func newStringStoreField(store gtki.TreeStore, index int) *stringStoreField {
	return &stringStoreField{
		store: store,
		index: index,
	}
}

func (s *stringStoreField) onError(f func(error)) {
	s.doOnError = f
}

func (s *stringStoreField) handlePotentialError(e error) {
	if e != nil && s.doOnError != nil {
		s.doOnError(e)
	}
}

func (s *stringStoreField) set(iter gtki.TreeIter, value string) {
	s.handlePotentialError(s.store.SetValue(iter, s.index, value))
}

func (s *stringStoreField) get(iter gtki.TreeIter) string {
	untypedResult, e := s.store.GetValue(iter, s.index)
	s.handlePotentialError(e)
	result, e2 := untypedResult.GetString()
	s.handlePotentialError(e2)
	return result
}
