package types

import (
	"bytes"
	"fmt"

	"github.com/mercadolibre/go-iso8583/serdes"
)

type List struct {
	Desc
	Items []Field
}

func (list List) Name() string {
	return "list"
}

func (list List) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	mapValue, ok := value.(serdes.Map)
	if !ok {
		msg := fmt.Sprintf("invalid value [%T], expected: %T", value, serdes.Map{})
		return nil, SerializerError{
			Message: msg, Value: value, Serdes: list,
		}
	}

	serializedData := new(bytes.Buffer)
	for _, field := range list.Items {
		var itemValue serdes.Value
		if field.Name == "" {
			itemValue = mapValue
		} else {
			v, itemExists := mapValue[field.Name]
			if !itemExists {
				continue
			}
			itemValue = v
		}

		data, err := field.SerDes.Serialize(itemValue)
		if err != nil {
			return nil, SerializerError{
				Message: "field serializer failed", Serdes: list, Field: field, Value: value, Cause: err,
			}
		}

		if _, err := serializedData.ReadFrom(data); err != nil {
			return nil, SerializerError{
				Message: "buffer failed", Serdes: list, Field: field, Value: value, Cause: err,
			}
		}
	}

	return serializedData, nil
}

func (list List) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	listValues := serdes.Map{}
	for _, field := range list.Items {
		if data.Len() == 0 {
			break
		}

		fieldValue, err := field.SerDes.Deserialize(data)
		if err != nil {
			return listValues, DeserializationError{
				Message: "field deserializer failed", Serdes: list, Field: field, Remaning: data.Len(), Cause: err,
			}
		}

		if field.Name != "" {
			listValues[field.Name] = fieldValue
			continue
		}

		mapValue, ok := fieldValue.(serdes.Map)
		if !ok {
			return listValues, DeserializationError{
				Message: "field deserializer failed, anonymous field requires a map value", Serdes: list, Field: field, Remaning: data.Len(), Cause: err,
			}
		}

		for mapKep, mapItem := range mapValue {
			listValues[mapKep] = mapItem
		}
	}

	return listValues, nil
}
