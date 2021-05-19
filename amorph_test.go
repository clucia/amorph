package amorph_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/slllc/amorph"
)

// twowaytest starts with two Amorphs, v0 and v1
// it calculates the diffs
// forward apply diffs to v0, verify against v1
// reverse apply diffs to that, verify against v0
func twowaytest(v0, v1 amorph.Amorph) (err error) {
	diffs := amorph.Diff(v0, v1)
	if diffs == nil {
		return //
	}
	data0patched, err := amorph.ApplyFwd(diffs, v0)
	if err != nil {
		return //
	}
	res := amorph.DeepEqual(data0patched, v1)
	if !res {
		return fmt.Errorf("failed match")
	}
	data0repatched, err := amorph.ApplyRev(diffs, data0patched)
	if err != nil {
		return //
	}
	res = amorph.DeepEqual(v0, data0repatched)
	if !res {
		return fmt.Errorf("failed match")
	}
	return nil
}

// WalkStringer uses Walk to iterate over an Amorph and convert it to a string representation.
func WalkStringer(in interface{}) (s string) {
	amorph.Walk(in, func(iter amorph.WalkIter) error {
		top := iter.Top()
		key := top.Key()
		rval := top.Value()

		// value := top.Dereference(key)

		l := iter.Len()
		s += `                                                       `[:l*2]

		switch trval := rval.(type) {
		case nil:
			// s += fmt.Sprintln("key ", key, ", trval = ", trval, " nil ")
		case []interface{}:
			s += fmt.Sprintln(key, " slice ") // , value)
		case map[string]interface{}:
			s += fmt.Sprintln(key, " map ") // , value)
		case string:
			s += fmt.Sprintln(key, ",", trval, " string ") // , value)
		case float64:
			s += fmt.Sprintln(key, ",", trval, " string ") // , value)
		default:
			s += fmt.Sprintln(rval)
		}
		return nil
	})
	return //

}

// file processing
func Test000(t *testing.T) {
	data, err := amorph.NewAmorphFromFile("test.json")
	if err != nil {
		t.FailNow()
	}
	node0 := data.([]interface{})[0]
	node1 := data.([]interface{})[1]

	err = twowaytest(node0, node1)
	if err != nil {
		t.FailNow()
	}

}

// test maps
func Test001(t *testing.T) {
	data0 := map[string]interface{}{
		"foo": "123",
		"bar": "456",
		"bax": "789",
	}
	data1 := map[string]interface{}{
		"foo": "999",
		"bar": "456",
		"tur": "333",
	}
	err := twowaytest(data0, data1)
	if err != nil {
		t.FailNow()
	}
}

// test slice
func Test002(t *testing.T) {
	data0 := []interface{}{
		"123",
		"456",
		"789",
	}
	data1 := []interface{}{
		"123",
		"xxx",
	}
	err := twowaytest(data0, data1)
	if err != nil {
		t.FailNow()
	}
}

// Slice containing map/slice
func Test003(t *testing.T) {
	data0 := []interface{}{
		"123",
		"456",
		[]interface{}{
			"string", 88, "test",
		},
		map[string]interface{}{
			"dooo": "jdkjdk",
			"farp": "jdkjdk",
			"flam": "jdkjdk",
		},
	}
	data1 := []interface{}{
		"123",
		"xxx",
		[]interface{}{
			"string", 91, 77, "nerf",
		},
		map[string]interface{}{
			"dooo": "jdkjdk",
			"farp": "newval",
			"rdin": "jdkj",
		},
	}
	err := twowaytest(data0, data1)
	if err != nil {
		t.FailNow()
	}
}

// map containing map/slice
func Test004(t *testing.T) {
	data0 := []interface{}{
		"123",
		"456",
		[]interface{}{
			"string", 88, "test",
		},
		map[string]interface{}{
			"dooo": "jdkjdk",
			"farp": "jdkjdk",
			"flam": "jdkjdk",
		},
	}
	data1 := []interface{}{
		"123",
		"xxx",
		[]interface{}{
			"string", 91, 77, "nerf",
		},
		map[string]interface{}{
			"dooo": "jdkjdk",
			"farp": "newval",
			"rdin": "jdkj",
		},
	}
	err := twowaytest(data0, data1)
	if err != nil {
		t.FailNow()
	}
}

func Test005(t *testing.T) {
	data0 := http.Request{}

	str := WalkStringer(data0)

	fmt.Println(str)
}

type MyRequest struct {
	Host   string
	Method string
}

// map containing map/slice
func Test006(t *testing.T) {

	data0 := []interface{}{
		"123",
		"456",
		[]interface{}{
			"string", 88, "test",
		},
		MyRequest{Host: "fred", Method: "GET"},
		map[string]interface{}{
			"dooo": "jdkjdk",
			"farp": "jdkjdk",
			"rdin": "jdkjdk",
		},
	}
	data1 := []interface{}{
		"123",
		"xxx",
		[]interface{}{
			"string", 91, 77, "nerf",
		},
		MyRequest{Host: "barney", Method: "GET"},
		map[string]interface{}{
			"dooo": "jdkjdk",
			"farp": "newval",
			"rdin": "jdkj",
		},
	}
	err := twowaytest(data0, data1)
	if err != nil {
		t.FailNow()
	}
	data0alt := amorph.DeepCopy(data0)
	fmt.Println("data0 = ", data0)
	fmt.Println("data0alt = ", data0alt)
	data1alt := amorph.DeepCopy(data1)
	fmt.Println("data1 = ", data1)
	fmt.Println("data1alt = ", data1alt)
	err = twowaytest(data0alt, data1)
	if err != nil {
		t.FailNow()
	}

}

type test0 struct {
	count int
	val   int
}

// Walk test
func Test007(t *testing.T) {
	rdr := bytes.NewReader(
		[]byte(`{
				"employee":{ 
					"name":"John", 
					"age":30, 
					"city":"New York", 
					"country": "USA",
					"burgers": [ "memphis", "johnny cash", "roy orbison", "burford"]
				},
				"test": ""
			}`,
		),
	)
	adata0, err := amorph.NewAmorphFromReader(rdr)
	if err != nil {
		t.FailNow()
	}
	adata0.(map[string]interface{})["test"] = MyRequest{Host: "fred", Method: "GET"}

	// fmt.Println("walkstringer = ", WalkStringer(adata0.Unref()))
	err = amorph.Walk(adata0, func(iter amorph.WalkIter) error {
		fmt.Println("iter = ", iter)
		top := iter.Top()
		key := top.Key()
		rval := top.Value()

		switch trval := rval.(type) {
		// case nil:
		//	fmt.Println("key ", key, ", trval = ", trval, " nil ")
		case []interface{}:
			fmt.Println(key, " slice ") // , value)
		case map[string]interface{}:
			fmt.Println(key, " map ") // , value)
		case string:
			fmt.Println(key, ",", trval, " string ") // , value)
		case float64:
			fmt.Println(key, ",", trval, "  float64 ") // , value)
		default:
			fmt.Println("other type: ", rval)
		}
		return nil
	})
	if err != nil {
		t.FailNow()
	}
	fmt.Println("walkstringer = ", WalkStringer(adata0))
}

// test using an 'outside' type
func Test008(t *testing.T) {
	rdr := bytes.NewReader(
		[]byte(`{
			"employee":{ 
				"name":"John", 
				"age":30, 
				"city":"New York", 
				"country": "USA",
				"burgers": [ "memphis", "johnny cash", "roy orbison", "burford"]
			}
		}`),
	)
	adata0, err := amorph.NewAmorphFromReader(rdr)
	if err != nil {
		t.FailNow()
	}
	adata0.(map[string]interface{})["booga"] = &test0{1, 2}

	rdr = bytes.NewReader(
		[]byte(`{
			"employee":{
				"name":"bob",
				"age":30,
				"city":"New York",
				"burgers": [ "paris", "johnny cash", "roy orbison"],
				"friends": [ "monica", "chandler", "phoebe", "joey"]
			}
		}`),
	)
	adata1, err := amorph.NewAmorphFromReader(rdr)
	if err != nil {
		t.FailNow()
	}
	adata1.(map[string]interface{})["booga"] = &test0{4, 3}

	diffs := amorph.Diff(adata0, adata1)
	if diffs == nil {
		t.FailNow()
	}
	fmt.Println("diffs = ", amorph.PatchStringer(diffs))

	twowaytest(adata0, adata1)
}

func Test009(t *testing.T) {
	adata0, err := amorph.NewAmorphFromString(
		`{
			"employee":{ 
				"name":"John" 
			}
		}`,
	)
	if err != nil {
		t.FailNow()
	}
	adata1, err := amorph.NewAmorphFromString(
		`{
			"employee":{
				"name":"Johny"
			}
		}`,
	)
	if err != nil {
		t.FailNow()
	}
	diffs := amorph.Diff(adata0, adata1)
	if diffs == nil {
		t.FailNow()
	}
	fmt.Println("diffs = ", amorph.PatchStringer(diffs))
}

func Test010(t *testing.T) {
	adata0, err := amorph.NewAmorphFromString(
		`[ "foo", {"why": "yknot", "where": "wherever" }, "bar", "baz" ]`,
	)
	if err != nil {
		t.FailNow()
	}
	adata1, err := amorph.NewAmorphFromString(
		`[ "foo", {"why": " frayed knot", "where": "wherever"}, "bar", "xxx", "baz" ]`,
	)
	if err != nil {
		t.FailNow()
	}
	diffs := amorph.Diff(adata0, adata1)
	if diffs == nil {
		t.FailNow()
	}
	fmt.Println("diffs = ", amorph.PatchStringer(diffs))
}
