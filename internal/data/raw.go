package data

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Asset struct {
	Class           string
	Package         string
	Name            string
	Data            AssetData
	ObjectHierarchy []IndexInfo
}

func (a Asset) UnmarshalJSON(bytes []byte) error {
	var raw rawAsset
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	a.Class = raw.Class
	a.Package = raw.Package
	a.Name = raw.Name
	a.ObjectHierarchy = raw.ObjectHierarchy

	data := getDataTypeForAssetClass(raw.Class)
	err = json.Unmarshal(raw.Data, &data)
	if err != nil {
		return errors.Wrap(err, "could not decode the AssetSerializedData")
	}
	a.Data = AssetData{&data}

	return nil
}

func getDataTypeForAssetClass(assetClass string) any {
	switch assetClass {
	case "Blueprint":
		return BlueprintData{}
	case "UserDefinedEnum":
		return EnumData{}
	case "UserDefinedStruct":
		return StructData{}
	}
	panic("unsupported asset class: " + assetClass)
	return nil
}

type rawAsset struct {
	Class           string
	Package         string
	Name            string
	Data            json.RawMessage
	ObjectHierarchy []IndexInfo
}

type AssetData struct {
	v any
}

func (i *AssetData) IsBlueprint() bool {
	_, ok := i.v.(*BlueprintData)
	return ok
}

func (i *AssetData) AsBlueprint() *BlueprintData {
	return i.v.(*BlueprintData)
}

func (i *AssetData) IsEnum() bool {
	_, ok := i.v.(*EnumData)
	return ok
}

func (i *AssetData) AsEnum() *EnumData {
	return i.v.(*EnumData)
}

func (i *AssetData) AsStruct() *StructData {
	return i.v.(*StructData)
}

func (i *AssetData) IsStruct() bool {
	_, ok := i.v.(*StructData)
	return ok
}

type BlueprintData struct {
}

type EnumData struct {
}

type StructData struct {
}

type Index int

type IndexInfo struct {
	// Shared
	Type       string
	ObjectName string
	Outer      Index

	// Import
	ClassPackage string
	ClassName    string

	// Export
	ObjectClass Index
	ObjectFlags int64
	Properties  any
	ObjectMark  string
}

func (i *IndexInfo) IsExport() bool {
	return i.Type == "Export"
}

func (i *IndexInfo) AsExport() ExportIndexInfo {
	return ExportIndexInfo{
		ObjectName:  i.ObjectName,
		Outer:       i.Outer,
		ObjectClass: i.ObjectClass,
		ObjectFlags: i.ObjectFlags,
		Properties:  i.Properties,
	}
}

func (i *IndexInfo) IsImport() bool {
	return i.Type == "Import"
}

func (i *IndexInfo) AsImport() ImportIndexInfo {
	return ImportIndexInfo{
		ObjectName:   i.ObjectName,
		Outer:        i.Outer,
		ClassPackage: i.ClassPackage,
		ClassName:    i.ClassName,
	}
}

func (i IndexInfo) IsSelfPackage() bool {
	return i.ObjectMark != ""
}

func (i IndexInfo) IsSelfClass() bool {
	return i.ObjectName == ""
}

type ExportIndexInfo struct {
	ObjectName  string
	Outer       Index
	ObjectClass Index
	ObjectFlags int64
	Properties  any
}

type ImportIndexInfo struct {
	ObjectName   string
	Outer        Index
	ClassPackage string
	ClassName    string
}
