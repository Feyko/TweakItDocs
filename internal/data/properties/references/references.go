package references

type Reference struct {
	Index        int64  `json:"index"`
	ObjectName   string `json:"object_name"`
	SerialOffset int64  `json:"serial_offset"`
}

func NewReference(jsonMap map[string]any) Reference {
	objectName := "NULLREF"
	if jsonMap["reference"] != nil {
		objectName = jsonMap["reference"].(map[string]any)["object_name"].(string)
	}
	return Reference{
		Index:      int64(jsonMap["index"].(float64)),
		ObjectName: objectName,
	}
}
