package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/mercadolibre/go-iso8583/serdes"
)

var (
	ErrIndefiniteLength = errors.New("indefinite length is not supported")
	ErrInvalidLength    = errors.New("invalid length")
	ErrTagNotFound      = errors.New("tag not found")
)

// TagValue pair.
type TagValue struct {
	Tag     int
	Value   []byte
	SizeLen int // size of length in number of bytes.
}

type BerTLV struct {
	Desc
	SizeLen int // size of length in number of bytes.
	Items   []Field
}

func (t BerTLV) Name() string {
	return "bertlv"
}

func (t BerTLV) Serialize(data serdes.Value) (*bytes.Buffer, error) {
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

		tagRaw, err := Raw{}.Serialize(field.Name)
		if err != nil {
			return nil, SerializerError{
				Message: "tag serializer failed", Serdes: t, Field: field, Cause: err,
			}
		}

		tag, _, err := readTag(bytes.NewReader(tagRaw.Bytes()))
		if err != nil {
			return nil, SerializerError{
				Message: "invalid tag", Serdes: t, Field: field, Value: tagRaw, Cause: err,
			}
		}

		builtTlv := TagValue{
			Tag:     tag,
			SizeLen: t.SizeLen,
			Value:   data.Bytes(),
		}

		if _, err := serializedData.ReadFrom(builtTlv.bufferBytes()); err != nil {
			return nil, SerializerError{
				Message: "data struct serializer failed", Value: data, Serdes: t, Cause: err,
			}
		}
	}

	return serializedData, nil
}

func (t BerTLV) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	listValues := serdes.Map{}

	mapTLV, err := decode(t.SizeLen, data.Bytes())
	if err != nil {
		return nil, DeserializationError{
			Message: "data struct deserializer failed", Serdes: t, Remaning: data.Len(), Cause: err,
		}
	}

	for _, tlv := range mapTLV {
		tag, err := Raw{}.Deserialize(bytes.NewBuffer(encodeInt(tlv.Tag)))
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
			continue
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

func (t BerTLV) findField(tag string) (Field, error) {
	for _, field := range t.Items {
		if field.Name == tag {
			return field, nil
		}
	}
	return Field{}, fmt.Errorf("field %s not found", tag)
}

// Bytes encodes TagValue into a byte buffer.
func (tv TagValue) bufferBytes() *bytes.Buffer {
	return bytes.NewBuffer(tv.bytes())
}

// bytes encodes TagValue into a byte slice.
func (tv TagValue) bytes() []byte {
	result := append(tv.tag(), tv.len()...)
	return append(result, tv.Value...)
}

// tag returns encoded tag value (two bytes tags are supported).
func (tv TagValue) tag() []byte {
	if (tv.Tag>>8)&0x1F == 0 {
		return []byte{byte(tv.Tag)}
	}

	return []byte{byte(tv.Tag >> 8), byte(tv.Tag & 0xff)}
}

// len returns encoded length of the value.
func (tv TagValue) len() []byte {
	l := len(tv.Value)

	// build size with fixed length
	if tv.SizeLen > 0 {
		r := encodeInt(l)
		result := make([]byte, tv.SizeLen-len(r))
		return append(result, r...)
	}

	// the first byte is a final byte?
	if l <= 0x7f {
		return []byte{byte(l)}
	}

	r := encodeInt(l)
	numOctets := len(r)
	result := make([]byte, 1+numOctets)
	result[0] = 0x80 | byte(numOctets)

	copy(result[1:], r)

	return result
}

// readFrom implements io.ReaderFrom.
func (tv *TagValue) readFrom(r io.Reader) (n int64, err error) {
	tag, tagn, err := readTag(r)
	if err != nil {
		return int64(tagn), err
	}

	l, ln, err := readLen(tv.SizeLen, r)
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

// isConstructed returns true if the value is constructed type(contains other TLV records).
func (tv TagValue) isConstructed() bool {
	if tv.Tag <= 0xff {
		return tv.Tag&0x20 != 0
	}

	return (tv.Tag>>8)&0x20 != 0
}

// readTag reads length of a tag and return it along with length of a length in bytes or an error.
func readTag(r io.Reader) (tag int, n int, err error) {
	b := make([]byte, 1)

	// reading first byte of the tag
	n, err = r.Read(b)
	if err != nil {
		return 0, n, err
	}

	tag = int(b[0])

	// it's a two byte tag
	if b[0]&0x1F == 0x1F {
		tag <<= 8

		n, err = r.Read(b)
		if err != nil {
			return 0, n + 1, err
		}

		n = 2
		tag |= int(b[0])
	}

	return tag, n, nil
}

// readLen reads length of a tag and return it along with length of a length in bytes and/or an error.
func readLen(sizeTam int, r io.Reader) (length int, n int, err error) {
	// when size of len is fixed
	if sizeTam > 0 {
		b := make([]byte, sizeTam)

		nb, err := r.Read(b)
		if err != nil {
			return 0, n, err
		}

		lenb := append(make([]byte, 4-nb), b...)
		return int(binary.BigEndian.Uint32(lenb)), sizeTam, nil
	}

	b := make([]byte, 1)

	n, err = r.Read(b)
	if err != nil {
		return 0, n, err
	}

	if b[0] == 0x80 {
		return 0, 1, ErrIndefiniteLength
	}

	if b[0]&0x80 == 0 {
		return int(b[0]), 1, nil
	}

	nb := int(b[0] & 0x7f)
	if nb > 4 {
		return 0, 1, ErrInvalidLength
	}

	lenb := make([]byte, 4)
	n, err = r.Read(lenb[4-nb:])
	if err != nil {
		return 0, n + 1, err
	}

	return int(binary.BigEndian.Uint32(lenb)), n + 1, nil
}

// encodeInt encodes an integer to BER format.
func encodeInt(in int) []byte {
	result := make([]byte, 4)

	binary.BigEndian.PutUint32(result, uint32(in))

	var lz int
	for ; lz < 4; lz++ {
		if result[lz] != 0 {
			break
		}
	}

	return result[lz:]
}

// decode decodes TLV encoded byte slice into slice of TagValue structs.
func decode(sizeTam int, p []byte) ([]TagValue, error) {
	r := bytes.NewReader(p)

	tv := TagValue{SizeLen: sizeTam}
	var result []TagValue
	for {
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

// Find finds first tag (DFS) in the TLV structure represented by p.
func Find(sizeTam int, tag int, p []byte) (*TagValue, error) {
	tvs, err := decode(sizeTam, p)
	if err != nil {
		return nil, err
	}

	for _, tv := range tvs {
		if tv.Tag == tag {
			return &tv, nil
		}

		if !tv.isConstructed() {
			continue
		}

		tvv, err := Find(sizeTam, tag, tv.Value)
		if err != nil {
			return nil, err
		}

		if tvv != nil {
			return tvv, nil
		}
	}

	return nil, ErrTagNotFound
}
