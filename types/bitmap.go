package types

import (
	"bytes"

	"github.com/mercadolibre/iso-8583/serdes"
)

type Bitmap struct {
	BlockSize int
	NumBits   int
}

func (Bitmap) Name() string {
	return "bitmap"
}

func (bitmap Bitmap) Serialize(value serdes.Value) (*bytes.Buffer, error) {
	blockSizeInBytes := bitmap.BlockSize / 8
	rawValue, ok := value.([]byte)
	if !ok {
		return nil, SerializerError{
			Message: "invalid value type", Serdes: bitmap, Value: value,
		}
	}

	lenValue := len(rawValue)
	numBlocks := lenValue / blockSizeInBytes
	padding := blockSizeInBytes - (lenValue % blockSizeInBytes)
	if padding != blockSizeInBytes {
		numBlocks++
		rawValue = append(rawValue, make([]byte, padding)...)
	}

	for blockIndex := 0; blockIndex < numBlocks; blockIndex++ {
		offset := blockIndex * blockSizeInBytes

		if blockIndex < numBlocks-1 {
			// if this is not the last block, we need to set the msb bit to indicate there are more blocks.
			rawValue[offset] |= 0x80
		}
	}

	return bytes.NewBuffer(rawValue), nil
}

func (bitmap Bitmap) Deserialize(data *bytes.Buffer) (serdes.Value, error) {
	maxNumBlocks := bitmap.NumBits / bitmap.BlockSize
	blockSizeInBytes := bitmap.BlockSize / 8

	var value []byte
	moreBlocks := true

	for blockIndex := 0; blockIndex < maxNumBlocks && moreBlocks; blockIndex++ {
		if data.Len() < blockSizeInBytes {
			return nil, DeserializationError{
				Message: "data has no bytes enough to decode block", Serdes: bitmap, Remaning: data.Len(),
			}
		}

		block := make([]byte, blockSizeInBytes)
		if _, err := data.Read(block); err != nil {
			return nil, DeserializationError{
				Message: "error reading data", Serdes: bitmap, Remaning: data.Len(),
			}
		}

		moreBlocks = (block[0] & 0x80) > 0
		if moreBlocks && blockIndex < maxNumBlocks-1 {
			block[0] &= 0x7F
		}

		value = append(value, block...)
	}

	return value, nil
}
