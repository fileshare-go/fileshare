package service

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/pkg/db"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
	"gorm.io/gorm"
)

var srv service.Service

func InitService() {
	cfg := config.Cfg()

	orm := db.OpenClientDB(cfg.Database)

	srv = service.Service{
		Orm:   orm,
		Error: nil,
	}
}

func Orm() *gorm.DB {
	return srv.Orm
}
