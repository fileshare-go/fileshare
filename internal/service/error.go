package service

import (
	"errors"

	"gorm.io/gorm"
)

var notAffectedError = errors.New("No rows affected")

func oneRowAffected(db *gorm.DB) error {
	if db.RowsAffected != 1 {
		return notAffectedError
	}
	return db.Error
}
