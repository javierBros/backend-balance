package config

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/javierBros/backend-balance/application"
	"github.com/javierBros/backend-balance/application/controller"
	"github.com/javierBros/backend-balance/application/notifications"
	"github.com/javierBros/backend-balance/application/services"
)

func ChargeDependencies(event events.S3Event, sess *session.Session, envVariables *application.EnvironmentVariables) *controller.FileProcessingController {
	//var dbConnection = connection.InitPGDBConnection()

	// notifications
	emailSender := notifications.NewEmailSender(sess, envVariables)

	// Service
	summaryProcessingService := services.NewSummaryProcessingService(emailSender)

	// Controller
	fileProcessingController := controller.NewFileProcessingController(summaryProcessingService, envVariables, event, sess)

	return fileProcessingController
}
