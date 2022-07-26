package amorph

func Union(a0, a1 Amorph, ops ...int) (ar Amorph, err error) {
	options := 0
	for _, v := range ops {
		options = v | options
	}
	return union(a0, a1, options)
}

func amorphUnion(a0, a1 Amorph, options int) (Amorph, error) {
	switch {
	case OptUnionSliceResolveAmorph0&options > 0:
		return a0, nil
	case OptUnionSliceResolveAmorph1&options > 0:
		return a1, nil
	case OptUnionSliceAlways&options > 0:
		return []interface{}{a0, a1}, nil
	case OptUnionSliceNotEqual&options > 0 && a0 != a1:
		return []interface{}{a0, a1}, nil
	case OptUnionSliceNotEqual&options > 0 && a0 == a1:
		return a0, nil
	default:
		return []interface{}{a0, a1}, nil
	}
}

func stringUnion(a0, a1 string, options int) (Amorph, error) {
	switch {
	case OptUnionSliceResolveAmorph0&options > 0:
		return a0, nil
	case OptUnionSliceResolveAmorph1&options > 0:
		return a1, nil
	case OptUnionSliceAlways&options > 0:
		return []interface{}{a0, a1}, nil
	case OptUnionSliceNotEqual&options > 0 && a0 != a1:
		return []interface{}{a0, a1}, nil
	case OptUnionSliceNotEqual&options > 0 && a0 == a1:
		return a0, nil
	default:
		return []interface{}{a0, a1}, nil
	}
}

func float64Union(a0, a1 float64, options int) (Amorph, error) {
	switch {
	case OptUnionSliceResolveAmorph0&options > 0:
		return a0, nil
	case OptUnionSliceResolveAmorph1&options > 0:
		return a1, nil
	case OptUnionSliceAlways&options > 0:
		fallthrough
	case OptUnionSliceNotEqual&options > 0 && a0 != a1:
		fallthrough
	default:
		return []interface{}{a0, a1}, nil
	}
}

func mapUnion(a0, a1 Amorph, options int) (Amorph, error) {
	var err error
	keys := make(map[string]struct{})
	for k := range a0.(map[string]interface{}) {
		keys[k] = struct{}{}
	}
	for k := range a1.(map[string]interface{}) {
		keys[k] = struct{}{}
	}
	ar := make(map[string]interface{})
	for k := range keys {
		v0, ok0 := a0.(map[string]interface{})[k]
		v1, ok1 := a1.(map[string]interface{})[k]

		switch {
		case ok0 && !ok1:
			ar[k] = v0
		case !ok0 && ok1:
			ar[k] = v1
		case ok0 && ok1:
			ar[k], err = union(v0, v1, options)
			if err != nil {
				return nil, err
			}
		}
	}
	return ar, nil
}

func sliceUnion(a0, a1 Amorph, options int) (Amorph, error) {
	var err error
	a0slice := a0.([]interface{})
	a1slice := a1.([]interface{})
	long1 := a0slice
	short1 := a1slice

	if len(a0slice) < len(a1slice) {
		long1 = a1slice
		short1 = a0slice
	}
	ar := NewNullSlice(len(long1)).([]interface{})
	for i, v0 := range long1 {
		if i < len(short1) {
			v1 := a1.([]interface{})[i]
			ar[i], err = union(v0, v1, options)
			if err != nil {
				return nil, err
			}
		} else {
			ar[i] = v0
		}
	}
	return ar, nil
}

func union(a0, a1 Amorph, options int) (ar Amorph, err error) {
	switch ca0 := a0.(type) {
	case nullType:
		switch a1.(type) {
		case nullType:
			return NULL, nil
		default:
			return a1, nil
		}
	case nil:
		switch a1.(type) {
		case nullType:
			return a0, nil
		default:
			return amorphUnion(a0, a1, options)
		}
	case string:
		switch ca1 := a1.(type) {
		case nullType:
			return a0, nil
		case string:
			return stringUnion(ca0, ca1, options)
		default:
			return amorphUnion(a0, a1, options)
		}
	case float64:
		switch ca1 := a1.(type) {
		case nullType:
			return a0, nil
		case float64:
			return float64Union(ca0, ca1, options)
		default:
			return amorphUnion(a0, a1, options)
		}
	case map[string]interface{}:
		switch a1.(type) {
		case nullType:
			return a0, nil
		case map[string]interface{}:
			return mapUnion(a0, a1, options)
		default:
			return amorphUnion(a0, a1, options)
		}
	case []interface{}:
		switch a1.(type) {
		case nullType:
			return a0, nil
		case []interface{}:
			return sliceUnion(a0, a1, options)
		default:
			return amorphUnion(a0, a1, options)
		}
	default:
		return nil, ErrUnsupportedType
	}
}
