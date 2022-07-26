package amorph

import "fmt"

var ErrMustSubtract = fmt.Errorf("can't subtract")
var ErrUnionNoSlicify = fmt.Errorf("Cannot perform union without slicify")
var ErrIntersectionNoSlicify = fmt.Errorf("Cannot perform intersection without slicify")
var ErrUnsupportedType = fmt.Errorf("unsupported type")
