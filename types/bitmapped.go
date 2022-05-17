package types

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/mercadolibre/go-iso8583/serdes"
)

type BitMapped struct {
	Desc
	Bitmap  Bitmap
	Mapping map[int]serdes.Serdes
}

func (bitMapped BitMapped) Name() string {
	return "bitMapped"
}

func (bitMapped BitMapped) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	bitsValue, sortedBitsNumber, err := bitMapped.normalizeValue(value)
	if err != nil {
		return nil, err
	}

	numBits := len(sortedBitsNumber)
	if numBits == 0 {
		return &bytes.Buffer{}, nil
	}

	bitmap, err := bitMapped.serializeBitmap(sortedBitsNumber, numBits)
	if err != nil {
		return nil, err
	}

	fields, err := bitMapped.serializeFields(sortedBitsNumber, bitsValue)
	if err != nil {
		return nil, err
	}

	data := append(bitmap.Bytes(), fields.Bytes()...)
	return bytes.NewBuffer(data), nil
}

func (bitMapped BitMapped) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	values := make(serdes.Map)
	bitmapValue, err := bitMapped.Bitmap.Deserialize(data)
	if err != nil {
		return nil, DeserializationError{
			Message: "error decoding bitmap", Serdes: bitMapped, Remaning: data.Len(), Cause: err,
		}
	}

	bitmap, ok := bitmapValue.([]byte)
	if !ok {
		return nil, DeserializationError{
			Message: "bitmap was deserialized to an invalid type", Serdes: bitMapped, Remaning: data.Len(), Cause: err,
		}
	}

	numBits := len(bitmap) * 8
	for bitNumber := 1; bitNumber <= numBits; bitNumber++ {
		bitIndex := bitNumber - 1
		if !bitMapped.checkBit(bitmap, bitIndex) {
			continue
		}

		deserializer, exists := bitMapped.Mapping[bitNumber]
		if !exists || deserializer == nil {
			return nil, DeserializationError{
				Message: fmt.Sprintf("bit %d not found", bitNumber), Serdes: bitMapped, Remaning: data.Len(), Cause: err,
			}
		}

		value, err := deserializer.Deserialize(data)
		if err != nil {
			return nil, DeserializationError{
				Message: fmt.Sprintf("deserialize bit %d failed", bitNumber), Serdes: bitMapped, Remaning: data.Len(), Cause: err,
			}
		}

		bitKey := strconv.Itoa(bitNumber)
		values[bitKey] = value
	}

	return values, nil
}

func (bitMapped BitMapped) normalizeValue(value serdes.Value) (map[int]interface{}, []int, error) {
	mapStringValue, ok := value.(serdes.Map)
	if !ok {
		return nil, nil, SerializerError{
			Message: "invalid value", Value: value, Serdes: bitMapped,
		}
	}

	var index int
	var list []int
	mapIntValue := map[int]interface{}{}

	for key, value := range mapStringValue {
		bitNumber, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		mapIntValue[bitNumber] = value
		list = append(list, bitNumber)
		index++
	}

	sort.Ints(list)
	return mapIntValue, list, nil
}

func (bitMapped BitMapped) serializeBitmap(sortedBitsNumber []int, numBits int) (*bytes.Buffer, error) {
	lastBitNumber := sortedBitsNumber[numBits-1]
	nbytes := lastBitNumber / 8
	if lastBitNumber%8 > 0 {
		nbytes++
	}

	bitmap := make([]byte, nbytes)
	for _, bitNumber := range sortedBitsNumber {
		bitIndex := bitNumber - 1
		cbit := 7 - uint(bitIndex%8)
		bitmap[bitIndex/8] |= 0x1 << cbit
	}

	data, err := bitMapped.Bitmap.Serialize(bitmap)
	if err != nil {
		return nil, SerializerError{
			Message: "error serializing bitmap", Serdes: bitMapped, Cause: err,
		}
	}

	return data, nil
}

func (bitMapped BitMapped) serializeFields(sortedBitsNumber []int, bitsValue map[int]interface{}) (*bytes.Buffer, error) {
	out := new(bytes.Buffer)
	for _, bitNumber := range sortedBitsNumber {
		serializer, exists := bitMapped.Mapping[bitNumber]
		if !exists || serializer == nil {
			return nil, SerializerError{
				Message: fmt.Sprintf("bit number %d not found", bitNumber), Serdes: bitMapped, Value: bitsValue,
			}
		}

		value := bitsValue[bitNumber]
		bitData, err := serializer.Serialize(value)
		if err != nil {
			return nil, SerializerError{
				Message: "serializer failed", Serdes: bitMapped, Field: Field{Name: strconv.Itoa(bitNumber), SerDes: serializer}, Value: value, Cause: err,
			}
		}

		if _, err := out.ReadFrom(bitData); err != nil {
			return nil, SerializerError{
				Message: "error coping buffer", Serdes: bitMapped, Field: Field{Name: strconv.Itoa(bitNumber), SerDes: serializer}, Value: value, Cause: err,
			}
		}
	}

	return out, nil
}

func (bitMapped *BitMapped) checkBit(bitmap []byte, index int) bool {
	cbyte := index / 8
	cbit := 7 - (index % 8)

	status := bitmap[cbyte]&(0x1<<uint(cbit)) != 0
	return status
}
