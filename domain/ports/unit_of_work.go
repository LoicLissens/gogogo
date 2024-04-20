package ports

import (
	"jiva-guildes/domain/ports/repositories"
)

type UnitOfWork interface {
	GuildeRepository(connectionPool interface{}) *repositories.GuildeRepository
}

type UnitOfWorkManager interface {
	Start() *UnitOfWork
}
