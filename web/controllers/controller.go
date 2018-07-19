package controllers

import (
	"fmt"
	"github.com/thanksapp/sdk/integration"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Application struct {
	Fabric *integration.SdkSetup
}

func renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	layoutPath := filepath.Join("web", "views", "layouts", "main.html")
	templatePath := filepath.Join("web", "views", "contents", templateName)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(templatePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	resultTemplate, err := template.ParseFiles(templatePath, layoutPath)
	if err != nil {
		// Log the detailed error
		fmt.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := resultTemplate.ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}