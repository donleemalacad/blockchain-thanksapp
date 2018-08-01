package controllers

import (
	"fmt"
	"net/http"
)

func (app *Application) ReplenishController(w http.ResponseWriter, r *http.Request) {
	addPointsToAll, err := app.Fabric.ReplenishPoints()
	fmt.Print("addPointsToAll controleer")
	fmt.Print(addPointsToAll)
	fmt.Printf("err")
	fmt.Print(err)
	if err != nil {
		http.Error(w, "Unable to replenish data from blockchain", 500)
	}
	fmt.Printf("ReplenishController ")
	fmt.Print(addPointsToAll)
	renderTemplate(w, r, "replenish.html", addPointsToAll)
}
