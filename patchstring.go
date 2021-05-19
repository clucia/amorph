package amorph

// Copyright 2021 Charles J. Luciano and Scalability
// Labs LLC. All rights reserved. Use of this source
// code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"strconv"
)

// PatchStringer converts a patch to a text representation
// mainly for debugging purposes
func PatchStringer(patch Patch) string {
	return describe(patch, "") + "\n"
}

func describe(patch Patch, indent string) (s string) {
	var typ string
	var ok bool
	typ, _, _, _, ok = unpack(DirFwd, patch)
	if !ok {
		return "error in describe"
	}
	switch {
	case patch == nil:
		return indent + "nil\n"
	case typ == "raw":
		return rawDescribe(patch, indent)
	case typ == "float64":
		return float64Describe(patch, indent)
	case typ == "string":
		return stringDescribe(patch, indent)
	case typ == "slice":
		return sliceDescribe(patch, indent)
	case typ == "map":
		return mapDescribe(patch, indent)
	default:
		panic("Malformed patch")
	}
}

func mapDescribe(patch Patch, indent string) (s string) {
	var typ string
	var valFwd interface{}
	var ok bool
	typ, _, valFwd, _, ok = unpack(DirFwd, patch)
	if !ok {
		return "error in describe"
	}
	patchMap, ok := valFwd.(map[string]interface{})
	if !ok {
		panic("Malformed patch")
	}
	if len(patchMap) == 0 {
		s += indent + "Empty mapPatch\n"
		return //
	}
	s += indent + " typ = " + typ + "\n"
	for k, v := range patchMap {
		istr := indent + "." + k
		s += describe(v, istr)
	}
	return //
}

func sliceDescribe(patch Patch, indent string) (s string) {
	var typ string
	var valFwd interface{}
	var ok bool
	typ, _, valFwd, _, ok = unpack(DirFwd, patch)
	if !ok {
		return "error in describe"
	}
	_ = typ
	patchSlice, ok := valFwd.([]Patch)
	if !ok {
		panic("Malformed patch")
	}
	if len(patchSlice) == 0 {
		s += indent + "Empty slicePatch\n"
		return //
	}
	for i, v := range patchSlice {
		istr := indent + "[" + strconv.FormatInt(int64(i), 10) + "]"

		s += describe(v, istr)
	}
	return //
}

func float64Describe(patch Patch, indent string) (s string) {
	var typ0, typ1 string
	var valFwd, valRev interface{}
	var deleteFwd, deleteRev bool
	var ok bool
	typ0, _, valFwd, deleteFwd, ok = unpack(DirFwd, patch)
	if !ok {
		return "error in describe"
	}
	typ1, _, valRev, deleteRev, ok = unpack(DirRev, patch)
	if !ok {
		return "error in describe"
	}
	if typ0 != typ1 {
		typ0 = "type mismatch"
	}
	s += indent + "float64Patch::describe:" +
		" typ = " + typ0 +
		", deleteFwd = " + strconv.FormatBool(deleteFwd) +
		", deleteRev = " + strconv.FormatBool(deleteRev) +
		", valFwd = " + strconv.FormatFloat(valFwd.(float64), 'f', 0, 64) +
		", valRev = " + strconv.FormatFloat(valRev.(float64), 'f', 0, 64) +
		"\n"
	return //
}

func stringDescribe(patch Patch, indent string) (s string) {
	var typ0, typ1 string
	var valFwd, valRev interface{}
	var deleteFwd, deleteRev bool
	var ok bool
	typ0, _, valFwd, deleteFwd, ok = unpack(DirFwd, patch)
	if !ok {
		return "error in describe"
	}
	typ1, _, valRev, deleteRev, ok = unpack(DirRev, patch)
	if !ok {
		return "error in describe"
	}
	if typ0 != typ1 {
		typ0 = "type mismatch"
	}
	s += indent +
		" typ = " + typ0 +
		", deleteFwd = " + strconv.FormatBool(deleteFwd) +
		", deleteRev = " + strconv.FormatBool(deleteRev) +
		", valFwd = " + valFwd.(string) +
		", valRev = " + valRev.(string) +
		"\n"
	return //
}

func rawDescribe(patch Patch, indent string) (s string) {
	var typ0, typ1 string
	var valFwd, valRev interface{}
	var deleteFwd, deleteRev bool
	var ok bool
	typ0, _, valFwd, deleteFwd, ok = unpack(DirFwd, patch)
	if !ok {
		return "error in describe"
	}
	typ1, _, valRev, deleteRev, ok = unpack(DirRev, patch)
	if typ0 != typ1 {
		typ0 = "type mismatch"
	}
	if !ok {
		return "error in describe"
	}
	s += indent + "rawPatch::describe:" +
		" typ = " + typ0 +
		", deleteFwd = " + strconv.FormatBool(deleteFwd) +
		", deleteRev = " + strconv.FormatBool(deleteRev) +
		", valFwd = " + fmt.Sprintf("%v", valFwd) +
		", valRev = " + fmt.Sprintf("%v", valRev) +
		"\n"
	return //
}
