package types_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/mercadolibre/go-iso8583/types"

	"github.com/stretchr/testify/assert"
)

func Test_Word_Deserialize_Success(t *testing.T) {
	definition := types.Word{Order: binary.BigEndian}
	data := bytes.NewBuffer([]byte{0x33, 0x55})

	value, err := definition.Deserialize(data)
	assert.NoError(t, err)
	assert.Equal(t, "13141", value)
}

func Test_Word_Deserialize_Error(t *testing.T) {
	definition := types.Word{Order: binary.BigEndian}
	data := bytes.NewBuffer([]byte{})

	_, err := definition.Deserialize(data)
	assert.Error(t, err)
}

func Test_Word_Serialize_Success(t *testing.T) {
	definition := types.Word{Order: binary.BigEndian}

	value := "1082"
	data, err := definition.Serialize(value)
	assert.NoError(t, err)

	expected := bytes.NewBuffer([]byte{0x04, 0x3A})
	assert.Equal(t, expected, data)
}

func Test_Word_Serialize_Error(t *testing.T) {
	definition := types.Word{Order: binary.BigEndian}

	_, err := definition.Serialize(0)
	assert.Error(t, err)

	_, err = definition.Serialize("sasS")
	assert.Error(t, err)
}
