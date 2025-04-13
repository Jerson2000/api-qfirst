package routes

import (
	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/controllers"
	"github.com/jerson2000/api-qfirst/middlewares"
	"github.com/justinas/alice"
)

func EmailRoutes(router *mux.Router) {
	otp := "/mailer"

	chain := alice.New(middlewares.CSRFMiddleware)
	// POST
	router.Handle(otp+"/otp/send", chain.ThenFunc(controllers.MailerGenerateOTP)).Methods("POST")
	router.Handle(otp+"/otp/validate", chain.ThenFunc(controllers.MailerValidateOTP)).Methods("POST")

}
