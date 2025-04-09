package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/middlewares"
	"github.com/jerson2000/api-qfirst/routes"
)

func init() {
	config.ConfigLoadEnvironmentVariables()
	config.ConfigConnectToDatabase()
	config.ConfigJwtKey()
	config.ConfigRefreshJwtKey()
	config.ConfigCasbinEnforcer()
	config.ConfigCSRF()
	config.ConfigCacheInit()
}

func main() {
	log.Println("Starting server...")
	router := *mux.NewRouter()
	log.Println("Setting up routes...")

	// User Routes
	routes.UserRoutes(&router)
	// Auth Routes
	routes.AuthRoutes(&router)
	// services
	routes.ServicesRoutes(&router)
	// booking
	routes.BookingRoutes(&router)

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", middlewares.CacheMiddleware(&router)))
}
