package controllers

import (
	"net/http"
)
func (app *Application) AddPersonController(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Success bool
	}{
		Success: false,
	}
	if r.Method == "POST" {
		// Call ParseForm to parse the raw query
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", 500)
		}

		// Get Form field name value
		name := r.FormValue("name")

		// Save to Ledger
		_, err := app.Fabric.AddPerson(name)
		if err != nil {
			http.Error(w, "Unable to query chaincode", 500)
		}
		
		data.Success = true
	}
	
	// Render Template
	renderTemplate(w, r, "add-person.html", data)
}