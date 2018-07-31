package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AllUsersLedger struct {
	Name           string `json:"name"`
	PointsReceived int    `json:"pointsreceived"` // for history of received points, if you convert it to cash this is where it should be deducted(10pts)
	PointsSent     int    `json:"pointssent"`     // for history of sent points
	PointsCurrent  int    `json:"pointscurrent"`  // points you can give currently, technically should always be either 1 or 0
	Giver          string `json:"giver"`          // Other peers or system
	Message        string `json:"message"`        // Message given upon
	SentTo         string `json:"sentto"`         // Person the point is sent to
	Error          int    `json:"error"`          // Boolean error
}

func (app *Application) TransferController(w http.ResponseWriter, r *http.Request) {
	allUser, err := app.Fabric.GetAllUserDetailsInLedger()
	if err != nil {
		http.Error(w, "Unable to retrieve data from blockchain", 500)
	}

	// Convert GetAllUser Chaincode Function to Bytes
	byteConverted := []byte(allUser)

	// Using struct to parse byte converted result
	var details []AllUsersLedger
	if err := json.Unmarshal(byteConverted, &details); err != nil {

		http.Error(w, "Unable to parse blockchain result", 500)
	}
	fmt.Println("details")
	fmt.Println(details)
	// Map so that you can iterate in template using range struct
	mapping := map[string]interface{}{
		"Result": details,
	}

	if r.Method == "POST" {
		// Call ParseForm to parse the raw query
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", 500)
		}

		// Get Form field name value
		from := r.FormValue("from")
		to := r.FormValue("to")
		message := r.FormValue("message")

		// Save to Ledger
		success, err := app.Fabric.TransferPoint(to, from, message)
		fmt.Println(success)

		if err != nil {
			// data.Fail = true
			// data.FailMessage = err.Error()
			// data.Success = false
		}
		// http.Redirect(w, r, "/transaction/all/", http.StatusSeeOther)
	}
	fmt.Print("mapping")
	fmt.Print(mapping)
	// Render Template
	renderTemplate(w, r, "transfer.html", mapping)
}
