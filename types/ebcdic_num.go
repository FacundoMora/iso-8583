package types

import (
	"bytes"
	"fmt"

	"github.com/mercadolibre/iso-8583/serdes"
)

type EbcdicNumeric struct {
	Desc
	NumDigits int
}

func (ebcdic EbcdicNumeric) Name() string {
	return "ebcdic"
}

func (ebcdic EbcdicNumeric) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	valueStr, ok := value.(string)
	if !ok {
		return nil, SerializerError{
			Message: "invalid value type", Serdes: ebcdic, Value: value,
		}
	}

	valueLen := len(valueStr)
	numDigits := ebcdic.NumDigits
	if numDigits == 0 {
		numDigits = valueLen
	}

	if ebcdic.NumDigits > 0 && valueLen > numDigits {
		return nil, SerializerError{
			Message: "value too long", Serdes: ebcdic, Value: value,
		}
	}

	paddedValue := fmt.Sprintf("%0*s", numDigits, valueStr)
	data := new(bytes.Buffer)
	for _, c := range paddedValue {
		data.WriteByte(asciiToEbcdic[c])
	}

	return data, nil
}

func (ebcdic EbcdicNumeric) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	numDigits := ebcdic.NumDigits
	if numDigits == 0 {
		numDigits = data.Len()
	}

	if data.Len() < numDigits {
		return nil, DeserializationError{
			Message: "data does not has bytes enough", Serdes: ebcdic, Remaning: data.Len(),
		}
	}

	var out string
	for count := 0; count < numDigits; count++ {
		c, err := data.ReadByte()
		if err != nil {
			return nil, DeserializationError{
				Message: "error reading bytes", Serdes: ebcdic, Remaning: data.Len(),
			}
		}

		out += string(ebcdicToASCII[c])
	}

	return out, nil
}
