package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/*
Welcome to db.go this file handles the reading, writing, update & deleting items
from DynamoDB.

Remember that the AWS Region we are using is "ap-southeast-2" which is Sydney

FUnctions in this file:
- getItem: fetch an item from DynamoDB
- putItem: create and item in DynamoDB

TODO:
- updateItem
- deleteItem
*/

//Declare a new DynamoDB instance. This is safe for concurrent use
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-southeast-2"))

func getItem(employeeid string) (*employee, error) {
	// Prepare the input for query
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Employees"),
		Key: map[string]*dynamodb.AttributeValue{
			"EmployeeID": {
				S: aws.String(employeeid),
			},
		},
	}

	//retrieve item from dynamodb, if no matching item found return nil
	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	//the result.item object returned has underlying type of
	//map[string]*AttributeValue. we can use unmarshalMap helper to
	//parse this straight into fields of a struct

	emp := new(employee)
	err = dynamodbattribute.UnmarshalMap(result.Item, emp)
	if err != nil {
		return nil, err
	}

	return emp, nil
}

//add a book record to DynamoDb
func putItem(emp *employee) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Employees"),
		Item: map[string]*dynamodb.AttributeValue{
			"EmployeeID": {
				S: aws.String(emp.EmployeeID),
			},
			"FirstName": {
				S: aws.String(emp.FirstName),
			},
			"LastName": {
				S: aws.String(emp.LastName),
			},
			"EmployeeType": {
				S: aws.String(emp.EmployeeType),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}
