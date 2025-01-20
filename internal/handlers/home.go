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
		target, err := database.GetTarget(db, time.Now().UTC().AddDate(0, 0, -i))
		targets[i] = strconv.Itoa(target)
		if err != nil {
			http.Error(w, "Invalid calorie input", http.StatusInternalServerError)
			return
		}
	}

	sums := [7]int{}
	for i := 0; i < 7; i++ {
		sum, err := database.GetDailyConsumption(db, time.Now().UTC().AddDate(0, 0, -i))

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

