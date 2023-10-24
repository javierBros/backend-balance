package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/javierBros/backend-balance/application/config"
	"github.com/javierBros/backend-balance/application/constants"
	"github.com/javierBros/backend-balance/application/model"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

var envVariables *config.EnvironmentVariables

func downloadFile(sess *session.Session, bucket, key string) ([]byte, error) {
	s3Client := s3.New(sess)
	output, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Println("Error downloading file from S3:", err)
		return nil, err
	}

	if *output.ContentLength > int64(1*1024*1024) {
		log.Println("File exceeds the maximum allowed size.")
		return nil, fmt.Errorf("File size exceeds the maximum allowed size")
	}

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
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

		month := transactions[i].Date.Format(constants.DateFormatForMonth)
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

func sendSummaryEmail(sess *session.Session, summary string) error {
	svc := ses.New(sess)
	_, err := svc.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(envVariables.DestinationEmail)},
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
		Source: aws.String(envVariables.DestinationEmail),
	})

	return err
}

func processCSV(data []byte) ([]model.Transaction, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	transactions := []model.Transaction{}

	for {
		record, err := reader.Read()
		if err != nil {
			fmt.Printf(`err line csv: %v\n`, err.Error())
			break
		}
		fmt.Printf(`record: %v\n`, record)

		if len(record) != 3 {
			continue // Skip invalid rows
		}

		date, _ := time.Parse(constants.DateFormatMMdd, record[1])

		amount, _ := strconv.ParseFloat(strings.TrimPrefix(record[2], "+"), 64)

		transaction := model.Transaction{
			Date:     date,
			Amount:   amount,
			IsCredit: strings.HasPrefix(record[2], "+"),
		}
		fmt.Printf(`transaction: %v\n`, record)

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func ProcessFile(ctx context.Context, event events.S3Event) error {

	envVariables = config.NewEnvironmentVariables()

	fmt.Printf(`S3 event: %v`, event)
	// Initialize AWS SES session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(envVariables.SESRegion),
	})
	if err != nil {
		fmt.Printf(`Error reading session: %v`, err.Error())
		return err
	}

	s3Record := event.Records[0].S3
	bucket := s3Record.Bucket.Name
	key := s3Record.Object.Key

	fmt.Printf(`Bucket: %v, Key: %v\n`, bucket, key)

	data, err := downloadFile(sess, bucket, key)
	if err != nil {
		log.Println("Error downloading file from S3:", err)
		return err
	}

	transactions, err := processCSV(data)
	if err != nil {
		log.Println("Error processing CSV:", err)
		return err
	}

	summary, err := calculateSummary(transactions)
	if err != nil {
		log.Println("Error calculating summary:", err)
		return err
	}

	err = sendSummaryEmail(sess, summary)
	if err != nil {
		fmt.Printf(`Error sending email: %v`, err.Error())
		return err
	}

	return nil
}

func main() {
	lambda.Start(ProcessFile)
}
