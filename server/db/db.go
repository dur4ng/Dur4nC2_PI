package db

import "gorm.io/gorm"

// Client - Database Client
var Client = newDBClient()

// Session - Database session
func Session() *gorm.DB {
	return Client.Session(&gorm.Session{
		FullSaveAssociations: true,
	})
}
