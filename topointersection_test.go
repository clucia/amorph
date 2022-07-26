package amorph_test

import (
	"fmt"
	"testing"

	"github.com/clucia/amorph"
	"github.com/stretchr/testify/assert"
)

// OptSlicify(options)
// OptIgnoreLeafValues(options)

func TestTopoIntersectionPrimitiveTypeString(t *testing.T) {
	var data0 interface{}
	var data1 interface{}

	data0 = nil
	data1 = nil

	intr, err := amorph.TopoIntersection(data0, data1)
	assert.Equal(t, intr, []interface{}{nil, nil})
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)

	intr, err = amorph.TopoIntersection(data0, data1, amorph.OptTopoIntersectionSliceNotEqual)
	assert.Equal(t, intr, nil)
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)

	data0 = "test"
	intr, err = amorph.TopoIntersection(data0, data1)
	assert.True(t, amorph.DeepEqual(intr, []interface{}{"test", nil}))
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)
}

func TestTopoIntersectionMap(t *testing.T) {
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
		"key0": []interface{}{"value0", "value0"},
		"key3": []interface{}{"nothesame", "alsonotthesame"},
		"key4": []interface{}{nil, "combinednil"},
		"key5": []interface{}{nil, nil},
	}
	intr, err := amorph.TopoIntersection(data0, data1)
	assert.True(t, amorph.DeepEqual(intr, res0))
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)
	fmt.Println("res0 = ", intr)
	fmt.Println("comparisin = ", amorph.DeepEqual(intr, res0))

	res1 := map[string]interface{}{
		"key0": "value0",
		"key3": []interface{}{"nothesame", "alsonotthesame"},
		"key4": []interface{}{nil, "combinednil"},
		"key5": nil,
	}
	intr, err = amorph.TopoIntersection(data0, data1, amorph.OptTopoIntersectionSliceNotEqual)
	assert.True(t, amorph.DeepEqual(intr, res1))
	assert.Nil(t, err)

	fmt.Println("intr = ", intr, ", err = ", err)
}
func TestTopoIntersectionSlice(t *testing.T) {
	inner0 := []interface{}{
		"another",
		"slice",
	}
	data0 := []interface{}{
		amorph.NULL,
		amorph.NULL,
		"pos2",
		"123",
		"456",
		"789",
		nil,
		nil,
		inner0,
	}
	data1 := []interface{}{
		amorph.NULL,
		"pos1",
		amorph.NULL,
		"123",
		"xxx",
		nil,
		nil,
		"yyy",
		inner0,
	}
	res0 := []interface{}{
		amorph.NULL,
		amorph.NULL,
		amorph.NULL,
		[]interface{}{"123", "123"},
		[]interface{}{"456", "xxx"},
		[]interface{}{"789", nil},
		[]interface{}{nil, nil},
		[]interface{}{nil, "yyy"},
		[]interface{}{
			[]interface{}{"another", "another"},
			[]interface{}{"slice", "slice"},
		},
	}
	intr, err := amorph.TopoIntersection(data0, data1)
	assert.True(t, amorph.DeepEqual(intr, res0))
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)
	fmt.Println("res0 = ", res0)

	res1 := []interface{}{
		amorph.NULL,
		amorph.NULL,
		amorph.NULL,
		"123",
		[]interface{}{"456", "xxx"},
		[]interface{}{"789", nil},
		nil,
		[]interface{}{nil, "yyy"},
		[]interface{}{"another", "slice"},
	}
	intr, err = amorph.TopoIntersection(data0, data1, amorph.OptTopoIntersectionSliceNotEqual)
	assert.True(t, amorph.DeepEqual(intr, res1))
	assert.Nil(t, err)
	fmt.Println("intr = ", intr, ", err = ", err)
	fmt.Println("res1 = ", res1)
}
