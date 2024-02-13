package database

import (
	"fmt"
	"log"

	"mygram-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "10051998"
	dbPort   = "5432"
	dbname   = "simple-api"
	db       *gorm.DB
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, dbPort)
	dsn := config
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database;", err)
	}

	log.Println("sukses koneksi ke database")
	db.Debug().AutoMigrate(
		models.User{},
		models.Photo{},
		models.Comment{},
		models.SocialMedia{},
	)
}

func GetDB() *gorm.DB {
	return db
}
