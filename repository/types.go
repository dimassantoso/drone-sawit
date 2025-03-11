// This file contains types that are used in the repository layer.
package repository

import (
	"time"
)

// Filter model
type Filter struct {
	Limit   int    `json:"limit" default:"10" form:"limit"`
	Page    int    `json:"page" default:"1" form:"page"`
	Offset  int    `json:"-"`
	Search  string `json:"search,omitempty" form:"search"`
	OrderBy string `json:"order_by,omitempty" form:"order_by"`
	Sort    string `json:"sort,omitempty" default:"desc" lower:"true" form:"sort"`
	ShowAll bool   `json:"show_all" form:"show_all"`
}

// CalculateOffset method
func (f *Filter) CalculateOffset() int {
	f.Offset = (f.Page - 1) * f.Limit
	return f.Offset
}

// FilterEstate model
type FilterEstate struct {
	Filter
	ID string
}

// FilterEstateTree model
type FilterEstateTree struct {
	Filter
	ID       string
	EstateID string
	X        int
	Y        int
}

type BaseModel struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Estate model
type Estate struct {
	BaseModel
	Width  int
	Length int
}

// EstateTree model
type EstateTree struct {
	BaseModel
	EstateID string
	X        int
	Y        int
	Height   int
}

type EstateTreeStats struct {
	Min    int     `json:"min"`
	Max    int     `json:"max"`
	Median float32 `json:"median"`
}

type CoordinatePoint struct {
	X int
	Y int
}
