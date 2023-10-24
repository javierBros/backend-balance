package notifications

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/javierBros/backend-balance/application"
)

type EmailSender struct {
	session      client.ConfigProvider
	envVariables *application.EnvironmentVariables
}

func (d *EmailSender) SendSummaryEmail(summary string) error {
	svc := ses.New(d.session)
	_, err := svc.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(d.envVariables.DestinationEmail)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(summary),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Transaction Summary"),
			},
		},
		Source: aws.String(d.envVariables.DestinationEmail),
	})

	return err
}

func NewEmailSender(session client.ConfigProvider,
	envVariables *application.EnvironmentVariables) *EmailSender {
	return &EmailSender{
		session:      session,
		envVariables: envVariables,
	}
}
