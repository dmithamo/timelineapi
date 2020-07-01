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
	r.HandleFunc("/tasks/{projectID}", getTasksByProject).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{taskID}", updateTask).Methods(http.MethodPatch)
	r.HandleFunc("/tasks/{taskID}", deleteTask).Methods(http.MethodDelete)
}
