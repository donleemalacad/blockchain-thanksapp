package controllers

import (
	"net/http"
	"encoding/json"
)

type ViewUserLedgerDetails struct {
	Name           string `json:"name"`
	PointsReceived int    `json:"pointsreceived"` // for history of received points, if you convert it to cash this is where it should be deducted(10pts)
	PointsSent     int    `json:"pointssent"`     // for history of sent points
	PointsCurrent  int    `json:"pointscurrent"`  // points you can give currently, technically should always be either 1 or 0
	Giver          string `json:"giver"`          // Other peers or system
	Message        string `json:"message"`        // Message given upon
	SentTo         string `json:"sentto"`         // Person the point is sent to
}

func (app *Application) ViewUserTransactionsController(w http.ResponseWriter, r *http.Request) {
	// User
	user := r.URL.Query()

	allUserHistory, err := app.Fabric.GetSpecificUserDetails(user.Get("name"))
	if err != nil {
		http.Error(w, "Unable to retrieve data from blockchain", 500)
	}

	// Convert GetAllUser Chaincode Function to Bytes
	byteConverted := []byte(allUserHistory)

	// Using struct to parse byte converted result
	var details []ViewUserLedgerDetails
	if err := json.Unmarshal(byteConverted, &details); err != nil {
		http.Error(w, "Unable to parse blockchain result", 500)
	}

	// Map so that you can iterate in template using range struct
	mapping := map[string]interface{}{"Result" : details}

	// Render Template
	renderTemplate(w, r, "user-transaction-viewed.html", mapping)
}