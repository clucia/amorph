package amorph

const (
	OptionNone                  = 1 << iota
	OptUnionSliceResolveAmorph0 //  use the value from Amorph0 this will trump OptionSliceAmorph1Mask if both are set
	OptUnionSliceResolveAmorph1 // use the value from Amorph1
	OptUnionSliceAlways         // puts the two values in a slice
	OptUnionSliceNotEqual       // puts the two values in a slice if they're not equal

	OptTopoIntersectionSliceResolveAmorph0 //  use the value from Amorph0 this will trump OptionSliceAmorph1Mask if both are set
	OptTopoIntersectionSliceResolveAmorph1 // use the value from Amorph1
	OptTopoIntersectionSliceAlways         // puts the two values in a slice
	OptTopoIntersectionSliceNotEqual       // puts the two values in a slice if they're not equal

	OptDifferenceMustSubtract

	OptTopoDifferenceMustSubtract
)

const (
	OptMustSubtract         = OptDifferenceMustSubtract | OptTopoDifferenceMustSubtract
	OptResolveAmorph0       = OptUnionSliceResolveAmorph0 | OptTopoIntersectionSliceResolveAmorph0
	OptResolveAmorph1       = OptUnionSliceResolveAmorph1 | OptTopoIntersectionSliceResolveAmorph1
	OptResolveSliceAlways   = OptUnionSliceAlways | OptTopoIntersectionSliceAlways
	OptResolveSliceNotEqual = OptUnionSliceNotEqual | OptTopoIntersectionSliceNotEqual
)

type Options int
