package main

import (
	"log"
	"net/http"

	"github.com/cameronek/Calorific/internal/database"
	"github.com/cameronek/Calorific/internal/handlers"

	//"path/filepath"
) 

func main() {

	db, err := database.Initialize("./calorific.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data initializied successfully")
	defer db.Close()


	mux := http.NewServeMux()

	// Serve static files
	// Not going to use external static files for the time being
	//fs := http.FileServer(http.Dir("static"))
	//mux.Handle("/static/", http.StripPrefix("/static/", fs))


/*
	staticPath, err := filepath.Abs("internal/templates")
	log.Println(staticPath)

	if err != nil {
		log.Fatalf("Failed to resolve static directory: %v", err)
	} */

	// Home route
	mux.HandleFunc("GET /{$}", handlers.HomeHandler)

	// Image route
	images := http.FileServer(http.Dir("static"))

	// GET localhost/static/
	mux.Handle("/static/", http.StripPrefix("/static/", images))

	log.Println("Server start on localhost 8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
