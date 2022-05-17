package serdes

import "bytes"

type Serializer interface {
	Serialize(value Value) (*bytes.Buffer, error)
}

type Deserialize interface {
	Deserialize(data *bytes.Buffer) (Value, error)
}

type Named interface {
	Name() string
}

type Serdes interface {
	Named
	Serializer
	Deserialize
}
