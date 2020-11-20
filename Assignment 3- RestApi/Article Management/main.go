package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"id"`
    Title string `json:"title"`
    Desc string `json:"desc"`
    Content string `json:"content"`
}

var articles []Article

func getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}
  
func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var article Article
	_ = json.NewDecoder(r.Body).Decode(&article)

	articles = append(articles, article)

	json.NewEncoder(w).Encode(&article)
}
  
func getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range articles {
	    if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
	    }
	}
	json.NewEncoder(w).Encode(&Article{})
}
  
func updateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range articles {
	    if item.Id == params["id"] {
			articles = append(articles[:index], articles[index+1:]...)
			var article Article
			_ = json.NewDecoder(r.Body).Decode(&article)
			article.Id = params["id"]
			articles = append(articles, article)
			json.NewEncoder(w).Encode(&article)
			return
		}
	}
	json.NewEncoder(w).Encode(articles)
}
  
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range articles {
	  	if item.Id == params["id"] {
			articles = append(articles[:index], articles[index+1:]...)
			break
	  	}
	}
	json.NewEncoder(w).Encode(articles)
}
  
func main() {

	articles = []Article{
        Article{Id: "1", Title: "Article", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Python", Desc: "Python Description", Content: "Python Content"},
		Article{Id: "3", Title: "Java", Desc: "Java Description", Content: "Java Content"},
    }

	myRouter := mux.NewRouter()
	
	myRouter.HandleFunc("/article", getArticles)
	myRouter.HandleFunc("/article", createArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", getArticle)
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	http.ListenAndServe(":8000", myRouter)
}