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
	Id                string
	TransactionNumber string
	Status            string
	Origin            string
	Destination       string
	TransactionType   string
	CreatedAt         internals.NullTime
	SubmittedAt       internals.NullTime
	ApprovedAt        internals.NullTime
	CanceledAt        internals.NullTime
}

type LocationRepository interface {
	GetLocation(ctx context.Context, locationCode string, typeLocation LocationType) ([]LocationRow, error)
	NewDraftTransaction(ctx context.Context, transaction_number, transaction_type, origin, destination string) (TransactionInfo, error)
	CheckTransaction(ctx context.Context, transaction_number string) (TransactionInfo, error)
	SetStatusTransaction(ctx context.Context, transaction_number, status string) error
	// literary like a location send items
	AllocateItem(ctx context.Context, transaction_number, item string) error
	// literary like a location receive items
	ReceiveItem(ctx context.Context, transaction_number, destination, origin string, item string) error
}

type StockTransfer interface {
	LocationRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
