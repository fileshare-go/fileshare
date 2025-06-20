package service

import (
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	Orm   *gorm.DB
	Error error
}

func (s *Service) AddError(err error) error {
	if s.Error == nil {
		s.Error = err
	} else if err != nil {
		s.Error = fmt.Errorf("%v; %w", s.Error, err)
	}
	return s.Error
}
