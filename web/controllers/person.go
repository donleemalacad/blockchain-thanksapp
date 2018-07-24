package controllers

import (
	"net/http"
)
func (app *Application) AddPersonController(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Success     bool
		Fail        bool
		FailMessage string
	}{
		Success: false,
		Fail: false,
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
				
		data.Success = true
		
		if err != nil {
			data.Fail = true
			data.FailMessage = err.Error()
			data.Success = false
		}
	}
	
	// Render Template
	renderTemplate(w, r, "add-person.html", data)
}