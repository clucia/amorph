package amorph_test

import (
	"fmt"
	"testing"

	"github.com/clucia/amorph"
	"github.com/stretchr/testify/assert"
)

func TestUnion000(t *testing.T) {
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
		"key1": "value1",
		"key2": "value2",
		"key3": []interface{}{"nothesame", "alsonotthesame"},
		"key4": []interface{}{nil, "combinednil"},
		"key5": []interface{}{nil, nil},
	}
	u, err := amorph.Union(data0, data1)
	assert.Nil(t, err)
	assert.True(t, amorph.DeepEqual(u, res0))
	fmt.Println("u = ", u, ", err = ", err)
	fmt.Println("r = ", res0)

	res1 := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
		"key3": []interface{}{"nothesame", "alsonotthesame"},
		"key4": []interface{}{nil, "combinednil"},
		"key5": nil,
	}
	u, err = amorph.Union(data0, data1, amorph.OptUnionSliceNotEqual)
	assert.Nil(t, err)
	assert.True(t, amorph.DeepEqual(u, res1))
	fmt.Println("u = ", u, ", err = ", err)
	fmt.Println("r = ", res1)

	res2 := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
		"key3": "nothesame",
		"key4": nil,
		"key5": nil,
	}
	u, err = amorph.Union(data0, data1, amorph.OptUnionSliceResolveAmorph0)
	assert.Nil(t, err)
	assert.True(t, amorph.DeepEqual(u, res2))
	fmt.Println("u = ", u, ", err = ", err)
	fmt.Println("r = ", res2)

	res3 := map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
		"key3": "alsonotthesame",
		"key4": "combinednil",
		"key5": nil,
	}
	u, err = amorph.Union(data0, data1, amorph.OptUnionSliceResolveAmorph1)
	assert.Nil(t, err)
	assert.True(t, amorph.DeepEqual(u, res3))
	fmt.Println("u = ", u, ", err = ", err)
	fmt.Println("r = ", res3)

	res4 := map[string]interface{}{
		"key0": []interface{}{"value0", "value0"},
		"key1": "value1",
		"key2": "value2",
		"key3": []interface{}{"nothesame", "alsonotthesame"},
		"key4": []interface{}{nil, "combinednil"},
		"key5": []interface{}{nil, nil},
	}
	u, err = amorph.Union(data0, data1, amorph.OptUnionSliceAlways)
	assert.Nil(t, err)
	assert.True(t, amorph.DeepEqual(u, res4))
	fmt.Println("u = ", u, ", err = ", err)
	fmt.Println("r = ", res4)
}

func TestUnion001(t *testing.T) {
	data0 := []interface{}{
		"value0",
		"value1",
		nil,
		"nothesame",
		amorph.NULL,
		amorph.NULL,
	}
	data1 := []interface{}{
		"value0",
		nil,
		"value2",
		"alsonotthesame",
		"Value3",
		amorph.NULL,
	}
	_ = data0
	_ = data1
	u, err := amorph.Union(data0, data1)
	fmt.Println("u = ", u, ", err = ", err)

	u, err = amorph.Union(data0, data1, amorph.OptUnionSliceNotEqual)
	fmt.Println("u = ", u, ", err = ", err)
}
