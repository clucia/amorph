package amorph

// Copyright 2021 Charles J. Luciano and Scalability
// Labs LLC. All rights reserved. Use of this source
// code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"sort"
)

// A WalkPos contains a Position within an Amorph.
type WalkPos interface {
	Key() interface{}
	Value() interface{}
	Rereference(key interface{}, value interface{})
	Dereference(key interface{}) interface{}
}

type genericWalkPos struct {
	key   interface{}
	value interface{}
}

// NewWalkPos creates a new WalkPos.
func NewWalkPos(k, v interface{}) WalkPos {
	return &genericWalkPos{k, v}
}

// Key returns the Key from a WalkPos. It could be
// an int if the WalkPos refers to an []interface{},
// or a string if the WalkPos refers to a map[string]interface{}
func (wp *genericWalkPos) Key() interface{} {
	return wp.key
}

// Key returns the Value from a WalkPos.
func (wp *genericWalkPos) Value() interface{} {
	return wp.value
}

// Rereference allows a walker's user funtion to change the values in
// a map or slice found by the walker.
func (wp *genericWalkPos) Rereference(key interface{}, value interface{}) {
	switch elem := wp.value.(type) {
	case map[string]interface{}:
		elem[key.(string)] = value
	case []interface{}:
		elem[key.(int)] = value
	default:
		panic("Only supported for maps and slices")
	}
}
func (wp *genericWalkPos) Delete(key interface{}) {
	switch elem := wp.value.(type) {
	case map[string]interface{}:
		delete(elem, key.(string))
	default:
		panic("Only supported for maps")
	}
}

// Key returns the dereferences the object referred to by a WalkPos.
func (wp *genericWalkPos) Dereference(key interface{}) interface{} {
	switch elem := wp.value.(type) {
	case map[string]interface{}:
		return elem[key.(string)]
	case []interface{}:
		return elem[key.(int)]
	case string:
		return elem
	default:
		panic("Attempt to dereference bad type")
	}
}

// WalkIter contains the WalkPos from every layer as
// Walk descends the hierarchy.
type WalkIter interface {
	Top() WalkPos                 // Top gets the top (deepest) WalkPos from an iterator
	Pop() WalkIter                // Pop returns a WalkIter with one layer removed
	Len() int                     // Len returns the depth of the Iterator
	Append(elem WalkPos) WalkIter // Append creates a new iterator with a new layer appended
	Copy() WalkIter               // Copy duplicates a WalkIter
}

type sliceWalkIter []WalkPos

// Copy duplicates an iterator
func (wi *sliceWalkIter) Copy() WalkIter {
	nwi := make(sliceWalkIter, len(*wi))
	for i, v := range *wi {
		nwi[i] = v
	}
	return &nwi
}

// Pop returns a WalkIter with one layer removed
func (wi *sliceWalkIter) Pop() WalkIter {
	l := len(*wi)
	if l == 0 {
		panic("Cannot Pop a zero length iterator")
	}
	nwi := (*wi)[:l-1]
	return &nwi
}

// Len returns the depth of the Iterato
func (wi *sliceWalkIter) Len() int {
	return len(*wi)
}

// Top gets the top (deepest) WalkPos from an iterator
func (wi *sliceWalkIter) Top() WalkPos {
	l := len(*wi)
	if l == 0 {
		panic("Shouldn't Top zero length iterator")
	}
	return (*wi)[l-1]
}

// Append creates a new iterator with a new layer appended
func (wi *sliceWalkIter) Append(elem WalkPos) WalkIter {
	(*wi) = append(*wi, elem)
	return wi
}

// NewSliceWalkIter creates a new iterator
func NewSliceWalkIter() WalkIter {
	var new sliceWalkIter = make([]WalkPos, 0)
	return WalkIter(&new)
}

// Walk walks over a hierarchy of map/slice/string/float64 and calls a user
// function for each node.
func Walk(in interface{}, wfunc func(iter WalkIter) error) error {
	iter := NewSliceWalkIter()
	wp := NewWalkPos(nil, in)
	iter = iter.Append(wp)
	return walk(in, iter, wfunc)
}

func walk(val interface{}, iter WalkIter, wfunc func(iter WalkIter) error) (err error) {
	// fmt.Println("len(iter) = ", iter.Len())
	var niter WalkIter

	switch typedIn := val.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(typedIn))
		for k := range typedIn {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := typedIn[k]
			wp := NewWalkPos(k, v)
			niter = iter.Copy()
			niter = niter.Append(wp)
			err = wfunc(niter)
			if err != nil {
				return //
			}
			err = walk(v, niter, wfunc)
			if err != nil {
				return //
			}
		}
		return err
	case []interface{}:
		for i, v := range typedIn {
			wp := NewWalkPos(i, v)
			niter = iter.Copy()
			niter := niter.Append(wp)
			err = wfunc(niter)
			if err != nil {
				return //
			}
			err = walk(v, niter, wfunc)
			if err != nil {
				return //
			}
		}
		return //
	default:
		niter = NewSliceWalkIter().Append(NewWalkPos(nil, val))
		err = wfunc(niter)
		if err != nil {
			return //
		}
	}
	return //
}
