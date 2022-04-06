package imports

import (
	"TweakItDocs/internal/data"
	"log"
)

type ImportValues struct {
	data.Import

	// Not from the pak, filled by some logic
	Index   int    `json:"index,omitempty"`
	Package string `json:"package,omitempty"`
}

type ResolvedRecord struct {
	data.Record
	DeepestPackage            data.Import
	ImportsFromDeepestPackage []data.Import
}

func Filter(records []data.Record) []ResolvedRecord {
	r := make([]ResolvedRecord, len(records))
	for i, record := range records {
		deepestPackage := findLastPackageImport(record.Imports)
		importsFromDeepest := findImportsFromPackage(record.Imports, deepestPackage.ObjectName)
		newRecord := ResolvedRecord{record, deepestPackage, importsFromDeepest}
		r[i] = newRecord
	}
	return r
}

func findLastPackageImport(imports []data.Import) data.Import {
	deepest := 0
	for i, imp := range imports {
		if imp.ClassName == "Package" && i > deepest {
			deepest = i
		}
	}
	return imports[deepest]
}

func findImportsFromPackage(imports []data.Import, packageName string) []data.Import {
	r := make([]data.Import, 0, len(imports)/3)
	for _, imp := range imports {
		if imp.ClassPackage == packageName {
			r = append(r, imp)
		}
	}
	return r
}

func arrayIndexToPakStyleIndex(i int) int {
	return -i - 1
}

func pakStyleIndexToArrayIndex(i int) int {
	return arrayIndexToPakStyleIndex(i) // Logic is actually the same
}

func fillPackagesOfImportValuesArray(imports []ImportValues) []ImportValues {
	for i, imp := range imports {
		if imp.OuterIndex == 0 {
			continue
		}
		imports[i] = fillPackageOfImportValues(imp, imports)
	}
	return imports
}

func fillPackageOfImportValues(imp ImportValues, imports []ImportValues) ImportValues {
	outerPackage := findPackageOfImport(imp, imports)
	imp.Package = outerPackage.ObjectName
	return imp
}

func findPackageOfImport(imp ImportValues, imports []ImportValues) ImportValues {
	index := imp.OuterIndex
	next := imports[pakStyleIndexToArrayIndex(int(index))]
	for next.OuterIndex != 0 {
		if next.ClassName == "Package" {
			log.Fatalf("Found a package import that isn't the root: %v", next)
		}
		next = imports[pakStyleIndexToArrayIndex(int(next.OuterIndex))]
	}
	if next.ClassName != "Package" {
		log.Fatalf("Found a root import that isn't a package: %v", next)
	}
	return next
}
