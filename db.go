package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//Declare a new DynamoDB instance. This is safe for concurrent use

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-southeast-2"))

func getItem(isbn String) (*book, error) {
	// Prepare the input for query
	input := &dynamodb.GetItemInput{
		Tablename: aws.String("Books"),
		Key: map[string]*dynamodb.AttributeValue{
			"ISBN": {
				S: aws.String(isbn),
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

	bk := new(book)
	err = dynamodbattribute.UnmarshalMap(result.Item, bk)
	if err != nil {
		return nil, err
	}

	return bk, nil
}
