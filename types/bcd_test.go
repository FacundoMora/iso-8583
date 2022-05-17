package types_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/iso-8583/types"

	"github.com/stretchr/testify/assert"
)

func Test_Bcd_Serialize_Fixed_Size(t *testing.T) {
	ser := types.Bcd{NumDigits: 10}

	buffer, err := ser.Serialize("1D34567890")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x1D, 0x34, 0x56, 0x78, 0x90}, buffer.Bytes())

	buffer, err = ser.Serialize("123456789")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x23, 0x45, 0x67, 0x89}, buffer.Bytes())

	buffer, err = ser.Serialize("56789")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x00, 0x05, 0x67, 0x89}, buffer.Bytes())

}

func Test_Bcd_Serialize_Var_Size(t *testing.T) {
	ser := types.Bcd{}

	buffer, err := ser.Serialize("1234567890")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78, 0x90}, buffer.Bytes())

	buffer, err = ser.Serialize("123456789")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x23, 0x45, 0x67, 0x89}, buffer.Bytes())

	buffer, err = ser.Serialize("56789")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x05, 0x67, 0x89}, buffer.Bytes())

	buffer, err = ser.Serialize("456789")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x45, 0x67, 0x89}, buffer.Bytes())

	buffer, err = ser.Serialize("0071")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x71}, buffer.Bytes())
}

func Test_Bcd_Serialize_Errors(t *testing.T) {
	ser := types.Bcd{}

	_, err := ser.Serialize(332131)
	assert.Error(t, err)

	_, err = ser.Serialize("12345ABC6789")
	assert.Error(t, err)

	ser = types.Bcd{NumDigits: 10}
	_, err = ser.Serialize("3423423432432423432423")
	assert.Error(t, err)
}

func Test_Bcd_Deserialize_Fixed_Size(t *testing.T) {
	des := types.Bcd{NumDigits: 10}

	value, err := des.Deserialize(bytes.NewBuffer([]byte{0x12, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "1234567890", value)

	value, err = des.Deserialize(bytes.NewBuffer([]byte{0x02, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "234567890", value)

	value, err = des.Deserialize(bytes.NewBuffer([]byte{0x00, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "34567890", value)
}

func Test_Bcd_Deserialize_Fixed_Size_NotPadded(t *testing.T) {
	des := types.Bcd{NumDigits: 10, NotPadded: true}

	// Value with half byte padded
	value, err := des.Deserialize(bytes.NewBuffer([]byte{0x02, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "0234567890", value)

	// Value with full byte padding
	value, err = des.Deserialize(bytes.NewBuffer([]byte{0x00, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "0034567890", value)
}

func Test_Bcd_Deserialize_Fixed_Size_Odd_NotPadded(t *testing.T) {
	des := types.Bcd{NumDigits: 9, NotPadded: true}

	// Value with half byte padded
	value, err := des.Deserialize(bytes.NewBuffer([]byte{0x02, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "234567890", value)

	// Value with full byte padding
	value, err = des.Deserialize(bytes.NewBuffer([]byte{0x00, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "034567890", value)
}

func Test_Bcd_Deserialize_Var_Size_NotPadded(t *testing.T) {
	des := types.Bcd{NotPadded: true}

	// Value with half byte padded
	value, err := des.Deserialize(bytes.NewBuffer([]byte{0x02, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "0234567890", value)

	// Value with full byte padding
	value, err = des.Deserialize(bytes.NewBuffer([]byte{0x00, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "0034567890", value)
}

func Test_Bcd_Deserialize_Var_Size(t *testing.T) {
	des := types.Bcd{}

	value, err := des.Deserialize(bytes.NewBuffer([]byte{0x12, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "1234567890", value)

	value, err = des.Deserialize(bytes.NewBuffer([]byte{0x02, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "234567890", value)

	value, err = des.Deserialize(bytes.NewBuffer([]byte{0x00, 0x34, 0x56, 0x78, 0x90}))
	assert.NoError(t, err)
	assert.Equal(t, "34567890", value)
}

func Test_Bcd_Deserialize_Error(t *testing.T) {
	des := types.Bcd{NumDigits: 10}

	_, err := des.Deserialize(bytes.NewBuffer([]byte{0x56, 0x78, 0x90}))
	assert.Error(t, err)
}
