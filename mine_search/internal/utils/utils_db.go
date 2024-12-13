package utils

import (
	"errors"

	"gorm.io/gorm"
)

func CheckDBConnection(db *gorm.DB) (bool, error) {
	if db == nil {
		return false, errors.New("Can't connect database")
	}
	rawDB, err := db.DB()
	if rawDB == nil || err != nil {
		return false, err
	}
	if err := rawDB.Ping(); err != nil {
		return false, err
	}
	return true, nil
}
