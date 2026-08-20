package services

import (
	"scm-simple-luke.com/dir/internals/db"
	"scm-simple-luke.com/dir/internals/domain"
)

type AuthenticationService struct {
	q   domain.Queries
	uow domain.UnitOfWork
}

type WarehouseService struct {
	q   domain.Queries
	uow domain.UnitOfWork
}

type Services struct {
	*AuthenticationService
	*WarehouseService
}

func NewServices(db *db.PgDBInstance) *Services {
	return &Services{
		AuthenticationService: NewAuthentication(db),
		WarehouseService:      NewWarehouseService(db),
	}
}


func NewAuthentication(dbpgInstance *db.PgDBInstance) *AuthenticationService {
	return &AuthenticationService{
		q:   dbpgInstance,
		uow: dbpgInstance,
	}
}

func NewWarehouseService(dbpgInstance *db.PgDBInstance) *WarehouseService {
	return &WarehouseService{
		q:   dbpgInstance,
		uow: dbpgInstance,
	}
}
