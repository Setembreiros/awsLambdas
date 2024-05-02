package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	fmt.Printf("PostConfirmation of user %s with email %s in pool %s\n", event.UserName, event.Request.UserAttributes["email"], event.UserPoolID)

	return event, nil
}

func main() {
	lambda.Start(HandleRequest)
}
