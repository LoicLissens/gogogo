package adapters

import (
	"jiva-guildes/adapters/db/repositories"
	"jiva-guildes/domain/ports"
	portrepo "jiva-guildes/domain/ports/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
	conn *pgxpool.Pool
}

func (uow UnitOfWork) GuildeRepository() portrepo.GuildeRepository {
	if uow.conn == nil {
		panic("Connection pool not set")
	}
	g := repositories.NewGuildeRepository(uow.conn)
	return &g
}

type UnitOfWorkManager struct {
	conn *pgxpool.Pool
}

func NewUnitOfWorkManager(connectionPool *pgxpool.Pool) UnitOfWorkManager {
	return UnitOfWorkManager{conn: connectionPool}
}
func (uowm *UnitOfWorkManager) Setup(connectionPool *pgxpool.Pool) { //TODO: Use it when start and shutdown event are implemented
	uowm.conn = connectionPool
}

func (uowm UnitOfWorkManager) Start() ports.UnitOfWork {
	if uowm.conn == nil {
		panic("Connection pool not set")
	}
	uow := UnitOfWork{conn: uowm.conn}
	return &uow
}
