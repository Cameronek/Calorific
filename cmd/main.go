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
	log.Println("Data initialized successfully")
	defer db.Close()

	mux := http.NewServeMux()

	// Home route
	mux.HandleFunc("/", handlers.HomeHandler)

	// POST: Add food route
	mux.HandleFunc("/addFood", handlers.AddFoodHandler)

	// Static route
	static := http.FileServer(http.Dir("static"))

	// GET localhost/static/
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	log.Println("Server start on localhost 8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
