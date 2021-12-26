package gui

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

type baseStoreField struct {
	store     gtki.TreeStore
	index     int
	doOnError func(error)
}

func (s *baseStoreField) onError(f func(error)) {
	s.doOnError = f
}

func (s *baseStoreField) handlePotentialError(e error) {
	if e != nil && s.doOnError != nil {
		s.doOnError(e)
	}
}

func (s *baseStoreField) baseGet(iter gtki.TreeIter) interface{} {
	untypedResult, e := s.store.GetValue(iter, s.index)
	s.handlePotentialError(e)
	result, e2 := untypedResult.GoValue()
	s.handlePotentialError(e2)
	return result
}

func (s *baseStoreField) baseGetWithError(iter gtki.TreeIter) (interface{}, error) {
	untypedResult, e := s.store.GetValue(iter, s.index)
	if e != nil {
		return nil, e
	}
	result, e2 := untypedResult.GoValue()
	if e2 != nil {
		return nil, e2
	}
	return result, nil
}

type stringStoreField struct {
	*baseStoreField
}

func newBaseStoreField(store gtki.TreeStore, index int) *baseStoreField {
	return &baseStoreField{
		store: store,
		index: index,
	}
}

func newStringStoreField(store gtki.TreeStore, index int) *stringStoreField {
	return &stringStoreField{newBaseStoreField(store, index)}
}

func (s *stringStoreField) set(iter gtki.TreeIter, value string) {
	s.handlePotentialError(s.store.SetValue(iter, s.index, value))
}

func (s *stringStoreField) get(iter gtki.TreeIter) string {
	return s.baseGet(iter).(string)
}

func (s *stringStoreField) getWithError(iter gtki.TreeIter) (string, error) {
	res, e := s.baseGetWithError(iter)
	if e != nil {
		return "", e
	}
	return res.(string), nil
}

type intStoreField struct {
	*baseStoreField
}

func newIntStoreField(store gtki.TreeStore, index int) *intStoreField {
	return &intStoreField{newBaseStoreField(store, index)}
}

func (s *intStoreField) set(iter gtki.TreeIter, value int) {
	s.handlePotentialError(s.store.SetValue(iter, s.index, value))
}

func (s *intStoreField) get(iter gtki.TreeIter) int {
	return s.baseGet(iter).(int)
}

func (s *intStoreField) getWithError(iter gtki.TreeIter) (int, error) {
	res, e := s.baseGetWithError(iter)
	if e != nil {
		return 0, e
	}
	return res.(int), nil
}

type pixbufStoreField struct {
	*baseStoreField
}

func newPixbufStoreField(store gtki.TreeStore, index int) *pixbufStoreField {
	return &pixbufStoreField{newBaseStoreField(store, index)}
}

func (s *pixbufStoreField) set(iter gtki.TreeIter, value gdki.Pixbuf) {
	s.handlePotentialError(s.store.SetValue(iter, s.index, value))
}

func (s *pixbufStoreField) get(iter gtki.TreeIter) gdki.Pixbuf {
	return s.baseGet(iter).(gdki.Pixbuf)
}

func (s *pixbufStoreField) getWithError(iter gtki.TreeIter) (gdki.Pixbuf, error) {
	res, e := s.baseGetWithError(iter)
	if e != nil {
		return nil, e
	}
	return res.(gdki.Pixbuf), nil
}
