package backend

import (
	"jiva-guildes/adapters"
	"jiva-guildes/adapters/db"
	"jiva-guildes/domain/ports"
	"jiva-guildes/settings"
	"net/url"
)

var path, _ = url.Parse(settings.AppSettings.DATABASE_URI)
var connectionPool = db.MountDB(path)

var UnitOfWorkManager = adapters.NewUnitOfWorkManager(connectionPool)
var ServiceManager = ports.ServiceManager{UnitOfWorkManager: &UnitOfWorkManager}
