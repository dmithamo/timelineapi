//package middleware intercepts http requests and ... does something useful
package middleware

import (
	"fmt"
	"net/http"

	"github.com/dmithamo/timelineapi/pkg/security"
	"github.com/dmithamo/timelineapi/pkg/utils"
)

// CheckAuth checks that authorization header is present and valid
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.SendJSONResponse(w, http.StatusUnauthorized,
					&utils.GenericJSONRes{
						Message: "no valid authorization token",
						Data:    nil,
					})
				return
			}

			// other errs
			utils.SendJSONResponse(w, http.StatusInternalServerError,
				&utils.GenericJSONRes{
					Message: fmt.Sprintf("err reading authorization token: %v", err.Error()),
					Data:    nil,
				})
			return
		}
		// validate token
		claims, err := security.ValidateToken(cookie.Value)

		if err != nil {
			if err.Error() == utils.TOKEN_EXPIRED_ERR {
				//means a refreshToken is needed
				refreshToken, refreshErr := security.GenerateToken(claims)
				if refreshErr != nil {
					utils.SendJSONResponse(w, http.StatusInternalServerError,
						&utils.GenericJSONRes{
							Message: fmt.Sprintf("err refreshing auth token: %v", refreshErr.Error()),
							Data:    nil,
						})
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    *refreshToken,
					HttpOnly: true,
					Path:     "/",
					SameSite: 1,
				})
			} else {
				// other errs
				utils.SendJSONResponse(w, http.StatusInternalServerError,
					&utils.GenericJSONRes{
						Message: fmt.Sprintf("err validating authorization token: %v", err.Error()),
						Data:    nil,
					})
				return
			}
		}

		// if is authorized
		next.ServeHTTP(w, r)
	})
}
