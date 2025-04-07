package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConfigConnectToDatabase() {
	var err error

	dsn := os.Getenv("DB_URL_STRING")
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	if err != nil {
		log.Fatal(err)
	}

}
