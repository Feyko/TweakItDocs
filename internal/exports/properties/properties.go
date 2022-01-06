package properties

import (
	"TweakItDocs/internal/exports/properties/references"
	"TweakItDocs/internal/sjsonhelp"
	"fmt"
	"log"
	"strings"
)

func typeFromStrType(s string) Property {
	var r Property
	switch propTypeToValueType(s) {
	case "Object":
		r = ObjectProperty{}
	case "Text":
		r = TextProperty{}
	case "Array":
		r = ArrayProperty{}
	case "Struct":
		r = StructProperty{}
	case "Bool":
		r = BoolProperty{}
	case "SoftObject":
		r = tagProperty{} // TODO: Make a custom type for this
	case "Map":
		r = MapProperty{}
	case "Delegate", "MulticastInlineDelegate", "MulticastSparseDelegate":
		r = EventProperty{}
	case "Set":
		r = SetProperty{}
	case "Enum", "Byte", "Float", "Name", "Str", "FieldPath", "Interface",
		"Int8", "Int16", "Int", "Int64",
		"UInt16", "UInt32", "UInt64":
		r = tagProperty{}
	default:
		log.Fatalf("Could not find property type for %v\n", s)
		//r = tagProperty{}
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
			propertyType: propTypeToValueType(f.propertyType),
		},
	}
}

type BoolProperty struct {
	baseProperty
}

func (o BoolProperty) New(f fields) Property {
	value := f.tag
	if value == nil {
		value = f.tag_data
	}
	return BoolProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Bool",
			value:        value.(bool),
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

type EventProperty struct {
	baseProperty
}

func (o EventProperty) New(f fields) Property {
	return EventProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Event",
			value:        nil,
		},
	}
}

type ObjectProperty struct {
	baseProperty
}

func (o ObjectProperty) New(f fields) Property {
	return ObjectProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Object",
			value:        references.NewReference(f.tag.(map[string]interface{})),
		},
	}
}

type SetProperty struct {
	baseProperty
}

func (o SetProperty) New(f fields) Property {
	innerType := f.tag_data.(string)
	return SetProperty{
		baseProperty{
			name:         f.name,
			propertyType: fmt.Sprintf("Set of %vs", propTypeToValueType(innerType)),
			value:        nil,
		},
	}
}

type ArrayProperty struct {
	baseProperty
}

// TODO Change the container properties (Array, Map, Set) to have the inner type in their value field

func (o ArrayProperty) New(f fields) Property {
	innerType := f.tag_data.(string)
	propertyType := typeFromStrType(innerType)
	isInnerTypeStruct := isStructProperty(propertyType)
	values := f.tag.([]interface{})
	out := make([]interface{}, len(values))
	for i, v := range values {
		var tag_data interface{} = nil
		if isInnerTypeStruct {
			structValue := arrayStructValueToNormalStructValue(v)
			v = structValue["tag"]
			tag_data = structValue["tag_data"]
		}
		out[i] = propertyTypeToPropertyValue(propertyType, v, tag_data)
	}
	return ArrayProperty{
		baseProperty{
			name:         f.name,
			propertyType: fmt.Sprintf("Array of %vs", propTypeToValueType(innerType)),
			value:        out,
		},
	}
}

func makeAnonymousProperty(propertyType Property, tag, tag_data interface{}) Property {
	return propertyType.New(fields{
		name:     "NONE",
		tag:      tag,
		tag_data: tag_data,
	})
}

func arrayStructValueToNormalStructValue(v interface{}) map[string]interface{} {
	m := v.(map[string]interface{})
	innerTagData := m["inner_tag_data"].(map[string]interface{})
	innerTagData["tag"] = m["properties"]
	return innerTagData
}

type StructProperty struct {
	baseProperty
}

func (o StructProperty) New(f fields) Property {
	properties := structValueToPropertyMaps(f.tag)

	structType := ""
	if f.tag_data != nil {
		structType = f.tag_data.(map[string]interface{})["type"].(string)
	}

	return StructProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Struct",
			value: map[string]interface{}{
				"struct_type": structType,
				"properties":  properties,
			},
		},
	}
}

func isStructProperty(p Property) bool {
	_, ok := p.(StructProperty)
	return ok
}

type MapProperty struct {
	baseProperty
}

func (o MapProperty) New(f fields) Property {
	tag_data := f.tag_data.(map[string]interface{})

	keyPropertyTypeStr := tag_data["key_type"].(string)
	valuePropertyTypeStr := tag_data["value_type"].(string)
	keyPropertyType := typeFromStrType(keyPropertyTypeStr)
	valuePropertyType := typeFromStrType(valuePropertyTypeStr)

	out := make([]sjsonhelp.JsonMap, 0)

	var values []interface{}

	if f.tag == nil {
		goto skip // The value of a map can be nil. In that case, completely skip the value handling
	}

	values = f.tag.([]interface{})
	for _, v := range values {
		vMap := v.(map[string]interface{})

		key := vMap["key"]
		value := vMap["value"]

		key = propertyTypeToPropertyValue(keyPropertyType, key, nil)
		value = propertyTypeToPropertyValue(valuePropertyType, value, nil)

		out = append(out, sjsonhelp.JsonMap{"key": key, "value": value})
	}

skip:

	return MapProperty{
		baseProperty{
			name:         f.name,
			propertyType: fmt.Sprintf("Map of %vs to %vs", propTypeToValueType(keyPropertyTypeStr), propTypeToValueType(valuePropertyTypeStr)),
			value:        out,
		},
	}
}

func propertyTypeToPropertyValue(propertyType Property, value, tag_data interface{}) interface{} {
	property := makeAnonymousProperty(propertyType, value, tag_data)
	return property.Value()
}

func structValueToPropertyMaps(value interface{}) []sjsonhelp.JsonMap {
	var r []sjsonhelp.JsonMap
	switch v := value.(type) {
	case map[string]interface{}:
		innerValue := v["value"]
		if innerValue == nil {
			return make([]sjsonhelp.JsonMap, 0)
		}
		r = structMapToPropertyMaps(innerValue.(map[string]interface{}))
	case []interface{}:
		r = mapArrayOfProperties(v)
	case nil:
		r = make([]sjsonhelp.JsonMap, 0)
	default:
		log.Fatalf("Unsupported struct value format: %#v", value)
	}
	return r
}

func mapArrayOfProperties(properties []interface{}) []sjsonhelp.JsonMap {
	r := make([]sjsonhelp.JsonMap, len(properties))
	for i, p := range properties {
		property := New(p.(map[string]interface{}))
		r[i] = ToMap(property)
	}
	return r
}

func structMapToPropertyMaps(m sjsonhelp.JsonMap) []sjsonhelp.JsonMap {
	r := make([]sjsonhelp.JsonMap, 0, len(m))
	for k, v := range m {
		innerType := ""
		var value interface{}
		switch inner := v.(type) {
		// Have to have those two cases. Really annoying, would be nice to get rid of this
		case sjsonhelp.JsonMap:
			innerType = "Struct"
			value = structMapToPropertyMaps(inner)
		case map[string]interface{}:
			innerType = "Struct"
			value = structMapToPropertyMaps(inner)
		case int, int64:
			innerType = "Int"
			value = inner
		case float64:
			innerType = "Float"
			value = inner
		default:
			log.Fatalf("Unsupported value type in struct map: %#v", v)
		}
		r = append(r, makePropertyMap(
			k,
			innerType,
			value,
		))
	}
	return r
}

func New(json sjsonhelp.JsonMap) Property {
	values := fieldsFromUnknownMap(json)
	propertyType := typeFromStrType(values.propertyType)
	property := propertyType.New(values)
	return property
}

func ToMap(p Property) sjsonhelp.JsonMap {
	return makePropertyMap(
		p.Name(),
		p.Type(),
		p.Value(),
	)
}

func makePropertyMap(name, type_ string, value interface{}) sjsonhelp.JsonMap {
	return sjsonhelp.JsonMap{
		"name":  name,
		"type":  type_,
		"value": value,
	}
}

func propTypeToValueType(s string) string {
	return strings.TrimSuffix(s, "Property")
}
