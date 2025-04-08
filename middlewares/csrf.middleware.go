package middlewares

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/jerson2000/api-qfirst/config"
)

func CSRFMiddleware(next http.Handler) http.Handler {
	return csrf.Protect(
		config.CSRFKey,
		csrf.Secure(true),   // Enable for HTTPS in production
		csrf.Path("/"),      // Cookie scope (you can adjust this based on your needs)
		csrf.HttpOnly(true), // Prevent JavaScript from accessing the CSRF cookie
	)(next)
}
