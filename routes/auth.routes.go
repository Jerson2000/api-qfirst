package routes

import (
	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/controllers"
)

func AuthRoutes(router *mux.Router) {
	auth := "/auth"
	router.HandleFunc(auth+"/signup", controllers.AuthSignup).Methods("POST")
	router.HandleFunc(auth+"/login", controllers.AuthLogin).Methods("POST")
	router.HandleFunc(auth+"/refresh", controllers.AuthRefresh).Methods("POST")
}
