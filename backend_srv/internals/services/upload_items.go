package services

import (
	"context"
	"errors"

	"scm-simple-luke.com/dir/internals/domain"
)

func (w *WarehouseService) UploadBulkItems(ctx context.Context, batch_size int, rows []domain.ItemInfo) error {
	conn, err := w.q.BeginUpload(ctx)
	if err != nil {
		return err
	}

	if batch_size < 1 || len(rows) < 1 {
		return errors.New("batch size dan jumlah item harus mencukupi")
	}

	err = conn.UploadBulkNewItems(ctx, batch_size, rows)
	if err != nil {
		return err
	}

	return nil
}
