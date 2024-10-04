package ports

import (
	"jiva-guildes/domain/ports/repositories"
)

type UnitOfWork interface {
	GuildeRepository() repositories.GuildeRepository
}

type UnitOfWorkManager interface {
	Start() (UnitOfWork, func())
}
