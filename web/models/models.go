package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/types"
)

var db *gorm.DB

//SetDB establishes connection to database and saves its handler into db *sqlx.DB
func SetDB(dbConfig *config.DatabaseConfig) {
	gormConfig := gorm.Config{}
	var err error
	switch dbConfig.Type {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dbConfig.GetConnectionString()), &gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbConfig.SqliteDefaultPath()), &gormConfig)
	}
	if err != nil {
		panic(err)
	}
}

//GetDB returns database handler
func GetDB() *gorm.DB {
	return db
}

//AutoMigrate runs gorm auto migration
func AutoMigrate(providers []*types.Provider) error {
	var err error
	db.AutoMigrate(&Task{}, &types.Resource{}, &types.ResourceManager{})
	err = CreateAllResourceManagers(providers)
	return err
}
