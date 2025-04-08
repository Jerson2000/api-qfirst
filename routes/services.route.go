package routes

import (
	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/controllers"
	"github.com/jerson2000/api-qfirst/middlewares"
	"github.com/justinas/alice"
)

func ServicesRoutes(router *mux.Router) {
	services := "/services"
	// avoid unauthorized access
	userChain := alice.New(middlewares.JwtMiddleware, middlewares.RBACMiddleware)
	// POST
	router.Handle(services, userChain.ThenFunc(controllers.ServiceCreate)).Methods("POST")
	// PUT & PATCH
	router.Handle(services+"/{id}", userChain.ThenFunc(controllers.ServiceUpdate)).Methods("PATCH")

	// GET
	router.Handle(services, userChain.ThenFunc(controllers.ServiceList)).Methods("GET")
	router.Handle(services+"/{id}", userChain.ThenFunc(controllers.ServiceGetById)).Methods("GET")

	// DELETE
	router.Handle(services+"/{id}", userChain.ThenFunc(controllers.ServiceDelete)).Methods("DELETE")
}
