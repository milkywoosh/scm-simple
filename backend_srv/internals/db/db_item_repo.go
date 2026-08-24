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
