package domain

import (
	"errors"
	"fmt"
	"time"

	"scm-simple-luke.com/dir/internals"
)

// implement location repo

type Warehouse struct {
	Description    string `json:"description"`
	LocationCode   string `json:"location_code"`
	ParentLocation string `json:"parent_location"`
	Point          any    `json:"point"`
}

type WarehouseInfo interface {
}

type TransferDuration struct {
	Duration time.Duration // `json:"transfer_duration"` // must be interval type in golang
}

func (t *TransferDuration) ReadAsString() string {
	return t.Duration.String()
}

type WarehouseOutboundTracker struct {
	Id                  int                  `json:"id"`
	OutboundTransaction internals.NullString `json:"outbound_transaction"`
	InboundTransaction  internals.NullString `json:"inbound_transaction"`
	DeliveryAt          internals.NullTime   `json:"delivery_at"`
	ArrivedAt           internals.NullTime   `json:"arrived_at"`
	Files               internals.NullInt64  `json:"files"`
}

type ValidationWarehouseInbound struct{}

func (v *ValidationWarehouseInbound) CheckOutboundNumber(info_outbound_number TransactionInfo) error {

	if info_outbound_number.TransactionType != "warehouse_to_warehouse" {
		msg := fmt.Sprintf("Gagal, tipe outbound_number bukan warehouse to warehouse, tipenya adalah: %s", info_outbound_number.TransactionType)
		return errors.New(msg)
	}

	return nil

}
