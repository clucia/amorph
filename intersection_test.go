package amorph_test

import (
	"fmt"
	"testing"

	"github.com/clucia/amorph"
	"github.com/stretchr/testify/assert"
)

// OptSlicify(options)
// OptIgnoreLeafValues(options)

func TestIntersectionPrimitiveTypeString(t *testing.T) {
	var data0 interface{}
	var data1 interface{}

	data0 = nil
	data1 = nil

	intr, err := amorph.Intersection(data0, data1)
	assert.Nil(t, intr)
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)

	data1 = "test"

	intr, err = amorph.Intersection(data0, data1)
	assert.Equal(t, intr, amorph.NULL)
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)

	data0 = "test"
	intr, err = amorph.Intersection(data0, data1)
	assert.Equal(t, intr, "test")
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)
}

func TestIntersectionMap(t *testing.T) {
	data0 := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key3": "nothesame",
		"key4": nil,
		"key5": nil,
	}
	data1 := map[string]interface{}{
		"key0": "value0",
		"key2": "value2",
		"key3": "alsonotthesame",
		"key4": "combinednil",
		"key5": nil,
	}
	res0 := map[string]interface{}{
		"key0": "value0",
		"key5": nil,
	}
	intr, err := amorph.Intersection(data0, data1)
	assert.True(t, amorph.DeepEqual(intr, res0))
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)
	/*
		intr, err = amorph.Intersection(data0, data1, amorph.OptIntersectionSliceNotEqual)
		assert.Nil(t, intr)
		assert.Equal(t, err, amorph.ErrIntersectionNoSlicify)
		fmt.Println("intr = ", intr, ", err = ", err)
	*/
	fmt.Println("intr = ", intr, ", err = ", err)
}
func TestIntersectionSlice(t *testing.T) {
	inner0 := []interface{}{
		"another",
		"slice",
	}
	data0 := []interface{}{
		"pos0",
		"123",
		"456",
		"789",
		nil,
		nil,
		inner0,
		amorph.NULL,
	}
	data1 := []interface{}{
		amorph.NULL,
		"123",
		"xxx",
		nil,
		nil,
		"yyy",
		inner0,
		"djjdkjdk",
	}
	res0 := []interface{}{
		amorph.NULL,
		"123",
		amorph.NULL,
		amorph.NULL,
		nil,
		amorph.NULL,
		inner0,
		amorph.NULL,
	}
	intr, err := amorph.Intersection(data0, data1)
	assert.True(t, amorph.DeepEqual(intr, res0))
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)
	/*
		intr, err = amorph.Intersection(data0, data1, amorph.OptIntersectionIgnoreLeafValues)
		assert.Nil(t, intr)
		assert.Equal(t, err, amorph.ErrIntersectionNoSlicify)
		fmt.Println("intr = ", intr, ", err = ", err)

		res0 := []interface{}{
			"123",
			[]interface{}{
				"456", "xxx",
			},
			[]interface{}{
				"789", nil,
			},
			inner0,
		}

		intr, err = amorph.Intersection(data0, data1, amorph.OptIntersectionIgnoreLeafValues, amorph.OptIntersectionSliceNotEqual)
		assert.True(t, amorph.DeepEqual(intr, res0))
		assert.Nil(t, err)
		fmt.Println("intr = ", intr, ", err = ", err)

		res1 := []interface{}{
			"123",
			"456",
			"789",
			inner0,
		}
		intr, err = amorph.Intersection(data0, data1, amorph.OptIntersectionIgnoreLeafValues, amorph.OptIntersectionSliceAmorph0)
		assert.True(t, amorph.DeepEqual(intr, res1))
		assert.Nil(t, err)
		fmt.Println("intr = ", intr, ", err = ", err)

		fmt.Println("intr = ", intr, ", err = ", err)

		// OptSliceAmorph0 use the value from Amorph0
		// OptSliceAmorph1 use the value from Amorph1
		// OptSliceAlways puts the two values in a slice
		// OptSliceNotEqual puts the two values in a slice if they're not equal
	*/
}
