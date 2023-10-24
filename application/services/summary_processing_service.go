package services

import (
	"errors"
	"fmt"
	"github.com/javierBros/backend-balance/application"
	"github.com/javierBros/backend-balance/application/model"
	"github.com/javierBros/backend-balance/application/notifications"
)

type SummaryProcessingService struct {
	emailSender notifications.IEmailSender
}

func (d *SummaryProcessingService) ProcessSummary(transactions []model.Transaction) error {
	summary, err := calculateSummary(transactions)
	if err != nil {
		return err
	}
	err = d.emailSender.SendSummaryEmail(summary)
	if err != nil {
		return err
	}
	return nil
}

func calculateSummary(transactions []model.Transaction) (string, error) {
	if len(transactions) == 0 {
		return "", errors.New("No transactions to calculate")
	}

	totalBalance := 0.0
	creditSum := 0.0
	creditCount := 0
	debitSum := 0.0
	debitCount := 0
	transactionByMonth := make(map[string]int)

	for i := 0; i < len(transactions); i++ {
		totalBalance += transactions[i].Amount

		month := transactions[i].Date.Format(application.DateFormatForMonth)
		transactionByMonth[month]++
		if transactions[i].IsCredit {
			creditSum += transactions[i].Amount
			creditCount++
		} else {
			debitSum += transactions[i].Amount
			debitCount++
		}
	}

	emailBody := `
        <html>
            <head></head>
            <body>
                <h1>Transaction Summary</h1>
                <p>Total balance is %.2f</p>
                <p>Average credit amount: %.2f</p>
                <p>Average debit amount: %.2f</p>
	`
	for month, count := range transactionByMonth {
		emailBody += fmt.Sprintf("<p>Number of transactions in %s: %d</p>", month, count)
	}

	emailBody += `<footer>
					<p>Stori Company</p>
					<img src="https://stori-resources.s3.amazonaws.com/stori_logo.png" alt="Company Logo" width="150" height="50">
				</footer>
            </body>
        </html>
    `
	emailBody = fmt.Sprintf(emailBody, totalBalance, creditSum/float64(creditCount), debitSum/float64(debitCount))

	fmt.Println(emailBody)

	return emailBody, nil
}

func NewSummaryProcessingService(emailSender notifications.IEmailSender) *SummaryProcessingService {
	return &SummaryProcessingService{
		emailSender: emailSender,
	}
}
