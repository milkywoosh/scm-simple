package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"scm-simple-luke.com/dir/internals/domain"
)

func (a *WarehouseService) WarehouseInfo(ctx context.Context, location_code string) ([]domain.LocationRow, error) {
	locationRepo, err := a.q.WarehouseQueries(ctx)
	log.Printf("errr WarehouseInfo %v", err)
	if err != nil {
		return nil, err
	}
	return locationRepo.GetLocation(ctx, location_code, domain.WarehouseType)
}

type WarehouseTransferDraft struct {
	TransactionNumber string
}

type WarehouseReceiveDraft struct {
}

func (a *WarehouseService) CreateDraftTx(ctx context.Context, location_origin, location_destination string) (WarehouseTransferDraft, error) {
	tx, err := a.uow.BeginWarehouseToWarehouse(ctx)
	log.Printf("errr CreateDraftTx %v", err)
	n := WarehouseTransferDraft{}
	if err != nil {
		return n, err
	}

	currdate := time.DateTime
	transaction_number := fmt.Sprintf("wh_to_wh_%v", currdate)

	new_transaction, err := tx.NewDraftTransaction(ctx, transaction_number, "warehouse_to_warehouse", location_origin, location_destination)
	log.Printf("errr CreateDraftTx 1 %v", err)

	if err != nil {
		if err != nil {
			return n, err
		}
	}

	err = tx.Commit(ctx)
	log.Printf("errr CreateDraftTx Commit %v", err)

	if err != nil {
		return n, err
	}

	n.TransactionNumber = new_transaction.Transaction_number
	return n, nil
}

func (a *WarehouseService) SendItem(ctx context.Context, transaction_number, item string) error {
	_, err := a.uow.BeginWarehouseToWarehouse(ctx)
	return err

}
