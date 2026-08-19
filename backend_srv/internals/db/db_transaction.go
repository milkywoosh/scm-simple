package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"scm-simple-luke.com/dir/internals/domain"
)

type pgTransaction struct {
	tx pgx.Tx
	domain.UserRepository
	domain.LocationRepository
}

func newPgTx(tx pgx.Tx) *pgTransaction {

	newUserRepo := NewUserRepository(tx)
	newLocationRepo := NewDBLocationRepository(tx)

	return &pgTransaction{
		tx:                 tx,
		UserRepository:     newUserRepo,
		LocationRepository: newLocationRepo,
	}
}

func (ot *pgTransaction) Commit(ctx context.Context) error {
	return ot.tx.Commit(ctx)
}

func (ot *pgTransaction) Rollback(ctx context.Context) error {
	return ot.tx.Rollback(ctx)
}

func (ot *pgTransaction) UserRepo() domain.UserRepository {
	return ot.UserRepository
}

func (ot *pgTransaction) LocationRepo() domain.LocationRepository {
	return ot.LocationRepository
}
