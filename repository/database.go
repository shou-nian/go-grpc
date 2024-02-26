package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	connection *gorm.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{connection: db}, nil
}

func (db *Database) Close() error {
	sqlDB, err := db.connection.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
