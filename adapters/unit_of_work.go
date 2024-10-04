package adapters

import (
	"context"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/repositories"
	"jiva-guildes/domain/ports"
	portrepo "jiva-guildes/domain/ports/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
	db pgx.Tx
}

func (uow UnitOfWork) GuildeRepository() portrepo.GuildeRepository {
	if uow.db == nil {
		panic("Connection pool not set")
	}
	g := repositories.NewGuildeRepository(uow.db)
	return &g
}

type UnitOfWorkManager struct {
	db db.PsqlDB
}

func NewUnitOfWorkManager(db db.PsqlDB) UnitOfWorkManager {
	return UnitOfWorkManager{db: db}
}
func (uowm *UnitOfWorkManager) Setup(connectionPool *pgxpool.Pool) { //TODO: Use it when start and shutdown event are implemented
	uowm.db = connectionPool
}

func (uowm UnitOfWorkManager) Start() (ports.UnitOfWork, func()) {
	if uowm.db == nil {
		panic("Connection pool not set")
	}
	tx, err := uowm.db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	uow := UnitOfWork{db: tx}
	return &uow, func() {
		tx.Commit(context.Background())
	}
}
