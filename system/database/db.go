// database/db.go
package database

import "gorm.io/gorm"

var db *gorm.DB

func SetDB(d *gorm.DB) {
	db = d
}

func GetDB() *gorm.DB {
	return db
}
