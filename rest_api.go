package main

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

var todos []Todo

func GetTodos(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range todos {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(404)
}

func AddReplaceTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var newTodo Todo
	_ = json.NewDecoder(r.Body).Decode(&newTodo)
	newTodo.Id = params["id"]

	for i := 0; i < len(todos); i++ {
		todo := &todos[i]
		if todo.Id == params["id"] {
			*todo = newTodo
			log.Print("Replaced Todo: ", newTodo)
			todoChanged <- newTodo
			return
		}
	}

	todos = append(todos, newTodo)
	log.Print("Added Todo: ", newTodo)
	todoChanged <- newTodo
	w.WriteHeader(201)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i := 0; i < len(todos); i++ {
		todo := &todos[i]
		if todo.Id == params["id"] {
			log.Print("Deleting Todo: ", todo)
			todos = append(todos[:i], todos[i+1:]...)
			return
		}
	}

	w.WriteHeader(404)
}
