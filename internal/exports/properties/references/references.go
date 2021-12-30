package references

import (
	"TweakItDocs/internal/sjsonhelp"
)

type Reference struct {
	Index      int64
	ObjectName string
}

func NewReference(jsonMap sjsonhelp.JsonMap) Reference {
	objectName := "NULLREF"
	if jsonMap["reference"] != nil {
		objectName = jsonMap["reference"].(map[string]interface{})["object_name"].(string)
	}
	return Reference{
		Index:      jsonMap["index"].(int64),
		ObjectName: objectName,
	}
}
