package model

import (
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

var ErrDBNotInitialized = errors.New("db not initialized")

func InitDB(databaseName string) error {
	db, err := gorm.Open(sqlite.Open(databaseName+".db?cache=shared&mode=rwc"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&TeamUser{}, &Task{}, &Report{}, &BattleReport{})
	if err != nil {
		return err
	}

	Conn = db
	return nil
}

func DB() (*gorm.DB, error) {
	if Conn == nil {
		return nil, ErrDBNotInitialized
	}
	return Conn, nil
}
