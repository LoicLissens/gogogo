package services

import (
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/models/guilde"
)

func (sm ServiceManager) CreateGuildeHandler(cmd commands.CreateGuildeCommand) (guilde.Guilde, error) {
	validated := true // TODO: change regarding who create it
	g, err := guilde.CreateFromCommand(cmd, validated)
	if err != nil {
		return guilde.Guilde{}, err
	}
	uow, close := sm.UnitOfWorkManager.Start()
	defer close()
	savedGuilde, err := uow.GuildeRepository().Save(*g)
	if err != nil {
		// uow.Rollback() //todo: implement rollback
		return savedGuilde, err
	}
	return savedGuilde, nil
}
