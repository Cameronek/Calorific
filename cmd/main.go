package main

import (
	"log"
	"net/http"

	"github.com/cameronek/Calorific/internal/database"
	"github.com/cameronek/Calorific/internal/handlers"
) 

func main() {

	db, err := database.Initialize("./calorific.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data initialzied successfully")
	defer db.Close()


	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Home route
	mux.HandleFunc("/", handlers.HomeHandler)

	log.Println("Server start on localhost 8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
