package domain

import (
	"context"

	"scm-simple-luke.com/dir/internals"
)

type ItemInfo struct {
	Id                    string               `json:"id"`
	Serial_number         string               `json:"serial_number"`
	Factory_serial_number string               `json:"factory_serial_number"`
	Created_at            internals.NullTime   `json:"created_at"`
	Curr_status           string               `json:"curr_status"`
	Curr_transaction      internals.NullString `json:"curr_transaction"`
	Curr_location_code    string               `json:"curr_location_code"`
	Product_code          string               `json:"product_code"`
	Introduction_number   string               `json:"introduction_number"`
}

type ItemRepository interface {
	// note kalo error balikin aray kosong aja(?)
	GetItem(ctx context.Context, identifier string) ([]ItemInfo, error)
	UnlockItems(ctx context.Context, transaction_number string) error
	GetItemsOnTransaction(ctx context.Context, transaction_number string) ([]EachItemTransaction, error)
	AllocateItem(ctx context.Context, transaction_number, identifier string) error
	DisAllocateItem(ctx context.Context, transaction_number, item string) error

	ReceiveInboundItem(ctx context.Context, transaction_number, identifier string) error
	CheckTransaction(ctx context.Context, transaction_number string) (TransactionInfo, error)
}

type WarehouseOutboundInfo struct {
	TransactionInfo TransactionInfo       `json:"transaction_info"`
	ListItems       []EachItemTransaction `json:"list_items"`
}

type WarehouseInboundInfo struct {
	TransactionInfo TransactionInfo       `json:"transaction_info"`
	ListItems       []EachItemTransaction `json:"list_items"`
}

type EachItemTransaction struct {
	Id           int32              `json:"id"` // int32 adjust serial pgdata type as int4
	IdTransfer   int32              `json:"id_trans_item_transfer"`
	SerialNumber string             `json:"identifier_item"`
	AddedAt      internals.NullTime `json:"added_at"`
}
