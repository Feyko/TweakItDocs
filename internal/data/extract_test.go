package data

import (
	"log"
	"os"
	"testing"
)

var result []Package

func BenchmarkExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		content, err := os.ReadFile("data.json")
		if err != nil {
			log.Fatalf("Could not read the data: %v", err)
		}

		extracted, err := ExtractAll(content)
		if err != nil {
			log.Fatalf("Could not decode the data: %v\n", err)
		}
		result = extracted
	}
}

func BenchmarkExtractParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, err := ExtractSplit("data")
		if err != nil {
			log.Fatal(err)
		}
		result = r
	}
}
