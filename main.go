package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Todo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

var todos []Todo

func GetTodos(w http.ResponseWriter, r *http.Request) {
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
			return
		}
	}

	todos = append(todos, newTodo)
	log.Print("Added Todo: ", newTodo)
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

var upgrader = websocket.Upgrader{}

func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
	}
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
	router.HandleFunc("/todo/{id}", GetTodo).Methods("GET")
	router.HandleFunc("/todo/{id}", AddReplaceTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteTodo).Methods("DELETE")
	router.HandleFunc("/ws", ServeWebSocket)
	log.Fatal(http.ListenAndServe(port, router))
}
