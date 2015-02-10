package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Host       string // Maybe used for templates.
	Livereload string // For Debug/Development.
	Port       string
}

var (
	config = &Config{
		Port: os.Getenv("PORT"),
	}

	api    = mux.NewRouter()
	db     *sqlx.DB
	Public http.FileSystem = http.Dir("./public")
)

func init() {
	if config.Port == "" {
		config.Port = "8080"
	}
}

func main() {

	SetupDB()
	RegisterAPI()
	api.PathPrefix("/").Handler(http.FileServer(Public))
	endpoint := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Printf("Listening on %s", endpoint)

	log.Fatal(http.ListenAndServe(endpoint, api))
}
