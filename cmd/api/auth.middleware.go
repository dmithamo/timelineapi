package main

import (
	"fmt"
	"net/http"

	"github.com/dmithamo/timelineapi/pkg/utils"
)

// CheckAuth checks that authorization header is present and valid
func (a *application) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.SendJSONResponse(w, http.StatusUnauthorized, "no valid authorization token", nil)
				return
			}
			// other errs
			utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err reading authorization: %v", err.Error()), nil)
			return
		}

		token := cookie.Value
		// get session_token from cache
		sessionToken, err := a.cache.Do("GET", token)

		if err != nil {
			utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err reading authorization: %v", err.Error()), nil)
			return
		}

		if sessionToken == nil {
			utils.SendJSONResponse(w, http.StatusUnauthorized, "no valid authorization token", nil)
			return
		}

		// if is authorized
		next.ServeHTTP(w, r)
	})
}
