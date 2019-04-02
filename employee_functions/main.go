package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

/*
Welcome to the main.go file for our employee-related Lambda functions
This file handles HTTP requests, calling on db.go to read and write data
and then returning a response. The db.go file handles writing to DynamoDB


Functions in this file:
- router: handling http requests
- show: returning data about an employee
- create: creating a new emplouee
- serverError: log server errors to log
- return client errors as HTTP response

TODO:
- update: update employee details
- delete: delete and employee from the database
*/

//regex function to check EmployeeID is valid
var employeeidRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)

//log errors to standard error output
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

//Employee struct holds the data about the employee, we may need to add more later
//alternative we may remove EmployeeType, in which case db.go will have to be modified
type employee struct {
	EmployeeID   string `json:"employeeid"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	EmployeeType string `json:"employeetype"`
}

//Router the request to either fetch data "show" or input a new employee "create"
//TODO: add PUT for update and DELETE function
func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

//Show is our "fetch" function to show data on the requested employee
func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the `employeeid` query string parameter from the request and
	// validate it.
	employeeid := req.QueryStringParameters["employeeid"]
	if !employeeidRegexp.MatchString(employeeid) {
		return clientError(http.StatusBadRequest)
	} //fetch a specific employee record from dynamodb in this case

	// Fetch the employee record from the database based on the employeeid value.
	emp, err := getItem(employeeid)
	if err != nil {
		return serverError(err)
	}
	if emp == nil {
		return clientError(http.StatusNotFound)
	}

	// The APIGatewayProxyResponse.Body field needs to be a string, so
	// we marshal the employee record into JSON.
	js, err := json.Marshal(emp)
	if err != nil {
		return serverError(err)
	}

	// Return a response with a 200 OK status and the JSON employee record
	// as the body.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	emp := new(employee)
	err := json.Unmarshal([]byte(req.Body), emp)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	//check employeeid is valid number range
	if !employeeidRegexp.MatchString(emp.EmployeeID) {
		return clientError(http.StatusBadRequest)
	}
	//return bad request error if fields empty
	if emp.FirstName == "" || emp.LastName == "" || emp.EmployeeType == "" {
		return clientError(http.StatusBadRequest)
	}
	//creat the employee item and return 201 "created" response
	err = putItem(emp)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/employees?employeeid=%s", emp.EmployeeID)},
	}, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
