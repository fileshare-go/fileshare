package service

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/pkg/db"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
	"gorm.io/gorm"
)

var mgr ServiceMgr
var srv service.Service

func InitServiceMgr() {
	cfg := config.Cfg()

	orm := db.OpenClientDB(cfg.Database)

	srv = service.Service{
		Orm:   orm,
		Error: nil,
	}
	mgr = ServiceMgr{
		FileInfo:  FileInfo{Service: srv},
		ShareLink: ShareLink{Service: srv},
		Record:    Record{Service: srv},
	}
}

type ServiceMgr struct {
	FileInfo
	ShareLink
	Record
}

func Mgr() *ServiceMgr {
	return &mgr
}

func Orm() *gorm.DB {
	return srv.Orm
}
