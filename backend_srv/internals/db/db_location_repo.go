package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"scm-simple-luke.com/dir/internals/domain"
)

type DBLocationRepository struct {
	Conn DBTX
}

func NewDBLocationRepository(db DBTX) domain.LocationRepository {
	return &DBLocationRepository{
		Conn: db,
	}
}

func (d *DBLocationRepository) GetLocation(ctx context.Context, locationCode string, locationType domain.LocationType) ([]domain.LocationRow, error) {

	if locationType == domain.WarehouseType {
		query := `select location_code, description, parent_location_code, created_at, updated_at, location_point from warehouses w where w.location_code = $1`
		rows, err := d.Conn.Query(ctx, query, locationCode)
		if err != nil {
			log.Printf("cek GetLocation repo: %v", err)
			if errors.Is(err, pgx.ErrNoRows) {
				msg := fmt.Sprintf("kode lokasi %s tidak ditemukan", locationCode)
				return []domain.LocationRow{}, errors.New(msg)
			}
			return nil, err
		}

		defer rows.Close()

		var datas []domain.LocationRow
		for rows.Next() {
			var data domain.LocationRow

			err := rows.Scan(
				&data.LocationCode,
				&data.Description,
				&data.ParentLocation,
				&data.CreatedAt,
				&data.UpdatedAt,
				&data.Point,
			)
			if err != nil {
				return nil, err
			}

			datas = append(datas, data)
		}
		return datas, nil
	}

	msg := fmt.Sprintf("tipe lokasi %s belum dapat diakses", locationType)
	return nil, errors.New(msg)

}

func (d *DBLocationRepository) NewDraftTransaction(ctx context.Context, transaction_number, transaction_type, origin, destination string) (domain.TransactionInfo, error) {

	query :=
		`INSERT INTO transaction_item_transfers (
			transaction_number, 
			status,
			transaction_type,
			origin, 
			destination
		) VALUES (
		 	$1, $2, $3, $4
		)`

	_, err := d.Conn.Exec(
		ctx,
		query,
		transaction_number,
		"draft",
		transaction_type,
		origin,
		destination,
	)
	if err != nil {
		return domain.TransactionInfo{}, err
	}

	n := domain.TransactionInfo{}
	n.TransactionNumber = transaction_number
	n.TransactionType = transaction_type

	return n, nil

}

func (d *DBLocationRepository) SetStatusTransaction(ctx context.Context, transaction_number, transaction_type, origin, destination string) (domain.TransactionInfo, error) {

	query :=
		`INSERT INTO transaction_item_transfers (
			transaction_number, 
			transaction_type,
			origin, 
			destination
		) VALUES (
		 	$1, $2, $3, $4
		)`

	_, err := d.Conn.Exec(
		ctx,
		query,
		transaction_number,
		transaction_type,
		origin,
		destination,
	)
	if err != nil {
		return domain.TransactionInfo{}, err
	}

	n := domain.TransactionInfo{}
	n.TransactionNumber = transaction_number
	n.TransactionType = transaction_type

	return n, nil

}

func (d *DBLocationRepository) AllocateItem(ctx context.Context, transaction_number, item string) error {
	// update items allocated
	// insert detail history
	// insert detail transaction

	query1 := `
		update items i
			set i.curr_status = 'ALLOCATED'
		where i.serial_number = $1
			and i.curr_transaction is null
		returning i.id
	`
	row := d.Conn.QueryRow(ctx, query1, item)

	var item_id int
	err := row.Scan(&item_id)
	if err != nil {
		log.Printf("err update allocated: %v", err)
		return err
	}

	query2 := `
		insert into item_histories (
			item_id,
			transaction_number,
			status_history,
			created_at
		) values (
			$1, $2, $3, CURRENT_TIMESTAMP
		)
	`
	_, err = d.Conn.Exec(ctx, query2, item, transaction_number, "ALLOCATED")
	log.Printf("err insert into item_histories")
	if err != nil {
		return err
	}

	TransInfo, err := d.CheckTransaction(ctx, transaction_number)
	log.Printf("err insert into CheckTransaction Allocate item")
	if err != nil {
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
	_, err = d.Conn.Exec(ctx, query3, TransInfo.Id, item)
	log.Printf("err insert into item_histories")
	if err != nil {
		return err
	}

	return nil

}

func (d *DBLocationRepository) ReceiveItem(ctx context.Context, transaction_number, from string, to string, item string) error {
	return nil
}

func (d *DBLocationRepository) CheckTransaction(ctx context.Context, transaction_number string) (domain.TransactionInfo, error) {
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

	data := domain.TransactionInfo{}

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
