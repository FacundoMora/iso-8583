package types_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/go-iso8583/types"

	"github.com/stretchr/testify/assert"
)

func Test_Raw_Serialize_Fixed_Size(t *testing.T) {
	ser := types.Raw{NumBytes: 20}

	buffer, err := ser.Serialize("8485a2859989819389a985846da2a39989958734")
	assert.NoError(t, err)

	expected := []byte{
		0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95,
		0x87, 0x34,
	}

	assert.Equal(t, expected, buffer.Bytes())
}

func Test_Raw_Serialize_Var_Size(t *testing.T) {
	ser := types.Raw{}

	buffer, err := ser.Serialize("8485a2859989819389a985846da2a39989958734")
	assert.NoError(t, err)

	expected := []byte{
		0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95,
		0x87, 0x34,
	}

	assert.Equal(t, expected, buffer.Bytes())

}

func Test_Raw_Serialize_Errors(t *testing.T) {
	ser := types.Raw{}

	_, err := ser.Serialize(332131)
	assert.Error(t, err)

	ser = types.Raw{NumBytes: 3}
	_, err = ser.Serialize("3423423432432423432423")
	assert.NoError(t, err)

	ser = types.Raw{NumBytes: 3}
	_, err = ser.Serialize("34234")

	assert.Error(t, err)
}

func Test_Raw_Deserialize_Fixed_Size(t *testing.T) {
	des := types.Raw{NumBytes: 20}

	dataIn := []byte{
		0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95,
		0x87, 0x34,
	}

	value, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.NoError(t, err)
	assert.Equal(t, "8485a2859989819389a985846da2a39989958734", value)
}

func Test_Raw_Deserialize_Var_Size(t *testing.T) {
	des := types.Raw{}

	dataIn := []byte{
		0x84, 0x85, 0xa2, 0x85, 0x99, 0x89, 0x81, 0x93, 0x89, 0xa9, 0x85, 0x84, 0x6d, 0xa2, 0xa3, 0x99, 0x89, 0x95,
		0x87, 0x34,
	}

	value, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.NoError(t, err)
	assert.Equal(t, "8485a2859989819389a985846da2a39989958734", value)
}

func Test_Raw_Deserialize_Error(t *testing.T) {
	des := types.Raw{NumBytes: 10}

	dataIn := []byte{0x84, 0x85, 0xa2, 0x85}
	_, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.Error(t, err)
}
