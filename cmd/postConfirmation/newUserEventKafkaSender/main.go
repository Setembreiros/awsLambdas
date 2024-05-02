package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	fmt.Printf("PostConfirmation of user %s with email %s in pool %s\n", event.UserName, event.Request.UserAttributes["email"], event.UserPoolID)

	_, err := createKafkaEvent(event)
	if err != nil {
		fmt.Println("error creating json:", err.Error())
		return event, err
	}

	return event, nil
}

func createKafkaEvent(event events.CognitoEventUserPoolsPostConfirmation) (string, error) {
	type NewUserRegisteredEvent struct {
		UserId   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		UserType string `json:"user_type"`
		Region   string `json:"region"`
		FullName string `json:"full_name"`
	}

	var userType string
	if event.UserPoolID == "eu-west-3_hOUCCL4yo" {
		userType = "UA"
	} else {
		userType = "UE"
	}

	kafkaEvent := NewUserRegisteredEvent{
		UserId:   event.Request.UserAttributes["sub"],
		Username: event.UserName,
		Email:    event.Request.UserAttributes["email"],
		UserType: userType,
		Region:   event.Request.UserAttributes["custom:region"],
		FullName: event.Request.UserAttributes["name"],
	}
	b, err := json.Marshal(kafkaEvent)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func main() {
	lambda.Start(HandleRequest)
}
