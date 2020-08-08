package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	if db, err = sql.Open("sqlite3", "todo.db"); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS TODO (id integer primary key autoincrement, todo varchar(255), done boolean)"); err != nil {
		log.Fatal(err)
	}
}
func getTodos(w http.ResponseWriter) {
	todos := []Todo{}
	rows, err := db.Query("SELECT * from TODO")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Todo, &todo.Done); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		todos = append(todos, todo)
	}
	if err := json.NewEncoder(w).Encode(&todos); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if _, err := db.Exec("INSERT INTO TODO (todo, done) values (?,?)", todo.Todo, todo.Done); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
func main() {
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodos(w)
		case http.MethodPost:
			createTodo(w, r)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
