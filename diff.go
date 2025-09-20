package amorph

// Copyright 2021 Charles J. Luciano and Scalability
// Labs LLC. All rights reserved. Use of this source
// code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Diff compares two Amorphs and creates a Patch that describes
// the differences. The patch can be used to change one into the
// other.
//
// The differences from amorph0 to amorph1 are considered to be
// the Forward direction.
func Diff(amorph0, amorph1 Amorph) (patch Patch) {
	switch cvtd0 := amorph0.(type) {
	case nullType:
		return map[string]interface{}{
			"typ":       "raw",
			"deleteRev": amorph0 == nil,
			"deleteFwd": amorph1 == nil,
			"valFwd":    amorph1,
			"valRev":    amorph0,
		}
	case nil:
		return map[string]interface{}{
			"typ":       "raw",
			"deleteRev": true,
			"deleteFwd": amorph1 == nil,
			"valFwd":    amorph1,
		}
	case float64:
		return float64Diff(cvtd0, amorph1)
	case string:
		return stringDiff(cvtd0, amorph1)
	case []interface{}:
		return sliceDiff(cvtd0, amorph1)
	case map[string]interface{}:
		return mapDiff(cvtd0, amorph1)
	default:
		return map[string]interface{}{
			"typ":       "raw",
			"deleteFwd": amorph1 == nil,
			"deleteRev": amorph0 == nil,
			"valFwd":    amorph1,
			"valRev":    amorph0,
		}
	}
}

func Diff0(amorph0, amorph1 Amorph) (patch Patch) {
	switch cvtd0 := amorph0.(type) {
	case nil:
		return map[string]interface{}{
			"typ":       "raw",
			"deleteRev": true,
			"deleteFwd": amorph1 == nil,
			"valFwd":    amorph1,
		}
	case float64:
		return float64Diff(cvtd0, amorph1)
	case string:
		return stringDiff(cvtd0, amorph1)
	case []interface{}:
		return sliceDiff(cvtd0, amorph1)
	case map[string]interface{}:
		return mapDiff(cvtd0, amorph1)
	default:
		return map[string]interface{}{
			"typ":       "raw",
			"deleteFwd": amorph1 == nil,
			"deleteRev": amorph0 == nil,
			"valFwd":    amorph1,
			"valRev":    amorph0,
		}
	}
}

func float64Diff(float0 float64, amorph1 Amorph) (patch Patch) {
	if amorph1 == nil {
		return map[string]interface{}{
			"typ":       "raw",
			"deleteFwd": true,
			"valRev":    float0,
		}
	}
	float1, ok := amorph1.(float64)
	if !ok {
		return map[string]interface{}{
			"typ":    "raw",
			"valFwd": amorph1,
			"valRev": float0,
		}
	}
	if float0 == float1 {
		return nil
	}
	return map[string]interface{}{
		"typ":    "float64",
		"valFwd": float1,
		"valRev": float0,
	}
}

func stringDiff(str0 string, amorph1 Amorph) (patch Patch) {
	if amorph1 == nil {
		return map[string]interface{}{
			"typ":       "raw",
			"deleteFwd": true,
			"valRev":    str0,
		}
	}
	str1, ok := amorph1.(string)
	if !ok {
		return map[string]interface{}{
			"typ":    "raw",
			"valFwd": amorph1,
			"valRev": str0,
		}
	}
	if str0 == str1 {
		return nil
	}
	return map[string]interface{}{
		"typ":    "string",
		"valFwd": str1,
		"valRev": str0,
	}
}

func mapDiff(map0 map[string]interface{}, amorph1 Amorph) (patch Patch) {
	prune := true
	if map0 == nil {
		panic("Shouldn't happen")
	}
	if amorph1 == nil {
		return map[string]interface{}{
			"typ":       "raw",
			"deleteFwd": true,
			"valRev":    map0,
		}
	}
	map1, ok := amorph1.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"typ":    "raw",
			"valFwd": amorph1,
			"valRev": map0,
		}
	}
	keys := make(map[string]struct{})
	for k := range map0 {
		keys[k] = struct{}{}
	}
	for k := range map1 {
		keys[k] = struct{}{}
	}
	mapPatch := map[string]interface{}{
		"typ":    "map",
		"valFwd": make(map[string]interface{}),
	}
	for k := range keys {
		var elemPatch Patch
		_, ok0 := map0[k]
		_, ok1 := map1[k]
		switch {
		case !ok0 && !ok1:
			panic("Shouldn't happen")
		}
		if ok0 && ok1 {
			elemPatch = Diff(map0[k], map1[k])
			if elemPatch == nil {
				continue
			}
		} else {
			elemPatch = map[string]interface{}{
				"typ":       "raw",
				"deleteFwd": !ok1,
				"deleteRev": !ok0,
				"valFwd":    map1[k],
				"valRev":    map0[k],
			}
		}
		if elemPatch != nil {
			prune = false
		}
		mapPatch["valFwd"].(map[string]interface{})[k] = elemPatch
	}
	if prune {
		return nil
	}
	return mapPatch
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sliceDiff(slice0 []interface{}, amorph1 Amorph) (patch Patch) {
	prune := true
	if slice0 == nil {
		panic("Shouldn't happen")
	}
	if amorph1 == nil {
		return map[string]interface{}{
			"typ":       "raw",
			"deleteFwd": true,
			"valRev":    slice0,
			"lenRev":    len(slice0),
		}
	}
	slice1, ok := amorph1.([]interface{})
	if !ok {
		return map[string]interface{}{
			"typ":    "slice",
			"valFwd": amorph1,
			"valRev": slice0,
			"lenRev": len(slice0),
		}
	}
	l0 := len(slice0)
	l1 := len(slice1)
	lmax := max(l0, l1)
	slicePatch := map[string]interface{}{
		"typ":    "slice",
		"valFwd": make([]Patch, lmax),
		"lenFwd": l1,
		"lenRev": l0,
	}
	for i := 0; i < lmax; i++ {
		var elementPatch Patch
		switch {
		case i < l0 && i < l1:
			elementPatch = Diff(slice0[i], slice1[i])
		case i < l0:
			elementPatch = map[string]interface{}{
				"typ":       "raw",
				"valRev":    slice0[i],
				"deleteFwd": true,
			}
		case i < l1:
			elementPatch = map[string]interface{}{
				"typ":       "raw",
				"valFwd":    slice1[i],
				"deleteRev": true,
			}
		default:
			panic("Shouldn't happen")
		}
		if elementPatch != nil {
			prune = false
		}
		slicePatch["valFwd"].([]Patch)[i] = elementPatch
	}
	if prune {
		return nil
	}
	return slicePatch
}
