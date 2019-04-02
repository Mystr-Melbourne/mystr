/*
Welcome to the main.go file for our shift-related Lambda functions
This file handles HTTP requests, calling on db.go to read and write data
and then returning a response. The db.go file handles writing to DynamoDB


Functions in this file:

TODO:

TWILIO INTEGRATION this may or may need to be another separate
set of functions

- router: handling http requests
- show: returning data about an individual shift
	and a collected map of individual shifts for a week or day period
- create: creating a new shift
- serverError: log server errors to log
- return client errors as HTTP response
- update probably not required
- delete: delete a shift from the details
*/