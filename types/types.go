package types

import (
	"github.com/mercadolibre/go-iso8583/serdes"
)

type Field struct {
	Name   string
	SerDes serdes.Serdes
}

type Desc string
