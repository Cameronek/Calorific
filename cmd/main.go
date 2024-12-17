package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ/runtime/render"
	"github.com/cameronek/Calorific/internal/handlers"
) 

func main() {

	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("./static", http.StripPrefix("./static", fs))

	// Home route
	mux.HandleFunc("/", handlers.HomeHandler)

	log.Println("Server start on localhost 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
