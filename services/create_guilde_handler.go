package services

import (
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/models/guilde"
)

func (sm ServiceManager) CreateGuildeHandler(cmd commands.CreateGuildeCommand) (guilde.Guilde, error) {

	guilde := guilde.New(cmd.Name, cmd.Img_url, cmd.Page_url)
	uow := sm.UnitOfWorkManager.Start()
	savedGuilde, _ := uow.GuildeRepository().Save(*guilde)
	return savedGuilde, nil
}
