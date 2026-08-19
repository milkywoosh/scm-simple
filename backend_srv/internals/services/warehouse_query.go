package services

import (
	"context"
	"log"

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
