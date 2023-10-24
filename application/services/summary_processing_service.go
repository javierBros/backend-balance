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

	summary := fmt.Sprintf("Total balance is %.2f\n", totalBalance)
	summary += fmt.Sprintf("Average credit amount: %.2f\n", creditSum/float64(creditCount))
	summary += fmt.Sprintf("Average debit amount: %.2f\n", debitSum/float64(debitCount))

	for month, count := range transactionByMonth {
		summary += fmt.Sprintf("Number of transactions in %s: %d\n", month, count)
	}

	return summary, nil
}

func NewSummaryProcessingService(emailSender notifications.IEmailSender) *SummaryProcessingService {
	return &SummaryProcessingService{
		emailSender: emailSender,
	}
}
