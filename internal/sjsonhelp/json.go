package sjsonhelp

import (
	"github.com/minio/simdjson-go"
	"log"
)

type JsonMap map[string]interface{}

// TODO: Fix this. It may error on empty json arrays. I believe peeking first skips the first object, I don't know why.
func JsonArrayToArrayOfObjects(a *simdjson.Array) []*simdjson.Object {
	var r []*simdjson.Object
	a.ForEach(func(i simdjson.Iter) {
		obj, err := i.Object(nil)
		if err != nil {
			log.Fatalf("Could not transform supposed object to an object: %v", err)
		}
		r = append(r, obj)
		typ := i.PeekNext()
		if typ == simdjson.TypeNone {
			return
		}
	})
	return r
}

func ExtractArray(o *simdjson.Object, path ...string) *simdjson.Array {
	elem, err := o.FindPath(nil, path...)
	if err != nil {
		log.Fatalf("Could not extract array from object at path %v", path)
	}
	v, err := elem.Iter.Array(nil)
	if err != nil {
		log.Fatalf("Could not interpret supposed array as an array: %v", err)
	}
	return v
}

func ExtractObject(o *simdjson.Object, path ...string) *simdjson.Object {
	elem, err := o.FindPath(nil, path...)
	if err != nil {
		log.Fatalf("Could not extract object from object at path %v", path)
	}
	v, err := elem.Iter.Object(nil)
	if err != nil {
		log.Fatalf("Could not interpret supposed object as an object: %v", err)
	}
	return v
}

func ExtractString(o *simdjson.Object, path ...string) string {
	elem, err := o.FindPath(nil, path...)
	if err != nil {
		log.Fatalf("Could not extract string from object at path %v", path)
	}
	v, err := elem.Iter.StringCvt()
	if err != nil {
		log.Fatalf("Could not interpret supposed string as a string: %v", err)
	}
	return v
}
