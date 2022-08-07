package properties

import (
	"log"
	"strings"
)

type Property struct {
	PropertyType string `json:"property_type"`
	Name         string `json:"name"`
	Value        any    `json:"value"`
}

type RawProperty struct {
	Name    string `json:"name"`
	Type    string `json:"property_type"`
	Tag     any    `json:"tag"`
	TagData any    `json:"tag_data"`
}

func jsonToRawProperty(json map[string]any) RawProperty {
	return RawProperty{
		Name:    json["name"].(string),
		Type:    json["property_type"].(string),
		Tag:     json["tag"],
		TagData: json["tag_data"],
	}
}

func typeFromStrType(s string) UnknownProperty {
	var r UnknownProperty
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

type UnknownProperty interface {
	Type() string
	Value() any
	New(RawProperty) UnknownProperty
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

func (p tagProperty) New(f RawProperty) UnknownProperty {
	return tagProperty{
		baseProperty{
			name:         f.Name,
			value:        f.Tag,
			propertyType: propTypeToValueType(f.Type),
		},
	}
}

type BoolProperty struct {
	baseProperty
}

func (o BoolProperty) New(f RawProperty) UnknownProperty {
	value := f.Tag
	if value == nil {
		value = f.TagData
	}
	return BoolProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Bool",
			value:        value.(bool),
		},
	}
}

type TextProperty struct {
	baseProperty
}

func (o TextProperty) New(f RawProperty) UnknownProperty {
	return TextProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Text",
			value:        f.Tag.(map[string]any)["source_string"],
		},
	}
}

type EventProperty struct {
	baseProperty
}

func (o EventProperty) New(f RawProperty) UnknownProperty {
	return EventProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Event",
			value:        nil,
		},
	}
}

type SoftObjectProperty struct {
	baseProperty
}

func (o SoftObjectProperty) New(f RawProperty) UnknownProperty {
	return SoftObjectProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "SoftObject",
			value:        nil,
		},
	}
}

type ObjectProperty struct {
	baseProperty
}

func (o ObjectProperty) New(f RawProperty) UnknownProperty {
	v := f.Tag.(map[string]any)
	return ObjectProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Object",
			value:        Index{Index: int(v["index"].(float64)), Null: v["reference"] == nil},
		},
	}
}

// Redefining a simpler Index type here because properties depending on the real Index makes a cyclic dependency
type Index struct {
	Index int
	Null  bool
}

type SetProperty struct {
	baseProperty
}

func (o SetProperty) New(f RawProperty) UnknownProperty {
	innerType := f.TagData.(string)
	innerValueType := propTypeToValueType(innerType)
	return SetProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Set",
			value: map[string]any{
				"inner_type": innerValueType,
				"value":      nil,
			},
		},
	}
}

type ArrayProperty struct {
	baseProperty
}

func (o ArrayProperty) New(f RawProperty) UnknownProperty {
	innerType := f.TagData.(string)
	propertyType := typeFromStrType(innerType)
	isInnerTypeStruct := isStructProperty(propertyType)
	values := f.Tag.([]any)
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
	value := map[string]any{
		"inner_type": propTypeToValueType(innerType),
		"value":      out,
	}
	return ArrayProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Array",
			value:        value,
		},
	}
}

func makeAnonymousProperty(propertyType UnknownProperty, tag, tag_data any) UnknownProperty {
	return propertyType.New(RawProperty{
		Tag:     tag,
		TagData: tag_data,
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

func (o StructProperty) New(f RawProperty) UnknownProperty {
	properties := structValueToPropertyMaps(f.Tag)

	structType := ""
	if f.TagData != nil {
		structType = f.TagData.(map[string]any)["type"].(string)
	}

	return StructProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Struct",
			value: map[string]any{
				"struct_type": structType,
				"properties":  properties,
			},
		},
	}
}

func isStructProperty(p UnknownProperty) bool {
	_, ok := p.(StructProperty)
	return ok
}

type MapProperty struct {
	baseProperty
}

func (o MapProperty) New(f RawProperty) UnknownProperty {
	tag_data := f.TagData.(map[string]any)

	keyPropertyTypeStr := tag_data["key_type"].(string)
	valuePropertyTypeStr := tag_data["value_type"].(string)
	keyPropertyType := typeFromStrType(keyPropertyTypeStr)
	valuePropertyType := typeFromStrType(valuePropertyTypeStr)

	out := make([]map[string]any, 0)

	var values []any

	if f.Tag == nil {
		goto skip // The value of a map can be nil. In that case, completely skip the value handling
	}

	values = f.Tag.([]any)
	for _, v := range values {
		vMap := v.(map[string]any)

		key := vMap["key"]
		value := vMap["value"]

		key = propertyTypeToPropertyValue(keyPropertyType, key, nil)
		value = propertyTypeToPropertyValue(valuePropertyType, value, nil)

		out = append(out, map[string]any{"key": key, "value": value})
	}

skip:
	keyValueTypeStr := propTypeToValueType(keyPropertyTypeStr)
	valueValueTypeStr := propTypeToValueType(valuePropertyTypeStr)

	value := map[string]any{
		"inner_type_key":   keyValueTypeStr,
		"inner_type_value": valueValueTypeStr,
		"value":            out,
	}

	return MapProperty{
		baseProperty{
			name:         f.Name,
			propertyType: "Map",
			value:        value,
		},
	}
}

func propertyTypeToPropertyValue(propertyType UnknownProperty, value, tag_data any) any {
	property := makeAnonymousProperty(propertyType, value, tag_data)
	return property.Value()
}

func structValueToPropertyMaps(value any) []map[string]any {
	var r []map[string]any
	switch v := value.(type) {
	case map[string]any:
		innerValue := v["value"]
		if innerValue == nil {
			return make([]map[string]any, 0)
		}
		r = structMapToPropertyMaps(innerValue.(map[string]any))
	case []any:
		r = mapArrayOfProperties(v)
	case nil:
		r = make([]map[string]any, 0)
	default:
		log.Fatalf("Unsupported struct value format: %#v", value)
	}
	return r
}

func mapArrayOfProperties(properties []any) []map[string]any {
	r := make([]map[string]any, len(properties))
	for i, p := range properties {
		raw := jsonToRawProperty(p.(map[string]any))
		property := newProperty(raw)
		r[i] = ToMap(property)
	}
	return r
}

func structMapToPropertyMaps(m map[string]any) []map[string]any {
	r := make([]map[string]any, 0, len(m))
	for k, v := range m {
		innerType := ""
		var value any
		switch inner := v.(type) {
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

func newProperty(raw RawProperty) UnknownProperty {
	propertyType := typeFromStrType(raw.Type)
	property := propertyType.New(raw)
	return property
}

func ToMap(p UnknownProperty) map[string]any {
	return makePropertyMap(
		p.Name(),
		p.Type(),
		p.Value(),
	)
}

func makePropertyMap(name, type_ string, value any) map[string]any {
	return map[string]any{
		"name":  name,
		"type":  type_,
		"value": value,
	}
}

func propTypeToValueType(s string) string {
	return strings.TrimSuffix(s, "Property")
}

func resolve(p UnknownProperty) Property {
	return Property{
		p.Type(),
		p.Name(),
		p.Value(),
	}
}

func DataToProperty(data RawProperty) Property {
	p := newProperty(data)
	return resolve(p)
}

//TODO break down the big functions in this file & merge the internal and exported functions, this is a mess
