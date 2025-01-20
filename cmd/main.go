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
	log.Println("Data initialized successfully")
	defer db.Close()

	mux := http.NewServeMux()

	// Home route
	mux.HandleFunc("/", handlers.HomeHandler)

	// POST: Add food route
	mux.HandleFunc("/addFood", handlers.AddFoodHandler)

	// POST: Edit Calorie Target
	mux.HandleFunc("/editTarget", handlers.EditTargetHandler)

	// DELETE: Delete food route
	mux.HandleFunc("/deleteFood", handlers.DeleteFoodHandler)

	// POST: Add food to daily calorie consumption
	mux.HandleFunc("/addCals", handlers.AddCalsHandler)

	// DELETE: Delete consumed calories
	mux.HandleFunc("/deleteCals", handlers.DeleteCalsHandler)

	// Static route
	static := http.FileServer(http.Dir("static"))

	// GET localhost/static/
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	log.Println("Server start on localhost 8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
