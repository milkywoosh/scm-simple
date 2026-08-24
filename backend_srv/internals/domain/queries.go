package domain

import "context"

type Queries interface {
	// StockTransferQueries(ctx context.Context) (ItemTransfer, error)
	AuthQueries(ctx context.Context) (UserRepository, error)
	WarehouseQueries(ctx context.Context) (LocationRepository, error)
	ItemQueries(ctx context.Context) (ItemRepository, error)
}
