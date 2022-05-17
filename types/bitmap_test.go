package types_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/go-iso8583/types"

	"github.com/stretchr/testify/assert"
)

func Test_Bitmap_Serialize_Blocks(t *testing.T) {
	ser := types.Bitmap{BlockSize: 64, NumBits: 128}

	value64 := []byte{0x75, 0x65, 0x12, 0x76, 0xF5, 0x2A, 0x43, 0x46}
	data, err := ser.Serialize(value64)
	assert.NoError(t, err)
	assert.Equal(t, data.Bytes(), value64)

	value128 := []byte{0x75, 0x65, 0x12, 0x76, 0xF5, 0x2A, 0x43, 0x46, 0x75, 0x25, 0x82, 0x76, 0x55, 0x2A, 0xA3, 0x4F}
	data, err = ser.Serialize(value128)
	assert.NoError(t, err)
	expected := []byte{0xF5, 0x65, 0x12, 0x76, 0xF5, 0x2A, 0x43, 0x46, 0x75, 0x25, 0x82, 0x76, 0x55, 0x2A, 0xA3, 0x4F}
	assert.Equal(t, expected, data.Bytes())

	_, err = ser.Serialize("invalid")
	assert.Error(t, err)
}

func Test_Bitmap_Serialize_Errors(t *testing.T) {
	ser := types.Bitmap{BlockSize: 64, NumBits: 128}

	_, err := ser.Serialize("invalid")
	assert.Error(t, err)
}

func Test_Bitmap_Deserialize_Blocks(t *testing.T) {
	des := types.Bitmap{BlockSize: 64, NumBits: 128}

	value64 := []byte{0x75, 0x65, 0x12, 0x76, 0xF5, 0x2A, 0x43, 0x46}
	data := bytes.NewBuffer(value64)
	value, err := des.Deserialize(data)
	assert.NoError(t, err)
	assert.Equal(t, value64, value)

}

func Test_Bitmap_Deserialize_Error(t *testing.T) {
	des := types.Bitmap{BlockSize: 64, NumBits: 128}

	value64 := []byte{0x75, 0x65, 0x12, 0x76, 0xF5, 0x2A}
	data := bytes.NewBuffer(value64)
	_, err := des.Deserialize(data)
	assert.Error(t, err)
}
