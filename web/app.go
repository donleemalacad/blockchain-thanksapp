package web

import (
	"fmt"
	"github.com/thanksapp/web/controllers"
	"net/http"
)

func Serve(app *controllers.Application) {
	// Tell Where Assets Directory is located
	fs := http.FileServer(http.Dir("web/assets"))

	// Every time assets is requested, return the contents inside of web/assets/
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// If index is requested, use HomeController
	http.HandleFunc("/", app.HomeController)

	// If /transaction/all/ path is requested, use AllTransactionController
	http.HandleFunc("/transaction/all/", app.AllTransactionsController)

	// If /transaction/all/ path is requested, use AllTransactionController
	http.HandleFunc("/add/", app.AddPersonController)

	// Listen to port 3000
	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}