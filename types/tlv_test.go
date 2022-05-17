package types_test

import (
	"bytes"
	"testing"

	"github.com/mercadolibre/iso-8583/serdes"
	"github.com/mercadolibre/iso-8583/types"

	"github.com/stretchr/testify/assert"
)

func TestTLV_SerializeDeserialize(t *testing.T) {
	tests := []struct {
		name    string
		data    serdes.Value
		taipe   serdes.Serdes
		want    *bytes.Buffer
		wantErr bool
	}{
		{
			name: "serialize and deserialize with tcc",
			taipe: types.VarLength{Length: types.EbcdicNumeric{NumDigits: 3}, Desc: "test", Data: types.List{
				Items: []types.Field{
					{Name: "tcc", SerDes: types.Ebcdic{NumDigits: 1, Desc: "tcc"}},
					{Name: "se", SerDes: types.TLV{
						Items: []types.Field{
							{Name: "21", SerDes: types.Ebcdic{}},
							{Name: "61", SerDes: types.Ebcdic{}},
						}},
					},
				},
			}},
			data:    map[string]interface{}{"tcc": "R", "se": map[string]interface{}{"21": "01010", "61": "00001"}},
			want:    bytes.NewBuffer([]byte{0xf0, 0xf1, 0xf9, 0xd9, 0xf2, 0xf1, 0xf0, 0xf5, 0xf0, 0xf1, 0xf0, 0xf1, 0xf0, 0xf6, 0xf1, 0xf0, 0xf5, 0xf0, 0xf0, 0xf0, 0xf0, 0xf1}),
			wantErr: false,
		},
		{
			name: "serialize and deserialize with tcc and subelements def",
			taipe: types.VarLength{Length: types.EbcdicNumeric{NumDigits: 3}, Desc: "test", Data: types.List{
				Items: []types.Field{
					{Name: "tcc", SerDes: types.Ebcdic{NumDigits: 1, Desc: "tcc"}},
					{Name: "se", SerDes: types.TLV{
						Items: []types.Field{
							{Name: "21", SerDes: types.TLV{
								Items: []types.Field{
									{Name: "01", SerDes: types.Ebcdic{NumDigits: 1, Desc: "mPOS Acceptance Device Type"}},
									{Name: "02", SerDes: types.Ebcdic{NumDigits: 2, Desc: "Additional Terminal Capability Indicator"}},
								},
							}},
							{Name: "61", SerDes: types.List{
								Items: []types.Field{
									{Name: "1", SerDes: types.Ebcdic{NumDigits: 1, Desc: "Partial Approval Terminal Support Indicator"}},
									{Name: "2", SerDes: types.Ebcdic{NumDigits: 1, Desc: "Purchase Amount Only Terminal Support Indicator"}},
									{Name: "3", SerDes: types.Ebcdic{NumDigits: 1, Desc: "Real-time Substantiation Indicator"}},
									{Name: "4", SerDes: types.Ebcdic{NumDigits: 1, Desc: "Merchant Transaction Fraud Scoring Indicator"}},
									{Name: "5", SerDes: types.Ebcdic{NumDigits: 1, Desc: "Final Authorization Indicator"}},
								},
							}},
						}},
					},
				},
			}},
			data:    map[string]interface{}{"tcc": "R", "se": map[string]interface{}{"21": map[string]interface{}{"01": "0"}, "61": map[string]interface{}{"1": "0", "2": "0", "3": "0", "4": "0", "5": "1"}}},
			want:    bytes.NewBuffer([]byte{0xf0, 0xf1, 0xf9, 0xd9, 0xf2, 0xf1, 0xf0, 0xf5, 0xf0, 0xf1, 0xf0, 0xf1, 0xf0, 0xf6, 0xf1, 0xf0, 0xf5, 0xf0, 0xf0, 0xf0, 0xf0, 0xf1}),
			wantErr: false,
		},
		{
			name: "serialize and deserialize with tcc and subelements 37",
			taipe: types.VarLength{Length: types.EbcdicNumeric{NumDigits: 3}, Desc: "test", Data: types.List{
				Items: []types.Field{
					{Name: "tcc", SerDes: types.Ebcdic{NumDigits: 1, Desc: "tcc"}},
					{Name: "se", SerDes: types.TLV{
						Items: []types.Field{
							{Name: "37", SerDes: types.TLV{
								Desc: "Additional Merchant Data",
								Items: []types.Field{
									{Name: "01", SerDes: types.EbcdicNumeric{NumDigits: 11, Desc: "Payment Facilitator ID"}},
									{Name: "02", SerDes: types.EbcdicNumeric{NumDigits: 11, Desc: "Independent Sales Organization ID"}},
									{Name: "03", SerDes: types.Ebcdic{NumDigits: 15, Desc: "Sub-Merchant ID"}},
								},
							}},
						}},
					},
				},
			}},
			data:    map[string]interface{}{"tcc": "T", "se": map[string]interface{}{"37": map[string]interface{}{"01": "00000241408", "03": "51"}}},
			want:    bytes.NewBuffer([]byte{0xf0, 0xf3, 0xf9, 0xe3, 0xf3, 0xf7, 0xf3, 0xf4, 0xf0, 0xf1, 0xf1, 0xf1, 0xf0, 0xf0, 0xf0, 0xf0, 0xf0, 0xf2, 0xf4, 0xf1, 0xf4, 0xf0, 0xf8, 0xf0, 0xf3, 0xf1, 0xf5, 0xf5, 0xf1, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40}),
			wantErr: false,
		},
		{
			name: "serialize and deserialize",
			taipe: types.TLV{SizeLen: 2, Items: []types.Field{
				{
					Name: "92",
					SerDes: types.Ebcdic{
						Desc:      "CVC 2",
						NumDigits: 3,
					},
				},
			}},
			data:    map[string]interface{}{"92": "123"},
			want:    bytes.NewBuffer([]byte{0xf9, 0xf2, 0xf0, 0xf3, 0xf1, 0xf2, 0xf3}),
			wantErr: false,
		},
		{
			name: "serialize and deserialize multiple tags",
			taipe: types.TLV{SizeLen: 2, Items: []types.Field{
				{
					Name: "21",
					SerDes: types.Ebcdic{
						Desc:      "Acceptance Data",
						NumDigits: 5,
					},
				},
				{
					Name: "61",
					SerDes: types.Ebcdic{
						Desc:      "POS Data, Extended Condition Codes",
						NumDigits: 5,
					},
				},
				{
					Name: "92",
					SerDes: types.Ebcdic{
						Desc:      "CVC 2",
						NumDigits: 3,
					},
				},
			}},
			data:    map[string]interface{}{"21": "01010", "61": "00001", "92": "123"},
			want:    bytes.NewBuffer([]byte{0xf2, 0xf1, 0xf0, 0xf5, 0xf0, 0xf1, 0xf0, 0xf1, 0xf0, 0xf6, 0xf1, 0xf0, 0xf5, 0xf0, 0xf0, 0xf0, 0xf0, 0xf1, 0xf9, 0xf2, 0xf0, 0xf3, 0xf1, 0xf2, 0xf3}),
			wantErr: false,
		},
		{
			name: "serialize and deserialize raw emv tlv",
			taipe: types.TLV{SizeLen: 2, Desc: "icc", Items: []types.Field{
				{
					Name:   "01",
					SerDes: types.Raw{},
				},
			}},
			data:    map[string]interface{}{"01": "9f26085dfaeed4e8ed3f8f"},
			want:    bytes.NewBuffer([]byte{0xf0, 0xf1, 0xf1, 0xf1, 0x9f, 0x26, 0x8, 0x5d, 0xfa, 0xee, 0xd4, 0xe8, 0xed, 0x3f, 0x8f}),
			wantErr: false,
		},
		{
			name: "serialize and deserialize with size len = 3",
			taipe: types.TLV{
				Desc:    "additional data (national use)",
				SizeLen: 3, SizeTag: 3,
				Items: []types.Field{
					{Name: "001", SerDes: types.Ebcdic{NumDigits: 4, Desc: "Installment Payment Data"}},
					{Name: "002", SerDes: types.Ebcdic{Desc: "Installment Payment Response Data"}},
				},
			},
			data:    map[string]interface{}{"001": "2001"},
			want:    bytes.NewBuffer([]byte{0xf0, 0xf0, 0xf1, 0xf0, 0xf0, 0xf4, 0xf2, 0xf0, 0xf0, 0xf1}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.taipe.Serialize(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TLV.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)

			if got != nil {
				got2, err := tt.taipe.Deserialize(got)
				if (err != nil) != tt.wantErr {
					t.Errorf("TLV.Deserialize() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Log(got2)
				assert.Equal(t, tt.data, got2)
			}
		})
	}
}

func TestTLV_ErrorSerialize(t *testing.T) {
	tests := []struct {
		name    string
		data    serdes.Value
		taipe   serdes.Serdes
		want    *bytes.Buffer
		wantErr bool
	}{
		{
			name: "error serialize odd value",
			taipe: types.TLV{
				Items: []types.Field{
					{Name: "aa", SerDes: types.Raw{}},
				},
			},
			data:    map[string]interface{}{"aa": "12345"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error serialize empty field name",
			taipe: types.TLV{
				Items: []types.Field{
					{Name: "", SerDes: types.Raw{}},
				},
			},
			data:    map[string]interface{}{"": "12345"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error serialize invalid serdes data",
			taipe: types.TLV{
				Items: []types.Field{
					{Name: "01", SerDes: types.Raw{}},
				},
			},
			data:    99,
			want:    nil,
			wantErr: true,
		},
		{
			name: "error serialize two tag but only one def",
			taipe: types.TLV{
				Items: []types.Field{
					{Name: "10", SerDes: types.Raw{}},
				},
			},
			data:    map[string]interface{}{"10": "123465", "1f": "123465"},
			want:    bytes.NewBuffer([]byte{0xf1, 0xf0, 0xf0, 0xf3, 0x12, 0x34, 0x65}),
			wantErr: false,
		},
		{
			name: "error serialize invalid capacity length",
			taipe: types.TLV{
				SizeLen: 1, SizeTag: 1,
				Items: []types.Field{
					{Name: "1", SerDes: types.Ebcdic{NumDigits: 10}},
				},
			},
			data:    map[string]interface{}{"1": "1234567890"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.taipe.Serialize(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TLV.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTLV_ErrorDeserialize(t *testing.T) {
	tests := []struct {
		name    string
		data    *bytes.Buffer
		taipe   serdes.Serdes
		want    serdes.Value
		wantErr bool
	}{
		{
			name: "error deserialize two tag but only one def",
			taipe: types.TLV{
				Items: []types.Field{
					{Name: "11", SerDes: types.Raw{}},
				},
			},
			data:    bytes.NewBuffer([]byte{0xf1, 0xf0, 0xf0, 0xf3, 0xf1, 0xf2, 0xf3, 0xf1, 0xf1, 0xf0, 0xf3, 0x12, 0x34, 0x56}),
			want:    map[string]interface{}{"10": "123", "11": "123456"},
			wantErr: false,
		},
		{
			name: "continue deserialize invalid field name mapping",
			taipe: types.TLV{
				Items: []types.Field{
					{Name: "11", SerDes: types.Raw{}},
				},
			},
			data:    bytes.NewBuffer([]byte{0xf1, 0xf0, 0x00, 0x00}),
			want:    map[string]interface{}{"10": ""},
			wantErr: false,
		},
		{
			name: "error deserialize invalid length",
			taipe: types.TLV{
				Items: []types.Field{
					{Name: "10", SerDes: types.Raw{}},
				},
			},
			data:    bytes.NewBuffer([]byte{0xf1, 0xf0, 0x00, 0x00, 0x01}),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.taipe.Deserialize(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TLV.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
