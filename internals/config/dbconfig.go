package config

import (
	"fmt"
	"foodorderapi/internals/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB
var e error

func Databaseinit() {

	host := "localhost"
	user := "foodorderuser"
	password := "foodorderpassword"
	dbName := "foodorderdb"
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable ", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if e != nil {
		panic(e)

	}

	if e = database.AutoMigrate(
		&models.Menu{},
		&models.Order{},
		&models.Customer{},
		&models.Merchant{},
		&models.Admin{},
		&models.General{},
		&models.Category{},
	); e != nil {
		panic(e)
	}

}

func DB() *gorm.DB {
	return database
}
