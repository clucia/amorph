package amorph

// Copyright 2021 Charles J. Luciano and Scalability
// Labs LLC. All rights reserved. Use of this source
// code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "reflect"

// Tests two Amorphs for equality

func DeepEqual(amorph0, amorph1 interface{}) bool {
	switch {
	case amorph0 == nil && amorph1 == nil:
		return true
	case amorph0 == nil:
		return false
	case amorph1 == nil:
		return false
	}
	switch cvt0 := amorph0.(type) {
	case string:
		return stringCmp(cvt0, amorph1)
	case float64:
		return float64Cmp(cvt0, amorph1)
	case []interface{}:
		return sliceCmp(cvt0, amorph1)
	case map[string]interface{}:
		return mapCmp(cvt0, amorph1)
	default:
		return reflect.DeepEqual(amorph0, amorph1)
	}
}

func sliceCmp(slice0 []interface{}, amorph1 Amorph) bool {
	slice1, ok := amorph1.([]interface{})
	if !ok {
		return false
	}
	if len(slice0) != len(slice1) {
		return false
	}
	for i := range slice0 {
		res := DeepEqual(slice0[i], slice1[i])
		if !res {
			return false
		}
	}
	return true
}

func mapCmp(map0 map[string]interface{}, amorph1 Amorph) bool {
	map1, ok := amorph1.(map[string]interface{})
	if !ok {
		return false
	}
	if len(map0) != len(map1) {
		return false
	}
	for k := range map0 {
		_, ok0 := map0[k]
		_, ok1 := map1[k]
		if !ok0 || !ok1 {
			return false
		}
		if !DeepEqual(map0[k], map1[k]) {
			return false
		}
	}
	return true
}
func stringCmp(str0 string, amorph1 interface{}) bool {
	str1, ok := amorph1.(string)
	if !ok {
		return false
	}
	return str0 == str1
}

func float64Cmp(float0 float64, amorph1 interface{}) bool {
	float1, ok := amorph1.(float64)
	if !ok {
		return false
	}
	return float0 == float1
}
