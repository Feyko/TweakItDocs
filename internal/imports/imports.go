package imports

import (
	"TweakItDocs/internal/sjsonhelp"
	"github.com/fatih/structs"
	sjson "github.com/minio/simdjson-go"
	"log"
)

func ExtractImports(obj *sjson.Object) []sjsonhelp.JsonMap {
	jsonArray := sjsonhelp.ExtractArray(obj, "summary", "imports")
	objects := sjsonhelp.JsonArrayToArrayOfObjects(jsonArray)
	values := arrayOfJsonObjectsToImportValues(objects)
	values = fillPackagesOfImportValuesArray(values)
	return mapImportValuesArray(values)
}

type importValues struct {
	ClassName    string `structs:"class_name"`
	ClassPackage string `structs:"class_package"`
	ObjectName   string `structs:"object_name"`
	OuterIndex   int64  `structs:"outer_index"`

	// Not from the pak, filled by some logic
	Index   int    `structs:"index"`
	Package string `structs:"package"`
}

func arrayOfJsonObjectsToImportValues(a []*sjson.Object) []importValues {
	r := make([]importValues, len(a))
	for i, object := range a {
		values := jsonObjectToImportValues(object)
		values.Index = arrayIndexToPakStyleIndex(i)
		r[i] = values
	}
	return r
}

func arrayIndexToPakStyleIndex(i int) int {
	return -i - 1
}

func pakStyleIndexToArrayIndex(i int) int {
	return arrayIndexToPakStyleIndex(i) // Logic is actually the same
}

func jsonObjectToImportValues(object *sjson.Object) importValues {
	objectMap, err := object.Map(nil)
	if err != nil {
		log.Fatalf("Could not turn json object into import values because could not turn a json object into a map: %v", err)
	}
	return importValues{
		ClassName:    objectMap["class_name"].(string),
		ClassPackage: objectMap["class_package"].(string),
		ObjectName:   objectMap["object_name"].(string),
		OuterIndex:   objectMap["outer_index"].(int64),
	}
}

func fillPackagesOfImportValuesArray(imports []importValues) []importValues {
	for i, imp := range imports {
		if imp.OuterIndex == 0 {
			continue
		}
		imports[i] = fillPackageOfImportValues(imp, imports)
	}
	return imports
}

func fillPackageOfImportValues(imp importValues, imports []importValues) importValues {
	outerPackage := findPackageOfImport(imp, imports)
	imp.Package = outerPackage.ObjectName
	return imp
}

func findPackageOfImport(imp importValues, imports []importValues) importValues {
	index := imp.OuterIndex
	next := imports[pakStyleIndexToArrayIndex(int(index))]
	for next.OuterIndex != 0 {
		if next.ClassName == "Package" {
			log.Fatalf("Found a package import that isn't the root: %v", next)
		}
		next = imports[pakStyleIndexToArrayIndex(int(next.OuterIndex))]
	}
	if next.ClassName != "Package" {
		log.Fatalf("Found a root import that isn't a package: %v", next)
	}
	return next
}

func mapImportValuesArray(imports []importValues) []sjsonhelp.JsonMap {
	r := make([]sjsonhelp.JsonMap, len(imports))
	for i, imp := range imports {
		r[i] = structs.Map(imp)
	}
	return r
}
