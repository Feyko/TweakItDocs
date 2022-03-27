package data

type Record struct {
	Filename string   `json:"filename"`
	Exports  []Export `json:"exports"`
	Imports  []Import `json:"imports"`
}

type Export struct {
	ObjectName string     `json:"object_name"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name    string `json:"name"`
	Type    string `json:"property_type"`
	Tag     any    `json:"tag"`
	TagData any    `json:"tag_data"`
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
	Properties []Property `json:"properties"`
}

type rawExportExport struct {
	ObjectName string `json:"object_name"`
}

type rawSummary struct {
	Imports []Import
}
