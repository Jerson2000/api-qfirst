package middlewares

import (
	"log"
	"net/http"

	"slices"

	"github.com/jerson2000/api-qfirst/config"
)

type captureResponse struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (cres *captureResponse) WriteHeader(statusCode int) {
	cres.statusCode = statusCode
	cres.ResponseWriter.WriteHeader(statusCode)
}

func (cres *captureResponse) Write(b []byte) (int, error) {
	cres.body = append(cres.body, b...)
	return cres.ResponseWriter.Write(b)
}

func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Only cache GET requests
		if req.Method != http.MethodGet {
			next.ServeHTTP(res, req)
			return
		}

		// exclude paths
		excludedPaths := []string{
			"/auth/csrf",
			"/swagger/swagger-ui.css",
			"/swagger/swagger-ui-standalone-preset.js",
			"/swagger/swagger-ui-bundle.js",
			"/swagger/favicon-32x32.png",
			"/ws",
			"/",
		}

		// Check if the requested path is in the excluded list
		if slices.Contains(excludedPaths, req.URL.Path) {
			next.ServeHTTP(res, req)
			return
		}

		cacheKey := req.URL.Path

		// verify if cache from key exist
		if cached, exists := config.CacheClient.Get(cacheKey); exists {
			res.WriteHeader(http.StatusOK)
			res.Write(cached)
			log.Println("Cached response: ", cacheKey)
			return
		}

		// capture the response
		cres := &captureResponse{ResponseWriter: res}
		next.ServeHTTP(cres, req)
		// only cache successful responses
		if cres.statusCode >= 200 && cres.statusCode <= 299 && len(cres.body) > 0 {
			// cache the response body
			config.CacheClient.Set(cacheKey, cres.body)
			log.Println("Cached path: ", cacheKey)
		}
	})
}
