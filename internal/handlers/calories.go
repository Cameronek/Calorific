package handlers

import (
	"net/http"
	"strconv"
	"github.com/cameronek/Calorific/internal/database"
)

func AddCalsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
	http.Error(w, "Method not a post request", http.StatusMethodNotAllowed)
	return
	}

	name := r.FormValue("foodName")
	foodCals := r.FormValue("foodCals")
	cals, err := strconv.ParseInt(foodCals, 10, 64)
	if err != nil {
		http.Error(w, "Invalid calories", http.StatusBadRequest)
		return
	}

	db, err := database.Initialize("./calorific.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

	_, err = db.Exec("INSERT INTO dailyConsumption (name, calories, date) VALUES (?, ?, DATE('now'))", name, cals)
	if err != nil {
		http.Error(w, "Failed to delete food", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteCalsHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.Exec("DELETE FROM dailyConsumption WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete food", http.StatusInternalServerError)
		return
	}		
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditTargetHandler(w http.ResponseWriter, r *http.Request) {
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

	cals, err := strconv.Atoi(r.FormValue("kCals"))

	if err != nil {
		http.Error(w, "Invalid calorie input", http.StatusBadRequest)
		return
	}

	today := r.FormValue("date")

	db, err := database.Initialize("./calorific.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if len(today) == 1 {
		today = "0" + today
	}

	_, err = db.Exec("UPDATE dailyGoal SET goalCalories = ? WHERE id = (SELECT id FROM dailyGoal ORDER BY id DESC LIMIT 1) AND strftime('%d', date) = ?", cals, today)
	if err != nil {
		http.Error(w, "Error editing target", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)	
}

