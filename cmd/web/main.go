// package main spins up an http server using gorilla/mux for request multi-plexing
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	db "github.com/dmithamo/timelineapi/pkg/dbservice"
	"github.com/dmithamo/timelineapi/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// initialize router, register middleware, add routes
	r := mux.NewRouter()
	r.Use(middleware.SetCorsPolicy)
	r.Use(middleware.EnforceContentType)
	registerRoutes(r)

	// parse flags, connect db, create tables start server
	dsn := flag.String("dsn", "", "data source name for the db")
	addr := flag.String("addr", ":3001", "address where to serve application")
	rdb := flag.Bool("rdb", false, "set to true to drop all db tables and recreate them")
	flag.Parse()

	conn, err := db.ConnectDB(dsn)
	if err != nil {
		log.Fatal("connectdb [start]: ", err)
	}

	if *rdb {
		err := db.DropTables(conn)
		if err != nil {
			log.Fatal("drop tables [start]: ", err)
		}
	}

	err = db.CreateTables(conn)
	if err != nil {
		log.Fatal("create tables [start]: ", err)
	}
	log.Println("successfully connected to db")

	//serve!
	srv := &http.Server{
		Addr:         *addr,
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("starting server at http://127.0.0.1%v", *addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("serve [start]: ", err)
	}
}

func registerRoutes(r *mux.Router) {
	// /auth
	r.HandleFunc("/auth/register", registerUser).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", loginUser).Methods(http.MethodPost)

	// secure routes
	s := r.PathPrefix("").Subrouter()
	s.Use(middleware.CheckAuth)

	// auth - update user
	s.HandleFunc("/auth/register/{userID:[a-z0-9-]+}", updateUser).Methods(http.MethodPatch)

	// /projects
	s.HandleFunc("/projects", createProject).Methods(http.MethodPost)
	s.HandleFunc("/projects", getProjects).Methods(http.MethodGet)
	s.HandleFunc("/projects/{projectID}", getProject).Methods(http.MethodGet)
	s.HandleFunc("/projects/{projectID}", updateProject).Methods(http.MethodPatch)
	s.HandleFunc("/projects/{projectID}", deleteProject).Methods(http.MethodDelete)

	// /tasks
	s.HandleFunc("/projects", createTask).Methods(http.MethodPost)
	s.HandleFunc("/tasks", getTasks).Methods(http.MethodGet)
	s.HandleFunc("/tasks/{taskID}", getTask).Methods(http.MethodGet)
	s.HandleFunc("/tasks/{projectID}", getTasksByProject).Methods(http.MethodGet)
	s.HandleFunc("/tasks/{taskID}", updateTask).Methods(http.MethodPatch)
	s.HandleFunc("/tasks/{taskID}", deleteTask).Methods(http.MethodDelete)
}
