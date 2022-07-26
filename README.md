# Introduction

Decoding json into an interface{} produces an hierarchical arrangement of four data types: 
+ map[string]interface{}
+ []interface{}
+ float
+ string

I could not find a name for a structure composed of these four types, so I called it an `Amorph`.

This library: [github.com/clucia/amorph](github.com/clucia/amorph) provides a number of useful tools for manipulating `Amorphs`. 

In particular:
#### Diff/Patch Operations:
+ Diff - generate a representation of the differences between two Amorphs
+ PatchFwd - Apply a set of differences to an `Amorph`
+ PatchRev - Reverse apply a set of differences to an `Amorph`

#### Set Operations:
+ Union
+ Intersection
+ Difference (subtraction)
+ Topological Intersection (leaf values are ignored)
+ Topological Difference (leaf values are ignored)

#### Utility Operations:
+ DeepCopy - Duplicate an Amorph
+ DeepEqual - Compare two Amorphs for equality

#### Walker Operations
Used internally in the implementation of Diff and Patch, walking an Amorph is also availble
to users of this package.

#### Additional background
+ JSON
+ gojq - A query language for Amorphs
+ Amorph Creation - Examples
+ Walker Behavior - Used internally to Diff/Patch
---
---
---
# Amorph Diff and Patch operations
## `amorph.Diff()` generates a Patch containing the differences of two `Amorphs`.

For example, you have two Amorphs that represent 'before' and 'after' conditions.

    patch, err := amorph.Diff(before, after)
   
Applying the patch in the forward direction to 'before' will produce 'after'.

Applying the patch in the reverse direction to 'after' will produce 'before'.

## PatchFwd and PatchRev

Formerly known as ApplyFwd and ApplyRev which are still included for compatibility, but, 
are now deprecated.

The patch methods PatchFwd and PatchRev take an Amorph as input and generate a new Amorph as output. The output Amorph will contain the input Amorph with the differences from the Patch applied.

    result, err := amorph.PatchFwd(patch, before)
    
    // amorph.DeepEqual(after, result) will be true

A patch can also be reverse applied:

    result, err := amorph.PatchRev(patch, after)
    
    // amorph.DeepEqual(before, result) will be true
---
---
---


# Set Operations
+ Union - union of two Amorphs
+ Intersection - Intersection of two Amorphs
+ TopoIntersection - Topological Intersection of two Amorphs
+ Difference - Difference of two Amorphs
+ TopoDifference - Topological Difference of two Amorphs
-----
## Union

	Union(a0, a1 Amorph, ops ...int) (ar Amorph, err error)

Union combines two amorphs.

If there are two leaf nodes in the same topological position, this is a conflict. There are a number of options for dealing with conflicts.
### Conflict Resolution Options (Described below) are supported

-----
## Intersection

Intersection creates an amorph containing all of the values that Amorph0 and Amorph1 have in common.

	Intersection(a0, a1 Amorph, ops ...int) (ar Amorph, err error)

### Intersection has no options

## TopoIntersection
	TopoIntersection(a0, a1 Amorph, ops ...int) (ar Amorph, err error)

If there are two leaf nodes in the same topological position, this is a conflict. There are a number of options for dealing with conflicts.

### Conflict Resolution Options (Described below) are supported

-----
## Difference

Difference creates an Amorph with the values in subtrahend removed from the minuend.

	Difference(minuend, subtrahend Amorph, ops ...int) (Amorph, error)

### OptDifferenceMustSubtract
This option tells Difference to return an error if there is something in the subtraend not present in the minuend.

-----
## TopoDifference

TopoDifference creates an amorph containing all of the values in the minuend that do not have a value in the topologically same position in the subtrahend.

	TopoDifference(min, a1 Amorph, ops ...int) (Amorph, error)

### OptDifferenceMustSubtract
This option tells TopoDifference to return an error if there is something in the subtraend not present in the minuend.

----
## Conflict Resolution Options (for Union and TopoIntersection)
### OptUnionSliceResolveAmorph0
When a conflict occurs, the value in Amorph0 is in the result.

### OptUnionSliceResolveAmorph1
When a conflict occurs, the value in Amorph1 is in the result.

### OptUnionSliceAlways (also the default behavior)
When a conflict occurs, a slice containing Amorph0 and Amorph1 is added to the result.

### OptUnionSliceNotEqual
When a conflict occurs, if the values are equal, one of them is added to the result. If they are not equal, a slice containing Amorph0 and Amorph1 is added to the result.
# Utility Operations
## DeepCopy(a0, a1) - Create an Amorph from any arbitrary data object:

amorph.DeepCopy encodes the input in json and then decodes it into an interface{}. The output will include only A4 types.

    copyOfInputAmorph := amorph.DeepCopy(inputAmorph)

## Amorph DeepEqual(a0, a1 amorph.Amorph) bool

Compares two Amorphs, returns true or false. While 'walking' the data structure, any data type found one of the A4 types results in a call to reflect.DeepEqual.

    // data0 and data1 are from above examples
    eq := amorph.DeepEqual(inputAmorph0, inputAmorph1)
# Additional Background

# JSON

The json aspect of this README is only an example of where the data might come from, Amorphs can be constructed in other ways as well.

Other types are handled in a limited, and not efficient manner. This of course could be addressed, but would require the reflect package, which this library only uses as a last resort.

# gojq

`gojq` is a library that implements `jq` in a Go program and operates on Amorphs.

# Amorph Creation
## Literal Amorph:

    var simpleLiteralAmorph amorph.Amorph
	simpleLiteralAmorph = []interface{}{3.14159}
	
## A more sophisticated literal Amorph:

    var literalAmorph amorph.Amorph
	literalAmorph = []interface{}{
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


# Walker behavior

The Walker is used to walk an amorph. It is used in the implementation of Diff and Patch 
functions.
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
