package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cameronek/Calorific/internal/database"
	"github.com/cameronek/Calorific/internal/templates"
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

	targets := [7]string{}
	for i := 0; i < 7; i++ {
		//target, err := database.GetTarget(db, strconv.Itoa(time.Now().AddDate(0, 0, -i).Day()))
		target, err := database.GetTarget(db, time.Now().AddDate(0, 0, -i))
		targets[i] = strconv.Itoa(target)
		if err != nil {
			http.Error(w, "Invalid calorie input", http.StatusInternalServerError)
			return
		}
	}

	sums := [7]int{}
	for i := 0; i < 7; i++ {
		//sum, err := database.GetDailyConsumption(db, strconv.Itoa(time.Now().AddDate(0,0,-i).Day()))
		sum, err := database.GetDailyConsumption(db, time.Now().AddDate(0, 0, -i))

		sums[i] = sum
		if err != nil {
			http.Error(w, "Error getting daily calorie consumption", http.StatusInternalServerError)
			return
		}
	}

	dailyFoods, err := database.GetDailyFoods(db)
	if err != nil {
		http.Error(w, "Error getting daily food consumption", http.StatusInternalServerError)
	}

	streak, err := database.GetStreak(db)
	if err != nil {
		http.Error(w, "Error getting streak of days where target is reached", http.StatusInternalServerError)
	}

	ctx := context.WithValue(context.Background(), "foods", foods)
	ctx = context.WithValue(ctx, "dailyFoods", dailyFoods)
	ctx = context.WithValue(ctx, "streak", streak)

	for i := 0; i < len(targets); i++ {
		ctx = context.WithValue(ctx, "target" + strconv.Itoa(i), targets[i])
	}

	for i := 0; i < len(targets); i++ {
		ctx = context.WithValue(ctx, "sum" + strconv.Itoa(i), sums[i])
	}
	
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

	_, err = db.Exec("UPDATE dailyGoal SET goalCalories = ? WHERE id = (SELECT id FROM dailyGoal ORDER BY id DESC LIMIT 1) AND strftime('%d', date) = ?", cals, today)
	if err != nil {
		http.Error(w, "Error editing target", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)	
}

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