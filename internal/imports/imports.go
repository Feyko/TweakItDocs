package imports

import (
	"log"
)

type importValues struct {
	ClassName    string `structs:"class_name"`
	ClassPackage string `structs:"class_package"`
	ObjectName   string `structs:"object_name"`
	OuterIndex   int64  `structs:"outer_index"`

	// Not from the pak, filled by some logic
	Index   int    `structs:"index"`
	Package string `structs:"package"`
}

func arrayIndexToPakStyleIndex(i int) int {
	return -i - 1
}

func pakStyleIndexToArrayIndex(i int) int {
	return arrayIndexToPakStyleIndex(i) // Logic is actually the same
}

func fillPackagesOfImportValuesArray(imports []importValues) []importValues {
	for i, imp := range imports {
		if imp.OuterIndex == 0 {
			continue
		}
		imports[i] = fillPackageOfImportValues(imp, imports)
	}
	return imports
}

func fillPackageOfImportValues(imp importValues, imports []importValues) importValues {
	outerPackage := findPackageOfImport(imp, imports)
	imp.Package = outerPackage.ObjectName
	return imp
}

func findPackageOfImport(imp importValues, imports []importValues) importValues {
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
