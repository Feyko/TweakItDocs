package main

import (
	"TweakItDocs/internal/data"
	"TweakItDocs/internal/imports"
	"TweakItDocs/internal/translate"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
	"time"
)

var c int

func main() {
	start := time.Now()
	_ = translate.GetProto()
	fmt.Println(time.Since(start))
	//file, err := os.OpenFile("benchout", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	//if err != nil {
	//	log.Fatalf("Could not open/create the benchmarking file: %v", err)
	//}
	//err = pprof.StartCPUProfile(file)
	//defer pprof.StopCPUProfile()
	//if err != nil {
	//	log.Fatalf("Could not start CPU profiling: %v", err)
	//}
	start = time.Now()
	content, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Could not read the data: %w", err)
	}

	extracted, err := data.Extract(content)
	if err != nil {
		log.Fatalf("Could not decode the data: %w\n", err)
	}
	fmt.Println("Finished extraction")
	fmt.Println(time.Since(start))
	filtered := imports.Filter(extracted)

	//marshalled, err := sjsonhelp.MarshalIndent(r, "", " ")
	marshalled, err := json.Marshal(filtered)
	if err != nil {
		log.Fatalf("Could not marshal json: %v", err)
	}

	err = os.WriteFile("filtered.json", marshalled, 0664)
	if err != nil {
		log.Fatalf("Could not output the result to the file: %w", err)
	}
}

func clean(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "could not read the file")
	}
	b = []byte(strings.ReplaceAll(string(b), "\\u0000", ""))
	err = os.WriteFile(filename, b, 0644)
	return errors.Wrap(err, "could not write to the file")
}
