package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dmithamo/timelineapi/pkg/security"
)

// RequestLogger logs requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go func() {
			var userID interface{}
			tokenData, err := security.DecodeToken(r)
			if err != nil {
				userID = nil
			} else {
				userID = tokenData
			}

			logEntry := fmt.Sprintf(
				"%v %v %v %v %v %v\n",
				time.Now().Format("2006/01/02 Mon 03:04:05PM-0700 MST"),
				r.Header.Get("User-Agent"),
				r.Proto,
				userID,
				r.Method,
				r.URL,
			)

			// create log file at server root, discarding any errors
			filename := fmt.Sprintf("%v.log", time.Now().Format("Mon Jan 02 2006"))
			f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("err opening log file: ", err)
			}

			defer f.Close()
			_, err = fmt.Fprint(f, logEntry)

			if err != nil {
				fmt.Println("err writing log to file: ", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
