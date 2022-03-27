package data

import "TweakItDocs/internal/data/properties"

type Record struct {
	Filename string   `json:"filename"`
	Exports  []Export `json:"exports"`
	Imports  []Import `json:"imports"`
}

type Export struct {
	ObjectName string                `json:"object_name"`
	Properties []properties.Property `json:"properties"`
}

type Import struct {
	ClassName    string `json:"class_name"`
	ClassPackage string `json:"class_package"`
	ObjectName   string `json:"object_name"`
	OuterIndex   int64  `json:"outer_index"`
}

type rawRecord struct {
	ExportRecord exportRecord `json:"export_record"`
	Exports      []rawExport  `json:"exports"`
	Summary      rawSummary   `json:"summary"`
}

type exportRecord struct {
	FileName string `json:"file_name"`
}

type rawExport struct {
	Data   rawExportData   `json:"data"`
	Export rawExportExport `json:"export"`
}

type rawExportData struct {
	Properties []properties.RawProperty `json:"properties"`
}

type rawExportExport struct {
	ObjectName string `json:"object_name"`
}

type rawSummary struct {
	Imports []Import
}
