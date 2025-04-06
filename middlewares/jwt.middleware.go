package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")

			if authHeader == "" {
				models.ResponseWithError(res, http.StatusUnauthorized, "No token provided")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &models.JwtClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return config.JWTKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					models.ResponseWithError(res, http.StatusUnauthorized, "Invalid token signature")
					return
				}
				models.ResponseWithError(res, http.StatusBadRequest, "Invalid token")
				return
			}

			if !token.Valid {
				models.ResponseWithError(res, http.StatusBadRequest, "Invalid token")
				return
			}

			ctx := context.WithValue(req.Context(), "claims", claims)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
}
