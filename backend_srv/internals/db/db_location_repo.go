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
