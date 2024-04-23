package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	fmt.Printf("PreSignup of user: %s\n", event.UserName)
	event.Response.AutoConfirmUser = true
	return event, nil
}

func main() {
	lambda.Start(HandleRequest)
}
