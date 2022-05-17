package types

import (
	"bytes"
	"strconv"

	"github.com/mercadolibre/iso-8583/serdes"
)

type Byte struct {
	Desc
}

func (b Byte) Name() string {
	return "byte"
}

func (b Byte) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	valueInt, err := b.valueAsInt(value)
	if err != nil {
		return nil, err
	}

	data := []byte{byte(valueInt)}
	return bytes.NewBuffer(data), nil
}

func (b Byte) valueAsInt(value serdes.Value) (int, error) {
	valueStr, ok := value.(string)
	if !ok {
		return 0, SerializerError{
			Message: "invalid value", Value: value, Serdes: b,
		}
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, SerializerError{
			Message: "invalid data", Serdes: b, Value: value, Cause: err,
		}
	}
	return valueInt, nil
}

func (b Byte) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	if data.Len() == 0 {
		return nil, DeserializationError{
			Message: "does not has data enough", Serdes: b, Remaning: data.Len(),
		}
	}

	deserializedByte, err := data.ReadByte()
	if err != nil {
		return nil, DeserializationError{
			Message: "error reading byte", Serdes: b, Remaning: data.Len(),
		}
	}

	value := strconv.Itoa(int(deserializedByte))
	return value, nil
}
