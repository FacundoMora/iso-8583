package types

import (
	"errors"
	"testing"
)

func TestDeserializationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "just message",
			err: DeserializationError{
				Message: "error deserialize",
			},
			want: "error deserialize",
		},
		{
			name: "message and cause",
			err: DeserializationError{
				Message: "error deserialize",
				Cause:   errors.New("error cause."),
			},
			want: "error deserialize -> error cause.",
		},
		{
			name: "message and cause and serdes",
			err: DeserializationError{
				Serdes:  Byte{},
				Message: "error deserialize",
				Cause:   errors.New("error cause."),
			},
			want: "error deserialize: deserializer: byte. -> error cause.",
		},
		{
			name: "message and cause and serdes and field",
			err: DeserializationError{
				Serdes:  Byte{},
				Field:   Field{Name: "01", SerDes: Byte{}},
				Message: "error deserialize",
				Cause:   errors.New("error cause."),
			},
			want: "error deserialize: deserializer: byte. field name: 01 - field serdes name: byte. -> error cause.",
		},
		{
			name: "message and cause and serdes and field and data remaning",
			err: DeserializationError{
				Serdes:   Byte{},
				Field:    Field{Name: "01", SerDes: Byte{}},
				Message:  "error deserialize",
				Cause:    errors.New("error cause."),
				Remaning: 20,
			},
			want: "error deserialize: deserializer: byte. field name: 01 - field serdes name: byte. *data remaning: 20. -> error cause.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("DeserializationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "just message",
			err: SerializerError{
				Message: "error serialize",
			},
			want: "error serialize",
		},
		{
			name: "just cause",
			err: SerializerError{
				Cause: errors.New("error cause."),
			},
			want: " -> error cause.",
		},
		{
			name: "message and cause",
			err: SerializerError{
				Message: "error serialize",
				Cause:   errors.New("error cause."),
			},
			want: "error serialize -> error cause.",
		},
		{
			name: "message and cause and serdes",
			err: SerializerError{
				Serdes:  Byte{},
				Message: "error serialize",
				Cause:   errors.New("error cause."),
			},
			want: "error serialize: serializer: byte. -> error cause.",
		},
		{
			name: "message and cause and serdes and field",
			err: SerializerError{
				Serdes:  Byte{},
				Field:   Field{Name: "01", SerDes: Byte{}},
				Message: "error serialize",
				Cause:   errors.New("error cause."),
			},
			want: "error serialize: serializer: byte. field name: 01 - field serdes name: byte. -> error cause.",
		},
		{
			name: "message and cause and serdes and field and value",
			err: SerializerError{
				Serdes:  Byte{},
				Field:   Field{Name: "01", SerDes: Byte{}},
				Message: "error serialize",
				Cause:   errors.New("error cause."),
				Value:   30,
			},
			want: "error serialize: serializer: byte. field name: 01 - field serdes name: byte. value type: int. -> error cause.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("DeserializationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
