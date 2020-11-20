package main
import (
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
)


type Book struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var book Book
  _ = json.NewDecoder(r.Body).Decode(&book)
  books = append(books, book)
  json.NewEncoder(w).Encode(&book)
}

func getBook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  for _, item := range books {
    if item.Id == params["id"] {
      json.NewEncoder(w).Encode(item)
      return
    }
  }
  json.NewEncoder(w).Encode(&Book{})
}

func updateBook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  for index, item := range books {
    if item.Id == params["id"] {
      books = append(books[:index], books[index+1:]...)
      var book Book
      _ = json.NewDecoder(r.Body).Decode(&book)
      book.Id = params["id"]
      books = append(books, book)
      json.NewEncoder(w).Encode(&book)
      return
    }
  }
  json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  for index, item := range books {
    if item.Id== params["id"] {
      books = append(books[:index], books[index+1:]...)
      break
    }
  }
  json.NewEncoder(w).Encode(books)
}

func main() {
  router := mux.NewRouter()

  books = []Book{
	Book{Id: "1", Name: "Book1", Author: "Durjoy Dutta", Publication: "Durjoy Dutta Publication"},
	Book{Id: "2", Name: "Book2", Author: "Hemang Patel", Publication: "Hemang Patel Publication"},
	Book{Id: "3", Name: "Book3", Author: "Dev Panchal", Publication: "Dev Panchal Publication"},
	}

  router.HandleFunc("/books", getBooks).Methods("GET")
  router.HandleFunc("/books", createBook).Methods("POST")
  router.HandleFunc("/books/{id}", getBook).Methods("GET")
  router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
  router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
  http.ListenAndServe(":8000", router)
}
