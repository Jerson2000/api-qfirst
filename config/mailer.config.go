package config

import (
	"log"
	"os"
)

var (
	Mailer_Email    string
	Mailer_Password string
)

func ConfigMailerInit() {
	Mailer_Email = os.Getenv("EMAIL")
	if Mailer_Email == "" {
		log.Fatal("env EMAIL is not set")
	}

	Mailer_Password = os.Getenv("PASSWORD")
	if Mailer_Password == "" {
		log.Fatal("env PASSWORD is not set")
	}
}
