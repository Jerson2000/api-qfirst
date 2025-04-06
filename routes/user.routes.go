package routes

import (
	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/controllers"
	"github.com/jerson2000/api-qfirst/middlewares"
	"github.com/justinas/alice"
)

func UserRoutes(router *mux.Router) {
	users := "/users"
	userChain := alice.New(middlewares.JwtMiddleware)
	// POST
	router.Handle(users, userChain.ThenFunc(controllers.UserCreate)).Methods("POST")
	// PUT & PATCH
	router.Handle(users+"/{id}", userChain.ThenFunc(controllers.UserUpdate)).Methods("PATCH")

	// GET
	router.Handle(users, userChain.ThenFunc(controllers.UserList)).Methods("GET")
	router.Handle(users+"/{id}", userChain.ThenFunc(controllers.UserGetById)).Methods("GET")

	// DELETE
	router.Handle(users+"/{id}", userChain.ThenFunc(controllers.UserDelete)).Methods("DELETE")
}
