package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mercadolibre/go-iso8583/serdes"
)

type Ebcdic struct {
	Desc
	NumDigits int
}

var ebcdicToASCII = []byte{
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 91, 46, 60, 40, 43, 33,
	38, 32, 32, 32, 32, 32, 32, 32, 32, 32, 93, 36, 42, 41, 59, 94,
	45, 47, 32, 32, 32, 32, 32, 32, 32, 32, 124, 44, 37, 95, 62, 63,
	32, 32, 32, 32, 32, 32, 238, 160, 161, 96, 58, 35, 64, 39, 61, 34,
	230, 97, 98, 99, 100, 101, 102, 103, 104, 105, 164, 165, 228, 163, 229, 168,
	169, 106, 107, 108, 109, 110, 111, 112, 113, 114, 170, 171, 172, 173, 174, 175,
	239, 126, 115, 116, 117, 118, 119, 120, 121, 122, 224, 225, 226, 227, 166, 162,
	236, 235, 167, 232, 237, 233, 231, 234, 158, 128, 129, 150, 132, 133, 148, 131,
	123, 65, 66, 67, 68, 69, 70, 71, 72, 73, 149, 136, 137, 138, 139, 140,
	125, 74, 75, 76, 77, 78, 79, 80, 81, 82, 141, 142, 143, 159, 144, 145,
	92, 32, 83, 84, 85, 86, 87, 88, 89, 90, 146, 147, 134, 130, 156, 155,
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 135, 152, 157, 153, 151, 32,
}

// asciiToEbcdic is a table to convert EBCDIC encoding to ascii encoding.
var asciiToEbcdic = []byte{
	0, 1, 2, 3, 55, 45, 46, 47, 22, 5, 37, 11, 12, 13, 14, 15,
	16, 17, 18, 19, 60, 61, 50, 38, 24, 25, 63, 39, 28, 29, 30, 31,
	64, 79, 127, 123, 91, 108, 80, 125, 77, 93, 92, 78, 107, 96, 75, 97,
	240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 122, 94, 76, 126, 110, 111,
	124, 193, 194, 195, 196, 197, 198, 199, 200, 201, 209, 210, 211, 212, 213, 214,
	215, 216, 217, 226, 227, 228, 229, 230, 231, 232, 233, 74, 224, 90, 95, 109,
	121, 129, 130, 131, 132, 133, 134, 135, 136, 137, 145, 146, 147, 148, 149, 150,
	151, 152, 153, 162, 163, 164, 165, 166, 167, 168, 169, 192, 106, 208, 161, 7,
	32, 33, 34, 35, 36, 21, 6, 23, 40, 41, 42, 43, 44, 9, 10, 27,
	48, 49, 26, 51, 52, 53, 54, 8, 56, 57, 58, 59, 4, 20, 62, 225,
	65, 66, 67, 68, 69, 70, 71, 72, 73, 81, 82, 83, 84, 85, 86, 87,
	88, 89, 98, 99, 100, 101, 102, 103, 104, 105, 112, 113, 114, 115, 116, 117,
	118, 119, 120, 128, 138, 139, 140, 141, 142, 143, 144, 154, 155, 156, 157, 158,
	159, 160, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183,
	184, 185, 186, 187, 188, 189, 190, 191, 202, 203, 204, 205, 206, 207, 218, 219,
	220, 221, 222, 223, 234, 235, 236, 237, 238, 239, 250, 251, 252, 253, 254, 255,
}

func (ebcdic Ebcdic) Name() string {
	return "ebcdic"
}

func (ebcdic Ebcdic) Serialize(value serdes.Value) (*bytes.Buffer, error) {
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

	paddedValue := fmt.Sprintf("%-*s", numDigits, valueStr)
	data := new(bytes.Buffer)
	for _, c := range paddedValue {
		data.WriteByte(asciiToEbcdic[c])
	}

	return data, nil
}

func (ebcdic Ebcdic) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
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

	out = strings.TrimRight(out, " ")
	return out, nil
}
