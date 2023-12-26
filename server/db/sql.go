package db

import (
	"Dur4nC2/server/domain/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// newDBClient - Initialize the db client

func newMySQLDBClient() *gorm.DB {
	//dbConfig := configs.GetDatabaseConfig()

	var dbClient *gorm.DB
	dbClient = mySQLClient()
	dbClient = &gorm.DB{}

	err := dbClient.AutoMigrate(
		&models.Operator{},
		&models.Host{},
		&models.Beacon{},
		&models.BeaconTask{},
		&models.IOC{},
		&models.Loot{},
	)
	if err != nil {
		fmt.Println(err)
	}

	//return dbClient
	return &gorm.DB{}
}

func mySQLClient() *gorm.DB {
	//dsn, err := dbConfig.DSN()
	//if err != nil {
	//	panic(err)
	//}
	var dsn = "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true"
	dbClient, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return dbClient
}
