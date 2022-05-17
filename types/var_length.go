package types

import (
	"bytes"
	"strconv"

	"github.com/mercadolibre/go-iso8583/serdes"
)

type VarLength struct {
	Desc
	Length serdes.Serdes
	Data   serdes.Serdes
}

func (varLen VarLength) Name() string {
	return "var_length"
}

func (varLen VarLength) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	serializedData, err := varLen.Data.Serialize(value)
	if err != nil {
		return nil, SerializerError{
			Message: "error serializing data", Serdes: varLen, Value: value, Cause: err,
		}
	}

	deserializedLength := serializedData.Len()
	if _, ok := varLen.Data.(Bcd); ok {
		deserializedLength = deserializedLength * 2
		// Visa specify (pag. 76) that BCD types must indicate real data size, ignoring leading zeros.
		if strType, ok := value.(string); ok {
			deserializedLength = len(strType)
		}
	}

	deserializedLengthStr := strconv.Itoa(deserializedLength)
	serializedLength, err := varLen.Length.Serialize(deserializedLengthStr)
	if err != nil {
		return nil, SerializerError{
			Message: "error serializing length", Serdes: varLen, Value: value, Cause: err,
		}
	}

	out := append(serializedLength.Bytes(), serializedData.Bytes()...)
	return bytes.NewBuffer(out), nil
}

func (varLen VarLength) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	deserializedLength, err := varLen.Length.Deserialize(data)
	if err != nil {
		return nil, DeserializationError{
			Message: "error deserializing length", Serdes: varLen, Remaning: data.Len(), Cause: err,
		}
	}

	lengthIn, err := varLen.valueAsInt(deserializedLength, data)
	if err != nil {
		return nil, err
	}

	if _, ok := varLen.Data.(Bcd); ok {
		numBytes := lengthIn / 2
		if lengthIn > 0 && lengthIn%2 != 0 {
			numBytes++
		}
		lengthIn = numBytes
	}

	if data.Len() < lengthIn {
		return nil, DeserializationError{
			Message: "does not has bytes enough in data buffer", Serdes: varLen, Remaning: data.Len(), Cause: err,
		}
	}

	if lengthIn == 0 {
		return bytes.NewBuffer([]byte{}), nil
	}

	serializedData := bytes.NewBuffer(data.Next(lengthIn))
	deserializedData, err := varLen.Data.Deserialize(serializedData)
	if err != nil {
		return nil, DeserializationError{
			Message: "deserializer failed", Serdes: varLen, Remaning: data.Len(), Cause: err,
		}
	}

	return deserializedData, nil
}

func (varLen VarLength) valueAsInt(deserializedLength serdes.Value, data *bytes.Buffer) (int, error) {
	lengthStr, ok := deserializedLength.(string)
	if !ok {
		return 0, DeserializationError{
			Message: "length deserializer returned an invalid type", Serdes: varLen, Remaning: data.Len(),
		}
	}

	lengthIn, err := strconv.Atoi(lengthStr)
	if !ok {
		return 0, DeserializationError{
			Message: "length deserializer returned a non numeric string", Serdes: varLen, Remaning: data.Len(), Cause: err,
		}
	}

	return lengthIn, nil
}
