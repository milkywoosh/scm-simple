package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"scm-simple-luke.com/dir/internals/domain"
)

type DBItemRepository struct {
	Conn DBTX
}

func NewDBItemRepository(db DBTX) domain.ItemRepository {
	return &DBItemRepository{
		Conn: db,
	}
}

func (d *DBItemRepository) GetItem(ctx context.Context, identifier string) ([]domain.ItemInfo, error) {

	query := `
		select 
			i.id,
			i.serial_number,
			i.factory_serial_number,
			i.created_at,
			i.curr_status,
			i.curr_transaction,
			i.curr_location_code,
			i.product_code,
			i.introduction_number
		from items i
			where i.serial_number = $1
	`
	rows, err := d.Conn.Query(ctx, query, identifier)
	if err != nil {
		log.Printf("cek GetItem repo: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			msg := fmt.Sprintf("kode identifier %s tidak ditemukan", identifier)
			return []domain.ItemInfo{}, errors.New(msg)
		}
		return []domain.ItemInfo{}, err
	}

	defer rows.Close()

	var datas []domain.ItemInfo
	for rows.Next() {
		var data domain.ItemInfo

		err := rows.Scan(
			&data.Id,
			&data.Serial_number,
			&data.Factory_serial_number,
			&data.Created_at,
			&data.Curr_status,
			&data.Curr_transaction,
			&data.Curr_location_code,
			&data.Product_code,
			&data.Introduction_number,
		)
		if err != nil {
			return []domain.ItemInfo{}, err
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (d *DBItemRepository) GetItemsOnTransaction(ctx context.Context, transaction_number string) ([]domain.EachItemTransaction, error) {

	query := `
			select * from transaction_item_transfer_details t 
			where t.id_trans_item_transfer = (
				select id 
					from transaction_item_transfers t1 
				where 
					t1.transaction_number  = $1
			)
	`
	rows, err := d.Conn.Query(ctx, query, transaction_number)
	if err != nil {
		log.Printf("err query GetItemsOnTransaction 1: %s", err.Error())
		return []domain.EachItemTransaction{}, err
	}
	defer rows.Close()

	datas := []domain.EachItemTransaction{}

	for rows.Next() {

		var data domain.EachItemTransaction
		if err := rows.Scan(
			&data.Id,
			&data.IdTransfer,
			&data.SerialNumber,
			&data.AddedAt,
		); err != nil {
			log.Printf("err loop list items: %s", err.Error())
			return []domain.EachItemTransaction{}, err
		}

		datas = append(datas, data)
	}
	return datas, nil
}

func (d *DBItemRepository) UnlockItems(ctx context.Context, transaction_number string) error {

	// fetch transaction_item_transfer_details where (select id from transaction where trans_number = transnumber)
	// loop update items
	//	set curr_transaction = NULL
	// where trans_number = transnumber_param

	datas, err := d.GetItemsOnTransaction(ctx, transaction_number)
	if err != nil {
		return err
	}

	for _, val := range datas {

		query2 := `
			UPDATE items
				SET curr_transaction = NULL
			WHERE curr_transaction = $1
				and serial_number = $2
		`
		_, err := d.Conn.Exec(ctx, query2, transaction_number, val.SerialNumber)
		if err != nil {
			log.Printf("err unlock items 1: %s", err.Error())
			return err
		}
	}

	return nil

}

func (d *DBItemRepository) AllocateItem(ctx context.Context, transaction_number, identifier string) error {
	// update items allocated
	// insert detail history
	// insert detail transaction

	// update items allocated
	// insert detail history
	// insert detail transaction

	query1 := `
		insert into item_histories (
			item_id,
			transaction_number,
			previous_status,
			status_history,
			created_at
		) values (
			(select id from items i where i.serial_number = $1), 
			$2, 
			(select curr_status from items i where i.serial_number = $1), 
			$3 , 
			CURRENT_TIMESTAMP
		)
	`
	_, err := d.Conn.Exec(ctx, query1, identifier, transaction_number, "ALLOCATED")
	if err != nil {
		log.Printf("err insert into item_histories: %v", err)
		return err
	}

	query2 := `
		update items
			set curr_status = 'ALLOCATED',
				curr_transaction = $2
		where serial_number = $1
			and curr_transaction is null
		returning id
	`
	row := d.Conn.QueryRow(ctx, query2, identifier, transaction_number)

	var item_id int
	err = row.Scan(&item_id)
	if err != nil {
		log.Printf("err update AllocateItem: %v", err)
		return err
	}

	TransInfo, err := d.CheckTransaction(ctx, transaction_number)
	if err != nil {
		log.Printf("err insert into CheckTransaction Allocate item, %v", err)
		return err
	}

	query3 := `
		INSERT INTO transaction_item_transfer_details (
			id_trans_item_transfer,
			identifier_item,
			added_at
		) VALUES (
			$1, $2, CURRENT_TIMESTAMP
		)
	`
	_, err = d.Conn.Exec(ctx, query3, TransInfo.Id, identifier)
	if err != nil {
		log.Printf("err insert into item_histories: %v", err)
		return err
	}

	return nil

}

func (d *DBItemRepository) CheckTransaction(ctx context.Context, transaction_number string) (domain.TransactionInfo, error) {
	query := `
		select
			id,
			transaction_number,
			status,
			transaction_type,
			origin,
			destination,
			created_at,
			submitted_at,
			approved_at
		from transaction_item_transfers t
			where t.transaction_number = $1
	`

	var data domain.TransactionInfo

	row := d.Conn.QueryRow(ctx, query, transaction_number)
	if err := row.Scan(
		&data.Id,
		&data.TransactionNumber,
		&data.Status,
		&data.TransactionType,
		&data.Origin,
		&data.Destination,
		&data.CreatedAt,
		&data.SubmittedAt,
		&data.ApprovedAt,
	); err != nil {
		log.Printf("cek CheckTransaction repo 1: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			msg := fmt.Sprintf("transaction_number %s tidak ditemukan", transaction_number)
			return data, errors.New(msg)
		}
		return data, err
	}

	return data, nil

}

func (d *DBItemRepository) ReceiveInboundItem(ctx context.Context, transaction_number, identifier string) error {

	query1 := `
		insert into item_histories (
			item_id,
			transaction_number,
			previous_status,
			status_history,
			created_at
		) values (
			(select id from items i where i.serial_number = $1), 
			$2, 
			(select curr_status from items i where i.serial_number = $1), 
			$3 , 
			CURRENT_TIMESTAMP
		)
	`
	_, err := d.Conn.Exec(ctx, query1, identifier, transaction_number, "AVAILABLE")
	if err != nil {
		log.Printf("err insert into item_histories: %v", err)
		return err
	}

	query2 := `
		update items
			set curr_status = 'AVAILABLE',
				curr_transaction = $2,
				curr_location_code = (
					select y.destination from 
						transaction_item_transfers y
					where y.transaction_number = $2
				)
		where serial_number = $1
			and curr_transaction = (
				select outbound_transaction
					from transaction_transfer_tracker x 
				where x.inbound_transaction = $2
			)
		returning id
	`
	row := d.Conn.QueryRow(ctx, query2, identifier, transaction_number)

	var item_id int
	err = row.Scan(&item_id)
	if err != nil {
		log.Printf("err update ReceiveInboundItem : %v", err)
		return err
	}

	TransInfo, err := d.CheckTransaction(ctx, transaction_number)
	if err != nil {
		log.Printf("err insert into ReceiveInboundItem 2, %v", err)
		return err
	}

	query3 := `
		INSERT INTO transaction_item_transfer_details (
			id_trans_item_transfer,
			identifier_item,
			added_at
		) VALUES (
			$1, $2, NOW()
		)
	`
	_, err = d.Conn.Exec(ctx, query3, TransInfo.Id, identifier)
	if err != nil {
		log.Printf("err insert into item_histories: %v", err)
		return err
	}

	return nil

}

func (d *DBItemRepository) DisAllocateItem(ctx context.Context, transaction_number, identifier string) error {

	query2 := `
		update items
			set curr_status = (
					select 
						previous_status 
					from item_histories ih 
						where ih.item_id = (select id from items where serial_number = $1) 
							and (ih.transaction_number = $2)
				),
				curr_transaction = null
		where serial_number = $1
			and curr_transaction = $2
		returning id
	`
	row := d.Conn.QueryRow(ctx, query2, identifier, transaction_number)

	var item_id int
	err := row.Scan(&item_id)
	if err != nil {
		log.Printf("err update disAllocateItem: %v", err)
		return err
	}

	query1 := `
		delete from item_histories 
			where item_id = (select id from items where serial_number = $1) 
				and (transaction_number = $2)
	`
	_, err = d.Conn.Exec(ctx, query1, identifier, transaction_number)
	if err != nil {
		log.Printf("err delete from item_histories disallocate: %v", err)
		return err
	}


	TransInfo, err := d.CheckTransaction(ctx, transaction_number)
	if err != nil {
		log.Printf("err insert into CheckTransaction Allocate item, %v", err)
		return err
	}

	query3 := `
		delete from transaction_item_transfer_details 
			where (id_trans_item_transfer = $1)
				and (identifier_item = $2)
	`
	_, err = d.Conn.Exec(ctx, query3, TransInfo.Id, identifier)
	if err != nil {
		log.Printf("err insert into item_histories: %v", err)
		return err
	}

	return nil
}
