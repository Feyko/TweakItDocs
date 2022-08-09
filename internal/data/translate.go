package data

import "github.com/pkg/errors"

func TranslateAll(packages []Package) ([]Data, error) {
	r := make([]Data, len(packages))
	for i, pkg := range packages {
		data, err := TranslateOne(pkg)
		if err != nil {
			return nil, errors.Wrapf(err, "error on package %v", pkg.Filename)
		}
		r[i] = data
	}
	return r, nil
}

func TranslateOne(pkg Package) (Data, error) {
	data := Data{}
	for _, export := range pkg.Exports {
		if !export.SuperIndex.Null {
			if !export.OuterIndex.Null {
				panic("bro wtf")
			}
			data.Classes = append(data.Classes, Class{
				Path:          pkg.ClassPath(export.ObjectName),
				ParentClass:   export.SuperIndex.Path(),
				NewProperties: nil,
				NewDefaults:   nil,
				Components:    nil,
				Functions:     nil,
			})
		}
	}
	return data, nil
}
