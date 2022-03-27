package exports

//TODO Interpret raw export values into the below info. Also figure out if we want to follow the TODOs below
type ExportInfo struct {
	Components       []string //TODO Make this an ObjectProperty array
	OverriddenValues []string //TODO Make this a Property array
	NewValues        []string //TODO Make this a Property array
	ImportedClasses  []string
}
