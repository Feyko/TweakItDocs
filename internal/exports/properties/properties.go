package properties

import (
	"TweakItDocs/internal/exports/properties/references"
	"TweakItDocs/internal/sjsonhelp"
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
	case "Map":
		r = MapProperty{}
	case "Delegate", "MulticastInlineDelegate", "MulticastSparseDelegate":
		r = EventProperty{}
	case "Set":
		r = SetProperty{}
	case "Enum", "Byte", "Float", "Name", "Str", "FieldPath", "Interface", "SoftObject",
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
	tag          any
	tag_data     any
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
	Value() any
	New(fields) Property
	Name() string
}

type baseProperty struct {
	name         string
	propertyType string
	value        any
}

func (p baseProperty) Type() string {
	return p.propertyType
}

func (o baseProperty) Value() any {
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
			value:        f.tag.(map[string]any)["source_string"],
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

type SoftObjectProperty struct {
	baseProperty
}

func (o SoftObjectProperty) New(f fields) Property {
	return SoftObjectProperty{
		baseProperty{
			name:         f.name,
			propertyType: "SoftObject",
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
			value:        references.NewReference(f.tag.(map[string]any)),
		},
	}
}

type SetProperty struct {
	baseProperty
}

func (o SetProperty) New(f fields) Property {
	innerType := f.tag_data.(string)
	innerValueType := propTypeToValueType(innerType)
	return SetProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Set",
			value: sjsonhelp.JsonMap{
				"inner_type": innerValueType,
				"value":      nil,
			},
		},
	}
}

type ArrayProperty struct {
	baseProperty
}

func (o ArrayProperty) New(f fields) Property {
	innerType := f.tag_data.(string)
	propertyType := typeFromStrType(innerType)
	isInnerTypeStruct := isStructProperty(propertyType)
	values := f.tag.([]any)
	out := make([]any, len(values))
	for i, v := range values {
		var tag_data any = nil
		if isInnerTypeStruct {
			structValue := arrayStructValueToNormalStructValue(v)
			v = structValue["tag"]
			tag_data = structValue["tag_data"]
		}
		out[i] = propertyTypeToPropertyValue(propertyType, v, tag_data)
	}
	value := sjsonhelp.JsonMap{
		"inner_type": propTypeToValueType(innerType),
		"value":      out,
	}
	return ArrayProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Array",
			value:        value,
		},
	}
}

func makeAnonymousProperty(propertyType Property, tag, tag_data any) Property {
	return propertyType.New(fields{
		tag:      tag,
		tag_data: tag_data,
	})
}

func arrayStructValueToNormalStructValue(v any) map[string]any {
	m := v.(map[string]any)
	innerTagData := m["inner_tag_data"].(map[string]any)
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
		structType = f.tag_data.(map[string]any)["type"].(string)
	}

	return StructProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Struct",
			value: map[string]any{
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
	tag_data := f.tag_data.(map[string]any)

	keyPropertyTypeStr := tag_data["key_type"].(string)
	valuePropertyTypeStr := tag_data["value_type"].(string)
	keyPropertyType := typeFromStrType(keyPropertyTypeStr)
	valuePropertyType := typeFromStrType(valuePropertyTypeStr)

	out := make([]sjsonhelp.JsonMap, 0)

	var values []any

	if f.tag == nil {
		goto skip // The value of a map can be nil. In that case, completely skip the value handling
	}

	values = f.tag.([]any)
	for _, v := range values {
		vMap := v.(map[string]any)

		key := vMap["key"]
		value := vMap["value"]

		key = propertyTypeToPropertyValue(keyPropertyType, key, nil)
		value = propertyTypeToPropertyValue(valuePropertyType, value, nil)

		out = append(out, sjsonhelp.JsonMap{"key": key, "value": value})
	}

skip:
	keyValueTypeStr := propTypeToValueType(keyPropertyTypeStr)
	valueValueTypeStr := propTypeToValueType(valuePropertyTypeStr)

	value := sjsonhelp.JsonMap{
		"inner_type_key":   keyValueTypeStr,
		"inner_type_value": valueValueTypeStr,
		"value":            out,
	}

	return MapProperty{
		baseProperty{
			name:         f.name,
			propertyType: "Map",
			value:        value,
		},
	}
}

func propertyTypeToPropertyValue(propertyType Property, value, tag_data any) any {
	property := makeAnonymousProperty(propertyType, value, tag_data)
	return property.Value()
}

func structValueToPropertyMaps(value any) []sjsonhelp.JsonMap {
	var r []sjsonhelp.JsonMap
	switch v := value.(type) {
	case map[string]any:
		innerValue := v["value"]
		if innerValue == nil {
			return make([]sjsonhelp.JsonMap, 0)
		}
		r = structMapToPropertyMaps(innerValue.(map[string]any))
	case []any:
		r = mapArrayOfProperties(v)
	case nil:
		r = make([]sjsonhelp.JsonMap, 0)
	default:
		log.Fatalf("Unsupported struct value format: %#v", value)
	}
	return r
}

func mapArrayOfProperties(properties []any) []sjsonhelp.JsonMap {
	r := make([]sjsonhelp.JsonMap, len(properties))
	for i, p := range properties {
		property := New(p.(map[string]any))
		r[i] = ToMap(property)
	}
	return r
}

func structMapToPropertyMaps(m sjsonhelp.JsonMap) []sjsonhelp.JsonMap {
	r := make([]sjsonhelp.JsonMap, 0, len(m))
	for k, v := range m {
		innerType := ""
		var value any
		switch inner := v.(type) {
		// Have to have those two cases. Really annoying, would be nice to get rid of this
		case sjsonhelp.JsonMap:
			innerType = "Struct"
			value = structMapToPropertyMaps(inner)
		case map[string]any:
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

func makePropertyMap(name, type_ string, value any) sjsonhelp.JsonMap {
	return sjsonhelp.JsonMap{
		"name":  name,
		"type":  type_,
		"value": value,
	}
}

func propTypeToValueType(s string) string {
	return strings.TrimSuffix(s, "Property")
}
