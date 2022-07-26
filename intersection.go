package amorph

// Intersection produces an Amorph that contains everything common to the
// two Amorphs
func mapIntersection(a0, a1 Amorph, options int) (Amorph, error) {
	ar := make(map[string]interface{})
	a0map := a0.(map[string]interface{})
	a1map := a1.(map[string]interface{})

	for a0key, a0elem := range a0map {
		a1elem, a1ok := a1map[a0key]
		if !a1ok {
			continue
		}
		arelem, err := intersection(a0elem, a1elem, options)
		if err != nil {
			return nil, err
		}
		if arelem == NULL {
			continue
		}
		ar[a0key] = arelem
	}
	return ar, nil
}

func sliceIntersection(a0, a1 Amorph, options int) (Amorph, error) {
	a0Slice := a0.([]interface{})
	a1Slice := a1.([]interface{})
	l0 := len(a0Slice)
	l1 := len(a1Slice)
	min := l1
	if l0 < l1 {
		min = l0
	}

	var err error
	ar := NewNullSlice(min).([]interface{})
	for i := 0; i < min; i++ {
		ar[i], err = intersection(a0Slice[i], a1Slice[i], options)
		if err != nil {
			return nil, err
		}
	}
	return ar, nil
}

func Intersection(a0, a1 Amorph, ops ...int) (ar Amorph, err error) {
	options := 0
	for _, v := range ops {
		options = v | options
	}
	return intersection(a0, a1, options)
}

func intersection(a0, a1 Amorph, options int) (Amorph, error) {
	switch ca0 := a0.(type) {
	case nullType:
		return NULL, nil
	case nil:
		switch a1.(type) {
		case nil:
			return nil, nil
		default:
			return NULL, nil
		}
	case string:
		switch ca1 := a1.(type) {
		case string:
			if ca0 == ca1 {
				return ca0, nil
			}
			return NULL, nil
		default:
			return NULL, nil
		}
	case float64:
		switch ca1 := a1.(type) {
		case float64:
			if ca0 == ca1 {
				return ca0, nil
			}
			return NULL, nil
		default:
			return NULL, nil
		}
	case map[string]interface{}:
		switch a1.(type) {
		case map[string]interface{}:
			return mapIntersection(a0, a1, options)
		default:
			return NULL, nil
		}
	case []interface{}:
		switch a1.(type) {
		case []interface{}:
			return sliceIntersection(a0, a1, options)
		default:
			return NULL, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return nil, nil
}
