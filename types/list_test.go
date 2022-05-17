package types_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/mercadolibre/iso-8583/serdes"
	"github.com/mercadolibre/iso-8583/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_List_Serialize_Success(t *testing.T) {
	serializer := &serdes.Mock{}

	serializer.On("Serialize", "v1_des").Return(bytes.NewBufferString("|v1_ser"), nil)
	serializer.On("Serialize", "v2_des").Return(bytes.NewBufferString("|v2_ser"), nil)
	serializer.On("Serialize", "v3_des").Return(bytes.NewBufferString("|v3_ser"), nil)

	definitions := types.List{
		Items: []types.Field{
			{Name: "f1", SerDes: serializer}, {Name: "f2", SerDes: serializer}, {Name: "f3", SerDes: serializer},
		},
	}

	values := serdes.Map{
		"f1": "v1_des",
		"f2": "v2_des",
		"f3": "v3_des",
	}

	data, err := definitions.Serialize(values)
	assert.NoError(t, err)
	assert.Equal(t, data.Bytes(), []byte("|v1_ser|v2_ser|v3_ser"))
	assert.Equal(t, "list", definitions.Name())
}

func Test_List_Serialize_Invalid_Value_Type(t *testing.T) {
	serializer := &serdes.Mock{}

	definitions := types.List{
		Items: []types.Field{
			{Name: "f1", SerDes: serializer}, {Name: "f2", SerDes: serializer}, {Name: "f3", SerDes: serializer},
		},
	}

	_, err := definitions.Serialize("Invalid type")
	assert.Error(t, err)
}

func Test_List_Serialize_Field_Serializer_Failed(t *testing.T) {

	serializerErrorInput := "v1_serializer_error"

	serializer := &serdes.Mock{}
	serializer.On("Serialize", "v1_des").Return(bytes.NewBufferString("|v1_ser"), nil)
	serializer.On("Serialize", serializerErrorInput).Return(bytes.NewBufferString(""), errors.New("error"))

	definitions := types.List{
		Items: []types.Field{
			{Name: "f1", SerDes: serializer}, {Name: "f2", SerDes: serializer}, {Name: "f3", SerDes: serializer},
		},
	}

	values := serdes.Map{
		"f1": "v1_des",
		"f2": serializerErrorInput,
	}

	_, err := definitions.Serialize(values)
	assert.Error(t, err)

	values = serdes.Map{
		"f1": serializerErrorInput,
	}

	_, err = definitions.Serialize(values)
	assert.Error(t, err)
}

func Test_Deserializer_Deserialize_Success(t *testing.T) {
	deserializer := &serdes.Mock{}
	strIn := "|v1_ser|v2_ser|v3_ser"
	len := len(strIn) / 3

	consume := func(args mock.Arguments) {
		buffer := args.Get(0).(*bytes.Buffer)
		buffer.Next(len)
	}

	for numValues := 1; numValues <= 3; numValues++ {
		in := bytes.NewBufferString(strIn)
		in.Next(len * (numValues - 1))
		out := fmt.Sprintf("v%d_des", numValues)
		deserializer.On("Deserialize", in).Run(consume).Return(out, nil)
	}

	definitions := types.List{
		Items: []types.Field{
			{Name: "f1", SerDes: deserializer}, {Name: "f2", SerDes: deserializer}, {Name: "f3", SerDes: deserializer},
		},
	}

	dataInputDeserializeList := bytes.NewBufferString(strIn)
	values, err := definitions.Deserialize(dataInputDeserializeList)
	assert.NoError(t, err)

	expected := serdes.Map{
		"f1": "v1_des",
		"f2": "v2_des",
		"f3": "v3_des",
	}

	assert.Equal(t, expected, values)
}

func Test_Deserializer_Deserialize_Error(t *testing.T) {
	deserializer := &serdes.Mock{}
	strIn := "|v1_ser|v2_ser|v3_ser"

	deserializer.On("Deserialize", mock.Anything).Return("", errors.New("error"))

	definitions := types.List{
		Items: []types.Field{
			{Name: "f1", SerDes: deserializer}, {Name: "f2", SerDes: deserializer}, {Name: "f3", SerDes: deserializer},
		},
	}

	dataInputDeserializeList := bytes.NewBufferString(strIn)
	_, err := definitions.Deserialize(dataInputDeserializeList)
	assert.Error(t, err)
}

func Test_Serialize_Bit90_Error(t *testing.T) {
	layout := types.List{Desc: "Original Transaction Data",
		Items: []types.Field{
			{Name: "1", SerDes: types.Bcd{NumDigits: 4, Desc: "original message type"}},
			{Name: "2", SerDes: types.Bcd{NumDigits: 6, Desc: "original trace number"}},
			{Name: "3", SerDes: types.Bcd{NumDigits: 10, Desc: "original transaction date"}},
			{Name: "4", SerDes: types.Bcd{NumDigits: 22, Desc: "original acquirer id and forwarding institution id"}},
		},
	}

	values := map[string]interface{}{
		"1": "0100", "2": "042249", "3": "2403221416", "4": "400552",
	}

	expected := []byte{
		0x01, 0x00, 0x04, 0x22, 0x49, 0x24, 0x03, 0x22, 0x14, 0x16, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x40, 0x05, 0x52,
	}

	ser, err := layout.Serialize(values)
	assert.NoError(t, err)
	assert.Equal(t, expected, ser.Bytes())
}
