package services

import (
	"jiva-guildes/adapters"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	"jiva-guildes/settings"
	"testing"
)

// TODO Try to use test suite (see:https://medium.com/nerd-for-tech/setup-and-teardown-unit-test-in-go-bd6fa1b785cd)
// TODO: https://stackoverflow.com/questions/31794141/can-i-create-shared-test-utilities
// TODO: Instead of create and destroy the database for each test, we cn use transactions and rollback to keep the database clean

func SetupTest(tb testing.TB) (ServiceManager, func(tb testing.TB)) {
	tables.InitAllTables()
	var connectionPool = db.MountDB(settings.AppSettings.DATABASE_URI)
	var UnitOfWorkManager = adapters.NewUnitOfWorkManager(connectionPool)
	var TestServiceManager = ServiceManager{UnitOfWorkManager: &UnitOfWorkManager}

	return TestServiceManager, func(tb testing.TB) {
		tables.DropAllTables()
		db.Teardown(connectionPool)
	}
}
