package backend

import (
	"jiva-guildes/adapters"
	"jiva-guildes/adapters/db"
	"jiva-guildes/domain/ports/views"
	"jiva-guildes/services"
	"jiva-guildes/settings"

	"github.com/go-playground/validator/v10"
)

var connectionPool = db.MountDB(settings.AppSettings.DATABASE_URI)

var UnitOfWorkManager = adapters.NewUnitOfWorkManager(connectionPool)

var ServiceManager = services.ServiceManager{UnitOfWorkManager: &UnitOfWorkManager}

var viewsManager views.ViewsManager

var Validate *validator.Validate = validator.New()
