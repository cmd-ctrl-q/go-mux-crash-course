package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port = ":8080"
)

func main() {
	router := mux.NewRouter() // create new instance of a router

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("up and running..."))
		fmt.Fprintln(w, "Up and running...")
	})

	router.HandleFunc("/posts", getAllPosts).Methods("GET")
	router.HandleFunc("/posts", addPost).Methods("POST")

	fmt.Println("Server listening on port ", port)
	log.Fatal(http.ListenAndServe(":8080", router))
}
