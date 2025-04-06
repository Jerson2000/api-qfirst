package config

import (
	"log"
	"os"
)

var JWTKey []byte
var RefreshJWTKey []byte

func ConfigJwtKey() {
	JWTKey := []byte(os.Getenv(("JWT_SECRET")))
	if len(JWTKey) == 0 {
		log.Fatal("env jwt secret is not set")
	}

}

func ConfigRefreshJwtKey() {
	RefreshJWTKey := []byte(os.Getenv(("REFRESH_JWT_SECRET")))
	if len(RefreshJWTKey) == 0 {
		log.Fatal("Env refresh jwt secret is not set")
	}
}
