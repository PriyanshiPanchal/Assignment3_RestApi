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

type Product struct 
{
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Quality string `json:"quality"`
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		
	var products []Product

	result, err := db.Query("SELECT id, name, description, quality from product")
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
	  var product Product
	  err := result.Scan(&product.ID, &product.Name, &product.Description, &product.Quality)
	  if err != nil {
		panic(err.Error())
	  }
	  products = append(products, product)
	}
	json.NewEncoder(w).Encode(products)
}
func createProduct(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO product(name, description, quality) VALUES(?, ?, ?)")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	description:= keyVal["description"]
	quality:= keyVal["quality"]
	_, err = stmt.Exec(name, description, quality)
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}
func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, name, description, quality FROM product WHERE id = ?", params["id"])
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	var product Product
	for result.Next() {
	  err := result.Scan(&product.ID, &product.Name, &product.Description,&product.Quality)
	  if err != nil {
		panic(err.Error())
	  }
	}
	json.NewEncoder(w).Encode(product)
}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE product SET name = ?, description=?, quality=? WHERE id = ?")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName := keyVal["name"]
	description:= keyVal["description"]
	quality:= keyVal["quality"]
	_, err = stmt.Exec(newName, description, quality, params["id"])
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM product WHERE id = ?")
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
	
	router.HandleFunc("/product", getProducts).Methods("GET")
	router.HandleFunc("/product", createProduct).Methods("POST")
	router.HandleFunc("/product/{id}", getProduct).Methods("GET")
	router.HandleFunc("/product/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")
	
	http.ListenAndServe(":8000", router)
}