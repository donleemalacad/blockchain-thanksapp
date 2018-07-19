package controllers

import (
	"net/http"
)

func (app *Application) HomeController(w http.ResponseWriter, r *http.Request) {
	helloValue, err := app.Fabric.SampleDataUsingSDK()
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	data := &struct {
		Hello string
	}{
		Hello: helloValue,
	}
	renderTemplate(w, r, "index.html", data)
}