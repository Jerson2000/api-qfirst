package routes

import (
	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/controllers"
	"github.com/jerson2000/api-qfirst/middlewares"
	"github.com/justinas/alice"
)

func AuthRoutes(router *mux.Router) {
	auth := "/auth"
	csrfChain := alice.New(middlewares.CSRFMiddleware)
	userChain := alice.New(middlewares.JwtMiddleware, middlewares.CSRFMiddleware)
	router.HandleFunc(auth+"/signup", controllers.AuthSignup).Methods("POST")
	router.HandleFunc(auth+"/login", controllers.AuthLogin).Methods("POST")
	router.Handle(auth+"/refresh", csrfChain.ThenFunc(controllers.AuthRefresh)).Methods("POST")
	router.Handle(auth+"/csrf", csrfChain.ThenFunc(controllers.AuthRequestCSRFToken)).Methods("GET")
	router.Handle(auth+"/current", userChain.ThenFunc(controllers.AuthCurrent)).Methods("POST")
}
