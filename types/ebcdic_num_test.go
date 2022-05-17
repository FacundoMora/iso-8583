package types_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/iso-8583/types"

	"github.com/stretchr/testify/assert"
)

func Test_EbcdicNumeric_Serialize_Fixed_Size(t *testing.T) {
	ser := types.EbcdicNumeric{NumDigits: 10}

	buffer, err := ser.Serialize("123456")
	assert.NoError(t, err)

	expected := []byte{0xf0, 0xf0, 0xf0, 0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6}

	assert.Equal(t, expected, buffer.Bytes())
}

func Test_EbcdicNumeric_Serialize_Var_Size(t *testing.T) {
	ser := types.EbcdicNumeric{}

	buffer, err := ser.Serialize("123456")
	assert.NoError(t, err)

	expected := []byte{0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6}

	assert.Equal(t, expected, buffer.Bytes())
}

func Test_EbcdicNumeric_Serialize_Errors(t *testing.T) {
	ser := types.EbcdicNumeric{}

	_, err := ser.Serialize(332131)
	assert.Error(t, err)

	ser = types.EbcdicNumeric{NumDigits: 10}
	_, err = ser.Serialize("3423423432432423432423")
	assert.Error(t, err)
}

func Test_EbcdicNumeric_Deserialize_Fixed_Size(t *testing.T) {
	des := types.EbcdicNumeric{NumDigits: 10}

	dataIn := []byte{0xf0, 0xf0, 0xf0, 0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6}

	value, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.NoError(t, err)
	assert.Equal(t, "0000123456", value)
}

func Test_EbcdicNumeric_Deserialize_Var_Size(t *testing.T) {
	des := types.EbcdicNumeric{}

	dataIn := []byte{0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6}

	value, err := des.Deserialize(bytes.NewBuffer(dataIn))
	assert.NoError(t, err)
	assert.Equal(t, "123456", value)
}
