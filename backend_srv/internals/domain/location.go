package domain

import (
	"context"

	"scm-simple-luke.com/dir/internals"
)

type LocationType int

func (t LocationType) String() string {
	return [3]string{
		"customer",
		"technician",
		"warehouse",
	}[t]
}

const (
	CustomerType   LocationType = iota + 1 // 0+1
	TechnicianType                         // 1+1
	WarehouseType                          // 2+1
)

type LocationRow struct {
	LocationCode   string               `json:"location_code"`
	Description    string               `json:"description"`
	ParentLocation internals.NullString `json:"parent_location"`
	CreatedAt      internals.NullTime   `json:"created_at"`
	UpdatedAt      internals.NullTime   `json:"updated_at"`
	Point          any                  `json:"point"`
}

type LocationRepository interface {
	GetLocation(ctx context.Context, locationCode string, typeLocation LocationType) ([]LocationRow, error)
}
