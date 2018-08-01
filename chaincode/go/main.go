/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type data struct {
	Name           string `json:"name"`
	PointsReceived int    `json:"pointsreceived"` // for history of received points, if you convert it to cash this is where it should be deducted(10pts)
	PointsSent     int    `json:"pointssent"`     // for history of sent points
	PointsCurrent  int    `json:"pointscurrent"`  // points you can give currently, technically should always be either 1 or 0
	Giver          string `json:"giver"`          // Other peers or system
	Message        string `json:"message"`        // Message given upon
	SentTo         string `json:"sentto"`         // Person the point is sent to
	Timestamp      string `json:"timestamp"`      // Timestamp
	Error          int    `json:"error"`          // Bolean error
}

type hist struct {
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initiating Thanks Chaincode")

	var args, args2 []string
	args = append(args, "Donlee Malacad")
	args = append(args, "System")

	t.addPerson(stub, args)

	args2 = append(args2, "Alvin Cadacio")
	args2 = append(args2, "System")

	t.addPerson(stub, args2)

	return shim.Success(nil)
}

// Invoke function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	if function == "transfer" {
		// Transfer a point from A to B
		return t.transfer(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// get current status
		return t.query(stub, args)
	} else if function == "addPerson" {
		// Adding of new person and its value
		return t.addPerson(stub, args)
	} else if function == "getHistoryOfPerson" {
		// Getting complete history of person
		return t.getHistoryOfPerson(stub, args)
	} else if function == "getAllUsers" {
		return t.getAllUsers(stub)
	} else if function == "addPointToAll" {
		return t.addPointToAll(stub)
	}
	return shim.Error("Invalid invoke function name. Expecting \"transfer\" \"delete\" \"query\"")
}

func (t *SimpleChaincode) addPointToAll(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("add Point to all users")

	localTime, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(localTime)
	fmt.Print(now)

	if now.Day() == 1 {
		fmt.Print("first of the month")
		fmt.Print(now.Hour())
		if now.Hour() != 1 {
			fmt.Print("not 1st hour of the 1st of the month")
			return shim.Error("Not time of the month to replenish(not 1st hour of the 1st of the month)")
		}
	} else {
		fmt.Print("not 1st of the month")

		return shim.Error("Not time of the month to replenish(not 1st of the month)")
	}

	startKey := "A"
	endKey := "zzzzzzzzzzzzz"
	usersKey := make(map[int]string)
	x := 0

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	fmt.Println(resultsIterator)
	if err != nil {
		fmt.Printf("ERROR")
		fmt.Print(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			fmt.Printf("ERROR2")
			fmt.Print(err.Error())
			return shim.Error(err.Error())
		}

		usersKey[x] = string(response.GetKey())
		x++
	}
	fmt.Print(usersKey)

	for i := 0; i < len(usersKey); i++ {
		fmt.Print("\n", usersKey[i])
		// add point
		AddPointPersonBytes, err := stub.GetState(usersKey[i])
		fmt.Printf("\nAddPointPersonBytes\n")
		fmt.Print(AddPointPersonBytes)
		if err != nil {
			return shim.Error("\n\nFailed to get state: ToPerson\n\n")
		}
		if AddPointPersonBytes == nil {
			return shim.Error("AddPointPerson Entity not found")
		}
		AddPointPerson := data{}

		err = json.Unmarshal([]byte(AddPointPersonBytes), &AddPointPerson)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("\nAddPointPerson\n")
		fmt.Print(AddPointPerson)
		AddPointPerson.Name = usersKey[i]
		AddPointPerson.PointsCurrent = AddPointPerson.PointsCurrent + 1
		AddPointPerson.Giver = "System"
		AddPointPerson.Message = "System Generated"
		AddPointPerson.SentTo = ""
		AddPointPerson.Timestamp = time.Now().String()

		fmt.Printf("\nAddPointPerson Updated\n")
		fmt.Print(AddPointPerson)
		AddPointPersonJSONasBytes, _ := json.Marshal(AddPointPerson)
		stub.PutState(usersKey[i], AddPointPersonJSONasBytes)
		fmt.Printf("AddPointPersonJSONasBytes")

		fmt.Print(AddPointPersonJSONasBytes)
	}

	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) addPerson(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Adding of Person requires 2 parameters")
	}

	// check if person already exists
	startKey := "A"
	endKey := "zzzzzzzzzzzzzzzzzzz"
	usersKey := make(map[int]string)
	x := 0

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	println(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		usersKey[x] = string(response.GetKey())
		x++
	}
	fmt.Print(usersKey)

	for i := 0; i < len(usersKey); i++ {
		fmt.Print("\n", usersKey[i])
		if usersKey[i] == args[0] {
			fmt.Printf("ERROR: person %s already added, try new name", args[0])
			return shim.Error("ERROR: person already added, try new name")
		}
	}

	Name := args[0]
	PointsReceived := 0
	PointsSent := 0
	// TODO: change 5 to 1, for testing only -done
	PointsCurrent := 1

	Giver := args[1]
	Message := "Welcome the Thanks Application, this is your initial point"
	SentTo := ""
	Timestamp := time.Now().String()
	Error := 0

	fmt.Printf("Added new Person, Name = %s\n", Name)
	fmt.Printf("Points: Received = %d, Sent = %d, Current = %d\n", PointsReceived, PointsSent, PointsCurrent)
	fmt.Printf("\nGiver = %s Message = %s\n", Giver, Message)

	data := &data{Name, PointsReceived, PointsSent, PointsCurrent, Giver, Message, SentTo, Timestamp, Error}
	fmt.Println("data: ", data)
	dataJSONasBytes, err := json.Marshal(data)

	err = stub.PutState(Name, []byte(dataJSONasBytes))
	if err != nil {
		return shim.Error("Failed to execute addPerson function")
	}

	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) getHistoryOfPerson(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	person := args[0]

	fmt.Printf("- start getting history of Person transaction: %s\n", person)

	resultsIterator, err := stub.GetHistoryForKey(person)

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("\n--- Transaction History of %s:---\n%s\n", person, buffer.String())

	return shim.Success(buffer.Bytes())
}

// Transaction makes payment of X units from FromPerson to ToPerson
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Transfer of Point Commencing")
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	time.LoadLocation("Asia/Shanghai")

	ToPerson := args[0]
	FromPerson := args[1]

	// get history for previous point sent checking
	resultsIterator, err := stub.GetHistoryForKey(FromPerson)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	history := make(map[string]map[string]string)

	var dataCheck data

	var timeKey time.Time

	fmt.Println(resultsIterator)
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		timeKey = time.Unix(response.Timestamp.Seconds, 0)

		err = json.Unmarshal([]byte(response.Value), &dataCheck)
		if err != nil {
			return shim.Error(err.Error())
		}

		history = map[string]map[string]string{timeKey.String(): map[string]string{"sentto": dataCheck.SentTo}}

		if timeKey.Year() == time.Now().Year() {
			if timeKey.Month() == time.Now().Month()-1 {
				if history[timeKey.String()]["sentto"] == ToPerson {
					fmt.Print("Same Person sent to last month")
					return shim.Error("Error")
				}
			}
		}

	}

	// Get current transferee details
	ToPersonbytes, err := stub.GetState(ToPerson)

	if err != nil {
		return shim.Error("\n\nFailed to get state: ToPerson\n\n")
	}
	if ToPersonbytes == nil {
		return shim.Error("ToPerson Entity not found")
	}

	// Get current transferer details
	FromPersonbytes, err := stub.GetState(FromPerson)

	if err != nil {
		return shim.Error("\n\nFailed to get state: FromPersonbytes\n\n")
	}
	if FromPersonbytes == nil {
		return shim.Error("\n\nFromPerson Entity not found\n\n")
	}

	TransfererPerson := data{}

	err = json.Unmarshal([]byte(FromPersonbytes), &TransfererPerson)
	if err != nil {
		return shim.Error(err.Error())
	}

	TransferPerson := data{}

	err = json.Unmarshal([]byte(ToPersonbytes), &TransferPerson)
	if err != nil {
		return shim.Error(err.Error())
	}

	if TransfererPerson.PointsCurrent <= 0 {
		TransferPerson.Error = 1
		TransfererPerson.Error = 1

		return shim.Error("\n\nCurrent Point is not enough\n\n")
	}
	TransferPerson.Name = ToPerson
	TransferPerson.PointsReceived = TransferPerson.PointsReceived + 1
	TransferPerson.Giver = FromPerson
	TransferPerson.Message = args[2]
	TransferPerson.SentTo = ""
	TransferPerson.Timestamp = time.Now().String()
	TransferPerson.Error = 0

	TransferPersonJSONasByres, _ := json.Marshal(TransferPerson)
	err = stub.PutState(ToPerson, TransferPersonJSONasByres)

	if err != nil {
		return shim.Error(err.Error())
	}

	TransfererPerson.Name = FromPerson
	TransfererPerson.PointsSent = TransfererPerson.PointsSent + 1
	TransfererPerson.PointsCurrent = TransfererPerson.PointsCurrent - 1
	TransfererPerson.Giver = ""
	message := fmt.Sprintf("Message to %s: %s", ToPerson, args[2])
	TransfererPerson.Message = message
	TransfererPerson.SentTo = ToPerson
	TransfererPerson.Timestamp = time.Now().String()
	TransfererPerson.Error = 0

	TransfererPersonJSONasByres, _ := json.Marshal(TransfererPerson)
	err = stub.PutState(FromPerson, TransfererPersonJSONasByres)

	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("\nTransfer of Point from %s to %s is completed", FromPerson, ToPerson)

	return shim.Success(TransfererPersonJSONasByres)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) getAllUsers(stub shim.ChaincodeStubInterface) pb.Response {
	startKey := "A"
	endKey := "zzzzzzzzzzzz"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- get All Persons query result:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]
	fmt.Println("\nargs: ", args)
	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Get Current State::%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func addPointToAll() {
	fmt.Println("add Point to all users")

	return
}

func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
