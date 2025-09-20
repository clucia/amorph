package amorph_test

import (
	"fmt"
	"testing"

	"github.com/clucia/amorph"
	"github.com/stretchr/testify/assert"
)

func TestDifferencePrimitiveTypeNil(t *testing.T) {
	var data0 interface{}
	var data1 interface{}

	data0 = nil
	data1 = nil

	diff, err := amorph.Difference(data0, data1)
	assert.Equal(t, diff, amorph.NULL)
	assert.Equal(t, err, nil)
}

// OptTrimMaps
// OptIgnoreLeafValues
// OptMustSubtract

func TestDifferencePrimitiveTypeString(t *testing.T) {
	var data0 interface{}
	var data1 interface{}
	data0 = "test"
	data1 = "test"

	diff, err := amorph.Difference(data0, data1, amorph.OptDifferenceMustSubtract)
	assert.Equal(t, diff, amorph.NULL)
	assert.Nil(t, err)
	fmt.Println("diff = ", diff, ", err = ", err)

	data0 = "test0"
	data1 = "test1"

	diff, err = amorph.Difference(data0, data1)
	assert.Equal(t, diff, "test0")
	assert.Nil(t, err)
	fmt.Println("diff = ", diff, ", err = ", err)

	diff, err = amorph.Difference(data0, data1, amorph.OptDifferenceMustSubtract)
	assert.Equal(t, diff, amorph.NULL)
	assert.Equal(t, err, amorph.ErrMustSubtract)
	fmt.Println("diff = ", diff, ", err = ", err)
}

func TestDifferencePrimitiveTypeFloat64(t *testing.T) {
	var data0 interface{}
	var data1 interface{}
	data0 = 1.6
	data1 = 1.6

	diff, err := amorph.Difference(data0, data1, amorph.OptDifferenceMustSubtract)
	assert.Equal(t, diff, amorph.NULL)
	assert.Nil(t, err)
	fmt.Println("diff = ", diff, ", err = ", err)

	data1 = 3.2

	diff, err = amorph.Difference(data0, data1)
	assert.Equal(t, diff, 1.6)
	assert.Nil(t, err)
	fmt.Println("diff = ", diff, ", err = ", err)

	diff, err = amorph.Difference(data0, data1, amorph.OptDifferenceMustSubtract)
	assert.Equal(t, diff, amorph.NULL)
	assert.Equal(t, err, amorph.ErrMustSubtract)
	fmt.Println("diff = ", diff, ", err = ", err)

	fmt.Println("diff = ", diff, ", err = ", err)
}

func TestDifferenceSlice(t *testing.T) {
	data0 := []interface{}{
		"pos0",
		amorph.NULL,
		amorph.NULL,
		"value0",
		"value1",
		nil,
		"notthesame",
	}
	data1 := []interface{}{
		amorph.NULL,
		"pos1",
		amorph.NULL,
		"value0",
		nil,
		"value2",
		"alsonotthesame",
	}
	result1 := []interface{}{
		"pos0",
		amorph.NULL,
		amorph.NULL,
		amorph.NULL,
		"value1",
		nil,
		"notthesame",
	}
	diff, err := amorph.Difference(data0, data1)
	fmt.Println("diff = ", diff, ", err = ", err)

	assert.True(t, amorph.DeepEqual(diff, result1))
	assert.Nil(t, err)

	diff, err = amorph.Difference(data0, data1, amorph.OptDifferenceMustSubtract)
	assert.Nil(t, diff)
	assert.Equal(t, err, amorph.ErrMustSubtract)
	fmt.Println("diff = ", diff, ", err = ", err)

	result2 := []interface{}{
		"pos0",
		amorph.NULL,
		amorph.NULL,
		amorph.NULL,
		"value1",
		nil,
		"notthesame",
	}
	diff, err = amorph.Difference(data0, data1)
	assert.True(t, amorph.DeepEqual(diff, result2))
	assert.Nil(t, err)
	fmt.Println("diff = ", diff, ", err = ", err)

	fmt.Println("diff = ", diff, ", err = ", err)

}

func TestDifference001(t *testing.T) {
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
		"key4": "notnil",
		"key5": nil,
	}
	diff, err := amorph.Difference(data0, data1, amorph.OptDifferenceMustSubtract)
	fmt.Println("diff = ", diff, ", err = ", err)

	fmt.Println("diff = ", diff, ", err = ", err)

}

func TestNullHandling(t *testing.T) {
	// Test cases for DeepEqual with NULL
	tests := []struct {
		name     string
		a        amorph.Amorph
		b        amorph.Amorph
		expected bool
	}{
		{
			name:     "NULL equals NULL",
			a:        amorph.NULL,
			b:        amorph.NULL,
			expected: true,
		},
		{
			name:     "NULL not equals nil",
			a:        amorph.NULL,
			b:        nil,
			expected: false,
		},
		{
			name:     "NULL not equals string",
			a:        amorph.NULL,
			b:        "test",
			expected: false,
		},
		{
			name:     "NULL not equals empty slice",
			a:        amorph.NULL,
			b:        []interface{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := amorph.DeepEqual(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDiffWithNull(t *testing.T) {
	// Test cases for Diff with NULL
	tests := []struct {
		name string
		a    amorph.Amorph
		b    amorph.Amorph
	}{
		{
			name: "Diff NULL to value",
			a:    amorph.NULL,
			b:    "test",
		},
		{
			name: "Diff value to NULL",
			a:    "test",
			b:    amorph.NULL,
		},
		{
			name: "Diff NULL to NULL",
			a:    amorph.NULL,
			b:    amorph.NULL,
		},
		{
			name: "Diff slice with NULL elements",
			a:    []interface{}{amorph.NULL, "test", amorph.NULL},
			b:    []interface{}{"test", amorph.NULL, "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch := amorph.Diff(tt.a, tt.b)
			restored0, err := amorph.ApplyFwd(patch, tt.a)
			assert.NoError(t, err)
			assert.True(t, amorph.DeepEqual(restored0, tt.b))

			restored1, err := amorph.ApplyRev(patch, restored0)
			assert.NoError(t, err)
			assert.True(t, amorph.DeepEqual(restored1, tt.a))
		})
	}
}
