package domain

import "context"

type UnitOfWork interface {
	BeginAuth(ctx context.Context) (Authentication, error)
	BeginSetLocation(ctx context.Context) (LocationRepository, error)
	// BeginStockTransfer(ctx context.Context) (WarehouseSrv, error)
	// BeginStockReceive(ctx context.Context) (WarehouseSrv, error)
}
