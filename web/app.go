package web

import (
	"fmt"
	"net/http"

	"github.com/thanksapp/web/controllers"
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

	// If /add/ path is requested, use AddPersonController
	http.HandleFunc("/add/", app.AddPersonController)

	// If /transaction/select/ path is requested, use UserTransactionsController
	http.HandleFunc("/transaction/select/", app.UserTransactionsController)

	// If /transaction/view/ path is requested, use ViewUserTransactionsController
	http.HandleFunc("/transaction/view/", app.ViewUserTransactionsController)

	// If /transaction/view/ path is requested, use ViewUserTransactionsController
	http.HandleFunc("/transfer/", app.TransferController)

	// If /transaction/view/ path is requested, use ViewUserTransactionsController
	http.HandleFunc("/replenish/", app.ReplenishController)

	// Listen to port 3000
	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
