package services

import (
	"jiva-guildes/domain/ports"
)

type ServiceManager struct {
	UnitOfWorkManager ports.UnitOfWorkManager
}
