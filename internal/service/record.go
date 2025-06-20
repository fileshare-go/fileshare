package service

import (
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
)

type Record struct {
	service.Service
}

func (r Record) Get() ([]model.Record, error) {
	var recordList []model.Record
	db := r.Orm.Find(&recordList)
	return recordList, db.Error
}

func (r Record) Insert(record model.Record) error {
	db := r.Orm.Create(&record)
	if db.RowsAffected != 1 {
		return notAffectedError
	}
	return db.Error
}
