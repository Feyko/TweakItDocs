package main

import (
	"TweakItDocs/internal/data"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	//_ "net/http/pprof"
	"os"
)

func main() {
	//go func() {
	//	_ = http.ListenAndServe("0.0.0.0:8081", nil)
	//}()
	//file, err := os.OpenFile("membenchout", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	//if err != nil {
	//	log.Fatalf("Could not open/create the benchmarking file: %v", err)
	//}
	//err = pprof.StartCPUProfile(file)
	//defer pprof.StopCPUProfile()
	//if err != nil {
	//	log.Fatalf("Could not start CPU profiling: %v", err)
	//}
	//cleanData("data.json")
	//return

	//start := time.Now()
	//content, err := os.ReadFile("data.json")
	//if err != nil {
	//	log.Fatalf("Could not read the data: %w", err)
	//}
	//extracted, err := data.ExtractAll(content)
	//fmt.Println(time.Since(start))
	//start = time.Now()
	_ = data.FilterForFilename(nil, "") // include symbol
	extracted, err := data.ExtractAllFromDir("testdata")
	//fmt.Println(time.Since(start))
	if err != nil {
		log.Fatalf("Could not decode the data: %w\n", err)
	}
	fmt.Println("Finished extraction")

	translated, err := data.TranslateAll(extracted)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not translate"))
	}
	fmt.Println("Finished translation")

	//marshalled, err := sjsonhelp.MarshalIndent(r, "", " ")
	marshalled, err := json.Marshal(extracted)
	if err != nil {
		log.Fatalf("Could not marshal json: %v", err)
	}

	err = os.WriteFile("filtered.json", marshalled, 0664)
	if err != nil {
		log.Fatalf("Could not output the result to the file: %w", err)
	}

	marshalled, err = json.Marshal(translated)
	if err != nil {
		log.Fatalf("Could not marshal json: %v", err)
	}

	err = os.WriteFile("translated.json", marshalled, 0664)
	if err != nil {
		log.Fatalf("Could not output the result to the file: %w", err)
	}
}
