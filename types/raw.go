package types

import (
	"bytes"
	"encoding/hex"

	"github.com/mercadolibre/go-iso8583/serdes"
)

type Raw struct {
	Desc
	NumBytes int
}

func (raw Raw) Name() string {
	return "raw"
}

func (raw Raw) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	rawValue, err := raw.normalizeData(value)
	if err != nil {
		return nil, err
	}

	numBytes := raw.NumBytes
	if numBytes == 0 {
		numBytes = len(rawValue)
	}

	if raw.NumBytes > 0 && numBytes > raw.NumBytes {
		return nil, SerializerError{
			Message: "value too long", Serdes: raw, Value: value,
		}
	}

	return bytes.NewBuffer(rawValue), nil
}

func (raw Raw) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	numBytes := raw.NumBytes
	if numBytes == 0 {
		numBytes = data.Len()
	}

	if data.Len() < numBytes {
		return nil, DeserializationError{
			Message: "data does not has bytes enough", Serdes: raw, Remaning: data.Len(),
		}
	}

	rawValue := data.Next(numBytes)
	valueHex := hex.EncodeToString(rawValue)
	return valueHex, nil
}

func (raw Raw) normalizeData(value serdes.Value) ([]byte, error) {
	valueStr, ok := value.(string)
	if !ok {
		return nil, SerializerError{
			Message: "invalid value type", Serdes: raw, Value: value,
		}
	}

	if len(valueStr)%2 != 0 {
		return nil, SerializerError{
			Message: "invalid value, hex string with odd num of char", Serdes: raw, Value: value,
		}
	}

	rawValue, err := hex.DecodeString(valueStr)
	if err != nil {
		return nil, SerializerError{
			Message: "invalid value, error decoding hex str", Serdes: raw, Value: value, Cause: err,
		}
	}

	return rawValue, nil
}
