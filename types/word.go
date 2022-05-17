package types

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/mercadolibre/go-iso8583/serdes"
)

type Word struct {
	Desc
	Order binary.ByteOrder
}

func (w Word) Name() string {
	return "word"
}

func (w Word) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	valueWord, err := w.valueAsInt(value)
	if err != nil {
		return nil, err
	}

	data := new(bytes.Buffer)
	if err := binary.Write(data, w.Order, uint16(valueWord)); err != nil {
		return nil, SerializerError{
			Message: "error encoding value", Serdes: w, Value: value, Cause: err,
		}
	}

	return data, nil
}

func (w Word) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	if data.Len() < 2 {
		return nil, DeserializationError{
			Message: "does not has data enough", Serdes: w, Remaning: data.Len(),
		}
	}

	var valueWord uint16
	if err := binary.Read(data, w.Order, &valueWord); err != nil {
		return nil, DeserializationError{
			Message: "error decoding data", Serdes: w, Remaning: data.Len(), Cause: err,
		}
	}

	value := strconv.FormatUint(uint64(valueWord), 10)
	return value, nil
}

func (w Word) valueAsInt(value serdes.Value) (int, error) {
	valueStr, ok := value.(string)
	if !ok {
		return 0, SerializerError{
			Message: "invalid value", Value: value, Serdes: w,
		}
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, SerializerError{
			Message: "invalid data", Serdes: w, Value: value, Cause: err,
		}
	}
	return valueInt, nil
}
