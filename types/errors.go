package types

import (
	"fmt"

	"github.com/mercadolibre/iso-8583/serdes"
)

type SerializerError struct {
	Message string       `json:"message"`
	Serdes  serdes.Named `json:"serdes"`
	Field   Field        `json:"field"`
	Value   serdes.Value `json:"value"`
	Cause   error        `json:"cause"`
}

func (err SerializerError) Error() string {
	errSerdes := err.Serdes
	var serdesName string
	if errSerdes != nil {
		serdesName = fmt.Sprintf(" serializer: %s.", errSerdes.Name())
	}

	errField := err.Field
	var field string
	if errField.SerDes != nil {
		field = fmt.Sprintf(" field name: %s - field serdes name: %s.", errField.Name, errField.SerDes.Name())
	}

	errValue := err.Value
	var value string
	if errValue != nil {
		value = fmt.Sprintf(" value type: %T.", errValue)

	}

	var delimiter string
	if value != "" || field != "" || serdesName != "" {
		delimiter = ":"
	}

	if err.Cause != nil {
		return fmt.Sprintf("%s%s%s%s%s -> %+v",
			err.Message, delimiter, serdesName, field, value, err.Cause)
	}

	return fmt.Sprintf("%s%s%s%s%s",
		err.Message, delimiter, serdesName, field, value)

}

func (err SerializerError) Unwrap() error {
	return err.Cause
}

type DeserializationError struct {
	Message  string       `json:"message"`
	Serdes   serdes.Named `json:"serdes"`
	Field    Field        `json:"field"`
	Remaning int          `json:"remaning"`
	Cause    error        `json:"cause"`
}

func (err DeserializationError) Error() string {
	errSerdes := err.Serdes
	var serdesName string
	if errSerdes != nil {
		serdesName = fmt.Sprintf(" deserializer: %s.", errSerdes.Name())
	}

	errField := err.Field
	var field string
	if errField.SerDes != nil {
		field = fmt.Sprintf(" field name: %s - field serdes name: %s.", errField.Name, errField.SerDes.Name())

	}

	errReman := err.Remaning
	var remaning string
	if errReman > 0 {
		remaning = fmt.Sprintf(" *data remaning: %d.", errReman)
	}

	var delimiter string
	if field != "" || serdesName != "" || remaning != "" {
		delimiter = ":"
	}

	if err.Cause != nil {
		return fmt.Sprintf("%s%s%s%s%s -> %+v",
			err.Message, delimiter, serdesName, field, remaning, err.Cause)
	}

	return fmt.Sprintf("%s%s%s%s%s",
		err.Message, delimiter, serdesName, field, remaning)

}

func (err DeserializationError) Unwrap() error {
	return err.Cause
}
