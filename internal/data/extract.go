package data

import (
	"encoding/json"
	"fmt"
)

func Extract(data []byte) ([]Record, error) {
	raw, err := extractRaw(data)
	if err != nil {
		return nil, err
	}
	return rawRecordsToRecordSlice(raw), nil

}

func extractRaw(data []byte) ([]rawRecord, error) {
	var out []rawRecord
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal the json: %w", err)
	}
	return out, nil
}

func rawRecordsToRecordSlice(raw []rawRecord) []Record {
	return mapSlice(raw, rawRecordToRecord)
}

func mapSlice[T any, R any](s []T, f func(T) R) []R {
	out := make([]R, len(s))
	for i, elem := range s {
		out[i] = f(elem)
	}
	return out
}

func rawRecordToRecord(raw rawRecord) Record {
	return Record{
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
		ObjectName: raw.Export.ObjectName,
		Properties: raw.Data.Properties,
	}
}
