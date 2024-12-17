package handlers

import (
	"Calorific/internal/templates"
	"net/http"

	"github.com/a-h/templ/runtime/render"
	"github.com/cameronek/Calorific/templ-project/internal/templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.Index("Gopher")
	render.New(w).Render(r.Context(), component)
}


