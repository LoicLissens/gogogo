package services

import (
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/models/guilde"
)

func (sm ServiceManager) UpdateGuildeHandler(cmd commands.UpdateGuildeCommand) (guilde.Guilde, error) {

	uow, close := sm.UnitOfWorkManager.Start()
	defer close()
	guide, err := uow.GuildeRepository().GetByUUID(cmd.Uuid)
	if err != nil {
		return guilde.Guilde{}, err
	}
	err = guide.UpdateFromCommand(cmd)
	if err != nil {
		return guilde.Guilde{}, err
	}
	updatedGuide, err := uow.GuildeRepository().Update(guide)
	if err != nil {
		return guilde.Guilde{}, err
	}
	return updatedGuide, nil
}
