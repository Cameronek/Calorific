package handlers

import (
	"context"
	"net/http"
	"github.com/cameronek/Calorific/internal/templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.Index()
	component.Render(context.Background(), w)
}


