package types_test

import (
	"bytes"
	"errors"
	"strconv"
	"testing"

	"github.com/mercadolibre/go-iso8583/serdes"
	"github.com/mercadolibre/go-iso8583/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_VarLen_Serialize(t *testing.T) {
	lengthSerializer := &serdes.Mock{}

	deserializedValue := "deserialized_data"

	serializedValue := "serialized_value"
	serializedValueLength := strconv.Itoa(len(serializedValue))
	lengthSerializer.On("Serialize", serializedValueLength).Return(bytes.NewBufferString("*"+serializedValueLength+"*"), nil)

	dataSerializer := &serdes.Mock{}
	dataSerializer.On("Serialize", deserializedValue).Return(bytes.NewBufferString(serializedValue), nil)

	definitions := types.VarLength{
		Length: lengthSerializer, Data: dataSerializer,
	}

	serializedData, err := definitions.Serialize(deserializedValue)
	assert.NoError(t, err)

	expectedSerializedData := "*16*serialized_value"
	assert.Equal(t, []byte(expectedSerializedData), serializedData.Bytes())
}

func Test_VarLen_Serialize_even_Bcd(t *testing.T) {
	lengthSerializer := &serdes.Mock{}

	deserializedValue := "123456"
	serializedValue := []byte{0x12, 0x34, 0x56}

	serializedValueLength := strconv.Itoa(len(serializedValue) * 2)
	lengthSerializer.On("Serialize", serializedValueLength).Return(bytes.NewBufferString("*"+serializedValueLength+"*"), nil)

	definitions := types.VarLength{
		Length: lengthSerializer, Data: types.Bcd{},
	}

	serializedData, err := definitions.Serialize(deserializedValue)
	assert.NoError(t, err)

	expectedSerializedData := append([]byte("*6*"), serializedValue...)
	assert.Equal(t, expectedSerializedData, serializedData.Bytes())
}

func Test_VarLen_Serialize_odd_Bcd(t *testing.T) {
	lengthSerializer := &serdes.Mock{}

	deserializedValue := "23456"
	serializedValue := []byte{0x02, 0x34, 0x56}

	serializedValueLength := strconv.Itoa(len(deserializedValue))
	lengthSerializer.On("Serialize", serializedValueLength).Return(bytes.NewBufferString("*"+serializedValueLength+"*"), nil)

	definitions := types.VarLength{
		Length: lengthSerializer, Data: types.Bcd{},
	}

	serializedData, err := definitions.Serialize(deserializedValue)
	assert.NoError(t, err)

	expectedSerializedData := append([]byte("*5*"), serializedValue...)
	assert.Equal(t, expectedSerializedData, serializedData.Bytes())
}

func Test_VarLen_Serialize_Error(t *testing.T) {
	lengthSerializer := &serdes.Mock{}

	deserializedDataValueError := "deserialized_data_error"
	deserializedDataValueErrorLength := strconv.Itoa(len(deserializedDataValueError))
	lengthSerializer.On("Serialize", deserializedDataValueErrorLength).Return(bytes.NewBufferString(deserializedDataValueErrorLength), nil)

	deserializedLengthValueError := "deserialized_length_error"
	deserializedLengthValueErrorLength := strconv.Itoa(len(deserializedLengthValueError))
	lengthSerializer.On("Serialize", deserializedLengthValueErrorLength).Return(bytes.NewBufferString(deserializedLengthValueErrorLength), errors.New("error length"))

	dataSerializer := &serdes.Mock{}
	dataSerializer.On("Serialize", deserializedDataValueError).Return(bytes.NewBufferString(deserializedDataValueError), errors.New("error data"))
	dataSerializer.On("Serialize", deserializedLengthValueError).Return(bytes.NewBufferString(deserializedLengthValueError), nil)

	definitions := types.VarLength{
		Length: lengthSerializer, Data: dataSerializer,
	}

	_, err := definitions.Serialize(deserializedDataValueError)
	assert.Error(t, err)

	_, err = definitions.Serialize(deserializedLengthValueError)
	assert.Error(t, err)
}

func Test_VarLen_Deserialize(t *testing.T) {
	lengthDeserializer := &serdes.Mock{}
	strSerialized := "*16*serialized_value"
	dataIn := bytes.NewBufferString(strSerialized)

	lengthDeserializer.On("Deserialize", dataIn).Run(func(args mock.Arguments) {
		buffer := args.Get(0).(*bytes.Buffer)
		buffer.Next(4)
	}).Return("16", nil)

	dataDeserializer := &serdes.Mock{}

	dataInputData := bytes.NewBufferString("serialized_value")
	deserializedStr := "deserialized_value"
	dataDeserializer.On("Deserialize", dataInputData).Return(deserializedStr, nil)

	definitions := types.VarLength{
		Length: lengthDeserializer, Data: dataDeserializer,
	}

	value, err := definitions.Deserialize(dataIn)
	assert.NoError(t, err)
	assert.Equal(t, "deserialized_value", value)
}

func Test_VarLen_Deserialize_Bcd(t *testing.T) {
	lengthDeserializer := &serdes.Mock{}

	serializedValue := []byte{0x12, 0x34, 0x56}
	serializedData := append([]byte("*6*"), serializedValue...)

	dataIn := bytes.NewBuffer(serializedData)

	lengthDeserializer.On("Deserialize", dataIn).Run(func(args mock.Arguments) {
		buffer := args.Get(0).(*bytes.Buffer)
		buffer.Next(3)
	}).Return("6", nil)

	definitions := types.VarLength{
		Length: lengthDeserializer, Data: types.Bcd{},
	}

	value, err := definitions.Deserialize(dataIn)
	assert.NoError(t, err)
	assert.Equal(t, "123456", value)
}

func Test_VarLen_Deserialize_Errors(t *testing.T) {
	lengthDeserializer := &serdes.Mock{}

	dataInLengthError := bytes.NewBufferString("*29*serialized_value_length_error")
	dataInDataError := bytes.NewBufferString("*27*serialized_value_data_error")

	consume := func(args mock.Arguments) {
		buffer := args.Get(0).(*bytes.Buffer)
		buffer.Next(4)
	}

	lengthDeserializer.On("Deserialize", dataInDataError).Run(consume).Return("27", nil)
	lengthDeserializer.On("Deserialize", dataInLengthError).Run(consume).Return(nil, errors.New("error length"))

	dataDeserializer := &serdes.Mock{}
	serializedData := bytes.NewBufferString("serialized_value_data_error")
	dataDeserializer.On("Deserialize", serializedData).Return(nil, errors.New("error data"))

	definitions := types.VarLength{
		Length: lengthDeserializer, Data: dataDeserializer,
	}

	_, err := definitions.Deserialize(dataInLengthError)
	assert.Error(t, err)

	_, err = definitions.Deserialize(dataInDataError)
	assert.Error(t, err)
}
