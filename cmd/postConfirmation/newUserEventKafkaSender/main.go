package main

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type App struct {
	producer sarama.SyncProducer
	topic    string
}

func HandleRequest(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	fmt.Printf("PostConfirmation of user %s with email %s in pool %s\n", event.UserName, event.Request.UserAttributes["email"], event.UserPoolID)

	kafkaEvent, err := createKafkaEvent(event)
	if err != nil {
		return event, err
	}

	app, err := createKafkaProducer()
	if err != nil {
		return event, err
	}

	err = app.publish(kafkaEvent)
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		return event, err
	}

	return event, nil
}

func createKafkaEvent(event events.CognitoEventUserPoolsPostConfirmation) (string, error) {
	type NewRegisteredUserEvent struct {
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

	kafkaEvent := NewRegisteredUserEvent{
		UserId:   event.Request.UserAttributes["sub"],
		Username: event.UserName,
		Email:    event.Request.UserAttributes["email"],
		UserType: userType,
		Region:   event.Request.UserAttributes["custom:region"],
		FullName: event.Request.UserAttributes["name"],
	}
	b, err := json.Marshal(kafkaEvent)
	if err != nil {
		fmt.Println("error creating json:", err.Error())
		return "", err
	}

	return string(b), nil
}

func createKafkaProducer() (*App, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	kafkaConn := "172.31.36.175:9092"

	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		return nil, err
	}

	return &App{
		producer: prd,
		topic:    "UserWasRegisteredEvent",
	}, nil
}

func (app *App) publish(event string) error {
	msg := &sarama.ProducerMessage{
		Topic: app.topic,
		Value: sarama.StringEncoder(event),
	}
	p, o, err := app.producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		return err
	}

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
