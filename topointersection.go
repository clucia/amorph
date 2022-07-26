package amorph

// Intersection produces an Amorph that contains everything common to the
// two Amorphs
//
// By default, Intersection requires the leaf nodes to have the same value.
//
// The OptIgnoreLeafValues option tells Intersection to ignore the values.
// This can create a situation where there are two values at the same topological location.
//
// These options resolve these conflicts:
// OptSliceAmorph0 use the value from Amorph0
// OptSliceAmorph1 use the value from Amorph1
// OptSliceAlways puts the two values in a slice
// OptSliceNotEqual puts the two values in a slice if they're not equal

func amorphTopoIntersection(a0, a1 Amorph, options int) (Amorph, error) {
	switch {
	case OptTopoIntersectionSliceAlways&options > 0:
		return []interface{}{a0, a1}, nil
	case OptTopoIntersectionSliceResolveAmorph0&options > 0:
		return a0, nil
	case OptTopoIntersectionSliceResolveAmorph1&options > 0:
		return a1, nil
	default:
		return []interface{}{a0, a1}, nil
	}
}

func nilTopoIntersection(a0, a1 Amorph, options int) (Amorph, error) {
	switch {
	case OptTopoIntersectionSliceAlways&options > 0:
		return []interface{}{a0, a1}, nil
	case OptTopoIntersectionSliceResolveAmorph0&options > 0:
		return a0, nil
	case OptTopoIntersectionSliceResolveAmorph1&options > 0:
		return a1, nil
	case OptTopoIntersectionSliceNotEqual&options > 0 && a0 == a1:
		return a0, nil
	default:
		return []interface{}{a0, a1}, nil
	}
}

func stringTopoIntersection(a0, a1 string, options int) (Amorph, error) {
	switch {
	case OptTopoIntersectionSliceAlways&options > 0:
		return []interface{}{a0, a1}, nil
	case OptTopoIntersectionSliceResolveAmorph0&options > 0:
		return a0, nil
	case OptTopoIntersectionSliceResolveAmorph1&options > 0:
		return a1, nil
	case OptTopoIntersectionSliceNotEqual&options > 0 && a0 == a1:
		return a0, nil
	default:
		return []interface{}{a0, a1}, nil
	}
}

func float64TopoIntersection(a0, a1 float64, options int) (Amorph, error) {
	switch {
	case OptTopoIntersectionSliceAlways&options > 0:
		return []interface{}{a0, a1}, nil
	case OptTopoIntersectionSliceResolveAmorph0&options > 0:
		return a0, nil
	case OptTopoIntersectionSliceResolveAmorph1&options > 0:
		return a1, nil
	case OptTopoIntersectionSliceNotEqual&options > 0 && a0 == a1:
		return a0, nil
	default:
		return []interface{}{a0, a1}, nil
	}
}

func mapTopoIntersection(a0, a1 Amorph, options int) (Amorph, error) {
	ar := make(map[string]interface{})
	a0map := a0.(map[string]interface{})
	a1map := a1.(map[string]interface{})

	for a0key, a0elem := range a0map {
		a1elem, a1ok := a1map[a0key]
		if !a1ok {
			continue
		}
		arelem, err := topoIntersection(a0elem, a1elem, options)
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

func sliceTopoIntersection(a0, a1 Amorph, options int) (Amorph, error) {
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
		ar[i], err = topoIntersection(a0Slice[i], a1Slice[i], options)
		if err != nil {
			return nil, err
		}
	}
	return ar, nil
}

func TopoIntersection(a0, a1 Amorph, ops ...int) (ar Amorph, err error) {
	options := 0
	for _, v := range ops {
		options = v | options
	}
	return topoIntersection(a0, a1, options)
}

func topoIntersection(a0, a1 Amorph, options int) (Amorph, error) {
	switch a0.(type) {
	case nullType:
		return NULL, nil
	case nil:
		return nilTopoIntersection(a0, a1, options)
	case string:
		switch a1.(type) {
		case nullType:
			return NULL, nil
		case string:
			return stringTopoIntersection(a0.(string), a1.(string), options)
		default:
			return amorphTopoIntersection(a0, a1, options)
		}
	case float64:
		switch a1.(type) {
		case nullType:
			return NULL, nil
		case float64:
			return float64TopoIntersection(a0.(float64), a1.(float64), options)
		default:
			return amorphTopoIntersection(a0, a1, options)
		}
	case map[string]interface{}:
		switch a1.(type) {
		case nullType:
			return NULL, nil
		case map[string]interface{}:
			return mapTopoIntersection(a0, a1, options)
		}
		return nil, nil
	case []interface{}:
		switch a1.(type) {
		case nullType:
			return NULL, nil
		case []interface{}:
			return sliceTopoIntersection(a0, a1, options)
		}
	default:
		return nil, ErrUnsupportedType
	}
	return nil, nil
}
