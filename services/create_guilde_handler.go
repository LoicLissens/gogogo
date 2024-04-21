package services

import (
	"jiva-guildes/domain/models/guilde"
)

func (sm ServiceManager) CreateGuildeHandler() guilde.Guilde {
	guilde := guilde.New("test", "test", "test")
	uow := sm.UnitOfWorkManager.Start()
	savedGuilde, _ := uow.GuildeRepository().Save(*guilde)
	return savedGuilde
}
