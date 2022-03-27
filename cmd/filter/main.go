package main

import (
	"TweakItDocs/internal/data"
	"TweakItDocs/internal/exports"
	"TweakItDocs/internal/imports"
	"TweakItDocs/internal/sjsonhelp"
	"encoding/json"
	"fmt"
	sjson "github.com/minio/simdjson-go"
	"log"
	"os"
	"regexp"
	"strings"
)

var c int

func main() {
	//file, err := os.OpenFile("benchout", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	//if err != nil {
	//	log.Fatalf("Could not open/create the benchmarking file: %v", err)
	//}
	//err = pprof.StartCPUProfile(file)
	//defer pprof.StopCPUProfile()
	//if err != nil {
	//	log.Fatalf("Could not start CPU profiling: %v", err)
	//}

	content, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Could not read the data: %w", err)
	}

	extracted, err := data.Extract(content)
	if err != nil {
		log.Fatalf("Could not decode the data: %w\n", err)
	}
	fmt.Print("Finished extraction")

	//marshalled, err := sjsonhelp.MarshalIndent(r, "", " ")
	marshalled, err := json.Marshal(extracted)
	if err != nil {
		log.Fatalf("Could not marshal json: %v", err)
	}

	err = os.WriteFile("filtered.json", marshalled, 0664)
	if err != nil {
		log.Fatalf("Could not output the result to the file: %w", err)
	}
}

func formatRecord(obj *sjson.Object) map[string]interface{} {
	f := extractFilename(obj)
	if !isValidAssetFilename(f) {
		return nil
	}
	importList := imports.ExtractImports(obj)
	c++
	return map[string]interface{}{
		"i":        c,
		"filename": f,
		"exports":  exports.ExtractExports(obj),
		"imports":  importList,
		//"ClassName": getParentClassFromImports(importList),
	}
}

func extractFilename(obj *sjson.Object) string {
	s := sjsonhelp.ExtractString(obj, "export_record", "file_name")
	return strings.TrimRight(s, "\u0000")
}

func isValidAssetFilename(s string) bool {
	//match, err := regexp.MatchString("^FactoryGame.*(Build|Desc|Recipe|Schematic)_.*$", s)
	match, err := regexp.MatchString("^FactoryGame.*$", s)
	if err != nil {
		return false
	}
	return match
}
