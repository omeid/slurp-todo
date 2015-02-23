package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Public http.FileSystem = http.Dir("./public")


	config = struct {
		Host       string // May be used for templates.
		Livereload string // For Debug/Development.
		Port       string
	}{
		Port: os.Getenv("PORT"),
	}
)

func init() {
	if config.Port == "" {
		config.Port = "8080"
	}
}

func main() {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
	  log.Fatal(err)
	}

	SetupDB(db)

	api := mux.NewRouter()
	RegisterAPI(api, db)
	api.PathPrefix("/").Handler(http.FileServer(Public))

	endpoint := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Printf("Listening on %s", endpoint)
	log.Fatal(http.ListenAndServe(endpoint, api))
}
