package exports

import (
	"TweakItDocs/internal/exports/properties"
	"TweakItDocs/internal/sjsonhelp"
	"github.com/minio/simdjson-go"
	"log"
)

type ExportInfo struct {
	Components       []string //TODO Make this an ObjectProperty array
	OverriddenValues []string //TODO Make this a Property array
	NewValues        []string //TODO Make this a Property array
	ImportedClasses  []string
}

type exportValues struct {
	name       string
	properties []sjsonhelp.JsonMap
}

func ExtractExports(obj *simdjson.Object) []sjsonhelp.JsonMap {
	a := sjsonhelp.ExtractArray(obj, "exports")
	return formatMultipleExports(a)
}

func formatMultipleExports(a *simdjson.Array) []sjsonhelp.JsonMap {
	objects := sjsonhelp.JsonArrayToArrayOfObjects(a)
	r := make([]sjsonhelp.JsonMap, len(objects))
	for i, o := range objects {
		r[i] = formatExport(o)
	}
	return r
}

func formatExport(export *simdjson.Object) sjsonhelp.JsonMap {
	name := sjsonhelp.ExtractString(export, "export", "object_name")
	data := sjsonhelp.ExtractArray(export, "data", "properties")
	dataArray := sjsonhelp.JsonArrayToArrayOfObjects(data)
	props := make([]sjsonhelp.JsonMap, len(dataArray))
	for i, propObject := range dataArray {
		propMap, err := propObject.Map(nil)
		if err != nil {
			log.Fatalf("Could not interpret json object as map: %v", err)
		}
		prop := properties.New(propMap)
		props[i] = properties.ToMap(prop)
	}
	return sjsonhelp.JsonMap{"name": name, "properties": props}
}
