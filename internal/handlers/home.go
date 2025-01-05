package handlers

import (
	"context"
	"github.com/cameronek/Calorific/internal/templates"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.Index()
	component.Render(context.Background(), w)
}
