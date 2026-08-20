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

type TransactionInfo struct {
	Transaction_number string
	Created_by         string
	Submitted_at       string
	Approved_at        string
	Canceled_at        string
	Transaction_type   string
}

type LocationRepository interface {
	GetLocation(ctx context.Context, locationCode string, typeLocation LocationType) ([]LocationRow, error)
	NewDraftTransaction(ctx context.Context, transaction_number, transaction_type, origin, destination string) (TransactionInfo, error)
	CheckTransaction(transaction_number string) error
	// literary like a location send items
	SendItem(ctx context.Context, transaction_number, origin, destination, item string) error
	// literary like a location receive items
	ReceiveItem(ctx context.Context, transaction_number, destination, origin string, item string) error
}

type StockTransfer interface {
	LocationRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
