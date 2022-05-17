package types

import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/mercadolibre/go-iso8583/serdes"
)

type Bcd struct {
	Desc
	NumDigits int
	NotPadded bool
}

func (bcd Bcd) Name() string {
	return "bcd"
}

func (bcd Bcd) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	valueStr, err := bcd.normalizeValue(value)
	if err != nil {
		return nil, err
	}

	numDigits := len(valueStr)
	if numDigits, err = bcd.normalizeNumDigits(value, numDigits); err != nil {
		return nil, err
	}

	if numDigits == 0 {
		return &bytes.Buffer{}, nil
	}

	raw := bcd.convertNumToBcd(valueStr, numDigits)
	return bytes.NewBuffer(raw), nil
}

func (bcd Bcd) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	numDigits := bcd.NumDigits
	if numDigits == 0 {
		numDigits = data.Len() * 2
	}

	odd := false
	if numDigits > 0 && numDigits%2 != 0 {
		numDigits++
		odd = true
	}

	numBytes := numDigits / 2
	if data.Len() < numBytes {
		return nil, DeserializationError{
			Message: "data does not has bytes enough", Serdes: bcd, Remaning: data.Len(),
		}
	}

	raw := data.Next(numBytes)
	value := hex.EncodeToString(raw)

	if !bcd.NotPadded {
		value = strings.TrimLeft(value, "0")
	} else if odd {
		// removes half byte from value
		value = value[1:]
	}

	return value, nil
}

func (bcd Bcd) convertNumToBcd(valueStr string, numDigits int) []byte {
	numBytes := numDigits / 2
	valueLen := len(valueStr)

	data := make([]byte, numBytes)
	for indexDigit := 0; indexDigit < numDigits; indexDigit++ {
		var digit int
		valueIndex := indexDigit - (numDigits - valueLen)
		if valueIndex >= 0 {
			if valueStr[valueIndex] == 'D' {
				digit = 0xD
			} else {
				// Gets the digit in decimal from char, example: '9' (char) -> 9 (int)
				digit = int(valueStr[valueIndex]) - 0x30
			}
		}

		// Calculate where and put the digit on current byte
		// Example:
		//  - For digit index 0: byte in = 0x00, digit value = 8, byte out = 0x80
		//  - For digit index 1: byte in = 0x80, digit value = 4, byte out = 0x84

		nibbleOffset := 4 * ((indexDigit + 1) % 2)
		data[indexDigit/2] |= byte(digit << nibbleOffset)
	}
	return data
}

func (bcd Bcd) normalizeValue(value serdes.Value) (string, error) {
	valueStr, ok := value.(string)
	if !ok {
		return "", SerializerError{
			Message: "invalid value type", Serdes: bcd, Value: value,
		}
	}

	for _, c := range valueStr {
		if c == 'D' {
			continue
		}
		if c < '0' || c > '9' {
			return "", SerializerError{
				Message: "for bcd type, just numbers are allowed in string", Serdes: bcd, Value: value,
			}
		}
	}

	return valueStr, nil
}

func (bcd Bcd) normalizeNumDigits(value serdes.Value, numDigits int) (int, error) {
	if bcd.NumDigits > 0 && numDigits > bcd.NumDigits {
		return 0, SerializerError{
			Message: "value too long", Serdes: bcd, Value: value,
		}
	}

	if bcd.NumDigits > 0 {
		numDigits = bcd.NumDigits
	}

	if numDigits%2 != 0 {
		numDigits++
	}

	return numDigits, nil
}
