// package middleware defines middleware fns
package middleware

import (
	"encoding/json"
	"net/http"
)

// EnforceContentType checks that the request body is JSON-formatted,
// and sets the response content-type as JSON
func EnforceContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedContentType := "application/json"
		w.Header().Set("Content-Type", allowedContentType)

		if (r.Method == http.MethodPost || r.Method == http.MethodPatch) &&
			r.Header.Get("Content-Type") != allowedContentType {

			type errMsg struct {
				Message string
			}
			msg := errMsg{"Bad request. Format request body as JSON"}
			response, err := json.Marshal(&msg)

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(response))

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}

		next.ServeHTTP(w, r)
	})
}
