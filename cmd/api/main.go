// package main spins up an http server using gorilla/mux for request multi-plexing
package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/dmithamo/timelineapi/pkg/dbservice"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

// application collects all the <injectable> dependencies of the app
type application struct {
	db    *sql.DB
	cache redis.Conn
}

func main() {
	// initialize app, router, register middleware, add routes
	app := &application{}
	r := mux.NewRouter()
	r.Use(app.SetCorsPolicy)
	r.Use(app.EnforceContentType)
	registerRoutes(r, app)

	// parse flags, connect db, create tables start server
	dsn := flag.String("dsn", "", "data source name for the db")
	addr := flag.String("addr", ":3001", "address where to serve application")
	rdb := flag.Bool("rdb", false, "set to true to drop all db tables and recreate them")
	cdsn := flag.String("cdsn", "", "data source name for the cache storage service")
	flag.Parse()

	// connect to main db
	db, err := dbservice.ConnectDB(dsn)
	if err != nil {
		log.Fatal("connectdb [start]: ", err)
	}
	defer db.Close()

	// connect to redis db
	cache, err := dbservice.ConnectCache(cdsn)
	if err != nil {
		log.Fatal("connectcache [start]: ", err)
	}
	defer cache.Close()

	// inject db, cache (and other dependencies) into app
	app.db = db
	app.cache = cache

	if *rdb {
		err := dbservice.DropTables(db)
		if err != nil {
			log.Fatal("drop tables [start]: ", err)
		}
	}

	err = dbservice.CreateTables(db)
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

func registerRoutes(r *mux.Router, a *application) {
	// /auth
	r.HandleFunc("/auth/register", a.registerUser).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", a.loginUser).Methods(http.MethodPost)

	// secure routes
	s := r.PathPrefix("").Subrouter()
	s.Use(a.CheckAuth)

	// auth - update user
	s.HandleFunc("/auth/register/{userID:[a-z0-9-]+}", a.updateUser).Methods(http.MethodPatch)

	// /projects
	s.HandleFunc("/projects", a.createProject).Methods(http.MethodPost)
	s.HandleFunc("/projects", a.getProjects).Methods(http.MethodGet)
	s.HandleFunc("/projects/{projectID}", a.getProject).Methods(http.MethodGet)
	s.HandleFunc("/projects/{projectID}", a.updateProject).Methods(http.MethodPatch)
	s.HandleFunc("/projects/{projectID}", a.deleteProject).Methods(http.MethodDelete)

	// /tasks
	s.HandleFunc("/projects", a.createTask).Methods(http.MethodPost)
	s.HandleFunc("/tasks", a.getTasks).Methods(http.MethodGet)
	s.HandleFunc("/tasks/{taskID}", a.getTask).Methods(http.MethodGet)
	s.HandleFunc("/tasks/{projectID}", a.getTasksByProject).Methods(http.MethodGet)
	s.HandleFunc("/tasks/{taskID}", a.updateTask).Methods(http.MethodPatch)
	s.HandleFunc("/tasks/{taskID}", a.deleteTask).Methods(http.MethodDelete)
}