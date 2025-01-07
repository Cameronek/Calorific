package handlers

import (
	"context"
	"github.com/cameronek/Calorific/internal/database"
	"github.com/cameronek/Calorific/internal/templates"
	"net/http"
	"strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	db, err := database.Initialize("./calorific.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()


	foods, err := database.GetFoods(db)
	if err != nil {
		http.Error(w, "Error getting foods", http.StatusInternalServerError)
	}

	ctx := context.WithValue(context.Background(), "foods", foods)

	component := templates.Index()
	component.Render(ctx, w)
}


// Refactor:  Move this handler into its own handler (foodHandlers.go)
func AddFoodHandler(w http.ResponseWriter, r *http.Request) {
	// If method passed in isnt a post request, error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	foodName := r.FormValue("food")
	calories, err := strconv.Atoi(r.FormValue("kCals"))

	if err != nil {
		http.Error(w, "Invalid calorie input", http.StatusBadRequest)
		return
	}

	db, err := database.Initialize("./calorific.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO food (name, calories) VALUES (?, ?)", foodName, calories)
	if err != nil {
		http.Error(w, "Error saving food", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not a post request", http.StatusMethodNotAllowed)
		return
	}

	foodID := r.FormValue("foodID")
	id, err := strconv.ParseInt(foodID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid food ID", http.StatusBadRequest)
		return
	}

	db, err := database.Initialize("./calorific.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

	_, err = db.Exec("DELETE FROM food WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete food", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
