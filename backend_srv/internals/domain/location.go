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
	Id                string             `json:"id"`
	TransactionNumber string             `json:"transaction_number"`
	Status            string             `json:"status"`
	Origin            string             `json:"origin"`
	Destination       string             `json:"Destination"`
	TransactionType   string             `json:"transaction_type"`
	CreatedAt         internals.NullTime `json:"created_at"`
	SubmittedAt       internals.NullTime `json:"submitted_at"`
	ApprovedAt        internals.NullTime `json:"approved_at"`
	CanceledAt        internals.NullTime `json:"canceled_at"`
}

type LocationRepository interface {
	GetLocation(ctx context.Context, locationCode string, typeLocation LocationType) ([]LocationRow, error)
	NewDraftTransaction(ctx context.Context, transaction_number, transaction_type, origin, destination string) (TransactionInfo, error)
	CheckTransaction(ctx context.Context, transaction_number string) (TransactionInfo, error)
	SetStatusTransaction(ctx context.Context, transaction_number, status string) error
	// literary like a location send items
	AllocateItem(ctx context.Context, transaction_number, item string) error
	DisAllocateItem(ctx context.Context, transaction_number, item string) error
	// literary like a location receive items
	ReceiveItem(ctx context.Context, transaction_number, destination, origin string, item string) error
	GetTransactionInfo(ctx context.Context, transaction_number string) (TransactionInfo, error)
	GetItemsOnTransaction(ctx context.Context, transaction_number string) ([]EachItemTransaction, error)

	InputOutboundTracker(ctx context.Context, transaction_number string) error
	EditOutboundTracker(ctx context.Context, transaction_number string) error // set pending (if delivery_at and arrived_at null then pending)

	InputInboundTracker(ctx context.Context, transaction_number, outbound_number string) error
	CalculateDurationTransfer(ctx context.Context, outbound_number, inbound_number string) (TransferDuration, error)
}

type StockTransfer interface {
	LocationRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
