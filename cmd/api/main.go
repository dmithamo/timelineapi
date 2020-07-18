// package main spins up an http server using gorilla/mux for request multi-plexing
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dmithamo/timelineapi/pkg/dbservice"
	"github.com/dmithamo/timelineapi/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// application collects all the <injectable> dependencies of the app
type application struct {
	db *sql.DB
}

func main() {
	// initialize app, router, register middleware, add routes
	app := &application{}
	r := mux.NewRouter()

	go registerRoutesAndMiddleware(r, app)

	// parse flags, connect db, create tables start server
	dsn := flag.String("dsn", "", "data source name for the db")
	addr := flag.String("addr", ":3001", "address where to serve application")
	rdb := flag.Bool("rdb", false, "set to true to drop all db tables and recreate them")
	flag.Parse()

	// also load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("loadenv [start]: ", err)
	}

	// connect to main db
	db, err := dbservice.ConnectDB(dsn)
	if err != nil {
		log.Fatal("connectdb [start]: ", err)
	}
	defer db.Close()

	// inject db, cache (and other dependencies) into app
	app.db = db

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

	// gracefully shutdown on Ctrl-C signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go shutDownServer(&signals)

	log.Printf("starting server at http://127.0.0.1%v", *addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("serve [start]: ", err)
	}
}

// shutDownServer kills the server gracefully
func shutDownServer(signals *chan os.Signal) {
	if <-*signals == os.Interrupt {
		for i := 3; i > 0; i-- {
			fmt.Printf("\rShutting down in %v", i)
			time.Sleep(1 * time.Microsecond)
		}
		fmt.Println("\nBye!")
		os.Exit(0)
	}
}

func registerRoutesAndMiddleware(r *mux.Router, a *application) {
	// router-wide middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.SetCorsPolicy)
	r.Use(middleware.EnforceContentType)

	// /auth
	r.HandleFunc("/users", a.registerUser).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", a.loginUser).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", a.logoutUser).Methods(http.MethodPost)

	// secure routes
	s := r.PathPrefix("").Subrouter()
	s.Use(middleware.CheckAuth)

	// auth - update user
	s.HandleFunc("/auth/register/{userID:[0-9a-z-]+}", a.updateUser).Methods(http.MethodPatch)

	// /actions
	s.HandleFunc("/actions", a.createAction).Methods(http.MethodPost)
	s.HandleFunc("/actions", a.getActions).Methods(http.MethodGet)
	s.HandleFunc("/actions/{actionID:[0-9a-z-]+}", a.getAction).Methods(http.MethodGet)
	s.HandleFunc("/actions/{actionID:[0-9a-z-]+}", a.updateAction).Methods(http.MethodPatch)
	s.HandleFunc("/actions/{actionID:[0-9a-z-]+}", a.deleteAction).Methods(http.MethodDelete)

	// /outputs
	s.HandleFunc("/outputs", a.createOutput).Methods(http.MethodPost)
	s.HandleFunc("/outputs", a.getOutputs).Methods(http.MethodGet)
	s.HandleFunc("/outputs/{outputID:[0-9-]+}", a.getOutput).Methods(http.MethodGet)
	s.HandleFunc("/outputs/{actionID:[0-9-]+}", a.getOutputsByAction).Methods(http.MethodGet)
	s.HandleFunc("/outputs/{actionID:[0-9-]+}", a.updateOutput).Methods(http.MethodPatch)
	s.HandleFunc("/outputs/{actionID:[0-9-]+}", a.deleteOutput).Methods(http.MethodDelete)
}
