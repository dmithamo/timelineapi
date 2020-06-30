// package main spins up an http server using gorilla/mux for request multi-plexing
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// initialize router, add routess
	r := mux.NewRouter()

	//serve!
	log.Println("starting server at http://127.0.0.1:3001")
	if err := http.ListenAndServe(":3001", r); err != nil {
		log.Fatal("error starting server", err)
	}
}
