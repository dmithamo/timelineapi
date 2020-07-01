package middleware

import (
	"net/http"
	"strings"

	"github.com/dmithamo/timelineapi/pkg/utils"
)

// CheckAuth checks that authorization header is present and valid
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.SendJSONResponse(w, http.StatusUnauthorized, "No authorization header", nil)
			return
		}

		authHeaderSlice := strings.SplitN(authHeader, "Bearer ", 2)
		if len(authHeaderSlice) != 2 {
			utils.SendJSONResponse(w, http.StatusUnauthorized, "Malformed authorization header", nil)
			return
		}

		// verify token

		// if is authorized
		next.ServeHTTP(w, r)
	})
}
