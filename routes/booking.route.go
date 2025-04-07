package routes

import (
	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/controllers"
	"github.com/jerson2000/api-qfirst/middlewares"
	"github.com/justinas/alice"
)

func BookingRoutes(router *mux.Router) {
	booking := "/booking"
	// avoid unauthorized access
	userChain := alice.New(middlewares.JwtMiddleware)
	// POST
	router.Handle(booking, userChain.ThenFunc(controllers.BookingCreate)).Methods("POST")
	// PUT & PATCH
	router.Handle(booking+"/{id}", userChain.ThenFunc(controllers.BookingUpdate)).Methods("PATCH")

	// GET
	router.Handle(booking, userChain.ThenFunc(controllers.BookingList)).Methods("GET")
	router.Handle(booking+"/{id}", userChain.ThenFunc(controllers.BookingGetById)).Methods("GET")

	// DELETE
	router.Handle(booking+"/{id}", userChain.ThenFunc(controllers.BookingDelete)).Methods("DELETE")
}
