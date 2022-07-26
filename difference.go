package amorph

// Difference produces an Amorph that contains everything from the minuend
// that isn't subtracted by the subtrahend.
//
// By default, the subtrahend must match all the way to the leaf values for a subtraction
// to occur.
//
// By default, elements in the subtrahend that have no topological match in the minuend
// are ignored. The OptMustSubtract option causes Difference to treat this as an error.
//
// The OptIgnoreLeafValues option says that the subtraction may occur if there's ANY
// value in topologically the same place in the two Amorphs.
//
func Difference(min, a1 Amorph, ops ...int) (Amorph, error) {
	options := 0
	for _, v := range ops {
		options = options | v
	}
	return difference(min, a1, options)
}

func sliceDifference(min, sub Amorph, options int) (Amorph, error) {
	m := min.([]interface{})
	s := sub.([]interface{})

	ar := NewNullSlice(len(m)).([]interface{})

	if len(s) > len(m) && (OptDifferenceMustSubtract&options) > 0 {
		return nil, ErrMustSubtract
	}

	var err error
	for i, v := range m {
		if i < len(s) {
			ar[i], err = difference(v, s[i], options)
			if err != nil {
				return nil, err
			}
		} else {
			ar[i] = m[i]
		}
	}

	return ar, nil
}

func mapDifference(min, sub Amorph, options int) (Amorph, error) {
	m := min.(map[string]interface{})
	s := sub.(map[string]interface{})
	keys := make(map[string]struct{})
	ar := make(map[string]interface{})
	for k := range m {
		keys[k] = struct{}{}
	}
	for k := range s {
		keys[k] = struct{}{}
	}
	var err error
	for k := range keys {
		v0, ok0 := m[k]
		v1, ok1 := s[k]
		switch {
		case ok0 && !ok1:
			ar[k] = m[k]
		case !ok0 && ok1:
			if OptDifferenceMustSubtract&options > 0 {
				return nil, ErrMustSubtract
			}
		case ok0 && ok1:
			var res Amorph
			res, err = difference(v0, v1, options)
			if err != nil {
				return nil, err
			}
			if res == NULL {
				continue
			}
			ar[k] = res
		default:
			panic("")

		}
	}
	return ar, nil
}

func difference(min, sub Amorph, options int) (Amorph, error) {
	switch min.(type) {
	case nullType:
		return NULL, nil
	case nil:
		switch sub.(type) {
		case nil:
			return NULL, nil
		default:
			if OptDifferenceMustSubtract&options > 0 {
				return NULL, ErrMustSubtract
			}
			return nil, nil
		}
	case string:
		switch sub.(type) {
		case string:
			if min == sub {
				return NULL, nil
			}
		}
		if OptDifferenceMustSubtract&options > 0 {
			return NULL, ErrMustSubtract
		}
		return min, nil
	case float64:
		switch sub.(type) {
		case float64:
			if min == sub {
				return NULL, nil
			}
		}
		if OptDifferenceMustSubtract&options > 0 {
			return NULL, ErrMustSubtract
		}
		return min, nil
	case map[string]interface{}:
		switch sub.(type) {
		case map[string]interface{}:
			return mapDifference(min, sub, options)
		}
		if OptDifferenceMustSubtract&options > 0 {
			return nil, ErrMustSubtract
		}
		return min, nil
	case []interface{}:
		switch sub.(type) {
		case []interface{}:
			return sliceDifference(min, sub, options)
		}
		if OptDifferenceMustSubtract&options > 0 {
			return nil, ErrMustSubtract
		}
		return min, nil
	default:
		panic("")
	}
}
