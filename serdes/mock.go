package serdes

import (
	"bytes"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Name() string {
	return "mock"
}

func (m *Mock) Serialize(value Value) (*bytes.Buffer, error) {
	called := m.Called(value)
	return called.Get(0).(*bytes.Buffer), called.Error(1)
}

func (m *Mock) Deserialize(data *bytes.Buffer) (Value, error) {
	called := m.Called(data)
	return called.Get(0), called.Error(1)
}
