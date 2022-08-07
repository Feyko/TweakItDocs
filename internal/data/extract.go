package data

import (
	"TweakItDocs/internal/data/properties"
	"encoding/json"
	"fmt"
	"regexp"
	"runtime/debug"
)

func Extract(data []byte) ([]Package, error) {
	debug.SetMaxStack(1 << 38) // No, it's not an infinite loop.. I'm just using a lot of stack. Sorry Go

	raw, err := extractRaw(data)
	if err != nil {
		return nil, err
	}
	extracted := rawRecordsToRecordSlice(raw)
	extracted = filter(extracted)
	resolve(extracted)
	return extracted, nil
}

func resolve(packages []Package) {
	for _, p := range packages {
		p.Resolve()
	}
}

func filter(data []Package) []Package {
	r := make([]Package, 0, len(data))
	for _, e := range data {
		if isValidAssetFilename(e.Filename) {
			r = append(r, e)
		}
	}
	return r
}

func isValidAssetFilename(s string) bool {
	//match, err := regexp.MatchString("^FactoryGame.*(Build|Desc|Recipe|Schematic)_.*$", s)
	match, err := regexp.MatchString("^FactoryGame.*$", s)
	if err != nil {
		return false
	}
	return match
}

func extractRaw(data []byte) ([]rawRecord, error) {
	var out []rawRecord
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal the json: %w", err)
	}
	return out, nil
}

func rawRecordsToRecordSlice(raw []rawRecord) []Package {
	return mapSlice(raw, rawRecordToRecord)
}

func mapSlice[T any, R any](s []T, f func(T) R) []R {
	out := make([]R, len(s))
	for i, elem := range s {
		out[i] = f(elem)
	}
	return out
}

func rawRecordToRecord(raw rawRecord) Package {
	return Package{
		Filename: raw.ExportRecord.FileName,
		Exports:  rawExportsToExportSlice(raw.Exports),
		Imports:  raw.Summary.Imports,
	}
}

func rawExportsToExportSlice(raw []rawExport) []Export {
	return mapSlice(raw, rawExportToExport)
}

func rawExportToExport(raw rawExport) Export {
	return Export{
		ClassIndex:    raw.Export.ClassIndex.Convert(),
		SuperIndex:    raw.Export.SuperIndex.Convert(),
		TemplateIndex: raw.Export.TemplateIndex.Convert(),
		OuterIndex:    raw.Export.OuterIndex.Convert(),
		ObjectName:    raw.Export.ObjectName,
		Properties:    rawPropertiesToPropertySlice(raw.Data.Properties),
	}
}

func rawPropertiesToPropertySlice(raw []properties.RawProperty) []properties.Property {
	return mapSlice(raw, properties.DataToProperty)
}
