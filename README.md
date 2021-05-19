# Introduction

Decoding json into an interface{} produces an hierarchical arrangement of four data types: float64, string are 'primative types' and form 'leaf nodes' in the hierarchy, []interface, and map[string]interface{} are compound types and compose layers of the hierarchy around the primative and compound types.

This hierarchical arrangement has no name that I've been able to find, so I label it Amorph which is short for Amorphic. The four types listed above will be referred to as the A4 types.

The json aspect of this README is only an example of where the data might come from, Amorphs can be constructed in other ways as well.

Other types are handled in a limited, and not efficient manner. This of course could be addressed, but would require the reflect package, which this library only uses as a last resort.
# Value Proposition

amorph provides a Diff function that diffs two Amorphs and produces a Patch. A Patch can be applied to an Amorph and a new Amorph reflecting the differences is produced. Patches can be both forward and reverse applied.

The Walker provides a means to iterate over an Amorph and provides functions needed to fully manipulate it.

Amorphs also work with gojq which allows sophisticated queries to be run on an Amorph and producing results that are Amorphs as well.
# Amorph Creation
## Literal Amorph:

	simpleLiteralAmorph := []interface{}{3.14159}
	
## A more sophisticated literal Amorph:

	literalAmorph := []interface{}{
		"123",
		"456",
		[]interface{}{
			"string", 88, "test",
		},
		map[string]interface{}{
			"dooo": "foo",
			"farp": "bar",
			"flam": "baz",
		},
	}

## Create an Amorph from json in a file:
```
	amorphFromFile, err := amorph.NewAmorphFromFile("test.json")
```

## Create an Amorph from a json string:

	amorphExample0, err := amorph.NewAmorphFromString(`{
		"employee":{ 
			"name":"John", 
			"age":30, 
			"city":"New York", 
			"country": "USA",
			"burgers": [ "memphis", "johnny cash", "roy orbison", "burford"]
		}
	}`)
	
## Create an Amorph from json coming in on an io.Reader:

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
	amorphExample1, err := amorph.NewAmorphFromReader(rdr)



## DeepCopy - Create an Amorph from any arbitrary data object:

amorph.DeepCopy encodes the input in json and then decodes it into an interface{}. The output will include only A4 types.

    copyOfInputAmorph := amorph.DeepCopy(inputAmorph)

## Amorph DeepEqual

Compares two Amorphs, returns true or false. While 'walking' the data structure, any data type found one of the A4 types results in a call to reflect.DeepEqual.

    // data0 and data1 are from above examples
    eq := amorph.DeepEqual(inputAmorph0, inputAmorph1)

## Amorph Diff function

As an example, you have two Amorphs that represent 'before' and 'after' conditions.

amorph.Diff()  generates a Patch containing their differences.

    patch, err := amorph.Diff(before, after)
   
Applying the patch in the forward direction to 'before' will produce 'after'.

Applying the patch in the reverse direction to 'after' will produce 'before'.

## Patch ApplyFwd and ApplyRev

The patch methods ApplyFwd and ApplyRev take an Amorph as input and generate a new Amorph as output. The output Amorph will contain the input Amorph with the differences from the Patch applied.

    result, err := amorph.ApplyFwd(patch, before)
    
    // amorph.DeepEqual(after, result) will be true

A patch can also be reverse applied:

    result, err := amorph.ApplyRev(patch, after)
    
    // amorph.DeepEqual(before, result) will be true

# Walker behavior
## Walk Function

amorph.Walk iterates over an Amorph and calls the user supplied function at every node in the Amorph.

The user supplied function receives an iterator of type WalkIterator.

	err = amorph.Walk(adata0, func(iter amorph.WalkIter) error {
		top := iter.Top()
		key := top.Key()
		rval := top.Value()

		switch trval := rval.(type) {
		case nil:
			fmt.Println("key ", key, ", trval = ", trval, " nil ")
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

## WalkIterator

WalkIterator is a stack that describes the path from the top of the Amorph down the node being 'visited'. It contains all of the map keys and slice indexes down to the visited node.

Internally the WalkIterator is implemented as a slice of WalkPos elements (described below).
### Len
Len() returns the length of the WalkIterator

### Top

Top() returns the top WalkPos element of the WalkIterator.
    top := iter.Top()
    
panics if the len of the iterator is zero.
    
#### Pop

Pop() returns a copy of a WalkIterator with the top WalkIterator removed.

    newiter := iter.Pop()

#### Copy
Copy() returns a copy of a WalkIterator.

    newiter := iter.Copy()

### WalkPos

WalkPos the the type of the elements in WalkIter

### Key and Value

If WalkPos refers to a slice, the value returned from Key() can be type asserted to int.

If WalkPos refers to a map, the value returned from Key() can be type asserted to string.

If WalkPos refers to anything else, do not use Key(), it will return nil.

Key() and Value() allow you to access the key and value for that element.

### Dereference

If the WalkPos points to a slice:

    val := wp.Dereference(idx)
    
it returns the value at idx in the slice.

If the WalkPos points to a map:

    val := wp.Dereference(key)
    
it returns the value associate with key in the map.

panics if called on an unsupported type. The type can be determined from the iterator.

### Rereference

If the WalkPos points to a slice:

    wp.Rereference(idx, value)
    
sets the value at idx in the slice.

If the WalkPos points to a map:

    wp.Rereference(key, value)
    
sets the value associate with key in the map.

panics if called on an unsupported type. The type can be determined from the iterator.

### Delete

If the WalkPos points to a map:

    wp.Delete(key)
    
deletes the key from the map.

panics if called on an unsupported type. The type can be determined from the iterator.
