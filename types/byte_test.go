package types_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/go-iso8583/types"

	"github.com/stretchr/testify/assert"
)

func Test_Byte_Deserialize_Success(t *testing.T) {
	definition := types.Byte{}
	data := bytes.NewBuffer([]byte{0xFA})

	value, err := definition.Deserialize(data)
	assert.NoError(t, err)
	assert.Equal(t, "250", value)
}

func Test_Byte_Deserialize_Error(t *testing.T) {
	definition := types.Byte{}
	data := bytes.NewBuffer([]byte{})

	_, err := definition.Deserialize(data)
	assert.Error(t, err)
}

func Test_Byte_Serialize_Success(t *testing.T) {
	definition := types.Byte{}

	value := "128"
	data, err := definition.Serialize(value)
	assert.NoError(t, err)

	expected := bytes.NewBuffer([]byte{0x80})
	assert.Equal(t, expected, data)
}

func Test_Byte_Serialize_Error(t *testing.T) {
	definition := types.Byte{}

	_, err := definition.Serialize(0)
	assert.Error(t, err)

	_, err = definition.Serialize("sasS")
	assert.Error(t, err)
}
