package middlewares

import (
	"net/http"

	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
)

func RBACMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		claims := req.Context().Value("claims").(*models.JwtClaims)
		role := claims.Role

		// The object is the requested URL path and the action is the HTTP method
		path := req.URL.Path
		method := req.Method

		// Check permission
		ok, err := config.EnforcerCasbin.Enforce(string(role), path, method)
		if err != nil {
			// log.Printf("Casbin enforce error: %v", err)
			models.ResponseWithError(res, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if !ok {
			models.ResponseWithError(res, http.StatusForbidden, "You don't have permission to access this resource!")
			return
		}

		// log.Printf("Access granted for user: %s (role: %s) to %s with method %s", user, role, path, method)
		next.ServeHTTP(res, req)
	})

}
