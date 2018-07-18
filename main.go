package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var todoChanged chan Todo

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	port = ":" + port
	log.Printf("Listening on port %s", port)

	todoChanged = make(chan Todo)
	go Connect2AMQPAndSetupQueue(GetAMQPUriToUse(), todoChanged)

	// add one sample entry
	todos = append(todos, Todo{Id: "1", Title: "First Todo"})

	router := mux.NewRouter()
	router.HandleFunc("/todo", GetTodos).Methods("GET")
	router.HandleFunc("/todo/{id}", GetTodo).Methods("GET")
	router.HandleFunc("/todo/{id}", AddReplaceTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteTodo).Methods("DELETE")
	router.HandleFunc("/ws", ServeWebSocket)
	log.Fatal(http.ListenAndServe(port, router))
}
