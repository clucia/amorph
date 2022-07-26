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
func TopoDifference(min, a1 Amorph, ops ...int) (Amorph, error) {
	options := 0
	for _, v := range ops {
		options = options | v
	}
	return topoDifference(min, a1, options)
}

func sliceTopoDifference(min, sub Amorph, options int) (Amorph, error) {
	m := min.([]interface{})
	s := sub.([]interface{})

	ar := NewNullSlice(len(m)).([]interface{})

	if len(s) > len(m) && (OptTopoDifferenceMustSubtract&options) > 0 {
		return NULL, ErrMustSubtract
	}

	var err error
	for i, v := range m {
		if i < len(s) {
			ar[i], err = topoDifference(v, s[i], options)
			if err != nil {
				return NULL, err
			}
		} else {
			ar[i] = m[i]
		}
	}

	return ar, nil
}

func mapTopoDifference(min, sub Amorph, options int) (Amorph, error) {
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
		minv, minOk := m[k]
		subv, subOk := s[k]
		switch {
		case minOk && !subOk:
			ar[k] = m[k]
		case !minOk && subOk:
			if OptTopoDifferenceMustSubtract&options > 0 {
				return NULL, ErrMustSubtract
			}
		case minOk && subOk:
			var res Amorph
			res, err = topoDifference(minv, subv, options)
			if err != nil {
				return NULL, err
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

func topoDifference(min, sub Amorph, options int) (Amorph, error) {
	switch min.(type) {
	case nullType:
		if OptTopoDifferenceMustSubtract&options > 0 {
			return NULL, ErrMustSubtract
		}
		return NULL, nil
	case nil:
		switch sub.(type) {
		case nullType:
			return min, nil
		case nil:
			return NULL, nil
		default:
			if OptTopoDifferenceMustSubtract&options > 0 {
				return NULL, ErrMustSubtract
			}
			return NULL, nil
		}
	case string:
		switch sub.(type) {
		case nullType:
			return min, nil
		default:
			return NULL, nil
		}
	case float64:
		switch sub.(type) {
		case nullType:
			return min, nil
		default:
			return NULL, nil
		}
	case map[string]interface{}:
		switch sub.(type) {
		case nullType:
			return min, nil
		case map[string]interface{}:
			return mapTopoDifference(min, sub, options)
		default:
			return NULL, nil
		}
	case []interface{}:
		switch sub.(type) {
		case nullType:
			return min, nil
		case []interface{}:
			return sliceTopoDifference(min, sub, options)
		default:
			return NULL, nil
		}
	default:
		panic("")
	}
}
