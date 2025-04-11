package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	_ "github.com/jerson2000/api-qfirst/docs"
	"github.com/jerson2000/api-qfirst/middlewares"
	"github.com/jerson2000/api-qfirst/routes"
	"github.com/jerson2000/api-qfirst/websocket"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

// @title QFirst API Documentation
// @version 1.0
// @description QFirst is a robust and easy-to-use booking system designed to help businesses manage appointments and reservations efficiently.

// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

// @securityDefinitions.apikey CsrfToken
// @in header
// @name X-CSRF-Token
// @description CSRF protection token.

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

	router.StrictSlash(true)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.UIConfig(map[string]string{
			"showExtensions": "true",
		}),
		httpSwagger.PersistAuthorization(true),
	))

	// Serve swagger.json file
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// websocket
	router.HandleFunc("/ws", websocket.HandleWS)

	// serving static files
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static")
	})

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", middlewares.CacheMiddleware(&router)))
}
