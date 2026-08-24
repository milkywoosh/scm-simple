package db

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"scm-simple-luke.com/dir/internals/domain"
)

// satisfied by type struct pgxpool.Pool **read: https://github.com/jackc/pgx/wiki/Getting-started-with-pgx
type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)              // return rows
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row                     // return atmost 1 rows cannot >1
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) // exec query without NO return any row
	// Prepare(ctx context.Context, query string) (*sql.Stmt, error) // pgxpool.Pool have no Prepare() signature
}

type Conn struct {
	DBPool *pgxpool.Pool
}

func NewConn(ctx context.Context, connString string) (*Conn, error) {

	pgxPoolInit, errPgxPool := pgxpool.New(ctx, connString)
	if errPgxPool != nil {
		// log.Fatalf("err pgx pool init: %v", errPgxPool)
		return nil, errPgxPool
	}

	errPing := pgxPoolInit.Ping(ctx)
	if errPing != nil {
		log.Printf("error Ping DB: %s\n", errPing.Error())
		return nil, errPing
	}

	return &Conn{
		DBPool: pgxPoolInit,
	}, nil
}

// implement UnitOfWork, Queries Interface
type PgDBInstance struct {
	db *pgxpool.Pool
}

func NewPgInstance(db *pgxpool.Pool) *PgDBInstance {
	return &PgDBInstance{
		db: db,
	}
}

func (or *PgDBInstance) beginQuery() (*pgQuery, error) {
	if or.db == nil {
		return nil, errors.New("DB Instance doesnt exists")
	}
	return newPgQuery(or.db), nil
}

// implement UnitOfWork Interface
func (or *PgDBInstance) beginTx(ctx context.Context) (*pgTransaction, error) {

	tx, err := or.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})

	return newPgTx(tx), err

}

// implement UnitOfWork Interface
func (or *PgDBInstance) BeginAuth(ctx context.Context) (domain.Authentication, error) {

	return or.beginTx(ctx)

}

func (or *PgDBInstance) BeginWarehouseToWarehouse(ctx context.Context) (domain.ItemTransfer, error) {
	// var itemTf domain.ItemTransfer
	return or.beginTx(ctx)

}

func (or *PgDBInstance) BeginSetLocation(ctx context.Context) (domain.LocationRepository, error) {
	pgTx, err := or.beginTx(ctx)
	return pgTx.LocationRepo(), err
}

/*

// implement UnitOfWork Interface
func (or *PgDBInstance) BeginStockTransfer(ctx context.Context) (domain.WarehouseSrv, error) {
	return or.begin(ctx)
}

// work it later
func (or *PgDBInstance) BeginStockReceive(ctx context.Context) (domain.WarehouseSrv, error) {
	return or.begin(ctx)
}

func (or *PgDBInstance) StockTransferQueries(ctx context.Context) (domain.ItemTransfer, error) {
	pgQuery, err := or.beginQuery()
	if err != nil {
		return nil, err
	}

	return pgQuery.itemTransfer, nil
}
*/

func (or *PgDBInstance) AuthQueries(ctx context.Context) (domain.UserRepository, error) {
	pgQuery, err := or.beginQuery()
	if err != nil {
		return nil, err
	}

	return pgQuery.AuthRepo(), nil
}
func (or *PgDBInstance) WarehouseQueries(ctx context.Context) (domain.LocationRepository, error) {
	pgQuery, err := or.beginQuery()
	if err != nil {
		return nil, err
	}

	return pgQuery.WarehouseRepo(), nil
}

func (or *PgDBInstance) ItemQueries(ctx context.Context) (domain.ItemRepository, error) {
	pgQuery, err := or.beginQuery()
	if err != nil {
		return nil, err
	}

	return pgQuery.itemRepo, nil
}
