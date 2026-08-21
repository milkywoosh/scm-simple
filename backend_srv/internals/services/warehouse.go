package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"scm-simple-luke.com/dir/internals"
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

type WarehouseTransferHeader struct {
	TransactionNumber string
	Status            string
}

type WarehouseReceiveDraft struct {
}

func (a *WarehouseService) CreateDraftTx(ctx context.Context, location_origin, location_destination string) (WarehouseTransferHeader, error) {
	tx, err := a.uow.BeginWarehouseToWarehouse(ctx)
	log.Printf("errr CreateDraftTx %v", err)
	n := WarehouseTransferHeader{}
	if err != nil {
		return n, err
	}

	defer tx.Rollback(ctx)

	currdate := time.Now().Format("20060102150405")
	randString := internals.RandomStringSuffix(5)
	transaction_number := fmt.Sprintf("WH_TO_WH_%s_%s", currdate, randString)

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

	n.TransactionNumber = new_transaction.TransactionNumber
	n.Status = new_transaction.Status
	return n, nil
}

func (a *WarehouseService) AllocateItem(ctx context.Context, transaction_number, identifier string) error {
	tx, err := a.uow.BeginWarehouseToWarehouse(ctx)

	if err != nil {
		log.Printf("*WarehouseService AllocateItem 1: %v", err)
		return err
	}

	defer tx.Rollback(ctx)

	TransInfo, err := tx.CheckTransaction(ctx, transaction_number)
	if err != nil {
		log.Printf("*WarehouseService CheckTransaction 1: %v", err)
		return err
	}

	curr_status := TransInfo.Status
	if curr_status != "draft" {
		log.Printf("*WarehouseService CheckTransaction 1: %v", err)
		msg := fmt.Sprintf("error status saat ini bukan draft! saat ini %v", curr_status)
		return errors.New(msg)
	}

	err = tx.AllocateItem(ctx, transaction_number, identifier)
	if err != nil {
		log.Printf("*WarehouseService AllocateItem 1: %v", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("errr AllocateItem Commit %v", err)
		return err
	}

	return nil
}

func (a *WarehouseService) SetSubmit(ctx context.Context, transaction_number string) error {

	tx, err := a.uow.BeginWarehouseToWarehouse(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	TransactionInfo, err := tx.CheckTransaction(ctx, transaction_number)
	if err != nil {
		return err
	}

	CurrStatus := TransactionInfo.Status
	if CurrStatus != "draft" {
		msg := fmt.Sprintf("status transaksi saat ini bukan draft tapi %s", CurrStatus)
		return errors.New(msg)
	}

	err = tx.SetStatusTransaction(ctx, transaction_number, "submitted")
	if err != nil {
		log.Printf("err SetStatusTransaction submit: %s", err.Error())
		log.Printf("err SetStatusTransaction submit: %s", err.Error())
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("err SetStatusTransaction submit commit 1: %s", err.Error())
		return err
	}
	return nil
}
