package types

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mercadolibre/go-iso8583/serdes"
)

// Mastercard Subelement Encoding Scheme type

// TagValueMas pair
type TagValueMas struct {
	Tag     []byte
	Value   []byte
	SizeLen int // size of length
	SizeTag int // size of tag
}

type TLV struct {
	Desc
	SizeLen int // size of length
	SizeTag int // size of tag
	Items   []Field
}

const (
	_defaultSizeLen = 2
	_defaultSizeTag = 2
)

func (t TLV) Name() string {
	return "tlv"
}

func (t TLV) Serialize(data serdes.Value) (*bytes.Buffer, error) {
	mapValue, ok := data.(serdes.Map)
	if !ok {
		return nil, SerializerError{
			Message: fmt.Sprintf("invalid value [%T], expected: %T", data, serdes.Map{}), Value: data, Serdes: t,
		}
	}

	serializedData := new(bytes.Buffer)
	for _, field := range t.Items {
		if field.Name == "" {
			return nil, SerializerError{
				Message: "field name not found", Serdes: t, Field: field,
			}
		}

		itemValue, ok := mapValue[field.Name]
		if !ok {
			continue
		}

		data, err := field.SerDes.Serialize(itemValue)
		if err != nil {
			return nil, SerializerError{
				Message: "value serializer failed", Serdes: t, Field: field, Cause: err,
			}
		}

		sizeTag := t.SizeTag
		if sizeTag == 0 {
			sizeTag = _defaultSizeTag
		}
		tag, err := EbcdicNumeric{NumDigits: sizeTag}.Serialize(field.Name)
		if err != nil {
			return nil, SerializerError{
				Message: "tag serializer failed", Serdes: t, Field: field, Cause: err,
			}
		}

		sizeLen := t.SizeLen
		if sizeLen == 0 {
			sizeLen = _defaultSizeLen
		}

		builtTlv := TagValueMas{
			Tag:     tag.Bytes(),
			Value:   data.Bytes(),
			SizeLen: sizeLen,
		}

		// check capacity
		if data.Len() >= intPow(10, sizeLen) {
			return nil, SerializerError{
				Message: "length serializer failed", Value: data, Serdes: t,
			}
		}

		if _, err := serializedData.ReadFrom(builtTlv.bufferBytes()); err != nil {
			return nil, SerializerError{
				Message: "data struct serializer failed", Value: data, Serdes: t, Cause: err,
			}
		}
	}

	return serializedData, nil
}

func (t TLV) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	listValues := serdes.Map{}

	sizeTag := t.SizeTag
	if sizeTag == 0 {
		sizeTag = _defaultSizeTag
	}

	sizeLen := t.SizeTag
	if sizeLen == 0 {
		sizeLen = _defaultSizeLen
	}

	mapTLV, err := decodeMas(sizeTag, sizeLen, data.Bytes())
	if err != nil {
		return nil, DeserializationError{
			Message: "data struct deserializer failed", Serdes: t, Remaning: data.Len(), Cause: err,
		}
	}

	for _, tlv := range mapTLV {
		tag, err := EbcdicNumeric{}.Deserialize(bytes.NewBuffer(tlv.Tag))
		if err != nil {
			return nil, DeserializationError{
				Message: "tag deserializer failed", Serdes: t, Cause: err,
			}
		}

		tagValue, ok := tag.(string)
		if !ok {
			return nil, DeserializationError{
				Message: "tag type is not string", Serdes: t,
			}
		}

		field, err := t.findField(tagValue)
		if err != nil {
			//continue
			// default serializer
			field = Field{Name: tagValue, SerDes: Ebcdic{}}
		}

		value, err := field.SerDes.Deserialize(bytes.NewBuffer(tlv.Value))
		if err != nil {
			return nil, DeserializationError{
				Message: "struct data deserializer failed", Serdes: t, Field: field, Cause: err,
			}
		}

		listValues[tagValue] = value
	}

	return listValues, nil
}

func (t TLV) findField(tag string) (Field, error) {
	for _, field := range t.Items {
		if field.Name == tag {
			return field, nil
		}
	}
	return Field{}, fmt.Errorf("field %s not found", tag)
}

// Bytes encodes TagValue into a byte buffer.
func (tv TagValueMas) bufferBytes() *bytes.Buffer {
	return bytes.NewBuffer(tv.bytes())
}

// bytes encodes TagValue into a byte slice.
func (tv TagValueMas) bytes() []byte {
	result := append(tv.tag(), tv.len()...)
	return append(result, tv.Value...)
}

// tag returns encoded tag value (two bytes tags are supported).
func (tv TagValueMas) tag() []byte {
	return tv.Tag
}

// len returns encoded length of the value.
func (tv TagValueMas) len() []byte {
	l := len(tv.Value)
	size := tv.SizeLen
	b := make([]byte, size)

	for i := 0; i < size-1; i++ {
		y := size - 1 - i
		exp := intPow(10, y)
		b[i] = byte(l / exp)
		l = l % exp
	}
	b[size-1] = byte(l)

	for i := range b {
		b[i] |= 0xf0
	}

	return b
}

func (tv *TagValueMas) readFrom(r io.Reader) (n int64, err error) {
	tag, tagn, err := readTagMas(tv.SizeTag, r)
	if err != nil {
		return int64(tagn), err
	}

	l, ln, err := readLenMas(tv.SizeLen, r)
	if err != nil {
		return int64(tagn) + int64(ln), fmt.Errorf("failed to read length: %v", err)
	}

	tv.Tag = tag

	if l == 0 {
		return int64(tagn) + int64(ln), nil
	}

	tv.Value = make([]byte, l)
	vn, err := r.Read(tv.Value)
	if vn < l {
		return int64(tagn) + int64(ln) + int64(vn), io.ErrUnexpectedEOF
	}

	if err != nil {
		return int64(tagn) + int64(ln) + int64(vn), fmt.Errorf("failed to read value: %v", err)
	}

	return int64(tagn) + int64(ln) + int64(vn), nil
}

// readTag reads tag and return it along with length of a length in bytes or an error.
func readTagMas(size int, r io.Reader) (tag []byte, n int, err error) {
	tag = make([]byte, size)

	n, err = r.Read(tag)
	if err != nil {
		return nil, n, err
	}

	return tag, n, nil
}

// readLen reads length of a tag and return it along with length of a length in bytes and/or an error.
func readLenMas(size int, r io.Reader) (length int, n int, err error) {
	b := make([]byte, size)

	n, err = r.Read(b)
	if err != nil {
		return 0, n, err
	}

	for i := range b {
		b[i] &= 0x0f
	}

	l := int(b[size-1])
	for i := 0; i < size-1; i++ {
		y := size - 1 - i
		exp := intPow(10, y)
		l += int(b[i]) * exp
	}

	return l, n, nil
}

// decode decodes TLV encoded byte slice into slice of TagValueMas structs.
func decodeMas(sizeTag int, sizeLen int, p []byte) ([]TagValueMas, error) {
	r := bytes.NewReader(p)

	var result []TagValueMas
	for {
		tv := TagValueMas{SizeTag: sizeTag, SizeLen: sizeLen}
		_, err := tv.readFrom(r)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		result = append(result, tv)
	}

	return result, nil
}

// intPow calculates x to the yth power
func intPow(x, y int) int {
	if y == 0 {
		return 1
	}
	result := x
	for i := 2; i <= y; i++ {
		result *= x
	}
	return result
}
