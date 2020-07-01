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

			response, err := formatResponseHelper("Bad request. Request body should be valid JSON")
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusUnprocessableEntity)
			_, err = w.Write(response)

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
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
			response, err := formatResponseHelper("Method not implemented")
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusNotImplemented)
			_, err = w.Write(response)

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// formatResponseHelper structures a response message as JSON
func formatResponseHelper(message string) ([]byte, error) {
	type errMessage struct {
		Message string
	}
	msg := errMessage{message}
	resp, err := json.Marshal(&msg)

	return resp, err
}
