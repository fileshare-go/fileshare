package service

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/pkg/db"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
	"gorm.io/gorm"
)

var srv service.Service

func InitClientService() {
	cfg := config.Cfg()

	orm := db.OpenClientDB(cfg.Database)

	srv = service.Service{
		Orm:   orm,
		Error: nil,
	}
}

func InitServerService() {
	cfg := config.Cfg()

	orm := db.OpenServerDB(cfg.Database)

	srv = service.Service{
		Orm:   orm,
		Error: nil,
	}
}

func Orm() *gorm.DB {
	return srv.Orm
}
