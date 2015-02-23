package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id   int64  `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type Todos []Todos

func SetupDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS todos (` +
		`id INTEGER PRIMARY KEY   AUTOINCREMENT, ` +
		`task VARCHAR(255) NOT NULL, ` +
		`done BOOLEAN NOT NULL)`)

	if err != nil {
		return err
	}

	//Some test data.
	for _, task := range []string{"Learn Slurp", "Be productive", "Be nice.", "Call parents.", "Help the poor"} {

		_, err := db.Exec(`INSERT INTO todos (task, done) VALUES ($1,$2)`, task, false)
		if err != nil {
			return err
		}

	}

	return nil
}

func RegisterAPI(api *mux.Router, db *sql.DB) {

	//Get All.
	api.Path("/api/todos").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT Id, Task, Done FROM todos")
		defer func() {
			if err != nil {
				log.Println(err)
				http.Error(w, "Well, this is embrassing. The server is broken. Try again?", http.StatusInternalServerError)
			}
		}()

		if err != nil {
			return
		}

		todos := []Todo{}
		for rows.Next() {
			todo := Todo{}
			if err = rows.Scan(&todo.Id, &todo.Task, &todo.Done); err != nil {
				return
			}
			todos = append(todos, todo)
		}

		json.NewEncoder(w).Encode(todos)
	})

	//New TODO.
	api.Path("/api/todos").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		todo := &Todo{}
		err := json.NewDecoder(r.Body).Decode(todo)
		if err != nil {
			log.Println(err)
			http.Error(w, "Malformed Request.", 422) //Unprocessable Entity.
			return
		}

		res, err := db.Exec(`INSERT INTO todos (task, done) VALUES (:task, :done)`, todo.Task, todo.Done)

		if err != nil {
			log.Println(err)
			http.Error(w, "Well, this is embrassing. The server is broken. Try again?", http.StatusInternalServerError)
			return
		}

		todo.Id, _ = res.LastInsertId()
		json.NewEncoder(w).Encode(todo)
	})

	//Update
	api.Path("/api/todos/{id}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		todo := &Todo{}
		err := json.NewDecoder(r.Body).Decode(todo)
		if err != nil {
			log.Println(err)
			http.Error(w, "Malformed Request.", 422) //Unprocessable Entity.
			return
		}

		_, err = db.Exec(`UPDATE todos SET task=$1, done=$2 WHERE id=$3`, todo.Task, todo.Done, mux.Vars(r)["id"])
		if err != nil {
			log.Println(err)
			http.Error(w, "Well, this is embrassing. The server is broken. Try again?", http.StatusInternalServerError)
			return
		}
	})
	//Delete
	api.Path("/api/todos/{id}").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := db.Exec(`DELETE FROM todos WHERE id=$1`, mux.Vars(r)["id"])
		if err != nil {
			log.Println(err)
			http.Error(w, "Well, this is embrassing. The server is broken. Try again?", http.StatusInternalServerError)
			return
		}
	})
}
