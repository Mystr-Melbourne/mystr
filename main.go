package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type book struct {
	ISBN	string	`json:"isbn"`
	Title	string	`json:"title"`
	Author	string	`json:"author"`
}

func show() (*book, error) {
	//fetch a specific book record from dynamodb in this case 
	//marcus aurelius
	bk, err := getItem("978-1292292838"}
	if err != nil {
		return nil, err
	}

	return bk, nil
}

func main() {
	lambda.Start(show)
}
