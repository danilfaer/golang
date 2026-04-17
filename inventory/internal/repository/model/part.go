package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Value interface {
	isValue()
}

type StringValue struct {
	StringValue string
}
type Int64Value struct {
	Int64Value int64
}
type DoubleValue struct {
	DoubleValue float64
}
type BoolValue struct {
	BoolValue bool
}

func (v StringValue) isValue() {}
func (v Int64Value) isValue()  {}
func (v DoubleValue) isValue() {}
func (v BoolValue) isValue()   {}

// Part — модель детали репозитория и документ MongoDB (поле _id = UUID).
type Part struct {
	UUID          string              `bson:"_id"`
	Name          string              `bson:"name"`
	Description   string              `bson:"description"`
	Price         float64             `bson:"price"`
	StockQuantity int64               `bson:"stock_quantity"`
	Category      Category            `bson:"category"`
	Dimensions    *Dimensions         `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer       `bson:"manufacturer,omitempty"`
	Tags          []string            `bson:"tags,omitempty"`
	Metadata      map[string]*Value   `bson:"-"`
	CreatedAt     time.Time           `bson:"created_at"`
	UpdatedAt     time.Time           `bson:"updated_at"`
}

// UnmarshalBSON читает документ MongoDB; metadata — вложенный документ с простыми типами.
func (p *Part) UnmarshalBSON(data []byte) error {
	type wire struct {
		UUID          string         `bson:"_id"`
		Name          string         `bson:"name"`
		Description   string         `bson:"description"`
		Price         float64        `bson:"price"`
		StockQuantity int64          `bson:"stock_quantity"`
		Category      Category       `bson:"category"`
		Dimensions    *Dimensions    `bson:"dimensions,omitempty"`
		Manufacturer  *Manufacturer  `bson:"manufacturer,omitempty"`
		Tags          []string       `bson:"tags,omitempty"`
		Metadata      bson.M         `bson:"metadata,omitempty"`
		CreatedAt     time.Time      `bson:"created_at"`
		UpdatedAt     time.Time      `bson:"updated_at"`
	}
	var w wire
	if err := bson.Unmarshal(data, &w); err != nil {
		return err
	}
	*p = Part{
		UUID:          w.UUID,
		Name:          w.Name,
		Description:   w.Description,
		Price:         w.Price,
		StockQuantity: w.StockQuantity,
		Category:      w.Category,
		Dimensions:    w.Dimensions,
		Manufacturer:  w.Manufacturer,
		Tags:          w.Tags,
		Metadata:      metadataBSONToModel(w.Metadata),
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
	return nil
}

// MarshalBSON сериализует Part для записи в MongoDB.
func (p *Part) MarshalBSON() ([]byte, error) {
	type wire struct {
		UUID          string         `bson:"_id"`
		Name          string         `bson:"name"`
		Description   string         `bson:"description"`
		Price         float64        `bson:"price"`
		StockQuantity int64          `bson:"stock_quantity"`
		Category      Category       `bson:"category"`
		Dimensions    *Dimensions    `bson:"dimensions,omitempty"`
		Manufacturer  *Manufacturer  `bson:"manufacturer,omitempty"`
		Tags          []string       `bson:"tags,omitempty"`
		Metadata      bson.M         `bson:"metadata,omitempty"`
		CreatedAt     time.Time      `bson:"created_at"`
		UpdatedAt     time.Time      `bson:"updated_at"`
	}
	w := wire{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      p.Category,
		Dimensions:    p.Dimensions,
		Manufacturer:  p.Manufacturer,
		Tags:          p.Tags,
		Metadata:      metadataModelToBSON(p.Metadata),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
	return bson.Marshal(w)
}

func metadataBSONToModel(m bson.M) map[string]*Value {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]*Value, len(m))
	for k, raw := range m {
		v, err := bsonValueToValue(raw)
		if err != nil {
			continue
		}
		out[k] = v
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func metadataModelToBSON(meta map[string]*Value) bson.M {
	if len(meta) == 0 {
		return nil
	}
	out := make(bson.M, len(meta))
	for k, ptr := range meta {
		if ptr == nil {
			continue
		}
		switch v := (*ptr).(type) {
		case StringValue:
			out[k] = v.StringValue
		case Int64Value:
			out[k] = v.Int64Value
		case DoubleValue:
			out[k] = v.DoubleValue
		case BoolValue:
			out[k] = v.BoolValue
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func bsonValueToValue(raw any) (*Value, error) {
	switch x := raw.(type) {
	case string:
		var v Value = StringValue{StringValue: x}
		return &v, nil
	case int32:
		var v Value = Int64Value{Int64Value: int64(x)}
		return &v, nil
	case int64:
		var v Value = Int64Value{Int64Value: x}
		return &v, nil
	case float64:
		var v Value = DoubleValue{DoubleValue: x}
		return &v, nil
	case bool:
		var v Value = BoolValue{BoolValue: x}
		return &v, nil
	default:
		return nil, fmt.Errorf("unsupported metadata type %T", raw)
	}
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
