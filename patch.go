package amorph

import "fmt"

// Copyright 2021 Charles J. Luciano and Scalability
// Labs LLC. All rights reserved. Use of this source
// code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Patch is an object that can change an Amorph
type Patch interface{}

const (
	DirFwd = "Fwd"
	DirRev = "Rev"
)

// ApplyFwd duplicates the input Amorph to the output Amorph with the
// differences in patch applied
func ApplyFwd(patch Patch, amorphIn Amorph) (absout Amorph, err error) {
	return apply(DirFwd, patch, amorphIn)
}

// ApplyFwd duplicates the input Amorph to the output Amorph with the
// differences in patch REVERSE applied
func ApplyRev(patch Patch, amorphIn Amorph) (amorphOut Amorph, err error) {
	return apply(DirRev, patch, amorphIn)
}

func apply(dir string, patch Patch, amorphIn Amorph) (amorphOut Amorph, err error) {
	typ, _, valX, _, ok := unpack(dir, patch)
	if !ok {
		return nil, fmt.Errorf("unpack " + dir + " error")
	}
	switch {
	case typ == "string":
		amorphOut = valX
		return //
	case typ == "float64":
		amorphOut = valX
		return //
	case typ == "slice":
		return sliceApply(dir, patch, amorphIn)
	case typ == "map":
		return mapApply(dir, patch, amorphIn)
	case typ == "raw":
		return rawApply(dir, patch, amorphIn)
	default:
		panic("Malformed Patch")
	}
}

func rawApply(dir string, patch Patch, amorphIn Amorph) (amorphOut Amorph, err error) {
	var valX interface{}
	var ok bool
	_, _, valX, _, ok = unpack(dir, patch)
	if !ok {
		return nil, fmt.Errorf("bad type")
	}
	return valX, nil
}

// unpack gets the type of patch, lengtn, value, and delete flag for a patch.
// only slice patches have a length
// map and slice patches both get their value from valFwd because valFwd
// contains the information needed to patch in both directions.
func unpack(dir string, ipatch Patch) (
	typ string,
	lenX int,
	valX interface{},
	deleteX bool,
	ok bool,
) {
	if ipatch == nil {
		return "nil", 0, nil, false, true
	}
	patch := ipatch.(map[string]interface{})

	t0, ok := patch["typ"]
	if !ok {
		return //
	}
	typ = t0.(string)

	f0, ok := patch["delete"+dir]
	if ok {
		deleteX = f0 == "true"
	}

	switch {
	case typ == "map":
		fallthrough
	case typ == "slice":
		f0, ok = patch["valFwd"]
	default:
		f0, ok = patch["val"+dir]
	}
	if ok {
		valX = f0
	}
	f2, ok := patch["len"+dir]
	if ok && f2 != nil {
		lenX = f2.(int)
	}
	ok = true
	return
}

// sliceApply duplicates the input Amorph (a slice) to the output Amorph
// with the changes applied from patch
func sliceApply(dir string, ipatch Patch, amorphIn Amorph) (amorphOut Amorph, err error) {
	var typ string
	var lenX int
	var valX interface{}
	var deleteX, ok bool

	typ, lenX, valX, deleteX, ok = unpack(dir, ipatch)
	if !ok {
		return nil, fmt.Errorf("bad type")
	}
	if typ != "slice" {
		panic("Can't happen")
	}
	var abs2 []interface{}

	if amorphIn == nil {
		abs2 = make([]interface{}, lenX)
	} else {
		abs2 = amorphIn.([]interface{})
		for len(abs2) < lenX {
			abs2 = append(abs2, nil)
		}
		abs2 = abs2[:lenX]
	}
	replace := valX.([]Patch)
	for i := 0; i < lenX; i++ {
		if replace[i] == nil {
			continue
		}
		elementPatch, ok := replace[i].(Patch)
		if !ok {
			panic("Malformed Patch")
		}
		if deleteX {
			if i < len(abs2) {
				abs2[i] = nil // TODO needed?
			}
			continue
		}
		if i < len(amorphIn.([]interface{})) {
			abs2[i], err = apply(dir, elementPatch, amorphIn.([]interface{})[i])
		} else {
			abs2[i], err = apply(dir, elementPatch, nil)
		}
		if err != nil {
			return //
		}
	}
	return abs2, nil
}

func mapApply(dir string, ipatch Patch, amorphIn Amorph) (amorphOut Amorph, err error) {
	var typ string
	var valX interface{}
	var deleteX, ok bool
	typ, _, valX, deleteX, ok = unpack(dir, ipatch)
	if !ok {
		return nil, fmt.Errorf("bad type")
	}
	if typ != "map" {
		panic("Malformed patch")
	}
	mapIn, ok := amorphIn.(map[string]interface{})
	if !ok {
		panic("Malformed patch")
	}
	var mapOut map[string]interface{}
	if deleteX {
		return nil, nil
	}

	if amorphIn == nil {
		mapOut = make(map[string]interface{})
	} else {
		mapOut = mapIn
	}
	patchMap, ok := valX.(map[string]interface{})
	if !ok {
		panic("Malformed patch")
	}
	for k, v := range patchMap {
		patch, ok := v.(Patch)
		if !ok {
			panic("Malformed patch")
		}
		del, _ := patch.(map[string]interface{})["delete"+dir].(bool)
		if del {
			delete(mapOut, k)
			continue
		}
		_, ok = mapOut[k]
		if ok {
			mapOut[k], err = apply(dir, patch, mapOut[k])
		} else {
			mapOut[k], err = apply(dir, patch, nil)
		}
		if err != nil {
			return //
		}
	}
	return mapOut, nil
}
