package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
  	"database/sql"
  	_"github.com/go-sql-driver/mysql"
	"net/http"
	"github.com/gorilla/mux"  
)

var db *sql.DB
var err error

type Todo struct 
{
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
}

func getToDos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		
	var todos []Todo

	result, err := db.Query("SELECT id, title, description from todo")
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
	  var todo Todo
	  err := result.Scan(&todo.ID, &todo.Title, &todo.Description)
	  if err != nil {
		panic(err.Error())
	  }
	  todos = append(todos, todo)
	}
	json.NewEncoder(w).Encode(todos)
}
func createToDo(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO todo(title, description) VALUES(?, ?)")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	description:= keyVal["description"]
	_, err = stmt.Exec(title, description)
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}
func getToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title, description FROM todo WHERE id = ?", params["id"])
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	var todo Todo
	for result.Next() {
	  err := result.Scan(&todo.ID, &todo.Title, &todo.Description)
	  if err != nil {
		panic(err.Error())
	  }
	}
	json.NewEncoder(w).Encode(todo)
}
func updateToDo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE todo SET title = ? WHERE id = ?")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
}
func deleteToDo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM todo WHERE id = ?")
	if err != nil {
	  panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
   if err != nil {
	  panic(err.Error())
	}
  fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}

func main() {
  	db, err = sql.Open("mysql", "root:priya@1298@tcp(127.0.0.1:3306)/go_account")
  	if err != nil {
    	panic(err.Error())
	}
	  
	defer db.Close()
	
	router := mux.NewRouter()
	
	router.HandleFunc("/todo", getToDos).Methods("GET")
	router.HandleFunc("/todo", createToDo).Methods("POST")
	router.HandleFunc("/todo/{id}", getToDo).Methods("GET")
	router.HandleFunc("/todo/{id}", updateToDo).Methods("PUT")
	router.HandleFunc("/todo/{id}", deleteToDo).Methods("DELETE")
	
	http.ListenAndServe(":8000", router)
}