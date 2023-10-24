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

/*func processCSV(data []byte) ([]model.Transaction, error) {
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

func downloadAndProcessFile(event events.S3Event, sess *session.Session) ([]model.Transaction, error) {

	s3Record := event.Records[0].S3
	bucket := s3Record.Bucket.Name
	key := s3Record.Object.Key

	fmt.Printf(`Bucket: %v, Key: %v\n`, bucket, key)

	data, err := downloadFile(sess, bucket, key)
	if err != nil {
		log.Println("Error downloading file from S3:", err)
		return nil, err
	}

	transactions, err := processCSV(data)
	if err != nil {
		log.Println("Error processiyng CSV:", err)
		return nil, err
	}

	return transactions, nil
}*/

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
	err = fileProcessingController.DownloadAndProcessFile()
	if err != nil {
		return err
	}
	/*
		fileProcessingController := config.ChargeDependencies(sess, envVariables)

		err = fileProcessingController.DownloadAndProcessFile()
		if err != nil {
			return err
		}*/

	return nil
}

func main() {
	lambda.Start(ProcessEvent)
}
