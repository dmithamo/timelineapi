package middleware

import (
	"net/http"

	"github.com/dmithamo/timelineapi/pkg/utils"
)

// EnforceContentType checks that the request body is JSON-formatted,
// and sets the response content-type as JSON
func EnforceContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedContentType := "application/json"
		w.Header().Set("Content-Type", allowedContentType)

		if (r.Method == http.MethodPost || r.Method == http.MethodPatch) &&
			r.Header.Get("Content-Type") != allowedContentType {
			utils.SendJSONResponse(w, http.StatusUnprocessableEntity, "Bad request. Request body should be valid JSON", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SetCorsPolicy set the cross origin request policy
func SetCorsPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")

		if r.Method == "OPTIONS" {
			utils.SendJSONResponse(w, http.StatusNotImplemented, "Not supported", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
