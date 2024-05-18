package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

type App struct {
	client     *cognito.Client
	userPoolId string
	username   string
	email      string
}

func HandleRequest(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	fmt.Printf("PreSignup of user %s with email %s in pool %s\n", event.UserName, event.Request.UserAttributes["email"], event.UserPoolID)

	app, err := createCognitoClient(event)
	if err != nil {
		return event, err
	}

	err = app.validateEmail()
	if err != nil {
		return event, err
	}

	err = app.validateUsername()
	if err != nil {
		return event, err
	}

	return event, nil
}

func createCognitoClient(event events.CognitoEventUserPoolsPreSignup) (*App, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	var userPoolId string
	if event.UserPoolID == "eu-west-3_hOUCCL4yo" {
		userPoolId = "eu-west-3_mplHHKWzh"
	} else {
		userPoolId = "eu-west-3_hOUCCL4yo"
	}

	return &App{
		client:     cognito.NewFromConfig(cfg),
		userPoolId: userPoolId,
		username:   event.UserName,
		email:      event.Request.UserAttributes["email"],
	}, nil
}

func (c *App) validateEmail() error {
	filter := "email = \"" + c.email + "\""
	request := c.createListUserRequest(filter)

	response, err := c.sendListUserRequest(request)
	if err != nil {
		return err
	}

	if len(response.Users) > 0 {
		return errors.New("EXISTING_EMAIL")
	}

	return err
}

func (c *App) validateUsername() error {
	filter := "username = \"" + c.username + "\""
	request := c.createListUserRequest(filter)

	response, err := c.sendListUserRequest(request)
	if err != nil {
		return err
	}

	if len(response.Users) > 0 {
		return errors.New("EXISTING_USERNAME")
	}

	return err
}

func (c *App) createListUserRequest(filter string) *cognito.ListUsersInput {
	return &cognito.ListUsersInput{
		UserPoolId: aws.String(c.userPoolId),
		Filter:     aws.String(filter),
	}
}

func (c *App) sendListUserRequest(request *cognito.ListUsersInput) (*cognito.ListUsersOutput, error) {
	response, err := c.client.ListUsers(context.Background(), request)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}
	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
