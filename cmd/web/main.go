// package main spins up an http server using gorilla/mux for request multi-plexing
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dmithamo/timelineapi/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// initialize router, register middleware, add routes
	r := mux.NewRouter()
	r.Use(middleware.EnforceContentType)

	// /auth
	r.HandleFunc("/auth/register", registerUser).Methods(http.MethodPost)
	r.HandleFunc("/auth/register/{userID:[a-z0-9-]+}", updateUser).Methods(http.MethodPatch)
	r.HandleFunc("/auth/login", loginUser).Methods(http.MethodPost)

	// /projects
	r.HandleFunc("/projects", createProject).Methods(http.MethodPost)
	r.HandleFunc("/projects", getProjects).Methods(http.MethodGet)
	r.HandleFunc("/projects/{projectID}", getProject).Methods(http.MethodGet)
	r.HandleFunc("/projects/{projectID}", updateProject).Methods(http.MethodPatch)
	r.HandleFunc("/projects/{projectID}", deleteProject).Methods(http.MethodDelete)

	// /tasks
	r.HandleFunc("/projects", createTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks", getTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{taskID}", getTask).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{projectID}", getTasksByProject).Methods(http.MethodGet) // *should this be under the /projects domain?
	r.HandleFunc("/tasks/{taskID}", updateTask).Methods(http.MethodPatch)
	r.HandleFunc("/tasks/{taskID}", deleteTask).Methods(http.MethodDelete)

	//serve!
	srv := &http.Server{
		Addr:         ":3001",
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Println("starting server at http://127.0.0.1:3001")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("error starting server", err)
	}
}
