package db

import "gorm.io/gorm"

func (or *OrmRepository) GetDB() *gorm.DB {
	return or.DB
}

func (or *OrmRepository) Transaction(f func(tx *gorm.DB) error) error {
	return or.GetDB().Transaction(f)
}
