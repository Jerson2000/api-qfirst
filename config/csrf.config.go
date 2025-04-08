package config

import (
	"log"
	"os"
)

var CSRFKey []byte

func ConfigCSRF() {

	CSRFKey = []byte(os.Getenv("CSRF_KEY"))
	if len(CSRFKey) == 0 {
		log.Fatal("env csrf key is not set")
	}
}
