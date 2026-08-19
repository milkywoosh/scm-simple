package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"scm-simple-luke.com/dir/internals/domain"
)

type pgQuery struct {
	db           *pgxpool.Pool
	userRepo     domain.UserRepository
	locationRepo domain.LocationRepository
}

func newPgQuery(db *pgxpool.Pool) *pgQuery {

	newUserRepo := NewUserRepository(db)
	newLocationRepo := NewDBLocationRepository(db)

	return &pgQuery{
		userRepo:     newUserRepo,
		locationRepo: newLocationRepo,
	}
}

func (ot *pgQuery) AuthRepo() domain.UserRepository {
	return ot.userRepo
}

func (ot *pgQuery) WarehouseRepo() domain.LocationRepository {
	return ot.locationRepo
}
