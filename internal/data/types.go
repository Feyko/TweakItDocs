package data

import "TweakItDocs/internal/data/properties"

type Data struct {
	Classes []Class
}

type Class struct {
	Path          string
	ParentClass   string
	NewProperties []properties.Property
	NewDefaults   []properties.Property
	Components    []Component
	Functions     []Function
}

type Component struct {
}

type Function struct {
}
