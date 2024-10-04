package services

import (
	"jiva-guildes/domain/commands"
)

func (sm ServiceManager) DeleteGuildeHandler(cmd commands.DeleteGuildeCommand) error {
	uow, close := sm.UnitOfWorkManager.Start()
	defer close()
	err := uow.GuildeRepository().Delete(cmd.Uuid)
	if err != nil {
		// uow.Rollback() //todo: implement rollback
		return err
	}
	return nil
}
