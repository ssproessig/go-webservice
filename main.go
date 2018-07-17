package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)


func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	port = ":" + port
	log.Printf("Listening on port %s", port)

	router := mux.NewRouter()
	log.Fatal(http.ListenAndServe(port, router))
}
