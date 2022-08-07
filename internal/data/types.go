package data

import "TweakItDocs/internal/data/properties"

type Class struct {
	Path          string
	ParentClass   string
	NewProperties []properties.Property
	NewDefaults   []properties.Property
}
