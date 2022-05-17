package serdes_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/go-iso8583/serdes"

	"github.com/stretchr/testify/assert"
)

func Test_Serializer_Serialize(t *testing.T) {
	ser := serdes.Mock{}

	fieldsIn := map[string]interface{}{
		"2":  "2312312322332",
		"11": "122344",
	}

	expectedData := new(bytes.Buffer)
	expectedData.Write([]byte("0123456789"))

	ser.On("Serialize", fieldsIn).Return(expectedData, nil)
	data, err := ser.Serialize(fieldsIn)

	assert.NoError(t, err)
	assert.Equal(t, expectedData.Bytes(), data.Bytes())
}

func Test_Deserializer_Deserialize(t *testing.T) {
	des := serdes.Mock{}

	expectedFields := map[string]interface{}{
		"2":  "2312312322332",
		"11": "122344",
	}

	data := bytes.NewBuffer([]byte("0123456789"))
	des.On("Deserialize", data).Return(expectedFields, nil)

	fieldsOut, err := des.Deserialize(data)

	assert.NoError(t, err)
	assert.Equal(t, expectedFields, fieldsOut)
}
