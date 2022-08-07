package data

import (
	"TweakItDocs/internal/data/properties"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
)

func newStack() IndexStack {
	return IndexStack{
		v: make([]int, 0),
	}
}

type IndexStack struct {
	v []int
}

func (s *IndexStack) Add(i int) {
	s.v = append(s.v, i)
}

func (s *IndexStack) Remove(number int) {
	s.v = s.v[:len(s.v)-number]
}

func (s *IndexStack) Exists(i int) bool {
	for _, v := range s.v {
		if v == i {
			return true
		}
	}
	return false
}

type Package struct {
	Filename string   `json:"filename"`
	Exports  []Export `json:"exports"`
	Imports  []Import `json:"imports"`
}

func (p *Package) Resolve() {
	for i := range p.Exports {
		(&p.Exports[i]).Resolve(p)
	}
	for i := range p.Imports {
		(&p.Imports[i]).Resolve(p)
	}
}

type Export struct {
	ClassIndex    Index                 `json:"class_index"`
	SuperIndex    Index                 `json:"super_index"`
	TemplateIndex Index                 `json:"template_index"`
	OuterIndex    Index                 `json:"outer_index"`
	ObjectName    string                `json:"object_name"`
	Properties    []properties.Property `json:"properties"`
}

func (e *Export) Resolve(p *Package) {
	e.ClassIndex.Resolve(p)
	e.SuperIndex.Resolve(p)
	e.TemplateIndex.Resolve(p)
	e.OuterIndex.Resolve(p)
	for i := range e.Properties {
		property := &e.Properties[i]
		if property.PropertyType == "Object" {
			partial, ok := property.Value.(properties.Index)
			if !ok {
				continue // If it is already a full index then it was resolved
			}
			index := indexFromPartialIndex(partial)
			property.Value = index
			index.Resolve(p)
		}
	}
}

type Import struct {
	ClassName    string `json:"class_name"`
	ClassPackage string `json:"class_package"`
	ObjectName   string `json:"object_name"`
	OuterIndex   int64  `json:"outer_index"`
	OuterPackage Index  `json:"outer_package"`
}

func (e *Import) Resolve(p *Package) {
	e.OuterPackage.Resolve(p)
}

type Index struct {
	Index     int
	Null      bool
	Reference any
}

type Reference interface {
	Resolve(p *Package)
}

func (e *Index) Resolve(p *Package) {
	if e.Null || e.Reference != nil {
		return
	}
	var ref Reference
	if e.Index >= 0 {
		ref = &p.Exports[e.Index]

	} else {
		ref = &p.Imports[-e.Index-1]
	}
	e.Reference = ref
	ref.Resolve(p)
}

func (i *Index) UnmarshalJSON(bytes []byte) error {
	var raw rawIndex
	err := json.Unmarshal(bytes, &raw)
	*i = raw.Convert()
	return err
}

func (i Index) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i.Index), 10)), nil
}

func (i Index) GetImport() *Import {
	return i.Reference.(*Import)
}

func (i Index) GetExport() *Import {
	return i.Reference.(*Import)
}

func indexFromPartialIndex(i properties.Index) *Index {
	return &Index{Index: i.Index, Null: i.Null}
}

type rawRecord struct {
	ExportRecord rawExportRecord `json:"export_record"`
	Exports      []rawExport     `json:"exports"`
	Summary      rawSummary      `json:"summary"`
}

type rawExportRecord struct {
	FileName string `json:"file_name"`
}

type rawExport struct {
	Data   rawExportData `json:"data"`
	Export rawExportInfo `json:"export"`
}

type rawExportData struct {
	Properties []properties.RawProperty `json:"properties"`
}

type rawExportInfo struct {
	ObjectName    string   `json:"object_name"`
	ClassIndex    rawIndex `json:"class_index"`
	SuperIndex    rawIndex `json:"super_index"`
	TemplateIndex rawIndex `json:"template_index"`
	OuterIndex    rawIndex `json:"outer_index"`
}

type rawSummary struct {
	Imports []Import
}

type rawIndex struct {
	Index     int
	Reference json.RawMessage
}

func (i rawIndex) Convert() Index {
	null := slices.Equal(i.Reference, []byte("null"))
	if len(i.Reference) < 20 && !null {
		fmt.Println(i.Reference, []byte("null"))
	}
	return Index{
		Index: i.Index,
		Null:  null,
	}
}
