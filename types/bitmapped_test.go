package types_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/mercadolibre/go-iso8583/serdes"
	"github.com/mercadolibre/go-iso8583/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_BitMapped_Serialize_Success(t *testing.T) {
	serializerFields := &serdes.Mock{}
	serializerFields.On("Serialize", "v_002_des").Return(bytes.NewBufferString("|v_002_des"), nil)
	serializerFields.On("Serialize", "v_066_des").Return(bytes.NewBufferString("|v_066_des"), nil)
	serializerFields.On("Serialize", "v_100_des").Return(bytes.NewBufferString("|v_100_des"), nil)

	serializerMap := types.BitMapped{
		Bitmap: types.Bitmap{BlockSize: 64, NumBits: 128},
		Mapping: map[int]serdes.Serdes{
			2:   serializerFields,
			66:  serializerFields,
			100: serializerFields,
		},
	}

	value := serdes.Map{
		"2":   "v_002_des",
		"66":  "v_066_des",
		"100": "v_100_des",
	}

	data, err := serializerMap.Serialize(value)
	assert.NoError(t, err)

	bitmap := []byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}
	expected := append(bitmap, []byte("|v_002_des|v_066_des|v_100_des")...)
	assert.Equal(t, expected, data.Bytes())
}

func Test_BitMapped_Serialize_Errors(t *testing.T) {
	serializerFields := &serdes.Mock{}
	serializerFields.On("Serialize", "v_002_des").Return(bytes.NewBufferString("|v_002_des"), nil)
	serializerFields.On("Serialize", "v_066_des").Return(bytes.NewBufferString("|v_066_des"), nil)
	serializerFields.On("Serialize", "v_070_invalid").Return(nil, errors.New("error"))
	serializerFields.On("Serialize", "v_100_des").Return(bytes.NewBufferString("|v_100_des"), nil)

	serializerMap := types.BitMapped{
		Bitmap: types.Bitmap{BlockSize: 64, NumBits: 128},
		Mapping: map[int]serdes.Serdes{
			2:   serializerFields,
			66:  serializerFields,
			100: serializerFields,
		},
	}

	valueUnkownField := serdes.Map{
		"2":   "v_002_des",
		"70":  "v_070_des_invalid",
		"100": "v_100_des",
	}

	_, err := serializerMap.Serialize(valueUnkownField)
	assert.Error(t, err)

	_, err = serializerMap.Serialize("invalid type")
	assert.Error(t, err)
}

func Test_BitMapped_Deserialize_Success(t *testing.T) {
	deserializer := &serdes.Mock{}

	strIn := "|v_002_ser|v_066_ser|v_100_ser"
	bitmap := []byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}
	dataIn := append(bitmap, []byte(strIn)...)
	streLen := len(strIn) / 3

	consume := func(args mock.Arguments) {
		buffer := args.Get(0).(*bytes.Buffer)
		buffer.Next(streLen)
	}

	fields := []int{2, 66, 100}
	for numValues := 1; numValues <= 3; numValues++ {
		in := bytes.NewBuffer(dataIn)
		in.Next(len(bitmap))

		next := streLen * (numValues - 1)
		if next > 0 {
			in.Next(next)
		}

		out := fmt.Sprintf("v_%03d_des", fields[numValues-1])
		deserializer.On("Deserialize", in).Run(consume).Return(out, nil)
	}

	definitions := types.BitMapped{
		Bitmap: types.Bitmap{BlockSize: 64, NumBits: 128},
		Mapping: map[int]serdes.Serdes{
			2:   deserializer,
			66:  deserializer,
			100: deserializer,
		},
	}

	data := bytes.NewBuffer(dataIn)
	value, err := definitions.Deserialize(data)
	assert.NoError(t, err)

	expected := serdes.Map{
		"2": "v_002_des", "66": "v_066_des", "100": "v_100_des",
	}

	assert.Equal(t, expected, value)
}

func Test_BitMapped_Deserialize_Error_Deserializing_Field(t *testing.T) {
	deserializer := &serdes.Mock{}

	strIn := "|v_002_ser|v_066_ser|v_100_ser"
	bitmap := []byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}
	dataIn := append(bitmap, []byte(strIn)...)

	in := bytes.NewBuffer(dataIn)
	in.Next(len(bitmap))

	deserializer.On("Deserialize", in).Return(nil, errors.New("error"))

	definitions := types.BitMapped{
		Bitmap: types.Bitmap{BlockSize: 64, NumBits: 128},
		Mapping: map[int]serdes.Serdes{
			2:   deserializer,
			66:  deserializer,
			100: deserializer,
		},
	}

	data := bytes.NewBuffer(dataIn)
	_, err := definitions.Deserialize(data)
	assert.Error(t, err)

}

func Test_BitMapped_Deserialize_Field_Unknown(t *testing.T) {
	deserializer := &serdes.Mock{}

	strIn := "|v_002_ser|v_066_ser|v_100_ser"
	bitmap := []byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}
	dataIn := append(bitmap, []byte(strIn)...)

	definitions := types.BitMapped{
		Bitmap: types.Bitmap{BlockSize: 64, NumBits: 128},
		Mapping: map[int]serdes.Serdes{
			66:  deserializer,
			100: deserializer,
		},
	}

	data := bytes.NewBuffer(dataIn)
	_, err := definitions.Deserialize(data)
	assert.Error(t, err)
}
