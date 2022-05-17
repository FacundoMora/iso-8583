package types_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/iso-8583/types"

	"github.com/stretchr/testify/assert"
)

func Test_Ebcdic_Serialize_Fixed_Size(t *testing.T) {
	ser := types.Ebcdic{NumDigits: 40}

	buffer, err := ser.Serialize("deserialized_string")
	assert.NoError(t, err)

	expected := []byte{0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95,
		0x87, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
		0x40, 0x40, 0x40, 0x40,
	}

	assert.Equal(t, expected, buffer.Bytes())
}

func Test_Ebcdic_Serialize_Var_Size(t *testing.T) {
	ser := types.Ebcdic{}

	buffer, err := ser.Serialize("deserialized_string")
	assert.NoError(t, err)

	expected := []byte{0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95, 0x87}
	assert.Equal(t, expected, buffer.Bytes())

}

func Test_Ebcdic_Serialize_Errors(t *testing.T) {
	ser := types.Ebcdic{}

	_, err := ser.Serialize(332131)
	assert.Error(t, err)

	ser = types.Ebcdic{NumDigits: 10}
	_, err = ser.Serialize("3423423432432423432423")
	assert.Error(t, err)
}

func Test_Ebcdic_Deserialize_Fixed_Size(t *testing.T) {
	des := types.Ebcdic{NumDigits: 40}

	dataIn := []byte{0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95,
		0x87, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
		0x40, 0x40, 0x40, 0x40,
	}

	value, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.NoError(t, err)
	assert.Equal(t, "deserialized_string", value)
}

func Test_Ebcdic_Deserialize_Var_Size(t *testing.T) {
	des := types.Ebcdic{}

	dataIn := []byte{0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95,
		0x87, 0x40, 0x40,
	}

	value, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.NoError(t, err)
	assert.Equal(t, "deserialized_string", value)
}

func Test_Ebcdic_Deserialize_Error(t *testing.T) {
	des := types.Bcd{NumDigits: 10}

	dataIn := []byte{0x84, 0x85, 0xa2, 0x85}
	_, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.Error(t, err)
}
