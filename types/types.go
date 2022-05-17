package types

import (
	"github.com/mercadolibre/iso-8583/serdes"
)

type Field struct {
	Name   string
	SerDes serdes.Serdes
}

type Desc string
