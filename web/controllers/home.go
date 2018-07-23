package controllers

import (
	"net/http"
)

func (app *Application) HomeController(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "index.html", nil)
}