package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)


type Todo struct {
	Id string `json:"id"`
	Title string `json:"title"`
}
var todos []Todo


func GetTodos(w http.ResponseWriter, r* http.Request) {
	json.NewEncoder(w).Encode(todos)
}


func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	port = ":" + port
	log.Printf("Listening on port %s", port)

	// add one sample entry
	todos = append(todos, Todo{Id: "1", Title: "First Todo"})

	router := mux.NewRouter()
	router.HandleFunc("/todo", GetTodos).Methods("GET")
	log.Fatal(http.ListenAndServe(port, router))
}
