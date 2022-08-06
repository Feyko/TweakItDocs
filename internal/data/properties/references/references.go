package references

type Reference struct {
	Index      int
	ObjectName string
}

func NewReference(jsonMap map[string]any) Reference {
	objectName := "NULLREF"
	if jsonMap["reference"] != nil {
		objectName = jsonMap["reference"].(map[string]any)["object_name"].(string)
	}
	return Reference{
		Index:      int(jsonMap["index"].(float64)),
		ObjectName: objectName,
	}
}
