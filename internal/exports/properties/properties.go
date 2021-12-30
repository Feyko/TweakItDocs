package properties

import (
	"TweakItDocs/internal/exports/properties/references"
	"TweakItDocs/internal/sjsonhelp"
	"fmt"
	"strings"
)

func typeFromStrType(s string) Property {
	var r Property
	switch s {
	case "ObjectProperty":
		r = ObjectProperty{}
	case "TextProperty":
		r = TextProperty{}
	case "ArrayProperty":
		r = ArrayProperty{}
	case "StructProperty":
		r = StructProperty{}
	case "EnumProperty", "ByteProperty", "FloatProperty", "NameProperty", "IntProperty":
		r = tagProperty{}
	default:
		//fmt.Printf("Could not find property type for %v. Defaulting to tag\n", s)
		r = tagProperty{}
	}
	return r
}

type fields struct {
	name         string
	propertyType string
	tag          interface{}
	tag_data     interface{}
}

func fieldsFromUnknownMap(json sjsonhelp.JsonMap) fields {
	return fields{
		name:         json["name"].(string),
		propertyType: json["property_type"].(string),
		tag:          json["tag"],
		tag_data:     json["tag_data"],
	}
}

type Property interface {
	Type() string
	Value() interface{}
	New(fields) Property
	Name() string
}

type baseProperty struct {
	name         string
	propertyType string
	value        interface{}
}

func (p baseProperty) Type() string {
	return p.propertyType
}

func (o baseProperty) Value() interface{} {
	return o.value
}

func (o baseProperty) Name() string {
	return o.name
}

type tagProperty struct {
	baseProperty
}

func (p tagProperty) New(f fields) Property {
	return tagProperty{
		baseProperty{
			name:         f.name,
			value:        f.tag,
			propertyType: strings.TrimRight(f.propertyType, "Property"),
		},
	}
}

type TextProperty struct {
	baseProperty
}

func (o TextProperty) New(f fields) Property {
	return TextProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Text",
			value:        f.tag.(map[string]interface{})["source_string"],
		},
	}
}

type ObjectProperty struct {
	baseProperty
}

func (o ObjectProperty) New(f fields) Property {
	//fmt.Println(f)
	return ObjectProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Object",
			value:        references.NewReference(f.tag.(map[string]interface{})),
		},
	}
}

type ArrayProperty struct {
	baseProperty
}

func (o ArrayProperty) New(f fields) Property {
	innerType := f.tag_data.(string)
	propertyType := typeFromStrType(innerType)
	values := f.tag.([]interface{})
	out := make([]interface{}, len(values))
	for i, v := range values {
		var p Property
		if _, ok := propertyType.(StructProperty); ok {
			p = makeStructProperty(v)
		} else {
			p = makeAnonymousProperty(propertyType, v)
		}
		out[i] = p.Value()
	}
	return ArrayProperty{
		baseProperty{
			name:         f.name,
			propertyType: fmt.Sprintf("Array of %vs", innerType),
			value:        out,
		},
	}
}

func makeAnonymousProperty(propertyType Property, v interface{}) Property {
	return propertyType.New(fields{
		name: "NONE",
		tag:  v,
	})
}

func makeStructProperty(v interface{}) Property {
	//fmt.Println(v)
	m := v.(map[string]interface{})
	innerTagData := m["inner_tag_data"].(map[string]interface{})
	innerTagData["tag"] = m["properties"]
	return StructProperty{}.New(fields{
		name:         "NONE",
		propertyType: "Struct",
		tag:          innerTagData,
		tag_data:     innerTagData["tag_data"],
	})
}

type StructProperty struct {
	baseProperty
}

func (o StructProperty) New(f fields) Property {
	properties := make(map[string]interface{})
	switch v := f.tag.(type) {
	case []interface{}:
		for _, p := range v {
			pmap := p.(map[string]interface{})
			propertyType := typeFromStrType(pmap["property_type"].(string))
			property := propertyType.New(fieldsFromUnknownMap(pmap))
			properties[property.Name()] = property.Value()
		}
	case map[string]interface{}:
		properties = v
	}

	structName := f.tag_data.(map[string]interface{})["type"].(string)

	return StructProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Struct",
			value: map[string]interface{}{
				"name":       structName,
				"properties": properties,
			},
		},
	}
}

func New(json sjsonhelp.JsonMap) Property {
	values := fieldsFromUnknownMap(json)
	propertyType := typeFromStrType(values.propertyType)
	//fmt.Println(propertyType)
	property := propertyType.New(values)
	return property
}

func ToMap(p Property) sjsonhelp.JsonMap {
	return sjsonhelp.JsonMap{
		"name":  p.Name(),
		"type":  p.Type(),
		"value": p.Value(),
	}
}
