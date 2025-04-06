package main

import (
	"fmt"
	"log"

	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
)

func init() {
	config.ConfigLoadEnvironmentVariables()
	config.ConfigConnectToDatabase()
}

func main() {
	err := migrateModels()
	if err != nil {
		log.Fatal(err)
	}
}

func migrateModels() error {
	err := config.Database.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}
