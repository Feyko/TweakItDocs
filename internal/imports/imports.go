package imports

import (
	"TweakItDocs/internal/cdo"
	"TweakItDocs/internal/sjsonhelp"
	"fmt"
	sjson "github.com/minio/simdjson-go"
	"log"
)

func ExtractImports(obj *sjson.Object) []importMap {
	return formatMultipleImports(sjsonhelp.ExtractArray(obj, "summary", "imports"))
}

func formatMultipleImports(a *sjson.Array) []importMap {
	objects := sjsonhelp.JsonArrayToArrayOfObjects(a)
	r := make([]importMap, len(objects))
	for i, o := range objects {
		r[i] = formatImport(o)
	}
	return r
}

func formatImport(imp *sjson.Object) importMap {
	mdata, err := imp.Map(nil)
	if err != nil {
		log.Fatalf("Could not parse imp data as map: %v", err)
	}
	return mdata
}

type importValues struct {
	ClassName  string
	ObjectName string
	OuterIndex int64
}

type importMap map[string]interface{}

func getParentClassFromImports(imports []importMap) string {
	importV := importsToImportValuesSlice(imports)
	cdos := getCDOsFromImports(importV)
	lowest := getLowestIndexes(cdos)
	lowest = mapImports(lowest, cdoImportToClassImport)
	nocdos := excludeCDOsFromImports(importV)
	if len(lowest) > 1 {
		lowest = filterImports(lowest, func(imp importValues) bool {
			return !inImports(nocdos, imp)
		})
	}
	if len(lowest) > 1 {
		fmt.Println("MORETHAN1")
		return "MORETHAN1: " + fmt.Sprint(lowest)
	}
	if len(lowest) < 1 {
		fmt.Println("AYO bro")
		return "NOFOUND"
	}
	return cdo.ToClassName(lowest[0].ObjectName)
}

func getLowestIndexes(imports []importValues) []importValues {
	r := make([]importValues, 0, len(imports))
	if len(imports) == 0 {
		return r
	}
	lowest := imports[0].OuterIndex
	for _, e := range imports {
		if e.OuterIndex < lowest {
			lowest = e.OuterIndex
			r = make([]importValues, 0, len(imports))
		}
		if e.OuterIndex == lowest {
			r = append(r, e)
		}
	}
	return r
}

func getCDOsFromImports(imports []importValues) []importValues {
	return filterImports(imports, func(imp importValues) bool {
		return cdo.Is(imp.ObjectName)
	})
}

func excludeCDOsFromImports(imports []importValues) []importValues {
	return filterImports(imports, func(imp importValues) bool {
		return !cdo.Is(imp.ObjectName)
	})
}

func importsToImportValuesSlice(imports []importMap) []importValues {
	r := make([]importValues, len(imports))
	for i, e := range imports {
		r[i] = importValuesFromImport(e)
	}
	return r
}

func importValuesFromImport(imp importMap) importValues {
	return importValues{
		ClassName:  imp["class_name"].(string),
		ObjectName: imp["object_name"].(string),
		OuterIndex: imp["outer_index"].(int64),
	}
}

func filterImports(imports []importValues, f func(imp importValues) bool) []importValues {
	r := make([]importValues, 0, len(imports))

	for _, e := range imports {
		if f(e) {
			r = append(r, e)
		}
	}

	return r
}

func inImports(imports []importValues, imp importValues) bool {
	for _, e := range imports {
		if imp.ObjectName == e.ClassName {
			//fmt.Printf("FOUND EQUIVALENCE:\n%v\n%v\n", e, imp)
			return true
		}
	}

	return false
}

func mapImports(imports []importValues, f func(imp importValues) importValues) []importValues {
	r := make([]importValues, len(imports))

	for i, e := range imports {
		r[i] = f(e)
	}

	return r
}

func cdoImportToClassImport(imp importValues) importValues {
	return importValues{
		ClassName:  "Class",
		ObjectName: cdo.ToClassName(imp.ObjectName),
		OuterIndex: imp.OuterIndex,
	}
}
