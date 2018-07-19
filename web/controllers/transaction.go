package controllers

import (
	"net/http"
)

func (app *Application) AllTransactionsController(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "all-transaction.html", nil)
}