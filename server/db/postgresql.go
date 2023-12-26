package db

import (
	"Dur4nC2/server/domain/models"
	"fmt"
	"gorm.io/gorm/logger"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	user = "postgres"
	pass = "1234"
	host = "127.0.0.1"
)

func newDBClient() *gorm.DB {
	var dbClient *gorm.DB
	dbClient = postgresqlClient()

	err := dbClient.AutoMigrate(
		&models.Operator{},
		&models.Host{},
		&models.Beacon{},
		&models.BeaconTask{},
		&models.IOC{},
		&models.Loot{},
		&models.ImplantConfig{},
		&models.Server{},
	)
	if err != nil {
		fmt.Println(err)
	}

	return dbClient
}

func postgresqlClient() *gorm.DB {
	dsn := "host=" + host + " user=" + user + " password=" + pass + " dbname=c2 port=5432 sslmode=disable"
	dbClient, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	if err != nil {
		fmt.Println("Could not connect to postgresql")
		syscall.Exit(2)
	}
	return dbClient
}
