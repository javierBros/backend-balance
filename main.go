package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/javierBros/backend-balance/application"
	"github.com/javierBros/backend-balance/application/config"
)

var envVariables *application.EnvironmentVariables

func ProcessEvent(ctx context.Context, event events.S3Event) error {

	envVariables = application.NewEnvironmentVariables()

	fmt.Printf(`S3 event: %v`, event)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(envVariables.SESRegion),
	})
	if err != nil {
		fmt.Printf(`Error reading session: %v`, err.Error())
		return err
	}

	fileProcessingController := config.ChargeDependencies(event, sess, envVariables)
	err = fileProcessingController.ProcessBalance()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(ProcessEvent)
}
